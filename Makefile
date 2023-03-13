VERSION=3.2.0
PROJECT=ztf

ifeq ($(OS),Windows_NT)
    PLATFORM="Windows"
else
    ifeq ($(shell uname),Darwin)
        PLATFORM="Mac"
    else
        PLATFORM="Unix"
    endif
endif

ifeq ($(PLATFORM),"Mac")
    QINIU_DIR=/Users/aaron/work/zentao/qiniu/
else
    QINIU_DIR=~/ztfZip
endif

QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/
SERVER_MAIN_FILE=cmd/server/main.go

COMMAND_MAIN_DIR=cmd/command/
COMMAND_MAIN_FILE=${COMMAND_MAIN_DIR}main.go

COMMAND_BIN_DIR=bin/
CLIENT_BIN_DIR=client/bin/
CLIENT_OUT_DIR=client/out/

BUILD_TIME=`git show -s --format=%cd`
GO_VERSION=`go version`
GIT_HASH=`git show -s --format=%H`
BUILD_CMD=go build -ldflags "-X 'main.AppVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}' -X 'main.GitHash=${GIT_HASH}'"
BUILD_CMD_WIN=go build -ldflags "-s -w -X 'main.AppVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}' -X 'main.GitHash=${GIT_HASH}'"

default: win64 win32 linux mac
server: server_win64 server_win32 server_linux server_mac

server_win64: prepare compile_server_mac copy_files_win64   zip_server_win64
server_win32: prepare compile_server_win32 copy_files_win32   zip_server_win32
server_linux: prepare compile_server_linux copy_files_linux   zip_server_linux
server_mac: prepare compile_server_mac copy_files_mac   zip_server_mac

win64: prepare compile_server_win64 package_gui_win64_client compile_launcher_win64 compile_command_win64 copy_files_win64 zip_server_win64 zip_client_win64
win32: prepare compile_server_win32 package_gui_win32_client compile_launcher_win32 compile_command_win32 copy_files_win32 zip_server_win32 zip_client_win32
linux: prepare compile_server_linux package_gui_linux_client                        compile_command_linux copy_files_linux zip_server_linux zip_client_linux
mac:   prepare compile_server_mac   package_gui_mac_client                          compile_command_mac   copy_files_mac   zip_server_mac   zip_client_mac

prepare: update_version prepare_res
update_version: update_version_in_config gen_version_file

update_version_in_config:
ifeq ($(PLATFORM),"Mac")
    @gsed -i "s/Version.*/Version = ${VERSION}/" conf/ztf.conf
else
	@sed -i "s/Version.*/Version = ${VERSION}/" conf/ztf.conf
endif

gen_version_file:
	@echo 'gen version'
	@mkdir -p ${QINIU_DIR}/${PROJECT}/
	@echo ${VERSION} > ${QINIU_DIR}/${PROJECT}/version.txt

compile_ui:
	@cd ui && yarn build --dest ../client/ui && cd ..

prepare_res:
	@echo 'start prepare res'
	@rm -rf res/res.go
	@go-bindata -o=res/res.go -pkg=res res/...

# launcher
compile_launcher_win64:
	@echo 'start compile win64 launcher'
	@cd cmd/launcher && \
        CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD} -x -v \
		-o ../../${COMMAND_BIN_DIR}win64/${PROJECT}-gui.exe && \
		cd ..

compile_launcher_win32:
	@echo 'start compile win32 launcher'
	@cd cmd/launcher && \
        CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD} -x -v \
		-o ../../${COMMAND_BIN_DIR}win32/${PROJECT}-gui.exe && \
        cd ..

# server
compile_server_win64:
	@echo 'start compile server win64'
	@rm -rf ${COMMAND_BIN_DIR}win64/${PROJECT}-server.exe
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD_WIN} -x -v \
		-o ${COMMAND_BIN_DIR}win64/${PROJECT}-server.exe ${SERVER_MAIN_FILE}

compile_server_win32:
	@echo 'start compile server win32'
	@rm -rf ${COMMAND_BIN_DIR}win32/${PROJECT}-server.exe
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD_WIN} -x -v \
		-o ${COMMAND_BIN_DIR}win32/${PROJECT}-server.exe ${SERVER_MAIN_FILE}

