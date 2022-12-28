curl https://sh.rustup.rs -sSf | sh
source $HOME/.cargo/env
rustup default nightly
rustup target add wasm32-unknown-unknown

yum install openssl-devel gcc wget -y

export LD_LIBRARY_PATH=$(pwd)

cd ./lib/subapi/api && cargo clean && cargo build --release
cd ..
\cp -rf $(pwd)/target/release/libsubapi.* ../
\cp -rf $(pwd)/target/release/libsubapi.* ../../../

wget https://studygolang.com/dl/golang/go1.18.9.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.18.9.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin
export GOPATH=/opt/go

