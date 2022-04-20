VERSION=3.0.0_beta
PROJECT=ztf
QINIU_DIR=/Users/aaron/work/zentao/qiniu/
QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/
MAIN_FILE=cmd/command/main.go

BIN_DIR=client/bin/

BUILD_TIME=`git show -s --format=%cd`
GO_VERSION=`go version`
GIT_HASH=`git show -s --format=%H`
BUILD_CMD=go build -ldflags "-X 'commConsts.appVersion=${VERSION}' -X 'commConsts.buildTime=${BUILD_TIME}' -X 'commConsts.goVersion=${GO_VERSION}' -X 'commConsts.gitHash=${GIT_HASH}'"

default: update_version prepare_res compile_win64 compile_win32 compile_linux compile_mac copy_files zip

win64: update_version prepare_res compile_win64 copy_files zip
win32: update_version prepare_res compile_win32 copy_files zip
linux: update_version prepare_res compile_linux copy_files zip
mac:   update_version prepare_res compile_mac copy_files zip

update_version: update_version_in_config gen_version_file
update_version_in_config:
	@gsed -i "s/Version.*/Version = ${VERSION}/" conf/ztf.conf
gen_version_file:
	@echo 'gen version'
	@mkdir -p ${QINIU_DIR}/${PROJECT}/
	@echo ${VERSION} > ${QINIU_DIR}/${PROJECT}/version.txt

prepare_res:
	@echo 'start prepare res'
	@rm -rf res/res.go
	@go-bindata -o=res/res.go -pkg=res res/...

compile_win64:
	@echo 'start compile win64'
	@rm -rf ./${BIN_DIR}/*
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${BIN_DIR}win32/${PROJECT}.exe ${MAIN_FILE}

compile_win32:
	@echo 'start compile win32'
	@rm -rf ./${BIN_DIR}/*
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD} -x -v -ldflags "-s -w" \
		-o ${BIN_DIR}win32/${PROJECT}.exe ${MAIN_FILE}

compile_linux:
	@echo 'start compile linux'
	@rm -rf ./${BIN_DIR}/*
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD} \
		-o ${BIN_DIR}linux/${PROJECT} ${MAIN_FILE}

compile_mac:
	@echo 'start compile mac'
	@rm -rf ./${BIN_DIR}/*
	@echo
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD} \
		-o ${BIN_DIR}darwin/${PROJECT} ${MAIN_FILE}

copy_files:
	@echo 'start copy files'
	@for platform in `ls ${BIN_DIR}`; \
		do cp -r {demo,runtime} "${BIN_DIR}$${platform}"; done

zip:
	@echo 'start zip'
	@find . -name .DS_Store -print0 | xargs -0 rm -f
	@for platform in `ls ${BIN_DIR}`; do mkdir -p ${QINIU_DIST_DIR}$${platform}; done

	@cd ${BIN_DIR} && \
		for platform in `ls ./`; \
		   do cd $${platform} && \
		   zip -r ${QINIU_DIST_DIR}$${platform}/${PROJECT}-cmd.zip ./* && \
		   md5sum ${QINIU_DIST_DIR}$${platform}/${PROJECT}-cmd.zip | awk '{print $$1}' | \
		          xargs echo > ${QINIU_DIST_DIR}$${platform}/${PROJECT}-cmd.zip.md5 && \
           cd ..; \
		done

upload_to:
	@echo 'upload...'
	@find ${QINIU_DIR} -name ".DS_Store" -type f -delete
	@qshell qupload2 --src-dir=${QINIU_DIR} --bucket=download --thread-count=10 --log-file=qshell.log \
					 --skip-path-prefixes=zz,zd,zmanager,driver --rescan-local --overwrite --check-hash