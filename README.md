# keycard-cli

`keycard` is a command line tool to manage [Status Keycards](https://github.com/status-im/status-keycard).

* [Dependencies](#dependencies)
* [Installation](#installation)
* [Continuous Integration](#continuous-integration)
* CLI Commands
  * [Card info](#card-info)
  * [Keycard applet installation](#keycard-applet-installation)
  * [Card initialization](#card-initialization)
  * [Deleting the applet](#deleting-the-applet)
  * [Keycard shell](#keycard-shell)

## Dependencies

On linux you need to install and run the [pcsc daemon](https://linux.die.net/man/8/pcscd).

## Installation

Download the binary for your platform from the [releases page](https://github.com/status-im/keycard-cli/releases).

## Continuous Integration

Jenkins builds provide:

* [PR Builds](https://ci.status.im/job/status-keycard/job/prs/job/keycard-cli/) - Run only the `test` and `build` targets.
* [Manual Builds](https://ci.status.im/job/status-keycard/job/keycard-cli/) - Create GitHub release draft with binaries for 3 platforms.

Successful PR builds are mandatory.

## Usage

### Card info

```bash
keycard info -l debug
```

The `info` command will print something like this:

```
Installed: true
Initialized: false
InstanceUID: 0x
PublicKey: 0x112233...
Version: 0x
AvailableSlots: 0x
KeyUID: 0x
```
### Keycard applet installation

The `install` command will install an applet to the card.
You can download the status `cap` file from the [status-im/status-keycard releases page](https://github.com/status-im/status-keycard/releases).

```bash
keycard install -l debug -a PATH_TO_CAP_FILE
```

In case the applet is already installed and you want to force a new installation you can pass the `-f` flag.


### Card initialization


```bash
keycard init -l debug
```

The `init` command initializes the card and generates the secrets needed to pair the card to a device.

```
PIN 123456
PUK 123456789012
Pairing password: RandomPairingPassword
```

### Deleting the applet

:warning: **WARNING! This command will remove the applet and all the keys from the card.** :warning:

```bash
keycard-cli delete -l debug
```

### Keycard shell

The shell can be used to interact with the KeyCard using `keycard-go`. You can start the shell with:

```
keycard-cli shell
```

Once in the shell, you may submit one command at a time, followed by Enter.

### Talking to the Card

Before you can communicate with the card, you need to initialize it. You can do this interactively in the shell or using `keycard init` specified in the above section.

> Once you initialize your card, **save the PIN, PUK, and Pairing Password fields that are generated**; you will need these to establish a connection with the card.

With the secrets in hand, run the following commands in the shell:

```
> keycard-select
> keycard-set-secrets <PIN> <PUK> <PairingPassword>
> keycard-pair
> keycard-open-secure-channel
> keycard-verify-pin <PIN>
```

If you don't get an error message, it means you are connected to the card! This connection will persist for the duration of your shell session.

> `keycard-pair` will print two values: `PAIRING KEY` and `PAIRING INDEX`. You should save these values if you wish to reuse this pairing with a new shell session. You can do this with the command: `>keycard-set-pairing <PAIRING KEY> <PAIRING INDEX>`

### Wallets

Once you have verified your PIN, you can create, import, and interact with your HD wallet. The keycard stores one seed (and thus one HD wallet) at a time.

#### Generating a Key

You can generate a key using the TRNG of the smartcard:

```
> keycard-generate-key
```

This will create and load a seed (and master keypair) that does not correspond to a mnemonic, so you cannot export it as a seed phrase. The generation utilizes the card's true random number generator.

##### Importing a Key or Seed

You can also import a seed or key:

```
> keycard-load-key <isSeed> <isExtended> <data>
```

If `isSeed=1`, your `data` should be a 64-byte hex string representing the master seed for the HD wallet you are creating.

> See `status-keycard` docs for formatting when `isSeed=0`. You should be including a keypair with an optional chaincode to indicate whether it is an extended keypair (`isExtended=1` for extended keypair imports). Importing keys directly is generally not recommended, as the data is less compact.


#### Setting a "Current" Key

The card applet works by loading a private key (based on a derivation path) into a "current" state. Once set as the current key, the subsequent signature request will be filled by that key.

You can make a key current by running:

```
> keycard-derive-key <DERIVATION_PATH>
```

Where `DERIVATION_PATH` is a BIP44 path, e.g. `m/44/0'/0'/0/0`

#### Exporting Keys

You can export both public and private keys based on a derivation path (or without a path, using the current key).

> Public keys can always be exported, but private keys have some restrictions, which you can read about [here].

When exporting a key, you have two parameters (`p1` and `p2`) to specify:

* `p1` - Derivation options
* `p2` - Type of data export

**`p1`**
|  Option    |   Description            |
|:-----------|:-------------------------|
| `0`     | Export Current Key       |
| `1`     | Derive                   |
| `2`     | Derive and make current key |

**`p2`**

|  Option    |   Description                          |
|:-----------|:---------------------------------------|
| `0`        | Export public AND private key          |
| `1`        | Export public key                      |
| `2`        | Export public key and chaincode        |


> **IMPORTANT NOTE**: GridPlus has disabled private key exports, while status does not allow chaincode exports

Example:

```
> keycard-export-key 2 2 m/44'/0'/0'/0/0

Exorted Public Key: <PUB_KEY>
Exported Chain Code: <CHAIN_CODE>
```

> All public keys are build on the secp256k1 curve and exported in uncompressed point format, i.e. `04{X-component}{Y-component}`.

#### Master Seed

The `masterSeed` contains the root entropy of the HD wallet. The seed may be exported if it is created/imported using an appropriate flag.

> A `flag` is a one-byte value. Currently the only allowable valures are 0 (non-exportable seed) and 1 (exportable seed), but more options may be added in the future.

**Generating a seed:**

```
> keycard-generate-key <flag>
```

If the second argument is not provided, this will generate a `masterSeed` which is *not* exportable. You may designate the seed as exportable with `keycard-generate-key 1`.

**Importing a Seed:**

```
> keycard-load-key 3 <flag> <seed>
```

Here the `3` indicates that we are loading a seed (64 byte hex string). The `flag` is the same as for generating: 1=exportable, 0=non-exportable.

**Exporting seed:**

You can export the master seed (if it is exportable) with the following command:

```
> keycard-export-seed
```

#### Signing

TODO

#### Other Functionality

TODO