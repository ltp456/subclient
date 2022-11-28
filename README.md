## Docker环境编译安装

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
    export LD_LIBRARY_PATH=/opt/chainmanagersv2/polkadex

    wget https://studygolang.com/dl/go1.16.9.linux-amd64.tar.gz

    tar -C /usr/local -xzf go1.16.9.linux-amd64.tar.gz

    export PATH=$PATH:/usr/local/go/bin

    go build


tips
    
    编译需要较大的内存，在docker环境有可能编译不通过，在测试服务器centos7上可以编译通过

    https://github.com/paritytech/substrate/issues/10857

rust版本

    rustc 1.60.0-nightly (09cb29c64 2022-02-15)






## SubstClient

    HTTP Endpoint: http://localhost:9933/
    Websocket Endpoint: ws://localhost:9944/


## Todo
* 关注metaData数据？
* 稳定性？
* 其它功能用go重写？


## refer
* https://docs.substrate.io/v3/runtime/custom-rpcs/
* https://github.com/itering/scale.go
* https://docs.substrate.io/v3/integration/client-libraries/




   

