.PHONY: tidy deploy_gui_win rs build compress

tidy:
	go mod tidy

deploy: rs build compress

rs:
	rsrc -manifest ./cmd/anto/anto.manifest -ico favicon.ico -o ./cmd/anto/rsrc.syso

build:
	go build -gcflags='-l -N' -ldflags='-w -s -H=windowsgui' -o ./bin/anto.win.exe anto/cmd/anto

compress:
	upx -9 ./bin/anto.win.exe
