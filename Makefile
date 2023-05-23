.PHONY: tidy install_fyne install_fyne_old

tidy:
	go mod tidy

install_fyne: # 安装fyne后, 建议执行下tidy
	go get fyne.io/fyne/v2@latest& go install fyne.io/fyne/v2/cmd/fyne@latest

# 老版本go安装fyne指令(version < 1.16)
install_fyne_old:
	go get fyne.io/fyne/v2& go get fyne.io/fyne/v2/cmd/fyne

gui_compress:
	upx -9 ./bin/*.exe