compile_server_linux:
	@echo 'start compile server linux'
	@rm -rf ${COMMAND_BIN_DIR}linux/${PROJECT}-server
ifeq ($(PLATFORM),"Mac")
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}linux/${PROJECT}-server ${SERVER_MAIN_FILE}
else
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}linux/${PROJECT}-server ${SERVER_MAIN_FILE}
endif

compile_server_mac:
	@echo 'start compile mac'
	@rm -rf ${COMMAND_BIN_DIR}darwin/${PROJECT}-server
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}darwin/${PROJECT}-server ${SERVER_MAIN_FILE}

# gui
package_gui_win64_client:
	@echo 'start package gui win64'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir ${CLIENT_BIN_DIR}win32
	@cp -rf ${COMMAND_BIN_DIR}win64/${PROJECT}-server.exe ${CLIENT_BIN_DIR}win32/${PROJECT}.exe

	@cd client && npm run package-win64 && cd ..
	@rm -rf ${CLIENT_OUT_DIR}win64 && mkdir ${CLIENT_OUT_DIR}win64 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-x64 ${CLIENT_OUT_DIR}win64/gui

package_gui_win32_client:
	@echo 'start package gui win32'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir ${CLIENT_BIN_DIR}win32
	@cp -rf ${COMMAND_BIN_DIR}win64/${PROJECT}-server.exe ${CLIENT_BIN_DIR}win32/${PROJECT}.exe

	@cd client && npm run package-win32 && cd ..
	@rm -rf ${CLIENT_OUT_DIR}win32 && mkdir ${CLIENT_OUT_DIR}win32 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-ia32 ${CLIENT_OUT_DIR}win32/gui

package_gui_linux_client:
	@echo 'start package gui linux'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir ${CLIENT_BIN_DIR}linux
	@cp -rf ${COMMAND_BIN_DIR}linux/${PROJECT}-server ${CLIENT_BIN_DIR}linux/${PROJECT}

	@cd client && npm run package-linux && cd ..
	@rm -rf ${CLIENT_OUT_DIR}linux && mkdir ${CLIENT_OUT_DIR}linux && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-linux-x64 ${CLIENT_OUT_DIR}linux/gui

