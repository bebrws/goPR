GO_FILES := $(shell find config cmd internal -name '*.go')

all: goPR.app/Contents/MacOS/goPR

goPR.app/Contents/MacOS/goPR: $(GO_FILES)
	rm -rf goPR.app
	mkdir -p goPR.app/Contents/MacOS
	mkdir -p goPR.app/Contents/Resources
	cp icon.icns goPR.app/Contents/Resources
	cp Info.plist goPR.app/Contents
	GO111MODULE=on go build -o goPR.app/Contents/MacOS/goPR cmd/main.go

# sign:
# 	@echo "No longer needed"
# 	sudo codesign -s - --deep goPR.app || true

clean:
	rm -rf goPR.app

run: goPR.app/Contents/MacOS/goPR
	goPR.app/Contents/MacOS/goPR

install:
	@echo "Will install a LaunchAgent to run this every 1200 seconds by default"
	@echo "Run: goPR.app/Contents/MacOS/goPR install numSeconds"
	@echo "to change the interval"
	@echo "To remove the LaunchAgent run: goPR.app/Contents/MacOS/goPR clean"
	@echo "Or: make clean"
	@echo "from this directory"
	goPR.app/Contents/MacOS/goPR install 1200

uninstall:
	@echo "Deleting launch LaunchAgent"
	goPR.app/Contents/MacOS/goPR clean

.PHONY: clean run install clean uninstall # sign