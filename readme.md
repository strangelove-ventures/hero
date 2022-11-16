# hero
**hero** is a blockchain built using Cosmos SDK and Tendermint

## Get started

Install [go](https://go.dev/dl/)

## Build and install to go bin path

```
make install
```

## Initialize config

Come up with a moniker for your node, then run:

```
herod init $MONIKER
```
 
 
 
## Launch with genesis file or run as standalone chain

To launch as a consumer chain, download and save shared genesis file to `~/.hero/config/genesis.json`. Additionally add peering information (`persistent_peers` or `seeds`) to `~/.hero/config/config.toml`

To instead launch as a standalone, single node chain, run:

```
herod add-consumer-section
```

## Launch node

```
herod start
```
