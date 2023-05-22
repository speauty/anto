.PHONY: tidy install_fyne install_fyne_old

tidy:
	go mod tidy

install_fyne: # 安装fyne后, 建议执行下tidy
	go get fyne.io/fyne/v2@latest& go install fyne.io/fyne/v2/cmd/fyne@latest

# 老版本go安装fyne指令(version < 1.16)
install_fyne_old:
	go get fyne.io/fyne/v2& go get fyne.io/fyne/v2/cmd/fyne


.PHONY: deploy_gui_win gui_win_rs gui_win_build gui_win_compress

deploy_gui_win: gui_win_rs gui_win_build gui_win_compress

gui_win_rs:
	rsrc -manifest ./cmd/gui_win/gui_win.manifest -ico favicon.ico -o ./cmd/gui_win/rsrc.syso

gui_win_build:
	go build -gcflags='-l -N' -ldflags='-w -s -H=windowsgui' -o ./bin/anto.win.exe anto/cmd/gui_win

gui_win_compress:
	upx -9 ./bin/anto.win.exe

