.PHONY: tidy

tidy:
	go mod tidy

.PHONY: deploy_gui_win gui_win_rs gui_win_build gui_win_compress
deploy_gui_win: gui_win_rs gui_win_build gui_win_compress

gui_win_rs:
	rsrc -manifest ./cmd/gui_win/gui_win.manifest -ico favicon.ico -o ./cmd/gui_win/rsrc.syso

gui_win_build:
	go build -gcflags='-l -N' -ldflags='-w -s -H=windowsgui' -o ./bin/anto.win.exe anto/cmd/gui_win

gui_win_compress:
	upx -9 ./bin/anto.win.exe
