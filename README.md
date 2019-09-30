<p align="center">
  <img alt="open-rpc logo" src="https://github.com/open-rpc/design/blob/master/png/open-rpc-logo-320x320.png?raw=true" />
</p>

# [OpenRPC 1.0](https://spec.open-rpc.org/)

This package contains a golang implementation of the OpenRPC specification. It can currently generate server-side stubs based on a compliant
JSON document, but please be aware that the tool is still in its infancy and may not generate correctly. If you happen across any bugs or
realise any possible improvements, please consider submitting a PR!

```
go test ./...
go run example/example-proxy-server.go
```