# go-coral/contract

The seam between [go-coral](https://github.com/wxh06/go-coral), the Coral client, and whatever
seals its requests. It is types only, with no dependency beyond the standard library, and a module
of its own so that the client and a signing provider each depend on it without depending on the
other.

    go get github.com/wxh06/go-coral/contract

Implement `Provider` to plug a signer into the client, whether it runs in-process or calls a remote
service; consumers of the client do not import this module directly. The contract itself is in
the package documentation.

Tagged as `contract/vX.Y.Z`, independently of the client. Fields may be added to `Request` in a
minor version; construct it with keyed literals.

Licensed under the [BSD 3-Clause License](LICENSE).
