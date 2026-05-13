package deposit

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/lndclient"
	"github.com/lightninglabs/loop/staticaddr/address"
	"github.com/lightninglabs/loop/staticaddr/script"
	"github.com/lightninglabs/loop/staticaddr/version"
	"github.com/lightninglabs/loop/swap"
	"github.com/lightninglabs/loop/test"
	"github.com/lightningnetwork/lnd/chainntnfs"
	"github.com/lightningnetwork/lnd/input"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRecoveryKeyFamiliesDedupes(t *testing.T) {
	require.Equal(t, []keychain.KeyFamily{
		keychain.KeyFamily(swap.StaticMultiAddressKeyFamily),
		keychain.KeyFamily(swap.StaticAddressChangeKeyFamily),
	}, recoveryKeyFamilies(
		keychain.KeyFamily(swap.StaticMultiAddressKeyFamily),
	))
}

func TestFindRecoveryAddressParamsMatchesFamilies(t *testing.T) {
	ctx := context.Background()

	for _, family := range []int32{
		swap.StaticAddressKeyFamily,
		swap.StaticMultiAddressKeyFamily,
		swap.StaticAddressChangeKeyFamily,
	} {
		t.Run(fmt.Sprintf("%d", family), func(t *testing.T) {
			wallet := &familyAwareWalletKit{}
			seed := recoveryAddressParams(
				t, wallet, swap.StaticAddressKeyFamily, 0,
			)
			target := recoveryAddressParams(t, wallet, family, 3)

			addrMgr := new(mockAddressManager)
			addrMgr.On("GetParameters", target.PkScript).Return(nil)
			addrMgr.On("GetStaticAddressParameters", mock.Anything).
				Return(seed, nil)

			manager := NewManager(&ManagerConfig{
				AddressManager: addrMgr,
				WalletKit:      wallet,
			})

			params, err := manager.findRecoveryAddressParams(
				ctx, &RecoveryRequest{
					PkScript:   target.PkScript,
					HeightHint: 50,
					ScanLimit:  5,
				},
			)
			require.NoError(t, err)
			require.Equal(t, target.PkScript, params.PkScript)
			require.EqualValues(t, family, params.KeyLocator.Family)
			require.EqualValues(t, 3, params.KeyLocator.Index)
		})
	}
}

func TestRecoverDepositCreatesNewDeposit(t *testing.T) {
	ctx := context.Background()
	manager, req, store := newRecoveryManagerTest(t, nil, nil)

	result, err := manager.RecoverDeposit(ctx, req)
	require.NoError(t, err)
	require.True(t, result.RecoveredDeposit)
	require.True(t, result.RecoveredAddress)
	require.Equal(t, req.TxID.String()+":1", result.OutPoint.String())

	require.NotNil(t, store.created)
	require.Equal(t, req.TxID, store.created.Hash)
	require.EqualValues(t, 1, store.created.Index)
	require.EqualValues(t, 77, store.created.AddressID)
	require.NotNil(t, store.created.AddressParams)
	require.Equal(t, req.PkScript, store.created.AddressParams.PkScript)
	require.True(t, store.created.IsInState(Deposited))
}

func TestRecoverDepositReactivatesReplacedDeposit(t *testing.T) {
	ctx := context.Background()
	existing := &Deposit{
		state:                Replaced,
		Value:                125_000,
		TimeOutSweepPkScript: []byte{0x51},
	}
	manager, req, store := newRecoveryManagerTest(t, existing, nil)
	existing.OutPoint = wire.OutPoint{
		Hash:  req.TxID,
		Index: req.VOut,
	}

	result, err := manager.RecoverDeposit(ctx, req)
	require.NoError(t, err)
	require.True(t, result.RecoveredDeposit)
	require.Same(t, existing, store.updated)
	require.EqualValues(t, 77, existing.AddressID)
	require.Equal(t, req.PkScript, existing.AddressParams.PkScript)
	require.True(t, existing.IsInState(Deposited))
}

func TestRecoverDepositMismatchedExistingValueErrors(t *testing.T) {
	ctx := context.Background()
	existing := &Deposit{
		state: Replaced,
		Value: 1,
	}
	manager, req, _ := newRecoveryManagerTest(t, existing, nil)
	existing.OutPoint = wire.OutPoint{
		Hash:  req.TxID,
		Index: req.VOut,
	}

	_, err := manager.RecoverDeposit(ctx, req)
	require.ErrorContains(t, err, "has value")
}

func TestRecoverDepositAlreadySpentErrors(t *testing.T) {
	ctx := context.Background()
	spendTx := wire.NewMsgTx(2)
	manager, req, _ := newRecoveryManagerTest(t, nil, spendTx)
	spendTx.AddTxIn(&wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  req.TxID,
			Index: req.VOut,
		},
	})

	_, err := manager.RecoverDeposit(ctx, req)
	require.ErrorContains(t, err, "already spent")
}

type recoveryStore struct {
	existing *Deposit
	created  *Deposit
	updated  *Deposit
}

func (s *recoveryStore) CreateDeposit(_ context.Context, d *Deposit) error {
	s.created = d
	return nil
}

func (s *recoveryStore) UpdateDeposit(context.Context, *Deposit) error {
	return nil
}

func (s *recoveryStore) UpdateRecoveredDeposit(_ context.Context,
	d *Deposit) error {

	s.updated = d
	return nil
}

func (s *recoveryStore) GetDeposit(context.Context, ID) (*Deposit, error) {
	return nil, ErrDepositNotFound
}

