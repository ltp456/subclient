cd ./lib/ccapsapi/api && cargo build --release
cd ..
\cp -rf $(pwd)/target/release/libccapsapi.* ../
\cp -rf $(pwd)/target/release/libccapsapi.* ../../../

