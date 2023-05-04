.PHONY: deploy build compress rs tidy run

# 发布win应用
deploy: rs build compress

build:
	go build -gcflags='-l -N' -ldflags='-w -s -H=windowsgui' -o .\bin\anto.exe anto

compress:
	upx -9 .\bin\anto.exe

rs:
	rsrc -manifest anto.manifest -ico favicon.ico -o rsrc.syso

tidy:
	go mod tidy

run: build
	cd .\bin&& anto.exe
