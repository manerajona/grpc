# grpc

## Install protocol buffer on linux

```sh
sudo apt install -y protobuf-compiler
```
## Install protocol buffer from source code

1. Download the latest release as zip based on your os/arch:

> https://github.com/protocolbuffers/protobuf/releases

3. Extract zip:

```sh
unzip protoc-<version>-<os>-<arch>.zip -d $HOME/.protoc
```
3. Update environment’s path variable:

```sh
export PATH="$PATH:$HOME/.protoc/bin"
```
4. check version:

```sh
protoc --version
```

## Install plugins

1. Install the protocol compiler plugins for Go:

(Check version tags [here](https://grpc.io/docs/languages/go/quickstart))

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

2. Update your PATH so that the protoc compiler can find the plugins:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Generate gRPC code

1. Move to the directory containing the .proc file:
```sh
cd <directory>
```

2. Run the command bellow
```sh
protoc -I . --go_out . --go_opt module=github.com/manerajona/grpc \
--go-grpc_out require_unimplemented_servers=false:. --go-grpc_opt module=github.com/manerajona/grpc *.proto
```

## Generate ssl cert

1. Check if you have OpenSSL installed:

```shell
openssl version
```

2. Check if the .sh file has exec permission and launch it:

```sh
chmod +x ssl.sh
./sslgen.sh
```
