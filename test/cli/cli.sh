#/bin/bash
if [ ! -d "/data/" ];then
    git clone https://github.com/easysoft/zentaoatf.git
    go env -w GOPROXY=https://goproxy.cn,direct
    go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps
else
    cd zentaoatf
    git pull
fi
cd test/cli
go run ./main.go -zentaoVersion 17.6.1