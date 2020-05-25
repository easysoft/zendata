rm -rf build
mkdir build
mkdir build/log
cp -r conf build/
cp -r data build/
cp -r demo build/

/Users/aaron/go/bin/go-bindata -o=res/res.go -pkg=res res/ res/doc

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/zd-x86.exe src/zd.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/zd-amd64.exe src/zd.go

GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/zd-linux src/zd.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/zd-mac src/zd.go

cd build

cp zd-x86.exe zd.exe
zip -r zd-win-x86-2.2.zip zd.exe conf data demo
rm zd.exe

cp zd-amd64.exe zd.exe
zip -r zd-win-amd64-2.2.zip zd.exe conf data demo
rm zd.exe

cp zd-linux zd
tar -zcvf zd-linux-2.2.tar.gz zd conf data demo
rm zd

cp zd-mac zd
zip -r zd-mac-2.2.zip zd conf data demo
rm zd

cd ..