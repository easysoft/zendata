rm -rf build
mkdir build
mkdir build/log
cp -r data build/
cp -r demo build/

go-bindata -o=res/res.go -pkg=res res/ res/doc

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/zd-x86.exe src/zd.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/zd-amd64.exe src/zd.go

cd build

cp zd-x86.exe zd.exe
zip -r zd-win-x86-1.0.zip zd.exe data demo
rm zd.exe

cp zd-amd64.exe zd.exe
zip -r zd-win-amd64-1.0.zip zd.exe data demo
rm zd.exe

cd ..