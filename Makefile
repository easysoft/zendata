VERSION=3.0
PROJECT=zd
QINIU_DIR=/Users/aaron/work/zentao/qiniu/
QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/
PACKAGE=${PROJECT}-${VERSION}
BINARY=zd
BIN_DIR=bin
BIN_ZIP_DIR=${BIN_DIR}/zip/${PROJECT}/${VERSION}/
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
BIN_WIN64=${BIN_OUT}win64/
BIN_WIN32=${BIN_OUT}win32/
BIN_LINUX=${BIN_OUT}linux/
BIN_MAC=${BIN_OUT}darwin/

default: update_version_in_config gen_version_file prepare_res compile_all copy_files package

win64: update_version_in_config gen_version_file prepare_res compile_win64 copy_files package
win32: update_version_in_config gen_version_file prepare_res compile_win32 copy_files package
linux: update_version_in_config gen_version_file prepare_res compile_linux copy_files package
mac: update_version_in_config gen_version_file prepare_res compile_mac copy_files package
upload: upload_to

prepare_res:
	@echo 'start prepare res'
	@cp res/zh/sample.yaml demo/default.yaml
	@rm -rf ${BIN_DIR}

compile_all: compile_win64 compile_win32 compile_linux compile_mac

build_ui:
	@echo 'compile ui'
	@cd ui && yarn build && cd ..

compile_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" -o ${BIN_WIN64}${BINARY}.exe cmd/command/main.go

compile_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -x -v -ldflags "-s -w" -o ${BIN_WIN32}${BINARY}.exe cmd/command/main.go

compile_linux:
	@echo 'start compile linux'
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ go build -o ${BIN_LINUX}${BINARY} cmd/command/main.go

compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ${BIN_MAC}${BINARY} cmd/command/main.go

copy_files:
	@echo 'start copy files to ${BIN_DIR}'
	@cp -r {.zd.conf,data,yaml,users,demo,runtime} ${BIN_DIR}
	@rm -rf ${BIN_DIR}/demo/out ${BIN_DIR}/yaml/article/chinese/slang/out ${BIN_DIR}/runtime/protobuf/out

	@mkdir -p ${BIN_DIR}/tmp/cache && sqlite3 tmp/cache/.data.db ".backup '${BIN_DIR}/tmp/cache/.data.db'"
	@sqlite3 '${BIN_DIR}/tmp/cache/.data.db' ".read 'xdoc/clear-data.txt'"

	@for platform in `ls ${BIN_OUT}`; do cp -r {.zd.conf,bin/data,bin/runtime,bin/yaml,bin/users,bin/demo,bin/tmp} "${BIN_OUT}$${platform}"; done

	@rm -rf ${BIN_OUT}linux/runtime/php \
		${BIN_OUT}linux/runtime/protobuf/bin/mac \
		${BIN_OUT}linux/runtime/protobuf/bin/win*
	@rm -rf ${BIN_OUT}darwin/runtime/php \
		${BIN_OUT}darwin/runtime/protobuf/bin/linux \
		${BIN_OUT}darwin/runtime/protobuf/bin/win*
	@rm -rf ${BIN_OUT}win32/runtime/protobuf/bin/mac \
		${BIN_OUT}win32/runtime/protobuf/bin/linux \
		${BIN_OUT}win32/runtime/protobuf/bin/win64
	@rm -rf ${BIN_OUT}win64/runtime/protobuf/bin/mac \
		${BIN_OUT}win64/runtime/protobuf/bin/linux \
		${BIN_OUT}win64/runtime/protobuf/bin/win32

package:
	@echo 'start package'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for platform in `ls ${BIN_OUT}`; do mkdir -p ${QINIU_DIST_DIR}$${platform}; done

	@cd ${BIN_OUT} && \
		for platform in `ls ./`; \
			do  cd $${platform} && \
				zip -ry ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip ./* && \
				md5sum ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip | awk '{print $$1}' | \
					xargs echo > ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip.md5 && \
				cd ..; \
			done

update_version_in_config:
	@gsed -i "s/Version.*/Version = ${VERSION}/" .zd.conf

gen_version_file:
	@echo 'gen version'
	@mkdir -p ${QINIU_DIR}/${PROJECT}/
	@echo ${VERSION} > ${QINIU_DIR}/${PROJECT}/version.txt

upload_to:
	@echo 'upload...'
	@find ${QINIU_DIR} -name ".DS_Store" -type f -delete
	@qshell qupload2 --src-dir=${QINIU_DIR} --bucket=download --thread-count=10 --log-file=qshell.log \
                     --skip-path-prefixes=ztf,zv,zmanager,driver,deeptest --rescan-local --overwrite --check-hash
