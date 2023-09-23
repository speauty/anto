.PHONY: tidy deploy_gui_win rs build compress

tidy:
	go mod tidy

BinName=anto-v3.4.5-windows.exe

deploy: rs build

rs:
	rsrc -manifest ./cmd/anto/anto.manifest -ico ./resource/favicon.ico -o ./cmd/anto/rsrc.syso

build:
	go build -gcflags='-l -N' -ldflags='-w -s -H=windowsgui' -o "./bin/${BinName}" anto/cmd/anto

compress:
	upx -9 "./bin/${BinName}"
