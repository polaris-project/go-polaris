# Polaris

A distributed computations protocol for the next generation of the web.

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-%23516aa0.svg)](https://godoc.org/github.com/polaris-project/go-polaris)
[![Go Report Card](https://goreportcard.com/badge/github.com/polaris-project/go-polaris)](https://goreportcard.com/report/github.com/polaris-project/go-polaris)
[![Build Status](https://travis-ci.com/polaris-project/go-polaris.svg?branch=master)](https://travis-ci.com/polaris-project/go-polaris)
[![Gluten Status](https://img.shields.io/badge/gluten-free-brightgreen.svg)](https://img.shields.io/badge/gluten-free-brightgreen.svg)

## Installation

### Binary

```zsh
go install github.com/polaris-project/go-polaris
```

### Source Code

Go >= 1.11:

```zsh
go gopath-get -u github.com/polaris-project/go-polaris
```

Legacy Versions:

```zsh
go get -u github.com/polaris-project/go-polaris
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
