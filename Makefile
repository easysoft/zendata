VERSION=1.2.0
PROJECT=zendata
PACKAGE=${PROJECT}-${VERSION}
BINARY=zd
BIN_DIR=bin
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
BIN_WIN64=${BIN_OUT}win64/
BIN_WIN32=${BIN_OUT}win32/
BIN_LINUX=${BIN_OUT}linux/
BIN_MAC=${BIN_OUT}mac/

default: prepare_res compile_all copy_files package

win64: prepare_res compile_win64 copy_files package
win32: prepare_res compile_win32 copy_files package
linux: prepare_res compile_linux copy_files package
mac: prepare_res compile_mac copy_files package

prepare_res:
	@echo 'start prepare res'
	@cp res/zh/sample.yaml demo/default.yaml
	@go-bindata -o=res/res.go -pkg=res res/ res/en res/zh
	@rm -rf ${BIN_DIR}

compile_all: compile_win64 compile_win32 compile_linux compile_mac

compile_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64
	@go build -x -v -ldflags "-s -w" -o ${BIN_WIN64}zd.exe src/zd.go

compile_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386
	@go build -x -v -ldflags "-s -w" -o ${BIN_WIN32}zd.exe src/zd.go

compile_linux:
	@echo 'start compile linux'
	@GO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ${BIN_LINUX}zd src/zd.go

compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ${BIN_MAC}zd src/zd.go

copy_files:
	@echo 'start copy files'
	@cp -r {data,demo} bin && rm -rf ${BIN_DIR}/demo/out

	@for subdir in `ls ${BIN_OUT}`; do cp -r {bin/data,bin/demo} "${BIN_OUT}$${subdir}"; done

package:
	@echo 'start package'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@cd ${BIN_DIR} && zip -r ${PACKAGE}.zip ${PROJECT}
	@cd ${BIN_DIR} && rm -rf ${PROJECT}