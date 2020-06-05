rm -rf build
mkdir build
mkdir build/log
cp -r data build/
cp -r demo build/

go-bindata -o=res/res.go -pkg=res res/ res/doc

CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o build/zd-mac src/zd.go

cd build

cp zd-mac zd
zip -r zd-mac-1.0.zip zd data demo
rm zd

cd ..