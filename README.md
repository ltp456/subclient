## Docker Centos7 编译

    docker pull centos:centos7

    docker run -it -v /Users/abel/bworkspace/docker/:/opt/ --name centos7  -d centos:centos7 
   
    curl https://sh.rustup.rs -sSf | sh 
    // 有个安装选项直接默认回车即可
    source $HOME/.cargo/env

    rustup default nightly
    rustup target add wasm32-unknown-unknown  

    yum install openssl-devel -y
    yum install gcc -y
    yum install wget -y

    // 进入到subclient目录
    sh build.sh

    // 将so文件拷贝到当前的目录，并配置当前的环境变量 
    export LD_LIBRARY_PATH=/opt/subclient

    wget https://studygolang.com/dl/golang/go1.18.9.linux-amd64.tar.gz

    tar -C /usr/local -xzf go1.18.9.linux-amd64.tar.gz

    vim /etc/environment

    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=/opt/go

    source /etc/environment

    go build

## Example

    var networkId = []byte{0}
    var wsEndpoint = "wss://rpc.polkadot.io"
    var httpEndpoint = "wss://rpc.polkadot.io"

    option := types.ClientOption{
            HttpEndpoint: httpEndpoint,
            WsEndpoint:   wsEndpoint,
            NetworkId:    networkId,
            WsSwitch:     true,
        }
    client, err := NewClient(option)
    if err != nil {
            panic(err)
    }

    head, err := client.chainGetHead()
	if err != nil {
		panic(err)
	}
	fmt.Println(head)

        


   

