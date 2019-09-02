# zentaoatf
ZentaoATF is an automation testing framework written in Golang.

## Features
1. Support most popular programming languages like Python, Ruby, Lua, Tcl, PHP, Shell, GO and Windows Bat;
2. Reduce the invasive of existing testing scripts
3. Integration with ZenTao - an open source project management system;
3. Easy to use with the help of UI window

## QuickStart
### Run from release file
1. Download last release file from [here](https://github.com/easysoft/zentaoatf/releases);
2. Type 'atf-2.0.0.alpha.exe help' to get the doc.

### Run from Golang codes
1. Use 'git clone https://github.com/easysoft/zentaoatf.git' to get the source codes;
2. Overwrite edit.go and view.go files from https://github.com/rocket049/gocui to fix the Chinese related bug;
3. Type `go get -u all' to get all dependencies;
4. Type 'go run src/atf.go help' to get the doc;

## Licenses
All source code is licensed under the [Z PUBLIC LICENSE](LICENSE.md).