package_gui_mac_client:
	@echo 'start package gui mac'
	@rm -rf ${CLIENT_BIN_DIR}/* && mkdir ${CLIENT_BIN_DIR}darwin
	@cp -rf ${COMMAND_BIN_DIR}darwin/${PROJECT}-server ${CLIENT_BIN_DIR}darwin/${PROJECT}

	@cd client && npm run package-mac && cd ..
	@rm -rf ${CLIENT_OUT_DIR}darwin && mkdir ${CLIENT_OUT_DIR}darwin && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-darwin-x64 ${CLIENT_OUT_DIR}darwin/gui && \
		mv ${CLIENT_OUT_DIR}darwin/gui/ztf.app ${CLIENT_OUT_DIR}darwin/ztf.app && rm -rf ${CLIENT_OUT_DIR}darwin/gui

# command line
compile_command_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD} -x -v \
		-o ${COMMAND_BIN_DIR}win64/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD} -x -v \
		-o ${COMMAND_BIN_DIR}win32/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_linux:
	@echo 'start compile linux'
ifeq ($(PLATFORM),"Mac")
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}linux/${PROJECT} ${COMMAND_MAIN_FILE}
else
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}linux/${PROJECT} ${COMMAND_MAIN_FILE}
endif

compile_command_linux_arm64:
	@echo 'start compile linux'

	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}linux/${PROJECT}_arm64 ${COMMAND_MAIN_FILE}

compile_command_mac:
	@echo 'start compile darwin'
	@echo
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}darwin/${PROJECT} ${COMMAND_MAIN_FILE}

copy_files_win64:
	@echo 'start copy files win64'
	@cp -r demo "${CLIENT_OUT_DIR}win64"
	@cp -r demo "${COMMAND_BIN_DIR}win64"
	@cp ${COMMAND_BIN_DIR}win64/ztf.exe "${CLIENT_OUT_DIR}win64"
	@cp ${COMMAND_BIN_DIR}win64/ztf-gui.exe "${CLIENT_OUT_DIR}win64"

copy_files_win32:
	@echo 'start copy files win32'
	@cp -r demo "${CLIENT_OUT_DIR}win32"
	@cp -r demo "${COMMAND_BIN_DIR}win32"
	@cp ${COMMAND_BIN_DIR}win32/ztf.exe "${CLIENT_OUT_DIR}win32"
	@cp ${COMMAND_BIN_DIR}win32/ztf-gui.exe "${CLIENT_OUT_DIR}win32"

copy_files_linux:
	@echo 'start copy files linux'
	@cp -r demo "${CLIENT_OUT_DIR}linux"
	@cp -r demo "${COMMAND_BIN_DIR}linux"
	@cp ${COMMAND_BIN_DIR}linux/ztf "${CLIENT_OUT_DIR}linux"

copy_files_mac:
	@echo 'start copy files darwin'
	@cp -r demo "${CLIENT_OUT_DIR}darwin"
	@cp -r demo "${COMMAND_BIN_DIR}darwin"
	@cp ${COMMAND_BIN_DIR}darwin/ztf "${CLIENT_OUT_DIR}darwin"

# zip server
zip_server_win64:
	@cd ${COMMAND_BIN_DIR}win64 && zip -ry ${QINIU_DIST_DIR}win64/${PROJECT}-server.zip ./demo ./${PROJECT}-server.exe && cd ../..
	@md5sum ${QINIU_DIST_DIR}win64/${PROJECT}-server.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}win64/${PROJECT}-server.zip.md5

zip_server_win32:
	@cd ${COMMAND_BIN_DIR}win32 && zip -ry ${QINIU_DIST_DIR}win32/${PROJECT}-server.zip ./demo ./${PROJECT}-server.exe && cd ../..
	@md5sum ${QINIU_DIST_DIR}win32/${PROJECT}-server.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}win32/${PROJECT}-server.zip.md5

zip_server_linux:
	@cd ${COMMAND_BIN_DIR}linux && zip -ry ${QINIU_DIST_DIR}linux/${PROJECT}-server.zip ./demo ./${PROJECT}-server && cd ../..
	@md5sum ${QINIU_DIST_DIR}linux/${PROJECT}-server.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}linux/${PROJECT}-server.zip.md5

zip_server_mac:
	@cd ${COMMAND_BIN_DIR}darwin && zip -ry ${QINIU_DIST_DIR}darwin/${PROJECT}-server.zip ./demo ./${PROJECT}-server && cd ../..
	@md5sum ${QINIU_DIST_DIR}darwin/${PROJECT}-server.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}darwin/${PROJECT}-server.zip.md5

# zip client
zip_client_win64:
	@echo 'start zip win64'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}win64 && rm -rf ${QINIU_DIST_DIR}win64/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}win64 && \
		zip -ry ${QINIU_DIST_DIR}win64/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}win64/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}win64/${PROJECT}.zip.md5 && \
        cd ../..; \

zip_client_win32:
	@echo 'start zip win32'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}win32 && rm -rf ${QINIU_DIST_DIR}win32/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}win32 && \
		zip -ry ${QINIU_DIST_DIR}win32/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}win32/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}win32/${PROJECT}.zip.md5 && \
        cd ../..; \

zip_client_linux:
	@echo 'start zip linux'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}linux && rm -rf ${QINIU_DIST_DIR}linux/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}linux && \
		zip -ry ${QINIU_DIST_DIR}linux/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}linux/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}linux/${PROJECT}.zip.md5 && \
        cd ../..; \

zip_client_mac:
	@echo 'start zip darwin'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}darwin && rm -rf ${QINIU_DIST_DIR}darwin/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}darwin && \
		zip -ry ${QINIU_DIST_DIR}darwin/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}darwin/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}darwin/${PROJECT}.zip.md5 && \
        cd ../..; \

upload_to:
	@echo 'upload...'
	@find ${QINIU_DIR} -name ".DS_Store" -type f -delete
	@qshell qupload2 --src-dir=${QINIU_DIR} --bucket=download --thread-count=10 --log-file=qshell.log \
					 --skip-path-prefixes=zd,zv,zmanager,driver,deeptest --rescan-local --overwrite --check-hash
