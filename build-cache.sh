cd ./lib/subapi/api && cargo build --release
cd ..
\cp -rf $(pwd)/target/release/libsubapi.* ../
\cp -rf $(pwd)/target/release/libsubapi.* ../../../
cd ../../../
