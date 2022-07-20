VERSION=3.0.0_beta1
PROJECT=ztf
QINIU_DIR=/Users/aaron/work/zentao/qiniu/
QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/
SERVER_MAIN_FILE=cmd/server/main.go
COMMAND_MAIN_FILE=cmd/command/main.go

COMMAND_BIN_DIR=bin/
CLIENT_BIN_DIR=client/bin/
CLIENT_OUT_DIR=client/out/

BUILD_TIME=`git show -s --format=%cd`
GO_VERSION=`go version`
GIT_HASH=`git show -s --format=%H`
BUILD_CMD=go build -ldflags "-X 'commConsts.appVersion=${VERSION}' -X 'commConsts.buildTime=${BUILD_TIME}' -X 'commConsts.goVersion=${GO_VERSION}' -X 'commConsts.gitHash=${GIT_HASH}'"

default: win64 win32 linux mac

win64: prepare build_gui_win64 compile_command_win64 copy_files_win64 create_shortcut_win64 zip_win64
win32: prepare build_gui_win32 compile_command_win32 copy_files_win32 create_shortcut_win32 zip_win32
linux: prepare build_gui_linux compile_command_linux copy_files_linux create_shortcut_linux zip_linux
mac: prepare build_gui_mac compile_command_mac copy_files_mac create_shortcut_mac zip_mac

prepare: update_version prepare_res

update_version: update_version_in_config gen_version_file

update_version_in_config:
	@gsed -i "s/Version.*/Version = ${VERSION}/" conf/ztf.conf

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

# gui
build_gui_win64: compile_gui_win64 package_gui_win64_client
compile_gui_win64:
	@echo 'start compile win64'
	@rm -rf ./${CLIENT_BIN_DIR}/*
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${CLIENT_BIN_DIR}win32/${PROJECT}.exe ${SERVER_MAIN_FILE}
package_gui_win64_client:
	@cd client && npm run package-win64 && cd ..
	@rm -rf ${CLIENT_OUT_DIR}win64 && mkdir ${CLIENT_OUT_DIR}win64 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-x64 ${CLIENT_OUT_DIR}win64/gui

build_gui_win32: compile_gui_win32 package_gui_win32_client
compile_gui_win32:
	@echo 'start compile win32'
	@rm -rf ./${CLIENT_BIN_DIR}/*
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${CLIENT_BIN_DIR}win32/${PROJECT}.exe ${SERVER_MAIN_FILE}
package_gui_win32_client:
	@cd client && npm run package-win32 && cd ..
	@rm -rf ${CLIENT_OUT_DIR}win32 && mkdir ${CLIENT_OUT_DIR}win32 && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-win32-ia32 ${CLIENT_OUT_DIR}win32/gui

build_gui_linux: compile_gui_linux package_gui_linux_client
compile_gui_linux:
	@echo 'start compile linux'
	@rm -rf ./${CLIENT_BIN_DIR}/*
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD} \
		-o ${CLIENT_BIN_DIR}linux/${PROJECT} ${SERVER_MAIN_FILE}
package_gui_linux_client:
	@cd client && npm run package-linux && cd ..
	@rm -rf ${CLIENT_OUT_DIR}linux && mkdir ${CLIENT_OUT_DIR}linux && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-linux-x64 ${CLIENT_OUT_DIR}linux/gui

build_gui_mac: compile_gui_mac package_gui_mac_client
compile_gui_mac:
	@echo 'start compile mac'
	@rm -rf ./${CLIENT_BIN_DIR}/*
	@echo
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD} \
		-o ${CLIENT_BIN_DIR}darwin/${PROJECT} ${SERVER_MAIN_FILE}
package_gui_mac_client:
	@cd client && npm run package-mac && cd ..
	@rm -rf ${CLIENT_OUT_DIR}darwin && mkdir ${CLIENT_OUT_DIR}darwin && \
		mv ${CLIENT_OUT_DIR}${PROJECT}-darwin-x64 ${CLIENT_OUT_DIR}darwin/gui

# command line
compile_command_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${COMMAND_BIN_DIR}win64/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${COMMAND_BIN_DIR}win32/${PROJECT}.exe ${COMMAND_MAIN_FILE}

compile_command_linux:
	@echo 'start compile linux'
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}linux/${PROJECT} ${COMMAND_MAIN_FILE}

compile_command_mac:
	@echo 'start compile darwin'
	@echo
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD} \
		-o ${COMMAND_BIN_DIR}darwin/${PROJECT} ${COMMAND_MAIN_FILE}

copy_files_win64:
	@echo 'start copy files win64'
	cp -r demo "${CLIENT_OUT_DIR}win64"
	cp ${COMMAND_BIN_DIR}win64/* "${CLIENT_OUT_DIR}win64"

copy_files_win32:
	@echo 'start copy files win32'
	cp -r demo "${CLIENT_OUT_DIR}win32"
	cp ${COMMAND_BIN_DIR}win32/* "${CLIENT_OUT_DIR}win32"

copy_files_linux:
	@echo 'start copy files linux'
	cp -r demo "${CLIENT_OUT_DIR}linux"
	cp ${COMMAND_BIN_DIR}linux/* "${CLIENT_OUT_DIR}linux"

copy_files_mac:
	@echo 'start copy files darwin'
	cp -r demo "${CLIENT_OUT_DIR}darwin"
	cp ${COMMAND_BIN_DIR}darwin/* "${CLIENT_OUT_DIR}darwin"

create_shortcut_win64:
	@echo 'create shortcut win64'
	cp xdoc/ztf-gui.cmd "${CLIENT_OUT_DIR}win64"

create_shortcut_win32:
	@echo 'create shortcut win32'
	cp xdoc/ztf-gui.cmd "${CLIENT_OUT_DIR}win32"

create_shortcut_linux:
	@echo 'create shortcut linux'
	cp xdoc/ztf-gui-linux.sh ${CLIENT_OUT_DIR}linux/ztf-gui.sh

create_shortcut_mac:
	@echo 'create shortcut mac'
	cp xdoc/ztf-gui-mac.sh ${CLIENT_OUT_DIR}darwin/ztf-gui.sh

zip_win64:
	@echo 'start zip win64'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}win64 && rm -rf ${QINIU_DIST_DIR}win64/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}win64 && \
		zip -ry ${QINIU_DIST_DIR}win64/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}win64/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}win64/${PROJECT}.zip.md5 && \
        cd ../..; \

zip_win32:
	@echo 'start zip win32'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}win32 && rm -rf ${QINIU_DIST_DIR}win32/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}win32 && \
		zip -ry ${QINIU_DIST_DIR}win32/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}win32/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}win32/${PROJECT}.zip.md5 && \
        cd ../..; \

zip_linux:
	@echo 'start zip linux'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@mkdir -p ${QINIU_DIST_DIR}linux && rm -rf ${QINIU_DIST_DIR}linux/${PROJECT}.zip
	@cd ${CLIENT_OUT_DIR}linux && \
		zip -ry ${QINIU_DIST_DIR}linux/${PROJECT}.zip ./* && \
		md5sum ${QINIU_DIST_DIR}linux/${PROJECT}.zip | awk '{print $$1}' | \
			xargs echo > ${QINIU_DIST_DIR}linux/${PROJECT}.zip.md5 && \
        cd ../..; \

zip_mac:
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
					 --skip-path-prefixes=zz,zd,zmanager,driver --rescan-local --overwrite --check-hash
