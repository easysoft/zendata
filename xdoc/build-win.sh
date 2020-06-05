rm -rf build
mkdir build
mkdir build/log
cp -r data build/
cp -r demo build/

go-bindata -o=res/res.go -pkg=res res/ res/doc

SET CGO_ENABLED=1
SET GOOS=windows

SET GOARCH=386
go build -o build\zd-x86.exe src/zd.go

SET GOARCH=amd64
go build -o build\zd-amd64.exe src/zd.go