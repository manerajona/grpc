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

## Create code in go from .proto file
```sh
$ protoc -I echo echo.proto --go_out=plugins=grpc:echo
$ protoc -I echo chat.proto --go_out=plugins=grpc:chat
```