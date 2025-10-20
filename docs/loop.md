## CLI interface - loop

control plane for your loopd.

Usage:

```bash
$ loop [GLOBAL FLAGS] [COMMAND] [COMMAND FLAGS] [ARGUMENTS...]
```

Global flags:

| Name                   | Description                                               | Type   |          Default value          |  Environment variables |
|------------------------|-----------------------------------------------------------|--------|:-------------------------------:|:----------------------:|
| `--rpcserver="…"`      | loopd daemon address host:port                            | string |        `localhost:11010`        |  `LOOPCLI_RPCSERVER`   |
| `--network="…"` (`-n`) | the network loop is running on e.g. mainnet, testnet, etc | string |            `mainnet`            |   `LOOPCLI_NETWORK`    |
| `--loopdir="…"`        | path to loop's base directory                             | string |            `~/.loop`            |   `LOOPCLI_LOOPDIR`    |
| `--tlscertpath="…"`    | path to loop's TLS certificate                            | string |   `~/.loop/mainnet/tls.cert`    |  `LOOPCLI_TLSCERTPATH` |
| `--macaroonpath="…"`   | path to macaroon file                                     | string | `~/.loop/mainnet/loop.macaroon` | `LOOPCLI_MACAROONPATH` |
| `--help` (`-h`)        | show help                                                 | bool   |             `false`             |         *none*         |
| `--version` (`-v`)     | print the version                                         | bool   |             `false`             |         *none*         |

### `out` command

perform an off-chain to on-chain swap (looping out).

Attempts to loop out the target amount into either the backing lnd's 	wallet, or a targeted address.  	The amount is to be specified in satoshis.  	Optionally a BASE58/bech32 encoded bitcoin destination address may be 	specified. If not specified, a new wallet address will be generated.

Usage:

```bash
$ loop [GLOBAL FLAGS] out [COMMAND FLAGS] amt [addr]
```

The following flags are supported:

| Name                         | Description                                                                                                                                                                                                                                                                    | Type     | Default value | Environment variables |
|------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|:-------------:|:---------------------:|
| `--addr="…"`                 | the optional address that the looped out funds should be sent to, if let blank the funds will go to lnd's wallet                                                                                                                                                               | string   |               |        *none*         |
| `--account="…"`              | the name of the account to generate a new address from. You can list the names of valid accounts in your backing lnd instance with "lncli wallet accounts list"                                                                                                                | string   |               |        *none*         |
| `--account_addr_type="…"`    | the address type of the extended public key specified in account. Currently only pay-to-taproot-pubkey(p2tr) is supported                                                                                                                                                      | string   |    `p2tr`     |        *none*         |
| `--amt="…"`                  | the amount in satoshis to loop out. To check for the minimum and maximum amounts to loop out please consult "loop terms"                                                                                                                                                       | uint     |      `0`      |        *none*         |
| `--htlc_confs="…"`           | the number of confirmations (in blocks) that we require for the htlc extended by the server before we reveal the preimage                                                                                                                                                      | uint     |      `1`      |        *none*         |
| `--conf_target="…"`          | the number of blocks from the swap initiation height that the on-chain HTLC should be swept within                                                                                                                                                                             | uint     |      `9`      |        *none*         |
| `--max_swap_routing_fee="…"` | the max off-chain swap routing fee in satoshis, if not specified, a default max fee will be used                                                                                                                                                                               | int      |      `0`      |        *none*         |
| `--fast`                     | indicate you want to swap immediately, paying potentially a higher fee. If not set the swap server might choose to wait up to 30 minutes before publishing the swap HTLC on-chain, to save on its chain fees. Not setting this flag therefore might result in a lower swap fee | bool     |    `false`    |        *none*         |
| `--payment_timeout="…"`      | the timeout for each individual off-chain payment attempt. If not set, the default timeout of 1 hour will be used. As the payment might be retried, the actual total time may be longer                                                                                        | duration |     `0s`      |        *none*         |
| `--asset_id="…"`             | the asset ID of the asset to loop out, if this is set, the loop daemon will require a connection to a taproot assets daemon                                                                                                                                                    | string   |               |        *none*         |
| `--asset_edge_node="…"`      | the pubkey of the edge node of the asset to loop out, this is required if the taproot assets daemon has multiple channels of the given asset id with different edge nodes                                                                                                      | string   |               |        *none*         |
| `--force`                    | Assumes yes during confirmation. Using this option will result in an immediate swap                                                                                                                                                                                            | bool     |    `false`    |        *none*         |
| `--label="…"`                | an optional label for this swap,limited to 500 characters. The label may not start with our reserved prefix: [reserved]                                                                                                                                                        | string   |               |        *none*         |
| `--verbose` (`-v`)           | show expanded details                                                                                                                                                                                                                                                          | bool     |    `false`    |        *none*         |
| `--channel="…"`              | the comma-separated list of short channel IDs of the channels to loop out                                                                                                                                                                                                      | string   |               |        *none*         |
| `--help` (`-h`)              | show help                                                                                                                                                                                                                                                                      | bool     |    `false`    |        *none*         |

