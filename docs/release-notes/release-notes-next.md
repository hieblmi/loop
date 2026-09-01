# Loop Client Release Notes

#### New Features

#### Breaking Changes

#### Bug Fixes

* Reject static Loop In server responses with missing or extra HTLC nonces
  before signing, preventing a client crash from an incomplete nonce list.

* A static Loop In that releases its deposits a second time after its payment
  deadline, at its HTLC expiry or when the server rejects the swap, now only
  releases deposits that are still looping in for it. Before, it could move a
  deposit that another Loop In, a withdrawal or a channel open had taken in the
  meantime back to `Deposited`, and retried forever once the deposit was spent.
  A static Loop In also no longer locks a deposit that another Loop In is
  already using, and a failed lock no longer resets deposits that a withdrawal
  or a channel open has taken.

* The static address Loop In timeout sweep now spends the exact HTLC output
  confirmed on chain, including server-published backup fee variants.
  Previously it always targeted the standard fee variant, so a confirmed high
  fee or extremely high fee HTLC transaction could not be swept back after the
  HTLC timed out. Sweeps that would leave a dust output are no longer signed.

* Preserve pending loop-in ownership of deposits during expiry and restart
  recovery. Release only deposits still owned by the swap, after a positive
  unspent-input or confirmed-conflict check rules out a confirmed HTLC.
  Recover abandoned initialization locks at startup, retry temporary wallet
  lookup failures, and finish timed-out swaps after their deposits are reused.
  Check live swap claims before releasing deposits with stale legacy mappings.
  New swaps retain unconfirmed funding transactions so a wallet-visible mined
  replacement can prove that the original HTLC cannot confirm after restart.
  Missing history alone does not release deposits; older swaps without a funding
  snapshot continue monitoring for recovery.

* Follow replacement static deposit expiry sweeps after a chain reorganization
  instead of waiting for the old transaction until loopd restarts.

* Static Address loop-in sweep requests are now rejected unless the server
  provides exactly one prevout for every sweep input. A duplicate or missing
  prevout previously crashed `loopd` while computing the sweep signature
  hashes.

#### Maintenance

#### Contributors (Alphabetical Order)
