VERSION=1.1

all: zip win

clean:
	rm -fr zentaoatf
	rm -fr *.zip
zip:
	mkdir zentaoatf
	cp -fr {zt,langs,hello} zentaoatf
	zip -r -9 ZenTaoATF.$(VERSION).zip zentaoatf
	rm -fr zentaoatf
win:        
	mkdir zentaoatf
	cp -fr {zt*,langs,hello} zentaoatf
	cp -fr php zentaoatf/
	zip -r -9 ZenTaoATF.$(VERSION).win.zip zentaoatf
	rm -fr zentaoatf
