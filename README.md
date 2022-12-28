
## Subclient
subclient is a JSON RPC client of the substrate, it can scan blocks, transfer transactions and so on

## Example

    var networkId = 0
    var wsEndpoint = "wss://rpc.polkadot.io"
    var httpEndpoint = "https://rpc.polkadot.io"


	option := types.ClientOption{
		HttpEndpoint: httpEndpoint,
		WsEndpoint:   wsEndpoint,
		NetworkId:    networkId,
		WsSwitch:     true,
		Debug:        false,
	}
	client, err = NewClient(option)
	if err != nil {
		panic(err)
	}

    head, err := client.GetFinalHeight()
	if err != nil {
		panic(err)
	}
	fmt.Println(head)



## Docker Centos7 编译

    docker pull centos:centos7

    docker run -it -v /Users/abel/bworkspace/docker/:/opt/ --name centos7  -d centos:centos7 
   
    curl https://sh.rustup.rs -sSf | sh
    source $HOME/.cargo/env

    rustup default nightly
    rustup target add wasm32-unknown-unknown  

    yum install openssl-devel  gcc wget -y

    sh build.sh

    export LD_LIBRARY_PATH=$(pwd)

    wget https://studygolang.com/dl/golang/go1.18.9.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.18.9.linux-amd64.tar.gz

    vim /etc/environment
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=/opt/go
    source /etc/environment

    go build
