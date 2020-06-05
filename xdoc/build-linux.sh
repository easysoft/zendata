rm -rf build
mkdir build
mkdir build/log
cp -r data build/
cp -r demo build/

go-bindata -o=res/res.go -pkg=res res/ res/doc

GO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/zd-linux src/zd.go

cd build

cp zd-linux zd
tar -zcvf zd-linux-1.0.tar.gz zd data demo
rm zd

cd ..