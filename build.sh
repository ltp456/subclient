cd ./lib/subapi/api && cargo clean && cargo build --release
cd ..
\cp -rf $(pwd)/target/release/libsubapi.* ../
\cp -rf $(pwd)/target/release/libsubapi.* ../../../
cd ../../../

