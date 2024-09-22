GO_FILES := $(shell find config cmd internal -name '*.go')

all: goPR.app/Contents/MacOS/goPR

goPR.app/Contents/MacOS/goPR: $(GO_FILES)
	GO111MODULE=on go build -o goPR.app/Contents/MacOS/goPR cmd/main.go

sign:
	@echo "No longer needed"
	sudo codesign -s - --deep goPR.app || true

clean:
	rm -rf goPR.app/Contents/MacOS/goPR

run: goPR.app/Contents/MacOS/goPR
	goPR.app/Contents/MacOS/goPR

.PHONY: clean sign run