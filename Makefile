VERSION=1.1

all: zip win32 win64

clean:
	rm -fr zentaoatf
	rm -fr *.zip
zip:
	mkdir zentaoatf
	cp -fr {zt,langs,hello} zentaoatf
	zip -r -9 ZenTaoATF.$(VERSION).zip zentaoatf
	rm -fr zentaoatf
win32:        
	mkdir zentaoatf
	cp -fr {zt,langs,hello} zentaoatf
	mkdir zentaoatf/runtime
	cp -fr runtime/ansicon/x32/* zentaoatf/runtime
	cp -fr runtime/php/* zentaoatf/runtime
	echo "@echo off" > zentaoatf/zt.bat
	echo "set console=ansicon" >> zentaoatf/zt.bat
	echo ".\runtime\ansicon.exe .\runtime\php.exe zt %*" >> zentaoatf/zt.bat
	zip -r -9 ZenTaoATF.$(VERSION).win32.zip zentaoatf
	rm -fr zentaoatf
win64:        
	mkdir zentaoatf
	cp -fr {zt,langs,hello} zentaoatf
	mkdir zentaoatf/runtime
	cp -fr runtime/ansicon/x64/* zentaoatf/runtime
	cp -fr runtime/php/* zentaoatf/runtime
	echo "@echo off" > zentaoatf/zt.bat
	echo "set console=ansicon" >> zentaoatf/zt.bat
	echo ".\runtime\ansicon.exe .\runtime\php.exe zt %*" >> zentaoatf/zt.bat
	zip -r -9 ZenTaoATF.$(VERSION).win64.zip zentaoatf
	rm -fr zentaoatf
