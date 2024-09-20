all: goPR.app/Contents/MacOS/goPR

goPR.app/Contents/MacOS/goPR: cmd/main.go 
	GO111MODULE=on go build -o goPR.app/Contents/MacOS/goPR cmd/main.go

sign:	
	sudo codesign -s - --deep goPR.app

clean:
	rm -rf goPR.app/Contents/MacOS/goPR

run: goPR.app/Contents/MacOS/goPR sign
	goPR.app/Contents/MacOS/goPR

.PHONY: clean sign run