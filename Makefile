init:
	go mod tidy

buildWin:
	go build "-ldflags=-w -s -H=windowsgui" -o .\bin\anto.exe anto && upx -9 .\bin\anto.exe

rsrs:
	rsrc -manifest anto.manifest -ico favicon.ico -o rsrc.syso