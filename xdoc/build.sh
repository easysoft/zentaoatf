rm -rf ztf
mkdir ztf
mkdir ztf/log
cp -r conf ztf/
cp -r demo ztf/

go-bindata -o=res/res.go -pkg=res res/ res/doc res/json res/template

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ztf/ztf-x86.exe src/ztf.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ztf/ztf-amd64.exe src/ztf.go

GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ztf/ztf-linux src/ztf.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ztf/ztf-mac src/ztf.go

cd ztf

cp ztf-x86.exe ztf.exe
zip -r ztf-win-x86-2.1.zip ztf.exe demo conf log
rm ztf.exe

cp ztf-amd64.exe ztf.exe
zip -r ztf-win-amd64-2.1.zip ztf.exe demo conf log
rm ztf.exe

cp ztf-linux ztf
tar -zcvf ztf-linux-2.1.tar.gz ztf demo conf log
rm ztf

cp ztf-mac ztf
zip -r ztf-mac-2.1.zip ztf demo conf log
rm ztf

cd ..