func (s *recoveryStore) DepositForOutpoint(context.Context,
	string) (*Deposit, error) {

	if s.existing == nil {
		return nil, ErrDepositNotFound
	}

	return s.existing, nil
}

func (s *recoveryStore) AllDeposits(context.Context) ([]*Deposit, error) {
	return nil, nil
}

func newRecoveryManagerTest(t *testing.T, existing *Deposit,
	spendTx *wire.MsgTx) (*Manager, *RecoveryRequest, *recoveryStore) {

	t.Helper()

	const (
		height     = int32(123)
		targetVout = uint32(1)
	)

	lnd := test.NewMockLnd()
	wallet := &familyAwareWalletKit{WalletKitClient: lnd.WalletKit}

	seed := recoveryAddressParams(
		t, wallet, swap.StaticAddressKeyFamily, 0,
	)
	target := recoveryAddressParams(
		t, wallet, swap.StaticMultiAddressKeyFamily, 2,
	)

	tx := wire.NewMsgTx(2)
	tx.AddTxOut(&wire.TxOut{Value: 1, PkScript: []byte{0x51}})
	tx.AddTxOut(&wire.TxOut{
		Value:    125_000,
		PkScript: target.PkScript,
	})
	txid := tx.TxHash()
	req := &RecoveryRequest{
		TxID:       txid,
		VOut:       targetVout,
		HeightHint: height - 10,
		PkScript:   target.PkScript,
		ScanLimit:  5,
	}

	confChan := make(chan *chainntnfs.TxConfirmation, 1)
	confErrChan := make(chan error, 1)
	confChan <- &chainntnfs.TxConfirmation{
		BlockHeight: uint32(height),
		Tx:          tx,
		Block:       blockWithTx(tx),
	}

	chainNotifier := new(MockChainNotifier)
	chainNotifier.On(
		"RegisterConfirmationsNtfn", mock.Anything, mock.Anything,
		mock.Anything, int32(1), req.HeightHint,
	).Return(confChan, confErrChan, nil)

	chainKit := new(MockChainKit)
	bestHeight := height
	if spendTx != nil {
		bestHeight = height + 1
		chainKit.On("GetBlockHash").Return(chainhash.Hash{9}, nil).Once()
		chainKit.On("GetBlock").Return(blockWithTx(spendTx), nil).Once()
	}
	chainKit.On("GetBestBlock", mock.Anything).Return(
		chainhash.Hash{}, bestHeight, nil,
	)

	staticAddress, err := target.TaprootAddress(&chaincfg.TestNet3Params)
	require.NoError(t, err)

	addrMgr := new(mockAddressManager)
	addrMgr.On("GetParameters", target.PkScript).Return(nil)
	addrMgr.On("GetStaticAddressParameters", mock.Anything).Return(seed, nil)
	addrMgr.On("RestoreAddress", mock.Anything).Run(
		func(args mock.Arguments) {
			params := args.Get(0).(*address.Parameters)
			params.ID = 77
		},
	).Return(staticAddress, true, nil)

	store := &recoveryStore{existing: existing}
	manager := NewManager(&ManagerConfig{
		AddressManager: addrMgr,
		ChainKit:       chainKit,
		Store:          store,
		WalletKit:      wallet,
		ChainNotifier:  chainNotifier,
		Signer:         lnd.Signer,
	})

	return manager, req, store
}

func blockWithTx(tx *wire.MsgTx) *wire.MsgBlock {
	return &wire.MsgBlock{
		Transactions: []*wire.MsgTx{tx},
	}
}

type familyAwareWalletKit struct {
	lndclient.WalletKitClient
}

func (w *familyAwareWalletKit) DeriveKey(_ context.Context,
	locator *keychain.KeyLocator) (*keychain.KeyDescriptor, error) {

	_, pubKey := familyAwareKey(locator.Family, locator.Index)
	return &keychain.KeyDescriptor{
		KeyLocator: *locator,
		PubKey:     pubKey,
	}, nil
}

func (w *familyAwareWalletKit) NextAddr(context.Context, string,
	walletrpc.AddressType, bool) (btcutil.Address, error) {

	return btcutil.NewAddressWitnessPubKeyHash(
		make([]byte, 20), &chaincfg.TestNet3Params,
	)
}

func familyAwareKey(family keychain.KeyFamily,
	index uint32) (*btcec.PrivateKey, *btcec.PublicKey) {

	var key [32]byte
	binary.BigEndian.PutUint32(key[24:], uint32(family))
	binary.BigEndian.PutUint32(key[28:], index+1)

	return btcec.PrivKeyFromBytes(key[:])
}

func recoveryAddressParams(t *testing.T, wallet *familyAwareWalletKit,
	family int32, index uint32) *address.Parameters {

	t.Helper()

	locator := keychain.KeyLocator{
		Family: keychain.KeyFamily(family),
		Index:  index,
	}
	clientKey, err := wallet.DeriveKey(context.Background(), &locator)
	require.NoError(t, err)

	staticAddress, err := script.NewStaticAddress(
		input.MuSig2Version100RC2, int64(defaultExpiry),
		clientKey.PubKey, defaultServerPubkey,
	)
	require.NoError(t, err)

	pkScript, err := staticAddress.StaticAddressScript()
	require.NoError(t, err)

	return &address.Parameters{
		ClientPubkey:     clientKey.PubKey,
		ServerPubkey:     defaultServerPubkey,
		Expiry:           defaultExpiry,
		PkScript:         pkScript,
		KeyLocator:       locator,
		ProtocolVersion:  version.ProtocolVersion_V0,
		InitiationHeight: 1,
	}
}

var _ lndclient.WalletKitClient = (*familyAwareWalletKit)(nil)
