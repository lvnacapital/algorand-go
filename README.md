# Algorand

Algorand is a CLI application that abstracts low-level Algorand node operations.

It uses the [Algorand Go SDK](https://github.com/algorand/go-algorand-sdk). For any questions please refer to the [developer documentation](https://developer.algorand.org/docs/go-sdk).

It also leverages [Cobra](https://github.com/spf13/cobra) for the CLI interface and [Viper](https://github.com/spf13/viper) for configuration.

## Build

Prerequisites:
* Install Algorand node using: [https://developer.algorand.org/docs/introduction-installing-node](https://developer.algorand.org/docs/introduction-installing-node)
* Install `git` from: [https://git-scm.com/downloads](https://git-scm.com/downloads)
* Install `go` from: [https://golang.org/dl/](https://golang.org/dl/)
* Install `gox` using: `go get -u github.com/mitchellh/gox`
* Install `golint` using: `go get -u golang.org/x/lint/golint`
* Install `dep` using: `curl https://raw.githubusercontent.com/golang/dep/masterinstall.sh | sh`
* Install `make` (on Windows) from: [http://gnuwin32.sourceforge.net/packages/make.htm](http://gnuwin32.sourceforge.net/packages/make.htm)

To build for all platforms on Linux use:

```
$ make -r build
```

To build for the current platform only:
```
$ make -r run
```

To build for all platforms on Windows use:
```
$ /c/Program\ Files\ \(x86\)/GnuWin32/bin/make.exe -r build
```

Refer to the [Makefile](Makefile) for other options.

## Configuration

To configure `algorand`, a `config.yml` file can be passed with the --config flag. Configuration read in through a file will overwrite the same configuration specified by a flag. If no config file is passed, and no flags are set, reasonable defaults will be used.

```yml
hostname: '127.0.0.1'                 # the Algorand node's IP
algod-port: '8080'                    # the `algod' process port
kmd-port: '7833'                      # the `kmd' process port
algod-token: '374de74fa794248762e5a'  # the `algod' process token
kmd-token: 'be84aa55f61665645ed680b'  # the `kmd' process token
```

These values are taken from:
* `$NODE/data/algod.net`
* `$NODE/data/algod.token`
* `$NODE/data/kmd-<VERSION>/kmd.net`
* `$NODE/data/kmd-<VERSION>/kmd.token`

Many commands assume that the node is set up as archival in `$NODE/data/config.json`:

```json
{
    "Archival": true
}
```

## Usage

Ensure the `algod` and `kmd` processes are started and that the node is synchronized to the network as described [here](https://developer.algorand.org/docs/introduction-installing-node).

```
$ $NODE/goal node start
Algorand node successfully started!
$ $NODE/goal kmd start
Successfully started kmd
$ $NODE/goal node status
```

To use simply run `algorand [command] [flags]`. Running `algorand` or `algorand -h` will show usage and a list of commands.

```
Usage:
  algorand [command] [flags]
```

## Examples

### status:

```
$ algorand status
Last committed block: 378432
Time since last block: 2.4s
Sync Time: 0.0s
Last consensus protocol: v4
Next consensus protocol: v4
Round for next consensus protocol: 378433
Next consensus protocol supported: true
```

### sign

```
$ algorand sign
Have 1 wallet(s):
[1] Name: JSWallet      ID: 895bad84b32bbffa28c2c069d6b49e9f
Select wallet [1]: 1
Please type in the password for 'JSWallet':

Have 1 address(es) in 'JSWallet':
[1] MYPI256EXJQIMTV3NHX2BNVAPL7WRXCOOJ67X4WE536RSUMKEZVSP4IBI4
Pick the account address to send from [1]: 1

Specify the account address to send to: KI6TMKHUQOGJ7EDZLOWFOGHBBWBIMBMKONMS565X7NSOFMAM6S2EK4GBHQ

Specify the amount to be transferred: 1000

Specify some note text (optional): Algorand!

Made transaction: {_struct:{} Type:pay Header:{_struct:{} Sender:MYPI256EXJQIMTV3NHX2BNVAPL7WRXCOOJ67X4WE536RSUMKEZVSP4IBI4 Fee:101 FirstValid:365680 LastValid:365682 Note:[169 65 108 103 111 114 97 110 100 33] GenesisID:testnet-v31.0} KeyregTxnFields:{_struct:{} VotePK:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] SelectionPK:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} PaymentTxnFields:{_struct:{} Receiver:KI6TMKHUQOGJ7EDZLOWFOGHBBWBIMBMKONMS565X7NSOFMAM6S2EK4GBHQ Amount:1000 CloseRemainderTo:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ}}

Made signed transaction using library: 82a3736967c440669145281bb559eb2a278aea14cda9946ec367dbd191b0a480acb21a9489c0d6361f8b946711fc8fc864ea066021cbd5c01b078f4d072e07d47f02913624fd07a374786e89a3616d74cd03e8a366656565a26676ce00059470a367656ead746573746e65742d7633312e30a26c76ce00059472a46e6f7465c40aa9416c676f72616e6421a3726376c420523d3628f4838c9f90795bac5718e10d8286058a73592efbb7fb64e2b00cf4b4a3736e64c420661e8d77c4ba60864ebb69efa0b6a07aff68dc4e727dfbf2c4eefd19518a266ba474797065a3706179

`kmd' made signed transaction with bytes: 82a3736967c440669145281bb559eb2a278aea14cda9946ec367dbd191b0a480acb21a9489c0d6361f8b946711fc8fc864ea066021cbd5c01b078f4d072e07d47f02913624fd07a374786e89a3616d74cd03e8a366656565a26676ce00059470a367656ead746573746e65742d7633312e30a26c76ce00059472a46e6f7465c40aa9416c676f72616e6421a3726376c420523d3628f4838c9f90795bac5718e10d8286058a73592efbb7fb64e2b00cf4b4a3736e64c420661e8d77c4ba60864ebb69efa0b6a07aff68dc4e727dfbf2c4eefd19518a266ba474797065a3706179

Signed transactions match!

Sent transaction with ID: tx-Y7XIQGCRU6IRLLCKSZJVTMQSH6HSDMDVW3WQFVOMNNBNFL5BO6KQ
```

### find:

```
$ algorand find
Find Transaction Using Transaction ID
-------------------------------------
Enter the transaction ID: tx-A6R7R6EL2I4QJRHBSRLE2B4AQ3N74MKRWQZARYCXQOR742HC3NGQ
No transactions in block: 362201
No transactions in block: 362200
...
No transactions in block: 346047
No transactions in block: 346046
Found transaction in block: 343498
Transaction: {
  "round": 343498,
  "fee": 101,
  "first-round": 343496,
  "from": "2VXBXLOZSLA5EXPYD3P2SS5ODNUDTMOWTIQPLEU2SZB2Z563IWIXMQKJKI",
  "last-round": 344496,
  "noteb64": "gqJiea9SaWNoYXJkIERhd2tpbnOkdGV4dNk7dGhlIHNvbHV0aW9uIG9mdGVuIHR1cm5zIG91dCBtb3JlIGJlYXV0aWZ1bCB0aGFuIHRoZSBwdXp6bGU=",
  "tx": "A6R7R6EL2I4QJRHBSRLE2B4AQ3N74MKRWQZARYCXQOR742HC3NGQ",
  "payment": {
    "amount": 100,
    "to": "NJY27OQ2ZXK6OWBN44LE4K43TA2AV3DPILPYTHAJAMKIVZDWTEJKZJKO4A"
  },
  "type": "pay"
}
Decoded type: map[interface {}]interface {}
Decoded byte: map[by:[82 105 99 104 97 114 100 32 68 97 119 107 105 110 115] text:[116 104 101 32 115 111 108 117 116 105 111 110 32 111 102 116 101 110 32 116 117 114 110 115 32 111 117 116 32 109 111 114 101 32 98 101 97 117 116 105 102 117 108 32 116 104 97 110 32 116 104 101 32 112 117 122 122 108 101]]
Decoded text:
        by: Richard Dawkins
        text: the solution often turns out more beautiful than the puzzle
```

```
$ algorand find -t Y7XIQGCRU6IRLLCKSZJVTMQSH6HSDMDVW3WQFVOMNNBNFL5BO6KQ
Find Transaction Using Transaction ID
-------------------------------------
Found transaction in block: 365682
Transaction: {
  "round": 365682,
  "fee": 101,
  "first-round": 365680,
  "from": "MYPI256EXJQIMTV3NHX2BNVAPL7WRXCOOJ67X4WE536RSUMKEZVSP4IBI4",
  "last-round": 365682,
  "noteb64": "qUFsZ29yYW5kIQ==",
  "tx": "Y7XIQGCRU6IRLLCKSZJVTMQSH6HSDMDVW3WQFVOMNNBNFL5BO6KQ",
  "payment": {
    "amount": 1000,
    "to": "KI6TMKHUQOGJ7EDZLOWFOGHBBWBIMBMKONMS565X7NSOFMAM6S2EK4GBHQ"
  },
  "type": "pay"
}
Decoded type: []uint8
Decoded byte: [65 108 103 111 114 97 110 100 33]
Decoded text: Algorand!
```