### `in` command

perform an on-chain to off-chain swap (loop in).

Send the amount in satoshis specified by the amt argument  		off-chain. 		 		By default the swap client will create and broadcast the  		on-chain htlc. The fee priority of this transaction can  		optionally be set using the conf_target flag.   		The external flag can be set to publish the on chain htlc  		independently. Note that this flag cannot be set with the  		conf_target flag.

Usage:

```bash
$ loop [GLOBAL FLAGS] in [COMMAND FLAGS] amt
```

The following flags are supported:

| Name                | Description                                                                                                             | Type   | Default value | Environment variables |
|---------------------|-------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--amt="…"`         | the amount in satoshis to loop in. To check for the minimum and maximum amounts to loop in please consult "loop terms"  | uint   |      `0`      |        *none*         |
| `--external`        | expect htlc to be published externally                                                                                  | bool   |    `false`    |        *none*         |
| `--conf_target="…"` | the target number of blocks the on-chain htlc broadcast by the swap client should confirm within                        | uint   |      `0`      |        *none*         |
| `--last_hop="…"`    | the pubkey of the last hop to use for this swap                                                                         | string |               |        *none*         |
| `--label="…"`       | an optional label for this swap,limited to 500 characters. The label may not start with our reserved prefix: [reserved] | string |               |        *none*         |
| `--force`           | Assumes yes during confirmation. Using this option will result in an immediate swap                                     | bool   |    `false`    |        *none*         |
| `--verbose` (`-v`)  | show expanded details                                                                                                   | bool   |    `false`    |        *none*         |
| `--route_hints="…"` | route hints that can each be individually used to assist in reaching the invoice's destination                          | string |     `[]`      |        *none*         |
| `--private`         | generates and passes routehints. Should be used if the connected node is only reachable via private channels            | bool   |    `false`    |        *none*         |
| `--help` (`-h`)     | show help                                                                                                               | bool   |    `false`    |        *none*         |

### `terms` command

Display the current swap terms imposed by the server.

Usage:

```bash
$ loop [GLOBAL FLAGS] terms [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `monitor` command

monitor progress of any active swaps.

Allows the user to monitor progress of any active swaps.

Usage:

```bash
$ loop [GLOBAL FLAGS] monitor [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `quote` command

get a quote for the cost of a swap.

Usage:

```bash
$ loop [GLOBAL FLAGS] quote [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `quote in` subcommand

get a quote for the cost of a loop in swap.

Allows to determine the cost of a swap up front.Either specify an amount or deposit outpoints.

Usage:

```bash
$ loop [GLOBAL FLAGS] quote in [COMMAND FLAGS] amt
```

The following flags are supported:

| Name                     | Description                                                                                                                                                                                                    | Type   | Default value | Environment variables |
|--------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--last_hop="…"`         | the pubkey of the last hop to use for the quote                                                                                                                                                                | string |               |        *none*         |
| `--conf_target="…"`      | the target number of blocks the on-chain htlc broadcast by the swap client should confirm within                                                                                                               | uint   |      `0`      |        *none*         |
| `--verbose` (`-v`)       | show expanded details                                                                                                                                                                                          | bool   |    `false`    |        *none*         |
| `--private`              | generates and passes routehints. Should be used if the connected node is only reachable via private channels                                                                                                   | bool   |    `false`    |        *none*         |
| `--route_hints="…"`      | route hints that can each be individually used to assist in reaching the invoice's destination                                                                                                                 | string |     `[]`      |        *none*         |
| `--deposit_outpoint="…"` | one or more static address deposit outpoints to quote for. Deposit outpoints are not to be used in combination with an amount. Eachadditional outpoint can be added by specifying --deposit_outpoint tx_id:idx | string |     `[]`      |        *none*         |
| `--help` (`-h`)          | show help                                                                                                                                                                                                      | bool   |    `false`    |        *none*         |

### `quote out` subcommand

get a quote for the cost of a loop out swap.

Allows to determine the cost of a swap up front.

Usage:

```bash
$ loop [GLOBAL FLAGS] quote out [COMMAND FLAGS] amt
```

The following flags are supported:

| Name                | Description                                                                                                                                                                                                                                                      | Type | Default value | Environment variables |
|---------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------|:-------------:|:---------------------:|
| `--conf_target="…"` | the number of blocks from the swap initiation height that the on-chain HTLC should be swept within in a Loop Out                                                                                                                                                 | uint |      `9`      |        *none*         |
| `--fast`            | Indicate you want to swap immediately, paying potentially a higher fee. If not set the swap server might choose to wait up to 30 minutes before publishing the swap HTLC on-chain, to save on chain fees. Not setting this flag might result in a lower swap fee | bool |    `false`    |        *none*         |
| `--verbose` (`-v`)  | show expanded details                                                                                                                                                                                                                                            | bool |    `false`    |        *none*         |
| `--help` (`-h`)     | show help                                                                                                                                                                                                                                                        | bool |    `false`    |        *none*         |

### `listauth` command

list all L402 tokens.

Shows a list of all L402 tokens that loopd has paid for.

Usage:

```bash
$ loop [GLOBAL FLAGS] listauth [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `fetchl402` command

fetches a new L402 authentication token from the server.

Fetches a new L402 authentication token from the server. This token is required to listen to notifications from the server, such as reservation notifications. If a L402 is already present in the store, this command is a no-op.

Usage:

```bash
$ loop [GLOBAL FLAGS] fetchl402 [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `listswaps` command

list all swaps in the local database.

Allows the user to get a list of all swaps that are currently stored in the database.

Usage:

```bash
$ loop [GLOBAL FLAGS] listswaps [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                  | Description                                                                                                             | Type   | Default value | Environment variables |
|-----------------------|-------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--loop_out_only`     | only list swaps that are loop out swaps                                                                                 | bool   |    `false`    |        *none*         |
| `--loop_in_only`      | only list swaps that are loop in swaps                                                                                  | bool   |    `false`    |        *none*         |
| `--pending_only`      | only list pending swaps                                                                                                 | bool   |    `false`    |        *none*         |
| `--label="…"`         | an optional label for this swap,limited to 500 characters. The label may not start with our reserved prefix: [reserved] | string |               |        *none*         |
| `--channel="…"`       | the comma-separated list of short channel IDs of the channels to loop out                                               | string |               |        *none*         |
| `--last_hop="…"`      | the pubkey of the last hop to use for this swap                                                                         | string |               |        *none*         |
| `--max_swaps="…"`     | Max number of swaps to return after filtering                                                                           | uint   |      `0`      |        *none*         |
| `--start_time_ns="…"` | Unix timestamp in nanoseconds to select swaps initiated after this time                                                 | int    |      `0`      |        *none*         |
| `--help` (`-h`)       | show help                                                                                                               | bool   |    `false`    |        *none*         |

### `swapinfo` command

show the status of a swap.

Allows the user to get the status of a single swap currently stored in the database.

Usage:

```bash
$ loop [GLOBAL FLAGS] swapinfo [COMMAND FLAGS] id
```

The following flags are supported:

| Name            | Description        | Type | Default value | Environment variables |
|-----------------|--------------------|------|:-------------:|:---------------------:|
| `--id="…"`      | the ID of the swap | uint |      `0`      |        *none*         |
| `--help` (`-h`) | show help          | bool |    `false`    |        *none*         |

### `getparams` command

show liquidity manager parameters.

Displays the current set of parameters that are set for the liquidity manager.

Usage:

```bash
$ loop [GLOBAL FLAGS] getparams [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `setrule` command

set liquidity manager rule for a channel/peer.

Update or remove the liquidity rule for a channel/peer.

Usage:

```bash
$ loop [GLOBAL FLAGS] setrule [COMMAND FLAGS] {shortchanid | peerpubkey}
```

The following flags are supported:

| Name                       | Description                                                                                                            | Type   | Default value | Environment variables |
|----------------------------|------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--type="…"`               | the type of swap to perform, set to 'out' for acquiring inbound liquidity or 'in' for acquiring outbound liquidity     | string |     `out`     |        *none*         |
| `--incoming_threshold="…"` | the minimum percentage of incoming liquidity to total capacity beneath which to recommend loop out to acquire incoming | int    |      `0`      |        *none*         |
| `--outgoing_threshold="…"` | the minimum percentage of outbound liquidity that we do not want to drop below                                         | int    |      `0`      |        *none*         |
| `--clear`                  | remove the rule currently set for the channel/peer                                                                     | bool   |    `false`    |        *none*         |
| `--help` (`-h`)            | show help                                                                                                              | bool   |    `false`    |        *none*         |

### `suggestswaps` command

show a list of suggested swaps.

Displays a list of suggested swaps that aim to obtain the liquidity balance as specified by the rules set in the liquidity manager.

Usage:

```bash
$ loop [GLOBAL FLAGS] suggestswaps [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `setparams` command

update the parameters set for the liquidity manager.

Updates the parameters set for the liquidity manager. Note the parameters are persisted in db to save the trouble of setting them again upon loopd restart. To get the defaultvalues, use `getparams` before any `setparams`.

Usage:

```bash
$ loop [GLOBAL FLAGS] setparams [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                            | Description                                                                                                                                                                                                                                                                                  | Type     | Default value | Environment variables |
|---------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|:-------------:|:---------------------:|
| `--sweeplimit="…"`              | the limit placed on our estimated sweep fee in sat/vByte                                                                                                                                                                                                                                     | int      |      `0`      |        *none*         |
| `--feepercent="…"`              | the maximum percentage of swap amount to be used across all fee categories                                                                                                                                                                                                                   | float    |      `0`      |        *none*         |
| `--maxswapfee="…"`              | the maximum percentage of swap volume we are willing to pay in server fees                                                                                                                                                                                                                   | float    |      `0`      |        *none*         |
| `--maxroutingfee="…"`           | the maximum percentage of off-chain payment volume that we are willing to pay in routingfees                                                                                                                                                                                                 | float    |      `0`      |        *none*         |
| `--maxprepayfee="…"`            | the maximum percentage of off-chain prepay volume that we are willing to pay in routing fees                                                                                                                                                                                                 | float    |      `0`      |        *none*         |
| `--maxprepay="…"`               | the maximum no-show (prepay) in satoshis that swap suggestions should be limited to                                                                                                                                                                                                          | uint     |      `0`      |        *none*         |
| `--maxminer="…"`                | the maximum miner fee in satoshis that swap suggestions should be limited to                                                                                                                                                                                                                 | uint     |      `0`      |        *none*         |
| `--sweepconf="…"`               | the number of blocks from htlc height that swap suggestion sweeps should target, used to estimate max miner fee                                                                                                                                                                              | int      |      `0`      |        *none*         |
| `--failurebackoff="…"`          | the amount of time, in seconds, that should pass before a channel that previously had a failed swap will be included in suggestions                                                                                                                                                          | uint     |      `0`      |        *none*         |
| `--autoloop`                    | set to true to enable automated dispatch of swaps, limited to the budget set by autobudget                                                                                                                                                                                                   | bool     |    `false`    |        *none*         |
| `--destaddr="…"`                | custom address to be used as destination for autoloop loop out, set to "default" in order to revert to default behavior                                                                                                                                                                      | string   |               |        *none*         |
| `--account="…"`                 | the name of the account to generate a new address from. You can list the names of valid accounts in your backing lnd instance with "lncli wallet accounts list"                                                                                                                              | string   |               |        *none*         |
| `--account_addr_type="…"`       | the address type of the extended public key specified in account. Currently only pay-to-taproot-pubkey(p2tr) is supported                                                                                                                                                                    | string   |    `p2tr`     |        *none*         |
| `--autobudget="…"`              | the maximum amount of fees in satoshis that automatically dispatched loop out swaps may spend                                                                                                                                                                                                | uint     |      `0`      |        *none*         |
| `--autobudgetrefreshperiod="…"` | the time period over which the automated loop budget is refreshed                                                                                                                                                                                                                            | duration |     `0s`      |        *none*         |
| `--autoinflight="…"`            | the maximum number of automatically dispatched swaps that we allow to be in flight                                                                                                                                                                                                           | uint     |      `0`      |        *none*         |
| `--minamt="…"`                  | the minimum amount in satoshis that the autoloop client will dispatch per-swap                                                                                                                                                                                                               | uint     |      `0`      |        *none*         |
| `--maxamt="…"`                  | the maximum amount in satoshis that the autoloop client will dispatch per-swap                                                                                                                                                                                                               | uint     |      `0`      |        *none*         |
| `--htlc_conf="…"`               | the confirmation target for loop in on-chain htlcs                                                                                                                                                                                                                                           | int      |      `0`      |        *none*         |
| `--easyautoloop`                | set to true to enable easy autoloop, which will automatically dispatch swaps in order to meet the target local balance                                                                                                                                                                       | bool     |    `false`    |        *none*         |
| `--localbalancesat="…"`         | the target size of total local balance in satoshis, used by easy autoloop                                                                                                                                                                                                                    | uint     |      `0`      |        *none*         |
| `--asset_easyautoloop`          | set to true to enable asset easy autoloop, which will automatically dispatch asset swaps in order to meet the target local balance                                                                                                                                                           | bool     |    `false`    |        *none*         |
| `--asset_id="…"`                | If set to a valid asset ID, the easyautoloop and localbalancesat flags will be set for the specified asset                                                                                                                                                                                   | string   |               |        *none*         |
| `--asset_localbalance="…"`      | the target size of total local balance in asset units, used by asset easy autoloop                                                                                                                                                                                                           | uint     |      `0`      |        *none*         |
| `--fast`                        | if set new swaps are expected to be published immediately, paying a potentially higher fee. If not set the swap server might choose to wait up to 30 minutes before publishing swap HTLCs on-chain, to save on chain fees. Not setting this flag therefore might result in a lower swap fees | bool     |    `false`    |        *none*         |
| `--help` (`-h`)                 | show help                                                                                                                                                                                                                                                                                    | bool     |    `false`    |        *none*         |

### `getinfo` command

show general information about the loop daemon.

Displays general information about the daemon like current version, connection parameters and basic swap information.

Usage:

```bash
$ loop [GLOBAL FLAGS] getinfo [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `abandonswap` command

abandon a swap with a given swap hash.

This command overrides the database and abandons a swap with a given swap hash.  !!! This command might potentially lead to loss of funds if it is applied to swaps that are still waiting for pending user funds. Before executing this command make sure that no funds are locked by the swap.

Usage:

```bash
$ loop [GLOBAL FLAGS] abandonswap [COMMAND FLAGS] ID
```

The following flags are supported:

| Name                       | Description                                                                                                        | Type | Default value | Environment variables |
|----------------------------|--------------------------------------------------------------------------------------------------------------------|------|:-------------:|:---------------------:|
| `--i_know_what_i_am_doing` | Specify this flag if you made sure that you read and understood the following consequence of applying this command | bool |    `false`    |        *none*         |
| `--help` (`-h`)            | show help                                                                                                          | bool |    `false`    |        *none*         |

### `reservations` command (aliases: `r`)

manage reservations.

With loopd running, you can use this command to manage your 		reservations. Reservations are 2-of-2 multisig utxos that 		the loop server can open to clients. The reservations are used 		to enable instant swaps.

Usage:

```bash
$ loop [GLOBAL FLAGS] reservations [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `reservations list` subcommand (aliases: `l`)

list all reservations.

List all reservations.

Usage:

```bash
$ loop [GLOBAL FLAGS] reservations list [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `instantout` command

perform an instant off-chain to on-chain swap (looping out).

Attempts to instantly loop out into the backing lnd's wallet. The amount 	will be chosen via the cli.

Usage:

```bash
$ loop [GLOBAL FLAGS] instantout [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description                                                                                                      | Type   | Default value | Environment variables |
|-----------------|------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--channel="…"` | the comma-separated list of short channel IDs of the channels to loop out                                        | string |               |        *none*         |
| `--addr="…"`    | the optional address that the looped out funds should be sent to, if let blank the funds will go to lnd's wallet | string |               |        *none*         |
| `--help` (`-h`) | show help                                                                                                        | bool   |    `false`    |        *none*         |

### `listinstantouts` command

list all instant out swaps.

List all instant out swaps.

Usage:

```bash
$ loop [GLOBAL FLAGS] listinstantouts [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `static` command (aliases: `s`)

perform on-chain to off-chain swaps using static addresses.

Usage:

```bash
$ loop [GLOBAL FLAGS] static [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `static new` subcommand (aliases: `n`)

Create a new static loop in address.

Requests a new static loop in address from the server. Funds that are 	sent to this address will be locked by a 2:2 multisig between us and the 	loop server, or a timeout path that we can sweep once it opens up. The  	funds can either be cooperatively spent with a signature from the server 	or looped in.

Usage:

```bash
$ loop [GLOBAL FLAGS] static new [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `static listunspent` subcommand (aliases: `l`)

List unspent static address outputs.

List all unspent static address outputs.

Usage:

```bash
$ loop [GLOBAL FLAGS] static listunspent [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description                                                            | Type | Default value | Environment variables |
|-------------------|------------------------------------------------------------------------|------|:-------------:|:---------------------:|
| `--min_confs="…"` | The minimum amount of confirmations an output should have to be listed | int  |      `0`      |        *none*         |
| `--max_confs="…"` | The maximum number of confirmations an output could have to be listed  | int  |      `0`      |        *none*         |
| `--help` (`-h`)   | show help                                                              | bool |    `false`    |        *none*         |

### `static listdeposits` subcommand

Displays static address deposits. A filter can be applied to only show deposits in a specific state.

Usage:

```bash
$ loop [GLOBAL FLAGS] static listdeposits [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description                                                                                                                                                                                                                                                                                                                                      | Type   | Default value | Environment variables |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--filter="…"`  | specify a filter to only display deposits in the specified state. Leaving out the filter returns all deposits. The state can be one of the following:  deposited withdrawing withdrawn looping_in looped_in opening_channel channel_published publish_expired_deposit sweep_htlc_timeout htlc_timeout_swept wait_for_expiry_sweep expired failed | string |               |        *none*         |
| `--help` (`-h`) | show help                                                                                                                                                                                                                                                                                                                                        | bool   |    `false`    |        *none*         |

### `static listwithdrawals` subcommand

Display a summary of past withdrawals.

Usage:

```bash
$ loop [GLOBAL FLAGS] static listwithdrawals [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `static listswaps` subcommand

Shows a list of finalized static address swaps.

Usage:

```bash
$ loop [GLOBAL FLAGS] static listswaps [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `static withdraw` subcommand (aliases: `w`)

Withdraw from static address deposits.

Withdraws from all or selected static address deposits by sweeping them 	to the internal wallet or an external address.

Usage:

```bash
$ loop [GLOBAL FLAGS] static withdraw [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                     | Description                                                                                                               | Type   | Default value | Environment variables |
|--------------------------|---------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--utxo="…"`             | specify utxos as outpoints(tx:idx) which willbe withdrawn                                                                 | string |     `[]`      |        *none*         |
| `--all`                  | withdraws all static address deposits                                                                                     | bool   |    `false`    |        *none*         |
| `--dest_addr="…"`        | the optional address that the withdrawn funds should be sent to, if let blank the funds will go to lnd's wallet           | string |               |        *none*         |
| `--sat_per_vbyte="…"`    | (optional) a manual fee expressed in sat/vbyte that should be used when crafting the transaction                          | uint   |      `0`      |        *none*         |
| `--amt="…"` (`--amount`) | the number of satoshis that should be withdrawn from the selected deposits. The change is sent back to the static address | uint   |      `0`      |        *none*         |
| `--help` (`-h`)          | show help                                                                                                                 | bool   |    `false`    |        *none*         |

### `static summary` subcommand (aliases: `s`)

Display a summary of static address related information.

Displays various static address related information about deposits,  	withdrawals, swaps and channel openings.

Usage:

```bash
$ loop [GLOBAL FLAGS] static summary [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `static in` subcommand

Loop in funds from static address deposits.

Requests a loop-in swap based on static address deposits. After the 	creation of a static address funds can be sent to it. Once the funds are 	confirmed on-chain they can be swapped instantaneously. If deposited 	funds are not needed they can we withdrawn back to the local lnd wallet.

Usage:

```bash
$ loop [GLOBAL FLAGS] static in [COMMAND FLAGS] [amt] [--all | --utxo xxx:xx]
```

The following flags are supported:

| Name                     | Description                                                                                                                                                             | Type     | Default value | Environment variables |
|--------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|:-------------:|:---------------------:|
| `--utxo="…"`             | specify the utxos of deposits as outpoints(tx:idx) that should be looped in                                                                                             | string   |     `[]`      |        *none*         |
| `--all`                  | loop in all static address deposits                                                                                                                                     | bool     |    `false`    |        *none*         |
| `--payment_timeout="…"`  | the maximum time in seconds that the server is allowed to take for the swap payment. The client can retry the swap with adjusted parameters after the payment timed out | duration |     `0s`      |        *none*         |
| `--amt="…"` (`--amount`) | the number of satoshis that should be swapped from the selected deposits. If thereis change it is sent back to the static address                                       | uint     |      `0`      |        *none*         |
| `--fast`                 | Usage: complete the swap faster by paying a higher fee, so the change output is available sooner                                                                        | bool     |    `false`    |        *none*         |
| `--last_hop="…"`         | the pubkey of the last hop to use for this swap                                                                                                                         | string   |               |        *none*         |
| `--label="…"`            | an optional label for this swap,limited to 500 characters. The label may not start with our reserved prefix: [reserved]                                                 | string   |               |        *none*         |
| `--route_hints="…"`      | route hints that can each be individually used to assist in reaching the invoice's destination                                                                          | string   |     `[]`      |        *none*         |
| `--private`              | generates and passes routehints. Should be used if the connected node is only reachable via private channels                                                            | bool     |    `false`    |        *none*         |
| `--force`                | Assumes yes during confirmation. Using this option will result in an immediate swap                                                                                     | bool     |    `false`    |        *none*         |
| `--verbose` (`-v`)       | show expanded details                                                                                                                                                   | bool     |    `false`    |        *none*         |
| `--help` (`-h`)          | show help                                                                                                                                                               | bool     |    `false`    |        *none*         |

### `static openchannel` subcommand

Open a channel to a an existing peer.

Attempt to open a new channel to an existing peer with the key  	node-key.  	The channel will be initialized with local-amt satoshis locally and 	push-amt satoshis for the remote node. Note that the push-amt is 	deducted from the specified local-amt which implies that the local-amt 	must be greater than the push-amt. Also note that specifying push-amt 	means you give that amount to the remote node as part of the channel 	opening. Once the channel is open, a channelPoint (txid:vout) of the 	funding output is returned.  	If the remote peer supports the option upfront shutdown feature bit 	(query listpeers to see their supported feature bits), an address to 	enforce payout of funds on cooperative close can optionally be provided. 	Note that if you set this value, you will not be able to cooperatively 	close out to another address.  	One can also specify a short string memo to record some useful 	information about the channel using the --memo argument. This is stored 	locally only, and is purely for reference. It has no bearing on the 	channel's operation. Max allowed length is 500 characters.

Usage:

```bash
$ loop [GLOBAL FLAGS] static openchannel [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                                    | Description                                                                                                                                                                                                                                                                                                           | Type   | Default value | Environment variables |
|-----------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--node_key="…"`                        | the identity public key of the target node/peer serialized in compressed format                                                                                                                                                                                                                                       | string |               |        *none*         |
| `--local_amt="…"`                       | the number of satoshis the wallet should commit to the channel                                                                                                                                                                                                                                                        | int    |      `0`      |        *none*         |
| `--base_fee_msat="…"`                   | the base fee in milli-satoshis that will be charged for each forwarded HTLC, regardless of payment size                                                                                                                                                                                                               | uint   |      `0`      |        *none*         |
| `--fee_rate_ppm="…"`                    | the fee rate ppm (parts per million) that will be charged proportionally based on the value of each forwarded HTLC, the lowest possible rate is 0 with a granularity of 0.000001 (millionths)                                                                                                                         | uint   |      `0`      |        *none*         |
| `--push_amt="…"`                        | the number of satoshis to give the remote side as part of the initial commitment state, this is equivalent to first opening a channel and sending the remote party funds, but done all in one step                                                                                                                    | int    |      `0`      |        *none*         |
| `--sat_per_vbyte="…"`                   | (optional) a manual fee expressed in sat/vbyte that should be used when crafting the transaction                                                                                                                                                                                                                      | int    |      `0`      |        *none*         |
| `--private`                             | make the channel private, such that it won't be announced to the greater network, and nodes other than the two channel endpoints must be explicitly told about it to be able to route through it                                                                                                                      | bool   |    `false`    |        *none*         |
| `--min_htlc_msat="…"`                   | (optional) the minimum value we will require for incoming HTLCs on the channel                                                                                                                                                                                                                                        | int    |      `0`      |        *none*         |
| `--remote_csv_delay="…"`                | (optional) the number of blocks we will require our channel counterparty to wait before accessing its funds in case of unilateral close. If this is not set, we will scale the value according to the channel size                                                                                                    | uint   |      `0`      |        *none*         |
| `--max_local_csv="…"`                   | (optional) the maximum number of blocks that we will allow the remote peer to require we wait before accessing our funds in the case of a unilateral close                                                                                                                                                            | uint   |      `0`      |        *none*         |
| `--close_address="…"`                   | (optional) an address to enforce payout of our funds to on cooperative close. Note that if this value is set on channel open, you will *not* be able to cooperatively close to a different address                                                                                                                    | string |               |        *none*         |
| `--remote_max_value_in_flight_msat="…"` | (optional) the maximum value in msat that can be pending within the channel at any given time                                                                                                                                                                                                                         | uint   |      `0`      |        *none*         |
| `--channel_type="…"`                    | (optional) the type of channel to propose to the remote peer ("tweakless", "anchors", "taproot")                                                                                                                                                                                                                      | string |               |        *none*         |
| `--zero_conf`                           | (optional) whether a zero-conf channel open should be attempted                                                                                                                                                                                                                                                       | bool   |    `false`    |        *none*         |
| `--scid_alias`                          | (optional) whether a scid-alias channel type should be negotiated                                                                                                                                                                                                                                                     | bool   |    `false`    |        *none*         |
| `--remote_reserve_sats="…"`             | (optional) the minimum number of satoshis we require the remote node to keep as a direct payment. If not specified, a default of 1% of the channel capacity will be used                                                                                                                                              | uint   |      `0`      |        *none*         |
| `--memo="…"`                            | (optional) a note-to-self containing some useful 				information about the channel. This is stored 				locally only, and is purely for reference. It 				has no bearing on the channel's operation. Max 				allowed length is 500 characters                                                                          | string |               |        *none*         |
| `--fundmax`                             | if set, the wallet will attempt to commit the maximum possible local amount to the channel. This must not be set at the same time as local_amt                                                                                                                                                                        | bool   |    `false`    |        *none*         |
| `--utxo="…"`                            | a utxo specified as outpoint(tx:idx) which will be used to fund a channel. This flag can be repeatedly used to fund a channel with a selection of utxos. The selected funds can either be entirely spent by specifying the fundmax flag or partially by selecting a fraction of the sum of the outpoints in local_amt | string |     `[]`      |        *none*         |
| `--help` (`-h`)                         | show help                                                                                                                                                                                                                                                                                                             | bool   |    `false`    |        *none*         |

