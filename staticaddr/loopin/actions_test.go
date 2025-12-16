package loopin

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/lndclient"
	"github.com/lightninglabs/loop/fsm"
	"github.com/lightninglabs/loop/staticaddr/address"
	"github.com/lightninglabs/loop/staticaddr/deposit"
	"github.com/lightninglabs/loop/staticaddr/script"
	"github.com/lightninglabs/loop/staticaddr/version"
	"github.com/lightninglabs/loop/test"
	"github.com/lightningnetwork/lnd/chainntnfs"
	"github.com/lightningnetwork/lnd/invoices"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/stretchr/testify/require"
)

// TestMonitorInvoiceAndHtlcTxReRegistersOnConfErr ensures that an error from
// the HTLC confirmation subscription triggers a re-registration. Without the
// regression fix, only the initial registration would be performed and the
// test would time out waiting for the second one.
func TestMonitorInvoiceAndHtlcTxReRegistersOnConfErr(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	mockLnd := test.NewMockLnd()

	clientKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	serverKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	swapHash := lntypes.Hash{1, 2, 3}

	loopIn := &StaticAddressLoopIn{
		SwapHash:              swapHash,
		HtlcCltvExpiry:        2_000,
		InitiationHeight:      uint32(mockLnd.Height),
		InitiationTime:        time.Now(),
		ProtocolVersion:       version.ProtocolVersion_V0,
		ClientPubkey:          clientKey.PubKey(),
		ServerPubkey:          serverKey.PubKey(),
		PaymentTimeoutSeconds: 3_600,
	}
	loopIn.SetState(MonitorInvoiceAndHtlcTx)

	// Seed the mock invoice store so LookupInvoice succeeds.
	mockLnd.Invoices[swapHash] = &lndclient.Invoice{
		Hash:  swapHash,
		State: invoices.ContractOpen,
	}

	cfg := &Config{
		AddressManager: &mockAddressManager{
			params: &address.Parameters{
				ClientPubkey:    clientKey.PubKey(),
				ServerPubkey:    serverKey.PubKey(),
				ProtocolVersion: version.ProtocolVersion_V0,
			},
		},
		ChainNotifier:  mockLnd.ChainNotifier,
		DepositManager: &noopDepositManager{},
		InvoicesClient: mockLnd.LndServices.Invoices,
		LndClient:      mockLnd.Client,
		ChainParams:    mockLnd.ChainParams,
	}

	f, err := NewFSM(ctx, loopIn, cfg, false)
	require.NoError(t, err)

	resultChan := make(chan fsm.EventType, 1)
	go func() {
		resultChan <- f.MonitorInvoiceAndHtlcTxAction(ctx, nil)
	}()

	// Capture the invoice subscription the action registers so we can feed
	// an update later and let the action exit.
	var invSub *test.SingleInvoiceSubscription
	select {
	case invSub = <-mockLnd.SingleInvoiceSubcribeChannel:
	case <-ctx.Done():
		t.Fatalf("invoice subscription not registered: %v", ctx.Err())
	}

	// The first confirmation registration should happen immediately.
	var firstReg *test.ConfRegistration
	select {
	case firstReg = <-mockLnd.RegisterConfChannel:
	case <-ctx.Done():
		t.Fatalf("htlc conf registration not received: %v", ctx.Err())
	}

	// Force the confirmation stream to error so the FSM re-registers.
	firstReg.ErrChan <- errors.New("test htlc conf error")

	// FSM registers again, otherwise it would time out.
	var secondReg *test.ConfRegistration
	select {
	case secondReg = <-mockLnd.RegisterConfChannel:
	case <-ctx.Done():
		t.Fatalf("htlc conf was not re-registered: %v", ctx.Err())
	}

	require.NotEqual(t, firstReg, secondReg)

	// Settle the invoice to let the action exit.
	invSub.Update <- lndclient.InvoiceUpdate{
		Invoice: lndclient.Invoice{
			Hash:  swapHash,
			State: invoices.ContractSettled,
		},
	}

	select {
	case event := <-resultChan:
		require.Equal(t, OnPaymentReceived, event)
	case <-ctx.Done():
		t.Fatalf("fsm did not return: %v", ctx.Err())
	}
}

