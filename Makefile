VERSION=3.1.0
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
BIN_DIR=bin
BIN_ZIP_DIR=${BIN_DIR}/zip/${PROJECT}/${VERSION}/
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
CLIENT_OUT_DIR=client/out/

CLIENT_BIN_DIR=client/bin/
CLIENT_OUT_DIR=client/out/
CLIENT_OUT_DIR_EXECUTABLE=${CLIENT_OUT_DIR}executable/
CLIENT_OUT_DIR_UPGRADE=${CLIENT_OUT_DIR}upgrade/

CLIENT_UI_DIR=client/ui/

COMMAND_MAIN_DIR=cmd/command/
COMMAND_MAIN_FILE=${COMMAND_MAIN_DIR}main.go

SERVER_MAIN_FILE=cmd/server/main.go

BUILD_TIME=`git show -s --format=%cd`
GO_VERSION=`go version`
GIT_HASH=`git show -s --format=%H`
BUILD_CMD_UNIX=go build -ldflags "-X 'main.AppVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}' -X 'main.GitHash=${GIT_HASH}'"
BUILD_CMD_WIN=go build -ldflags "-s -w -X 'main.AppVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}' -X 'main.GitHash=${GIT_HASH}'"

default: clear build_ui prepare_build compile_all copy_files package package_upgrade
clear:
	@rm -rf ${BIN_DIR}
	@rm -rf ${CLIENT_OUT_DIR}

prepare_build: clear update_version_in_config gen_version_file prepare_res

win64: prepare_build compile_launcher_win64 compile_server_win64 package_gui_win64_client compile_command_win64 copy_files package package_upgrade
win32: prepare_build compile_launcher_win32 compile_server_win32 package_gui_win32_client compile_command_win32 copy_files package package_upgrade
linux: prepare_build                        compile_server_linux package_gui_linux_client compile_command_linux copy_files package package_upgrade
mac: prepare_build                          compile_server_mac   package_gui_mac_client   compile_command_mac   copy_files package package_upgrade

compile_all: compile_win64 compile_win32 compile_linux compile_mac

compile_win64: compile_launcher_win64 compile_server_win64 package_gui_win64_client compile_command_win64
compile_win32: compile_launcher_win32 compile_server_win32 package_gui_win32_client compile_command_win32
compile_linux: compile_server_linux package_gui_linux_client compile_command_linux
compile_mac: compile_server_mac package_gui_mac_client compile_command_mac

upload: upload_to

prepare_res:
	@echo 'start prepare res'
	@cp res/zh/sample.yaml demo/default.yaml
	@rm -rf ${BIN_DIR}


build_ui:
	@echo 'compile ui'
	@cd ui && yarn build --dest ../client/ui && cd ..

compile_server_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD_WIN} -x -v \
 		-o ${BIN_DIR}/win64/server.exe ${SERVER_MAIN_FILE}

	@rm -rf "${CLIENT_OUT_DIR_UPGRADE}win64" && mkdir -p "${CLIENT_OUT_DIR_UPGRADE}win64" && \
  		cp ${BIN_DIR}/win64/server.exe "${CLIENT_OUT_DIR_UPGRADE}win64/"

compile_server_win32:
	@echo 'start compile server win32'
	@rm -rf ${BIN_DIR}/win32/server.exe
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD_WIN} -x -v \
		-o ${BIN_DIR}/win32/server.exe ${SERVER_MAIN_FILE}

	@rm -rf "${CLIENT_OUT_DIR_UPGRADE}win32" && mkdir -p "${CLIENT_OUT_DIR_UPGRADE}win32" && \
  		cp ${BIN_DIR}/win32/server.exe "${CLIENT_OUT_DIR_UPGRADE}win32/"

compile_server_linux:
	@echo 'start compile server linux'
	@rm -rf ${BIN_DIR}/linux/server
ifeq ($(OS),"Mac")
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}/linux/server ${SERVER_MAIN_FILE}
else
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}/linux/server ${SERVER_MAIN_FILE}
endif

	@rm -rf "${CLIENT_OUT_DIR_UPGRADE}linux" && mkdir -p "${CLIENT_OUT_DIR_UPGRADE}linux" && \
  		cp ${BIN_DIR}/linux/server "${CLIENT_OUT_DIR_UPGRADE}linux/"

compile_server_mac:
	@echo 'start compile server mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=gcc CXX=g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}/darwin/server ${SERVER_MAIN_FILE}

	@rm -rf "${CLIENT_OUT_DIR_UPGRADE}darwin" && mkdir -p "${CLIENT_OUT_DIR_UPGRADE}darwin" && \
  		cp ${BIN_DIR}/darwin/server "${CLIENT_OUT_DIR_UPGRADE}darwin/"

# gui
package_gui_win64_client:
	@echo 'start package gui win64'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir -p ${CLIENT_BIN_DIR}win32
	@cp -rf ${BIN_DIR}/win64/server.exe ${CLIENT_BIN_DIR}win32/server.exe

	@cd client  && npm run package-win64 && cd ..
	@rm -rf ${CLIENT_OUT_DIR_EXECUTABLE}win64 && mkdir -p ${CLIENT_OUT_DIR_EXECUTABLE}win64 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-x64 ${CLIENT_OUT_DIR_EXECUTABLE}win64/gui

package_gui_win32_client:
	@echo 'start package gui win32'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir -p ${CLIENT_BIN_DIR}win32
	@cp -rf ${BIN_DIR}/win32/server.exe ${CLIENT_BIN_DIR}win32/server.exe

	@cd client && npm run package-win32 && cd ..
	@rm -rf ${CLIENT_OUT_DIR_EXECUTABLE}win32 && mkdir -p ${CLIENT_OUT_DIR_EXECUTABLE}win32 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-ia32 ${CLIENT_OUT_DIR_EXECUTABLE}win32/gui

