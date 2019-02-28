# Polaris Go Implementation Specifications

## Third-Party Dependencies

### Database

For all applicable instances, [boltdb](https://github.com/boltdb/bolt) will deemed the working database engine and client.

### Sha3

For hashing via Sha3, Go's [/x/Sha3]("golang.org/x/crypto/sha3") package will be used.

### Dependency Management

To manage external dependencies, use Go's [module system](https://github.com/golang/go/wiki/Modules). To enable modules in Go...

```zsh
# enable go modules: https://github.com/golang/go/wiki/Modules
export GO111MODULE=on
```

To download all necessary dependencies into the vendor/ folder...

```zsh
# download dependencies to vendor folder
go mod vendor
```

## Code Standards

All code is documented in the [godoc-recognizable](https://blog.golang.org/godoc-documenting-go-code) format, that of which specifies each exported function should be preceded by a comment in the following format:

```Go
// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
```

Additionally, all package names--regardless of the file they're found in (but must be a .go file)--must be documented in the following format:

```Go
// Package sort provides primitives for sorting slices and user-defined
// collections.
package sort
```

Finally, as would be found in the godoc "documenting Go code spec", all known bugs should be documented in the following format:

```Go
// BUG(r): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
```

In the aforementioned example, "r" would be replaced with the name or username of an individual responsible for or knowing of the bug.

## User Accounts

Each user account is composed of an ecdsa private key. However, one should note that an "account" is not the same as an "address." An address is simply the 0x-prefix encoded sha3 hash of the byte value of a serialized ecdsa public key. Generally, an account's public key is derived from the account's private key (generated via `crypto/ecdsa`, `elliptic.p521()`, `crypto/rand`).

An ecdsa key pair is represented by the `Account` struct, that of which stores a private key in the form of a pointer to an `ecdsa.PrivateKey`. Additionally, an account should also store a field called `SerializedPrivateKey`, that of which should be empty at all times except when the account is written to persistent memory.

### Generation

As was mentioned earlier, each account's private key is derived from the Golang standard library ecdsa package `GenerateKey()` method. Generally, this `GenerateKey()` method will be used in conjunction with both a `crypto/rand` Reader (`rand.Reader` [make sure `crypto/rand` is imported, rather than `math/rand`]) and `elliptic.P521()`, though this can be substituted as seen fit.

### Serialization

Accounts will be serialized via JSON and stored in a keystore path defined in the applicable common package (a child of the root Polaris data path specified in common).

Whilst writing an account to persistent memory via JSON, the account's private key will be temporarily set to nil, and reset to its previous state after writing to memory. During the time that the ecdsa.privateKey pointer is set to nil, a field of `Account` will be set--`SerializedPrivateKey`. This `SerializedPrivateKey` field will take the value of an x509 encoded byte array value of the given private key. After having stored this serialized private key in the given account struct, the account will be written to persistent memory, and the value of `PrivateKey` subsequently reset to before writing to memory--`SerializedPrivateKey` should be set to nil after having reset `PrivateKey` to its previous state.

Whilst reading an account from persistent memory via JSON, the account's serialized private key will be recovered from the given `"account\_{address}.json"` file (i.e. `0x000` => `"account_0x000.json"`), that of which should be deserialized into a pointer to an ecdsa private key. After having recovered the private key pointer from the serialized private key, the serialized private key should be set to nil, and the `PrivateKey` field be set to the deserialized private key (actual `ecdsa.PrivateKey` instance pointer).

### Addresses

Account addresses will--as has been stated earlier--be derived from the account public key. To obtain the account address, one simply hashes the x509 encoded byte value of the account public key via Polaris's crypto package `Sha3` method.

### Code Structure

Additionally, all of the `Account` functionality will be written in its own package, rather than in `common` or `types` (preferably in a package called `accounts`). Therefore, as the accounts logic will be written in its own package, all of the account related logic should be written in a folder called "accounts".

## Dag

Unlike many other cryptographically secured digital currencies, Polaris is based around a directed-acyclic-graph data structure. This "dag" structure is composed of a list of transactions where each transaction contains the hash of its parents (previous transactions).

No transaction--with the exception of the genesis--can have a parent transaction value of nil.

Each entry into the acyclic graph will be treated as an entry into the dag's respective database. Additionally, all of the dag-related logic should take place in the types package. To ensure that the dag never reaches a size that is not indexable, the dag will not be treated as a strict slice of transaction pointers, but simply a key-value database instance, that of which will operate on [boltdb](https://github.com/boltdb/bolt).

The only piece of information that the `dag.go` `Dag` struct will serve and store will be the hash of the genesis transaction, the dag's `DagConfig` pointer (contains supply allocation information and other metadata), and the dag length (should be stored as a pointer to a big integer). The dag's config stores an "identifier" that will be used to open a new database, as well write to memory (i.e. db stored under folder with name equivalent to identifier). The working dag db instance will be stored as a global variable in `dag.go` (all methods needing access to the dag db should not open a new db, but use the `dag.go` `WorkingDagDB` variable as a pointer to the db instance). In addition, the `dag.go` initializing pseudo-constructor method should accept a `DagConfig` struct instance pointer, that of which will provide the dag version, the genesis transaction information (`Alloc` address => float64 map), and the string dag identifier.

## Transaction

As is the case with most other digital currency networks, a transaction is an atomic action on the network, that of which can represent the transfer of data or monetary value. In the case of Polaris, either may be true.

### Transaction Fields

A single `Transaction` consists of the following fields:

| Field              | Value                                                                                                                      | Type             |
| ------------------ | -------------------------------------------------------------------------------------------------------------------------- | ---------------- |
| AccountNonce       | Transaction index in account list of transactions.                                                                         | uint64           |
| Amount             | Transaction value.                                                                                                         | \*big.Int        |
| Sender             | Transaction sender address.                                                                                                | \*common.Address |
| Recipient          | Transaction recipient address.                                                                                             | \*common.Address |
| ParentTransactions | Parent transaction hashes (usually 1, but in the case of a poorly synchronized network, may be more).                      | []common.Hash    |
| GasPrice           | Amount of polaris willing to pay per single unit of gas (in increments of 0.000000001 polaris).                            | \*big.Int        |
| GasLimit           | Amount of gas willing to pay at max.                                                                                       | uint64           |
| Payload            | Data sent with transaction (i.e. contract bytecode, message, etc...)                                                       | []byte           |
| Signature          | ECDSA sender signature.                                                                                                    | \*Signature      |
| Hash               | Transaction hash including transaction signature (if set). To verify, exclude signature from tx hash as message to verify. | common.Hash      |

## Transaction Signatures

Transactions are signed via ECDSA. Furthermore, all of the signature-related logic has already been written, and is located in types/transaction_signature.go. Finally, each transaction contains a pointer to a signature struct instance.

### Signature Fields

A single `Signature` consists of the following fields:

| Field | Value                                                                                  | Type      |
| ----- | -------------------------------------------------------------------------------------- | --------- |
| V     | Signature value (in the case of a transaction, the hash before setting the signature). | []byte    |
| R     | Signature recovery value.                                                              | \*big.Int |
| S     | Signature recovery value.                                                              | \*big.Int |
