rm -rf build
mkdir build
mkdir build/log
cp -r data build/
cp -r demo build/

go-bindata -o=res/res.go -pkg=res res/ res/doc

CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o build/zd-mac src/zd.go