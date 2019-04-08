Algorand is a CLI application that abstracts low-level Algorand node operations.

## Usage

Simply run `algorand <COMMAND> <FLAGS>`. Running `algorand` or `algorand -h` will show usage and a list of commands.

```
Usage:
  algorand [command] [flags]
```

## Configuration

To configure `algorand`, a `config.yml` file can be passed with the --config flag. Configuration read in through a file will overwrite the same configuration specified by a flag. If no config file is passed, and no flags are set, reasonable defaults will be used.

```yml
hostname: '127.0.0.1'                               # the Algorand node's IP
algod-port: '8080'                                  # the `algod' process port
kmd-port: '7833'                                    # the `kmd' process port
algod-token: '374de74fa794248762e5ac17c8b39f19a05'  # the `algod' process token
kmd-token: 'be84aa55f61665645ed680b27ee15f11653'    # the `kmd' process token
log-level: 'INFO'                                   # the logging level
```

## Examples

### status:
```
$ algorand status -c ./config.yml
Made an algod client
Made a kmd client
algod: algod.Client, kmd: kmd.Client
algod last round: 331484
algod time since last round: 2432713600
algod catchup: 0
algod latest version: v4
```
