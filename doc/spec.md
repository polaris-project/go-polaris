# Polaris Go Implementation Specifications

## Specs

### Third-Party Dependencies

#### Database

For all applicable instances, [boltdb](https://github.com/boltdb/bolt) will deemed the working database engine and client.

#### Sha3

For hashing via Sha3, Go's [/x/Sha3]("golang.org/x/crypto/sha3") package will be used.

#### Dependency Management

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

### Code Standards

#### TODO:

### User Accounts

Each user account is composed of an ecdsa private key. However, one should note that an "account" is not the same as an "address." An address is simply the 0x-prefix encoded byte value of a serialized ecdsa public key. Generally, an account's public key is derived from the account's private key (generated via `crypto/ecdsa`, `elliptic.p521()`, `crypto/rand`).

An ecdsa key pair is represented by the `Account` struct, that of which stores private in the form of a pointer to an `ecdsa.PrivateKey`.

#### Generation

As was mentioned earlier, each account's private key is derived from the Golang standard library ecdsa package `GenerateKey()` method. Generally, this `GenerateKey()` method will be used in conjunction with both a `crypto/rand` Reader (`rand.Reader` [make sure `crypto/rand` is imported, rather than `math/rand`]) and `elliptic.P521()`, though this can be substituted as seen fit.

#### Serialization

Accounts will be serialized via JSON and stored in a keystore path defined in the applicable common package (that of which is a child of the root Polaris data path specified in common).

Whilst writing an account to persistent memory via JSON, the account's private key will be temporarily set to nil, and reset after writing to memory. During the time that the ecdsa.privateKey pointer is set to nil, a field of `Account` will be set--`SerializedPrivateKey`. This `SerializedPrivateKey` field will take the value of an x509 encoded byte array value of the given private key. After having stored this serialized private key in the given account struct, the account will be written to persistent memory, and the value of `PrivateKey` subsequently reset to before writing to memory--`SerializedPrivateKey` should be set to nil after having reset `PrivateKey` to its previous state.

#### Addresses

Account addresses will--as has been stated earlier--be derived from the account public key. To obtain the account address, one simply hashes the x509 encoded byte value of the account public key via Polaris's crypto package `Sha3` method.

#### Code Structure

Additionally, all of the `Account` functionality will be written in its own package, rather than in `common` or `types` (preferably in a package called `accounts`).
