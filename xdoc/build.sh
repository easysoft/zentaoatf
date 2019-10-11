rm -rf ztf
mkdir ztf
mkdir ztf/log
cp -r conf ztf/
cp -r demo ztf/

go-bindata -o=res/res.go -pkg=res res/ res/doc res/json res/template

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ztf/ztf.exe src/ztf.go
GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ztf/ztf-linux src/ztf.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ztf/ztf-mac src/ztf.go

cd ztf
zip -r ztf-win-2.0.RC.zip ztf.exe demo conf log
cp -r ztf-linux ztf
zip -r ztf-linux-2.0.RC.zip ztf demo conf log
cp -r ztf-mac ztf
zip -r ztf-mac-2.0.RC.zip ztf demo conf log
rm ztf
cd ..