VERSION=$(shell head -n 1 VERSION)

all: tgz

clean:
	rm -fr zentaoatf
	rm -fr *.zip
tgz:
	mkdir zentaoatf
	cp -fr {zt,langs,hello} zentaoatf
	find zentaoatf/.git |xargs rm -fr
	zip -r -9 ZenTaoATF.$(VERSION).zip zentaoatf
	rm -fr zentaoatf
