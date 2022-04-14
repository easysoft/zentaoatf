VERSION=3.0.0
PROJECT=deeptest
PACKAGE=${PROJECT}-${VERSION}
BINARY=ztf
MAIN_FILE=cmd/server/main.go
BIN_DIR=client/bin/
BIN_ZIP_DIR=${BIN_DIR}/zip/${PROJECT}/${VERSION}/
BIN_ZIP_RELAT=../../../zip/${PROJECT}/${VERSION}/
BIN_OUT=${BIN_DIR}/
BIN_WIN64=${BIN_OUT}win64/
BIN_WIN32=${BIN_OUT}win32/
BIN_LINUX=${BIN_OUT}linux/
BIN_MAC=${BIN_OUT}darwin/

BUILD_TIME=`git show -s --format=%cd`
GO_VERSION=`go version`
GIT_HASH=`git show -s --format=%H`
BUILD_CMD=go build -ldflags "-X 'commConsts.appVersion=${VERSION}' -X 'commConsts.buildTime=${BUILD_TIME}' -X 'commConsts.goVersion=${GO_VERSION}' -X 'commConsts.gitHash=${GIT_HASH}'"

default: prepare_res compile_all copy_files package

win64: prepare_res compile_win64 copy_files package
win32: prepare_res compile_win32 copy_files package
linux: prepare_res compile_linux copy_files package
mac: prepare_res compile_mac copy_files package

prepare_res:
	@echo 'start prepare res'
	@go-bindata -o=res/res.go -pkg=res res/...
	@rm -rf ${BIN_DIR}

compile_all: build_win64 compile_win32 compile_linux compile_mac

build_win64: compile_win64 package_win64_client
compile_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${BIN_WIN32}${BINARY}.exe ${MAIN_FILE}
package_win64_client:
	cd client && npm run package-win64 && cd ..

build_win32: compile_win32 package_win64_client
compile_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${BIN_WIN32}${BINARY}.exe ${MAIN_FILE}
package_win32_client:
	cd client && npm run package-win32 && cd ..

build_linux: compile_linux package_linux_client
compile_linux:
	@echo 'start compile linux'
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD}
		-o ${BIN_LINUX}${BINARY} ${MAIN_FILE}
package_linux_client:
	cd client && npm run package-linux && cd ..

build_mac: compile_mac package_mac_client
compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD} \
		-o ${BIN_MAC}${BINARY} ${MAIN_FILE}
package_mac_client:
	cd client && npm run package-mac && cd ..

copy_files:
	@echo 'start copy files'
	@cp -r {cmd/server/server.yml,cmd/server/perms.yml,cmd/server/rbac_model.conf} bin
	@for subdir in `ls ${BIN_OUT}`; \
	    do cp -r {bin/server.yml,bin/perms.yml,bin/rbac_model.conf} "${BIN_OUT}$${subdir}/ztf"; done

package:
	@echo 'start package'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for subdir in `ls ${BIN_OUT}`; do mkdir -p ${BIN_DIR}/zip/${PROJECT}/${VERSION}/$${subdir}; done

	@cd ${BIN_OUT} && \
		for subdir in `ls ./`; do cd $${subdir} && zip -r ${BIN_ZIP_RELAT}$${subdir}/${BINARY}.zip "${BINARY}" && cd ..; done
