# zentaoatf
ZentaoATF is an automation testing framework written in Golang.

## Features
1. Support most popular programming languages like Python, Ruby, Lua, Tcl, PHP, Shell and GO;
2. Reduce the invasive of existing testing scripts
3. Integration with ZenTao - an open source project management system;
3. Easy to use with the help of UI window

## QuickStart
### Run
1. Download corresponding release file from [here](https://github.com/easysoft/zentaoatf/tree/master/release);
2. Type 'atf-2.0.0.alpha.exe' to get the help doc. 

### Dev
1. Use 'git clone https://github.com/easysoft/zentaoatf.git' to get the source codes;
2. Type `go get -u all' to get all dependencies;
3. Type 'go run src/atf.go' to get the help doc;
4. To open the CUI window, type 'go run src/atf.go cui'

### Test Suite

### Test Script


## Example:
#### Import test cases from remote Zentao system
go run src/atf.go gen -u http://ruiyinxin.test.zentao.net -t product -v 1 -l python -a autotest01 -p P2ssw0rd

#### Run test scripts in specified folder
go run src/atf.go run -d scripts -l python

#### Batch run with test suite
go run src/atf.go run -f scripts/all.suite -l python

#### Rerun failed test cases in specified result file
go run src/atf.go rerun -p logs/suite-all-2019-08-21T133157/result.txt

#### List test scripts
go run src/atf.go list -d scripts -l python

#### Brief test scripts in dir
go run src/atf.go view -d scripts -l python

#### View test scripts by path
go run src/atf.go view -f scripts/tc-1.py -f scripts/tc-2.py

#### Switch work dir to another path
go run src/atf.go switch -p /Users/aaron/dev/go/autotest/

#### Change tool language（en: English, zh: Simplified Chinese）
go run src/atf.go set -l zh

#### Open CUI Window
go run src/atf.go cui

#### Submit test result to remote Zentao system
![submit_result](xdoc/snapshot/submit_result.jpg)

#### Report bug for failed test case to remote Zentao system
![report_bug](xdoc/snapshot/report_bug.jpg)

## Licenses
All source code is licensed under the [Z PUBLIC LICENSE](LICENSE.md).