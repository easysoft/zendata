VERSION=2.3
PROJECT=zd

ifeq ($(OS),Windows_NT)
    OS="Windows"
else
    ifeq ($(shell uname),Darwin)
        OS="Mac"
    else
        OS="Unix"
    endif
endif

ifeq ($(OS),"Mac")
    QINIU_DIR=/Users/aaron/work/zentao/qiniu/
else
    QINIU_DIR=~/zentao/
endif

QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/
PACKAGE=${PROJECT}-${VERSION}
BINARY=zd
BIN_DIR=bin
BIN_ZIP_DIR=${BIN_DIR}/zip/${PROJECT}/${VERSION}/
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
CLIENT_OUT_DIR=client/out/

BIN_WIN64=${BIN_DIR}win64/
BIN_WIN32=${BIN_DIR}win32/
BIN_LINUX=${BIN_DIR}linux/
BIN_MAC=${BIN_DIR}darwin/

COMMAND_BIN_DIR=bin/
CLIENT_BIN_DIR=client/bin/
CLIENT_OUT_DIR=client/out/
CLIENT_OUT_DIR_EXECUTABLE=${CLIENT_OUT_DIR}executable/
CLIENT_OUT_DIR_UPGRADE=${CLIENT_OUT_DIR}upgrade/

COMMAND_MAIN_DIR=cmd/command/
COMMAND_MAIN_FILE=${COMMAND_MAIN_DIR}main.go

SERVER_MAIN_FILE=cmd/server/main.go

default: update_version_in_config gen_version_file prepare_res compile_all copy_files package

win64: update_version_in_config gen_version_file prepare_res compile_win64 package_gui_win64_client compile_command_win64 copy_files package
win32: update_version_in_config gen_version_file prepare_res compile_win32 package_gui_win32_client compile_command_win32 copy_files package
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

compile_client_ui:
	@cd ui && UI_IN_CLIENT=1 yarn build --dest ../client/ui && cd ..

compile_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" -o ${COMMAND_BIN_DIR}win64/${PROJECT}-server.exe ${SERVER_MAIN_FILE}

compile_win32:
	@echo 'start compile server win32'
	@rm -rf ${COMMAND_BIN_DIR}win32/${PROJECT}-server.exe
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -ldflags "-s -w" -x -v -o ${COMMAND_BIN_DIR}win32/${PROJECT}-server.exe ${SERVER_MAIN_FILE}

compile_linux:
	@echo 'start compile linux'
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ go build -o ${BIN_LINUX}${BINARY} cmd/command/main.go

compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ${BIN_MAC}${BINARY} cmd/command/main.go

# gui
package_gui_win64_client:
	@echo 'start package gui win64'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir ${CLIENT_BIN_DIR}win32
	@cp -rf ${COMMAND_BIN_DIR}win64/${PROJECT}-server.exe ${CLIENT_BIN_DIR}win32/${PROJECT}.exe

	@cd client  && npm run package-win64 && cd ..
	@rm -rf ${CLIENT_OUT_DIR}win64 && mkdir ${CLIENT_OUT_DIR}win64 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-x64 ${CLIENT_OUT_DIR}win64/gui

package_gui_win32_client:
	@echo 'start package gui win32'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir -p ${CLIENT_BIN_DIR}win32
	@cp -rf ${COMMAND_BIN_DIR}win32/${PROJECT}-server.exe ${CLIENT_BIN_DIR}win32/${PROJECT}.exe

	@cd client && npm run package-win32 && cd ..
	@rm -rf ${CLIENT_OUT_DIR}win32 && mkdir -p ${CLIENT_OUT_DIR}win32 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-ia32 ${CLIENT_OUT_DIR}win32/gui

# command line
compile_command_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		go build -x -v -ldflags "-s -w" -o ${COMMAND_BIN_DIR}win64/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		go build -x -v -ldflags "-s -w" -o ${COMMAND_BIN_DIR}win32/${PROJECT}.exe ${COMMAND_MAIN_FILE}

copy_files:
	@echo 'start copy files to ${BIN_DIR}'
	@cp -r .zd.conf ${BIN_DIR}
	@cp -r data ${BIN_DIR}
	@cp -r yaml ${BIN_DIR}
	@cp -r users ${BIN_DIR}
	@cp -r demo ${BIN_DIR}
	@cp -r runtime ${BIN_DIR}
	@rm -rf ${BIN_DIR}/demo/out ${BIN_DIR}/yaml/article/chinese/slang/out ${BIN_DIR}/runtime/protobuf/out

	@mkdir -p ${BIN_DIR}/tmp/cache && sqlite3 tmp/cache/.data.db ".backup '${BIN_DIR}/tmp/cache/.data.db'"
	@sqlite3 '${BIN_DIR}/tmp/cache/.data.db' ".read 'xdoc/clear-data.txt'"

	@for platform in `ls ${CLIENT_OUT_DIR}`;do \
		cp -r .zd.conf "${CLIENT_OUT_DIR}$${platform}"; \
		cp -r bin/data "${CLIENT_OUT_DIR}$${platform}"; \
		cp -r bin/runtime "${CLIENT_OUT_DIR}$${platform}"; \
		cp -r bin/users "${CLIENT_OUT_DIR}$${platform}"; \
		cp -r bin/demo "${CLIENT_OUT_DIR}$${platform}"; \
		cp -r bin/tmp "${CLIENT_OUT_DIR}$${platform}"; \
		true || cp ${COMMAND_BIN_DIR}$${platform}/zd.exe "${CLIENT_OUT_DIR}$${platform}"; \
		# cp ${COMMAND_BIN_DIR}$${platform}/zd-gui.exe "${CLIENT_OUT_DIR}$${platform}"; \
	done

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
	@for platform in `ls ${CLIENT_OUT_DIR}`; do mkdir -p ${QINIU_DIST_DIR}$${platform}; done

	@cd ${CLIENT_OUT_DIR} && \
		for platform in `ls ./`; \
			do  cd $${platform} && \
				zip -ry ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip ./* && \
				md5sum ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip | awk '{print $$1}' | \
					xargs echo > ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip.md5 && \
				cd ..; \
			done

update_version_in_config:
ifeq ($(OS),"Mac")
	@gsed -i "s/Version.*/Version = ${VERSION}/" .zd.conf
else
	@sed -i "s/Version.*/Version = ${VERSION}/" .zd.conf
endif

gen_version_file:
	@echo 'gen version'
	@mkdir -p ${QINIU_DIR}/${PROJECT}/
	@echo ${VERSION} > ${QINIU_DIR}/${PROJECT}/version.txt

upload_to:
	@echo 'upload...'
	@find ${QINIU_DIR} -name ".DS_Store" -type f -delete
	@qshell qupload2 --src-dir=${QINIU_DIR} --bucket=download --thread-count=10 --log-file=qshell.log \
                     --skip-path-prefixes=ztf,zv,zmanager,driver,deeptest --rescan-local --overwrite --check-hash
