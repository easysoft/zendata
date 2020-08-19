VERSION=1.2.0
PROJECT=zendata
PACKAGE=${PROJECT}-${VERSION}
BINARY=zd
BIN_DIR=bin
BIN_ZIP_DIR=${BIN_DIR}/zip/
BIN_ZIP_RELAT=../../../zip/
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
BIN_WIN64=${BIN_OUT}win64/zd/
BIN_WIN32=${BIN_OUT}win32/zd/
BIN_LINUX=${BIN_OUT}linux/zd/
BIN_MAC=${BIN_OUT}mac/zd/

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
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" -o ${BIN_WIN64}zd.exe src/zd.go

compile_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -x -v -ldflags "-s -w" -o ${BIN_WIN32}zd.exe src/zd.go

compile_linux:
	@echo 'start compile linux'
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ go build -o ${BIN_LINUX}zd src/zd.go

compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ${BIN_MAC}zd src/zd.go

copy_files:
	@echo 'start copy files'
	@cp -r {data,yaml,users,demo,tmp} bin && rm -rf ${BIN_DIR}/demo/output

	@for subdir in `ls ${BIN_OUT}`; do cp -r {bin/data,bin/yaml,bin/users,bin/demo,bin/tmp} "${BIN_OUT}$${subdir}/zd"; done

package:
	@echo 'start package'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for subdir in `ls ${BIN_OUT}`; do mkdir -p ${BIN_DIR}/zip/$${subdir}; done

	@cd ${BIN_OUT} && \
		for subdir in `ls ./`; do cd $${subdir} && zip -r ${BIN_ZIP_RELAT}$${subdir}/${BINARY}.zip "${BINARY}" && cd ..; done
	#@cd ${BIN_ZIP_DIR} && zip -r ${PACKAGE}.zip ./
	#@cd ${BIN_DIR} && rm -rf ${PROJECT}
