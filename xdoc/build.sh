rm -rf zdata
mkdir zdata
mkdir zdata/log
cp -r conf zdata/
cp -r runtime zdata/
cp -r demo zdata/

/Users/aaron/go/bin/go-bindata -o=res/res.go -pkg=res res/ res/doc

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o zdata/zdata-x86.exe src/zdata.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o zdata/zdata-amd64.exe src/zdata.go

GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zdata/zdata-linux src/zdata.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o zdata/zdata-mac src/zdata.go

cd zdata

cp zdata-x86.exe zdata.exe
zip -r zdata-win-x86-2.2.zip zdata.exe conf demo
rm zdata.exe

cp zdata-amd64.exe zdata.exe
zip -r zdata-win-amd64-2.2.zip zdata.exe conf demo
rm zdata.exe

cp zdata-linux zdata
tar -zcvf zdata-linux-2.2.tar.gz zdata conf demo
rm zdata

cp zdata-mac zdata
zip -r zdata-mac-2.2.zip zdata conf demo
rm zdata

cd ..