// TestMonitorInvoiceAndHtlcTxInvoiceErr asserts that invoice subscription
// errors surface through the FSM as OnError so we don't silently hang.
func TestMonitorInvoiceAndHtlcTxInvoiceErr(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	mockLnd := test.NewMockLnd()

	clientKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	serverKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	swapHash := lntypes.Hash{9, 9, 9}

	loopIn := &StaticAddressLoopIn{
		SwapHash:              swapHash,
		HtlcCltvExpiry:        3_000,
		InitiationHeight:      uint32(mockLnd.Height),
		InitiationTime:        time.Now(),
		ProtocolVersion:       version.ProtocolVersion_V0,
		ClientPubkey:          clientKey.PubKey(),
		ServerPubkey:          serverKey.PubKey(),
		PaymentTimeoutSeconds: 3_600,
	}
	loopIn.SetState(MonitorInvoiceAndHtlcTx)

	mockLnd.Invoices[swapHash] = &lndclient.Invoice{
		Hash:  swapHash,
		State: invoices.ContractOpen,
	}

	cfg := &Config{
		AddressManager: &mockAddressManager{
			params: &address.Parameters{
				ClientPubkey:    clientKey.PubKey(),
				ServerPubkey:    serverKey.PubKey(),
				ProtocolVersion: version.ProtocolVersion_V0,
			},
		},
		ChainNotifier:  mockLnd.ChainNotifier,
		DepositManager: &noopDepositManager{},
		InvoicesClient: mockLnd.LndServices.Invoices,
		LndClient:      mockLnd.Client,
		ChainParams:    mockLnd.ChainParams,
	}

	f, err := NewFSM(ctx, loopIn, cfg, false)
	require.NoError(t, err)

	resultChan := make(chan fsm.EventType, 1)
	go func() {
		resultChan <- f.MonitorInvoiceAndHtlcTxAction(ctx, nil)
	}()

	var invSub *test.SingleInvoiceSubscription
	select {
	case invSub = <-mockLnd.SingleInvoiceSubcribeChannel:
	case <-time.After(time.Second):
		t.Fatalf("invoice subscription not registered")
	}

	// Drain the initial HTLC confirmation registration so it doesn't block.
	select {
	case <-mockLnd.RegisterConfChannel:
	case <-time.After(time.Second):
		t.Fatalf("htlc conf registration not received")
	}

	// Inject an invoice subscription error and expect the FSM to surface an
	// error event instead of silently logging it.
	invSub.Err <- errors.New("test invoice sub error")

	select {
	case event := <-resultChan:
		require.Equal(t, fsm.OnError, event)
	case <-time.After(time.Second):
		t.Fatalf("expected invoice error to surface, got timeout")
	}
}

