mkdir -p build
cd build
#echo "brew --prefix curl"
#brew --prefix curl
#echo "brew ls --verbose curl"
#brew ls --verbose curl
# for normal install use following command
cmake -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ ..
make
sudo make install

cd ../..