package_gui_linux_client:
	@echo 'start package gui linux'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir -p ${CLIENT_BIN_DIR}linux
	@cp -rf ${BIN_DIR}/linux/server ${CLIENT_BIN_DIR}linux/server

	@cd client && npm run package-linux && cd ..
	@rm -rf ${CLIENT_OUT_DIR_EXECUTABLE}linux && mkdir -p ${CLIENT_OUT_DIR_EXECUTABLE}linux && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-linux-x64 ${CLIENT_OUT_DIR_EXECUTABLE}linux/gui

package_gui_mac_client:
	@echo 'start package gui mac'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir -p ${CLIENT_BIN_DIR}darwin
	@cp -rf ${BIN_DIR}/darwin/server ${CLIENT_BIN_DIR}darwin/server

	@cd client && npm run package-mac && cd ..
	@rm -rf ${CLIENT_OUT_DIR_EXECUTABLE}darwin && mkdir -p ${CLIENT_OUT_DIR_EXECUTABLE}darwin && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-darwin-x64 ${CLIENT_OUT_DIR_EXECUTABLE}darwin/gui && \
		mv ${CLIENT_OUT_DIR_EXECUTABLE}darwin/gui/zd.app ${CLIENT_OUT_DIR_EXECUTABLE}darwin/zd.app && rm -rf ${CLIENT_OUT_DIR_EXECUTABLE}darwin/gui


# launcher
compile_launcher_win64:
	@echo 'start compile win64 launcher'
	@cd cmd/launcher && \
        CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD_WIN} -x -v \
		-o ../../${BIN_DIR}/win64/${PROJECT}-gui.exe && \
		cd ..

compile_launcher_win32:
	@echo 'start compile win32 launcher'
	@cd cmd/launcher && \
        CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD_WIN} -x -v \
		-o ../../${BIN_DIR}/win32/${PROJECT}-gui.exe && \
        cd ..

# command line
compile_command_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD_WIN} -x -v \
 		-o ${BIN_DIR}/win64/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD_WIN} -x -v \
 		-o ${BIN_DIR}/win32/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_linux:
	@echo 'start compile linux'
	@rm -rf ${BIN_DIR}/linux/${PROJECT}
ifeq ($(OS),"Mac")
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}/linux/${PROJECT} ${COMMAND_MAIN_FILE}
else
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}/linux/${PROJECT} ${COMMAND_MAIN_FILE}
endif

compile_command_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}/darwin/${PROJECT} ${COMMAND_MAIN_FILE}


copy_files:
	@echo 'start copy files to ${BIN_DIR}'
	@cp -r .zd.conf ${BIN_DIR}
	@cp -r data ${BIN_DIR}
	@cp -r yaml ${BIN_DIR}
	@cp -r users ${BIN_DIR}
	@cp -r demo ${BIN_DIR}
	@cp -r runtime ${BIN_DIR}
	@rm -rf ${BIN_DIR}/demo/out ${BIN_DIR}/yaml/article/chinese/slang/out ${BIN_DIR}/runtime/protobuf/out

	@rm -rf ${BIN_DIR}/tmp
	@mkdir -p ${BIN_DIR}/tmp/cache && sqlite3 tmp/cache/.data.db ".backup '${BIN_DIR}/tmp/cache/.data.db'"
	@sqlite3 '${BIN_DIR}/tmp/cache/.data.db' ".read 'xdoc/clear-data.txt'"

	@for platform in `ls ${CLIENT_OUT_DIR_EXECUTABLE}`;do \
		cp -r .zd.conf "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp -r bin/data "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp -r bin/runtime "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp -r bin/yaml "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp -r bin/users "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp -r bin/demo "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp -r bin/tmp "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}"; \
		cp ${BIN_DIR}/$${platform}/${PROJECT}.exe "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}" || true; \
		cp ${BIN_DIR}/$${platform}/${PROJECT} "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}" || true; \
		cp ${BIN_DIR}/$${platform}/${PROJECT}-gui.exe "${CLIENT_OUT_DIR_EXECUTABLE}$${platform}" || true; \
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
	@for platform in `ls ${CLIENT_OUT_DIR_EXECUTABLE}`; do mkdir -p ${QINIU_DIST_DIR}$${platform}; done

	@cd ${CLIENT_OUT_DIR_EXECUTABLE} && \
		for platform in `ls ./`; \
			do  cd $${platform} && \
				pwd; \
				zip -ry ${QINIU_DIST_DIR}$${platform}/${PROJECT}.zip ./* && \
				md5sum ${QINIU_DIST_DIR}$${platform}/${PROJECT}.zip | awk '{print $$1}' | \
					xargs echo > ${QINIU_DIST_DIR}$${platform}/${PROJECT}.zip.md5 && \
				cd ../; \
			done

package_upgrade:
	@echo 'start package upgrade'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for platform in `ls ${CLIENT_OUT_DIR_UPGRADE}`; do mkdir -p ${QINIU_DIST_DIR}$${platform}; done

	@cd ${CLIENT_OUT_DIR_UPGRADE} && \
		for platform in `ls ./`; \
			do  cd $${platform} && \
				cp -r ../../../ui ./; \
				zip -ry ${QINIU_DIST_DIR}$${platform}/${PROJECT}-upgrade.zip ./* && \
				md5sum ${QINIU_DIST_DIR}$${platform}/${PROJECT}-upgrade.zip | awk '{print $$1}' | \
					xargs echo > ${QINIU_DIST_DIR}$${platform}/${PROJECT}-upgrade.zip.md5 && \
				cd ../; \
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
