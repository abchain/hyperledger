Use following line to update protos:
```

protoc --proto_path=<GOPATH>/src/hyperledger.abchain.org/protos --proto_path=<CWD> --go_out=. *.proto

```