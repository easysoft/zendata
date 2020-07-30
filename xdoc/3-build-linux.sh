rm -rf bin
mkdir bin

go-bindata -o=res/res.go -pkg=res res/ res/en res/zh

GO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/zd-linux src/zd.go
scp bin/zd-linux* aaron@172.16.13.1:/Users/aaron/rd/project/zentao/go/zd/bin