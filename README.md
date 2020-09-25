# grpc

## install protocol buffer linux
```sh
$ sudo apt install -y protobuf-compiler
$ protoc --version
```

## install protocol buffer go
```sh
$ export GOBIN=$HOME/go/bin
$ cd $GOBIN
$ go install github.com/golang/protobuf/protoc-gen-go
```