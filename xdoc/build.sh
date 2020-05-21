rm -rf zt
mkdir zt
mkdir zt/log
cp -r conf zt/
cp -r runtime zt/
cp -r demo zt/

/Users/aaron/go/bin/go-bindata -o=res/res.go -pkg=res res/ res/doc

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o zt/zt-x86.exe src/zt.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o zt/zt-amd64.exe src/zt.go

GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zt/zt-linux src/zt.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o zt/zt-mac src/zt.go

cd zt

cp zt-x86.exe zt.exe
zip -r zt-win-x86-2.2.zip zt.exe conf demo
rm zt.exe

cp zt-amd64.exe zt.exe
zip -r zt-win-amd64-2.2.zip zt.exe conf demo
rm zt.exe

cp zt-linux zt
tar -zcvf zt-linux-2.2.tar.gz zt conf demo
rm zt

cp zt-mac zt
zip -r zt-mac-2.2.zip zt conf demo
rm zt

cd ..