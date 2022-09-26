if [ ! -d "zentaoatf" ];then
    git clone https://github.com/easysoft/zentaoatf.git
    go env -w GOPROXY=https://goproxy.cn,direct
    cd zentaoatf
    git checkout ztf-42_zhaoke
    go mod tidy || true
    go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps
    cd ui
    npm install
    cd ../
else
    cd zentaoatf
    git pull
    go mod tidy || true
fi
cd test/cli
go run ./main.go -zentaoVersion 17.6.1