// TestMonitorInvoiceAndHtlcTxStaleConfirmation verifies we clear a stale
// htlcConfirmed flag across a re-registration after a confirmation stream
// error. Without the fix we would think the HTLC is confirmed and sweep on
// timeout even though the confirmation was lost (e.g. due to a reorg).
func TestMonitorInvoiceAndHtlcTxStaleConfirmation(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	mockLnd := test.NewMockLnd()
	// Start below the HTLC expiry so the first block notification doesn't
	// trigger timeout immediately.
	mockLnd.Height = 1

	clientKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	serverKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)

	swapHash := lntypes.Hash{7, 7, 7}

	loopIn := &StaticAddressLoopIn{
		SwapHash:              swapHash,
		HtlcCltvExpiry:        20,
		InitiationHeight:      uint32(mockLnd.Height),
		InitiationTime:        time.Now(),
		ProtocolVersion:       version.ProtocolVersion_V0,
		ClientPubkey:          clientKey.PubKey(),
		ServerPubkey:          serverKey.PubKey(),
		PaymentTimeoutSeconds: 3_600,
	}
	loopIn.SetState(MonitorInvoiceAndHtlcTx)

	mockLnd.Invoices[swapHash] = &lndclient.Invoice{
		Hash:  swapHash,
		State: invoices.ContractOpen,
	}

	cfg := &Config{
		AddressManager: &mockAddressManager{
			params: &address.Parameters{
				ClientPubkey:    clientKey.PubKey(),
				ServerPubkey:    serverKey.PubKey(),
				ProtocolVersion: version.ProtocolVersion_V0,
			},
		},
		ChainNotifier:  mockLnd.ChainNotifier,
		DepositManager: &noopDepositManager{},
		InvoicesClient: mockLnd.LndServices.Invoices,
		LndClient:      mockLnd.Client,
		ChainParams:    mockLnd.ChainParams,
	}

	f, err := NewFSM(ctx, loopIn, cfg, false)
	require.NoError(t, err)

	resultChan := make(chan fsm.EventType, 1)
	go func() {
		resultChan <- f.MonitorInvoiceAndHtlcTxAction(ctx, nil)
	}()

	// Wait for invoice and confirmation subscriptions to be registered.
	select {
	case <-mockLnd.SingleInvoiceSubcribeChannel:
	case <-ctx.Done():
		t.Fatalf("invoice subscription not registered: %v", ctx.Err())
	}

	var firstReg *test.ConfRegistration
	select {
	case firstReg = <-mockLnd.RegisterConfChannel:
	case <-ctx.Done():
		t.Fatalf("htlc conf registration not received: %v", ctx.Err())
	}

	// Deliver a confirmation directly on the registration channel to set
	// htlcConfirmed = true.
	firstReg.ConfChan <- &chainntnfs.TxConfirmation{
		Tx: &wire.MsgTx{
			TxOut: []*wire.TxOut{
				{
					PkScript: firstReg.PkScript,
				},
			},
		},
	}

	// Ensure the confirmation was consumed by the action before injecting
	// an error.
	require.Eventually(t, func() bool {
		return len(firstReg.ConfChan) == 0
	}, time.Second, time.Millisecond*10)

	// Error the conf stream to force a re-registration; htlcConfirmed must
	// be cleared, otherwise we'd wrongly sweep on timeout.
	firstReg.ErrChan <- errors.New("conf stream error")

	var secondReg *test.ConfRegistration
	select {
	case secondReg = <-mockLnd.RegisterConfChannel:
	case <-ctx.Done():
		t.Fatalf("htlc conf was not re-registered: %v", ctx.Err())
	}
	require.NotEqual(t, firstReg, secondReg)

	// Advance chain past the HTLC expiry. With stale htlcConfirmed this
	// would take the sweep branch; correct behavior is to time out.
	require.NoError(t, mockLnd.NotifyHeight(loopIn.HtlcCltvExpiry+1))

	select {
	case event := <-resultChan:
		require.Equal(t, OnSwapTimedOut, event)
	case <-time.After(time.Second):
		t.Fatalf("expected swap timeout due to stale conf, got timeout")
	}
}

// mockAddressManager is a minimal AddressManager implementation used by the
// test FSM setup.
type mockAddressManager struct {
	params *address.Parameters
}

// GetStaticAddressParameters returns the configured address parameters.
func (m *mockAddressManager) GetStaticAddressParameters(_ context.Context) (
	*address.Parameters, error) {

	return m.params, nil
}

// GetStaticAddress is unused for this test and returns nil.
func (m *mockAddressManager) GetStaticAddress(_ context.Context) (
	*script.StaticAddress, error) {

	return nil, nil
}

// noopDepositManager is a stub DepositManager used to satisfy FSM config.
type noopDepositManager struct{}

// GetAllDeposits implements DepositManager with a no-op.
func (n *noopDepositManager) GetAllDeposits(_ context.Context) (
	[]*deposit.Deposit, error) {

	return nil, nil
}

// AllStringOutpointsActiveDeposits implements DepositManager with a no-op.
func (n *noopDepositManager) AllStringOutpointsActiveDeposits(
	_ []string, _ fsm.StateType) ([]*deposit.Deposit, bool) {

	return nil, false
}

// TransitionDeposits implements DepositManager with a no-op.
func (n *noopDepositManager) TransitionDeposits(context.Context,
	[]*deposit.Deposit, fsm.EventType, fsm.StateType) error {

	return nil
}

// DepositsForOutpoints implements DepositManager with a no-op.
func (n *noopDepositManager) DepositsForOutpoints(context.Context, []string,
	bool) ([]*deposit.Deposit, error) {

	return nil, nil
}

// GetActiveDepositsInState implements DepositManager with a no-op.
func (n *noopDepositManager) GetActiveDepositsInState(fsm.StateType) (
	[]*deposit.Deposit, error) {

	return nil, nil
}
