rm -rf build
mkdir build

go-bindata -o=res/res.go -pkg=res res/ res/en res/zh

GO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/zd-linux src/zd.go
scp build/zd-linux* aaron@172.16.13.1:/Users/aaron/rd/project/zentao/go/zd/build