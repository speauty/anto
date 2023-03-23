init:
	go mod tidy

buildWin:
	go build "-ldflags=-w -s -H=windowsgui" -o .\bin\translator.exe translator && upx -9 .\bin\translator.exe

rsrs:
	rsrc -manifest translator.manifest -ico favicon.ico -o rsrc.syso