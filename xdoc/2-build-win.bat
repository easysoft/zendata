rd /s/q build
mkdir build

go-bindata -o=res/res.go -pkg=res res/ res/en res/zh

SET CGO_ENABLED=1
SET GOOS=windows

SET GOARCH=386
go build -o build\zd-x86.exe src\zd.go

SET GOARCH=amd64
go build -o build\zd-amd64.exe src\zd.go

scp build\zd-* aaron@172.16.13.1:/Users/aaron/rd/project/zentao/go/zd/build