VERSION=1.0.0
PROJECT=deeptest
PACKAGE=${PROJECT}-${VERSION}
BINARY=deeptest
BIN_DIR=bin
BIN_ZIP_DIR=${BIN_DIR}/zip/${PROJECT}/${VERSION}/
BIN_ZIP_RELAT=../../../zip/${PROJECT}/${VERSION}/
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
BIN_WIN64=${BIN_OUT}win64/zagent-server/
BIN_WIN32=${BIN_OUT}win32/zagent-server/
BIN_LINUX=${BIN_OUT}linux/zagent-server/
BIN_MAC=${BIN_OUT}mac/zagent-server/

default: prepare_res compile_all copy_files package

win64: prepare_res compile_win64 copy_files package
win32: prepare_res compile_win32 copy_files package
linux: prepare_res compile_linux copy_files package
mac: prepare_res compile_mac copy_files package

prepare_res:
	@echo 'start prepare res'
	@go-bindata -o=res/res.go -pkg=res res/...
	@rm -rf ${BIN_DIR}

compile_all: compile_win64 compile_win32 compile_linux compile_mac

compile_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BIN_WIN64}zagent-server.exe cmd/server/main.go

compile_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ${BIN_WIN32}zagent-server.exe cmd/server/main.go

compile_linux:
	@echo 'start compile linux'
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BIN_LINUX}zagent-server cmd/server/main.go

compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BIN_MAC}zagent-server cmd/server/main.go

copy_files:
	@echo 'start copy files'
	@cp -r {cmd/server/server.yml,cmd/server/perms.yml,cmd/server/rbac_model.conf} bin
	@for subdir in `ls ${BIN_OUT}`; \
	    do cp -r {bin/server.yml,bin/perms.yml,bin/rbac_model.conf} "${BIN_OUT}$${subdir}/zagent-server"; done

package:
	@echo 'start package'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for subdir in `ls ${BIN_OUT}`; do mkdir -p ${BIN_DIR}/zip/${PROJECT}/${VERSION}/$${subdir}; done

	@cd ${BIN_OUT} && \
		for subdir in `ls ./`; do cd $${subdir} && zip -r ${BIN_ZIP_RELAT}$${subdir}/${BINARY}.zip "${BINARY}" && cd ..; done
