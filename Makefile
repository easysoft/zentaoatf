VERSION=2.5
PROJECT=ztf
QINIU_DIR=/Users/aaron/work/zentao/qiniu/
QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/
PACKAGE=${PROJECT}-${VERSION}
BINARY=ztf
BIN_DIR=bin
BIN_ZIP_DIR=${BIN_DIR}/zip/${PROJECT}/${VERSION}/
BIN_OUT=${BIN_DIR}/${PROJECT}/${VERSION}/
BIN_WIN64=${BIN_OUT}win64/${BINARY}/
BIN_WIN32=${BIN_OUT}win32/${BINARY}/
BIN_LINUX=${BIN_OUT}linux/${BINARY}/
BIN_MAC=${BIN_OUT}mac/${BINARY}/
BIN_ARM=${BIN_OUT}arm/${BINARY}/

default: update_version_in_config gen_version_file prepare_res compile_all copy_files package

win64: update_version_in_config gen_version_file prepare_res compile_win64 copy_files package
win32: update_version_in_config gen_version_file prepare_res compile_win32 copy_files package
linux: update_version_in_config gen_version_file prepare_res compile_linux copy_files package
mac: update_version_in_config gen_version_file prepare_res compile_mac copy_files package
arm: update_version_in_config gen_version_file prepare_res compile_arm copy_files package
upload: upload_to

prepare_res:
	@echo 'start prepare res'
	@go-bindata -o=res/res.go -pkg=res res/...
	@rm -rf ${BIN_DIR} && mkdir -p ${BIN_DIR}

compile_all: compile_win64 compile_win32 compile_linux compile_mac compile_arm

compile_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BIN_WIN64}${BINARY}.exe src/main.go

compile_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ${BIN_WIN32}${BINARY}.exe src/main.go

compile_linux:
	@echo 'start compile linux'
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BIN_LINUX}${BINARY} src/main.go

compile_mac:
	@echo 'start compile mac'
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BIN_MAC}${BINARY} src/main.go

compile_arm:
	@echo 'start compile arm'
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o ${BIN_ARM}${BINARY} src/main.go

copy_files:
	@echo 'start copy files'

	@cp -r {conf,runtime,demo} ${BIN_DIR}
	@for platform in `ls ${BIN_OUT}`; \
		do cp -r {bin/conf,bin/runtime,bin/demo} "${BIN_OUT}$${platform}/${BINARY}"; done

package:
	@echo 'start package'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for platform in `ls ${BIN_OUT}`; do mkdir -p ${QINIU_DIST_DIR}$${platform}; done

	@cd ${BIN_OUT} && \
		for platform in `ls ./`; \
		   do cd $${platform} && \
		   zip -r ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip ${BINARY} && \
		   md5sum ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip | awk '{print $$1}' | \
		          xargs echo > ${QINIU_DIST_DIR}$${platform}/${BINARY}.zip.md5 && \
           cd ..; \
		done

update_version_in_config:
	@gsed -i "s/Version.*/Version = ${VERSION}/" conf/ztf.conf

gen_version_file:
	@echo 'gen version'
	@mkdir -p ${QINIU_DIR}/${PROJECT}/
	@echo ${VERSION} > ${QINIU_DIR}/${PROJECT}/version.txt

upload_to:
	@echo 'upload...'
	@find ${QINIU_DIR} -name ".DS_Store" -type f -delete
	@qshell qupload2 --src-dir=${QINIU_DIR} --bucket=download --thread-count=10 --log-file=qshell.log \
					 --skip-path-prefixes=zd,zmanager,driver --rescan-local --overwrite --check-hash
