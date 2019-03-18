# Polaris

A distributed computations protocol for the next generation of the web.

## Installation

### Binary

```zsh
go install github.com/polaris-project/go-polaris
```

### Source Code

```zsh
cd $GOPATH && go get -u github.com/polaris-project/go-polaris
```

## Usage

### Joining the Testnet

```zsh
go-polaris
```

### Creating a Private Network

Note: Before creating a private network, make sure to create a `data/config/config_network_name.json` file.

```zsh
go-polaris --network your_network_name --bootstrap-address
```
