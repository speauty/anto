package ui

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"gui.subtitle/src/logic/translate"
	mt2 "gui.subtitle/src/srv/mt"
	"gui.subtitle/src/srv/mt/aliyun"
	"gui.subtitle/src/srv/mt/bd"
	"gui.subtitle/src/util/lang"
	"os"
	"path/filepath"
	"strings"
)

var isSpecialMode = walk.NewMutableCondition()

var mtEngineValues = []string{"阿里云", "百度"}

func init() {
	MustRegisterCondition("isSpecialMode", isSpecialMode)
}

type AppWindow struct {
	*walk.MainWindow
	cfg *AppWindowCfg
}

type AppWindowCfg struct {
	Title string
}

func (aw *AppWindow) Start(cfg *AppWindowCfg) error {
	aw.cfg = cfg
	if err := aw.newMainWindow(); err != nil {
		return err
	}
	flag := win.GetWindowLong(aw.Handle(), win.GWL_STYLE)
	flag |= win.WS_OVERLAPPED  // always on top
	flag &= ^win.WS_THICKFRAME // fixed size
	win.SetWindowLong(aw.Handle(), win.GWL_STYLE, flag)
	aw.Run()
	return nil
}

func (aw *AppWindow) newMainWindow() error {
	var mtEngineComboBox *walk.ComboBox

	var aLiIdEdit *walk.TextEdit
	var aLiKeyEdit *walk.TextEdit
	var aLiLocationEdit *walk.TextEdit

	var bdIdEdit *walk.TextEdit
	var bdKeyEdit *walk.TextEdit

	var fromLanguageComboBox *walk.ComboBox
	var toLanguageComboBox *walk.ComboBox
	langVars := lang.ZH.GetMaps()

	var subtitleFileEdit *walk.TextEdit
	var subtitleFileSave *walk.TextEdit

	var logLabel *walk.TextEdit
	var stateLabel *walk.TextEdit

	ctx := context.Background()

	return MainWindow{
		AssignTo: &aw.MainWindow,
		Title:    aw.cfg.Title,
		Size:     Size{Width: 400, Height: 600},
		MinSize:  Size{Width: 400, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "***** 翻译小助手 *****", Alignment: AlignHCenterVCenter, Font: Font{Bold: true, PointSize: 10}},
			Label{Text: "版本号: 1.0.5  作者: speauty  邮箱: speauty@163.com", Alignment: AlignHCenterVCenter},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "翻译引擎", ToolTipText: "机器翻译引擎, 当前支持阿里云、百度"},
					ComboBox{
						Name:    "mtEngineComboBox",
						MinSize: Size{Width: 80}, MaxSize: Size{Width: 80, Height: 20},
						AssignTo:     &mtEngineComboBox,
						Model:        mtEngineValues,
						CurrentIndex: 0,
					},
					HSpacer{},
				},
			},
			GroupBox{
				MinSize: Size{Height: 110},
				MaxSize: Size{Height: 110},
				Layout:  VBox{},
				Visible: Bind("mtEngineComboBox.CurrentIndex == 0"),
				Children: []Widget{
					HSplitter{
						Children: []Widget{
							Label{Text: "访问身份", ToolTipText: "AccessKeyId"},
							TextEdit{AssignTo: &aLiIdEdit},
							HSpacer{},
						},
					},
					HSplitter{
						Children: []Widget{
							Label{Text: "访问密钥", ToolTipText: "AccessKeySecret"},
							TextEdit{AssignTo: &aLiKeyEdit},
							HSpacer{},
						},
					},
					HSplitter{
						Children: []Widget{
							Label{Text: "访问区域", ToolTipText: "AccessKeySecret"},
							TextEdit{AssignTo: &aLiLocationEdit, Text: "cn-hangzhou"},
							HSpacer{},
						},
					},
					VSpacer{},
				},
			},
			GroupBox{
				MinSize: Size{Height: 80},
				MaxSize: Size{Height: 80},
				Layout:  VBox{},
				Visible: Bind("mtEngineComboBox.CurrentIndex == 1"),
				Children: []Widget{
					HSplitter{
						Children: []Widget{
							Label{Text: "访问身份", ToolTipText: "AccessKeyId"},
							TextEdit{AssignTo: &bdIdEdit},
							HSpacer{},
						},
					},
					HSplitter{
						Children: []Widget{
							Label{Text: "访问密钥", ToolTipText: "AccessKeySecret"},
							TextEdit{AssignTo: &bdKeyEdit},
							HSpacer{},
						},
					},
					VSpacer{},
				},
			},
			GroupBox{
				MinSize: Size{Height: 100},
				MaxSize: Size{Height: 100},
				Layout:  VBox{},
				Children: []Widget{
					HSplitter{
						MinSize: Size{Height: 20}, MaxSize: Size{Height: 20},
						Children: []Widget{
							Label{Text: "来源语言", MinSize: Size{Width: 40}, MaxSize: Size{Width: 40}},
							ComboBox{
								Name:    "fromLanguageComboBox",
								MinSize: Size{Width: 112}, MaxSize: Size{Width: 112},
								AssignTo:     &fromLanguageComboBox,
								Model:        langVars,
								CurrentIndex: 0,
							},
							Label{Text: "目标语言", MinSize: Size{Width: 40}, MaxSize: Size{Width: 40}},
							ComboBox{
								Name:    "toLanguageComboBox",
								MinSize: Size{Width: 112}, MaxSize: Size{Width: 112},
								AssignTo:     &toLanguageComboBox,
								Model:        langVars,
								CurrentIndex: 1,
								OnCurrentIndexChanged: func() {
									/*if toLanguageComboBox.CurrentIndex() == 0 { // 目标语言禁止自动
										_ = toLanguageComboBox.SetCurrentIndex(1)
									}*/
								},
							},
							HSpacer{},
						},
					},
					HSplitter{
						MinSize: Size{Height: 20}, MaxSize: Size{Height: 20},
						Children: []Widget{
							Label{Text: "字幕文件"},
							PushButton{
								MinSize: Size{Width: 10}, MaxSize: Size{Width: 10},
								Text: "选择",
								OnClicked: func() {
									dlg := new(walk.FileDialog)
									dlg.Title = "选择字幕文件"
									dlg.Filter = "字幕文件 (*.srt)|*.srt"

									if ok, err := dlg.ShowOpen(aw); err != nil {
										walk.MsgBox(aw, "错误", fmt.Sprintf("选择字幕文件异常, 错误: %s", err), walk.MsgBoxIconError)
										return
									} else if !ok {
										if subtitleFileEdit.Text() == "" {
											walk.MsgBox(aw, "错误", "选择字幕文件失败, 请重新选择", walk.MsgBoxIconError)
										}
										return
									}
									_ = subtitleFileEdit.SetText(dlg.FilePath)
								},
							},
							TextEdit{AssignTo: &subtitleFileEdit, ReadOnly: true, HScroll: true},
						},
					},
					HSplitter{
						MinSize: Size{Height: 20}, MaxSize: Size{Height: 20},
						Children: []Widget{
							Label{Text: "保存路径"},
							PushButton{
								MinSize: Size{Width: 10}, MaxSize: Size{Width: 10},
								Text: "选择",
								OnClicked: func() {
									dlg := new(walk.FileDialog)
									dlg.Title = "选择保存路径"

									if ok, err := dlg.ShowBrowseFolder(aw); err != nil {
										walk.MsgBox(aw, "错误", fmt.Sprintf("选择保存路径异常, 错误: %s", err), walk.MsgBoxIconError)
										return
									} else if !ok {
										if subtitleFileSave.Text() == "" {
											walk.MsgBox(aw, "错误", "选择保存路径失败, 请重新选择", walk.MsgBoxIconError)
										}
										return
									}
									_ = subtitleFileSave.SetText(dlg.FilePath)
								},
							},
							TextEdit{AssignTo: &subtitleFileSave, ReadOnly: true, HScroll: true},
						},
					},
					VSpacer{},
				},
			},
			PushButton{
				Name: "submit-btn",
				Text: "提交",
				OnClicked: func() {
					if logLabel.Visible() {
						logLabel.SetVisible(false)
					}
					if stateLabel.Visible() {
						stateLabel.SetVisible(false)
					}
					_ = stateLabel.SetText("")
					_ = logLabel.SetText("")
					timeStart := carbon.Now()
					currentMTEngineName := getCurrentMTEngineName(mtEngineComboBox.CurrentIndex())

					if currentMTEngineName == "阿里云" {
						if aLiIdEdit.Text() == "" {
							walk.MsgBox(aw, "警告", "请设置访问身份", walk.MsgBoxIconWarning)
							defer func() { _ = aLiIdEdit.SetFocus() }()
							return
						}
						if aLiKeyEdit.Text() == "" {
							walk.MsgBox(aw, "警告", "请设置访问密钥", walk.MsgBoxIconWarning)
							defer func() { _ = aLiKeyEdit.SetFocus() }()
							return
						}
					} else if currentMTEngineName == "百度" {
						if bdIdEdit.Text() == "" {
							walk.MsgBox(aw, "警告", "请设置访问身份", walk.MsgBoxIconWarning)
							defer func() { _ = bdIdEdit.SetFocus() }()
							return
						}
						if bdKeyEdit.Text() == "" {
							walk.MsgBox(aw, "警告", "请设置访问密钥", walk.MsgBoxIconWarning)
							defer func() { _ = bdKeyEdit.SetFocus() }()
							return
						}
					} else {
						walk.MsgBox(aw, "提示", fmt.Sprintf("当前翻译引擎[%s]暂未接入, 尽情期待", currentMTEngineName), walk.MsgBoxIconWarning)
						return
					}

					if fromLanguageComboBox.CurrentIndex() == toLanguageComboBox.CurrentIndex() {
						walk.MsgBox(aw, "警告", fmt.Sprintf(
							"来源语言[%s]和目标语言[%s]不能一样, 请重新设置",
							lang.ZH.GetLangByIdx(fromLanguageComboBox.CurrentIndex()).GetCH(),
							lang.ZH.GetLangByIdx(toLanguageComboBox.CurrentIndex()).GetCH(),
						), walk.MsgBoxIconWarning)
						return
					}

					if subtitleFileEdit.Text() == "" {
						walk.MsgBox(aw, "警告", "请选择字幕文件", walk.MsgBoxIconWarning)
						return
					}
					if subtitleFileSave.Text() == "" {
						walk.MsgBox(aw, "警告", "请选择保存路径", walk.MsgBoxIconWarning)
						return
					}

					file, err := os.Open(subtitleFileEdit.Text())
					defer func() {
						_ = file.Close()
					}()
					if err != nil {
						walk.MsgBox(aw, "错误", fmt.Sprintf("打开目标字幕文件异常, 错误: %s", err.Error()), walk.MsgBoxIconError)
						return
					}
					sts, err := translate.Reader(file)
					if err != nil {
						walk.MsgBox(aw, "错误", fmt.Sprintf("载入目标字幕文件内容异常, 错误: %s", err.Error()), walk.MsgBoxIconError)
						return
					}
					basename := filepath.Base(subtitleFileEdit.Text())
					ext := filepath.Ext(subtitleFileEdit.Text())
					translatedFilePath := fmt.Sprintf("%s/%s.translated.srt", subtitleFileSave.Text(), strings.ReplaceAll(basename, ext, ""))
					create, err := os.Create(translatedFilePath)
					defer func() {
						_ = create.Close()
					}()
					if err != nil {
						walk.MsgBox(aw, "错误", fmt.Sprintf("创建翻译字幕文件异常, 错误: %s", err.Error()), walk.MsgBoxIconError)
						return
					}
					fromLanguage := lang.ZH.GetLangByIdx(fromLanguageComboBox.CurrentIndex())
					toLanguage := lang.ZH.GetLangByIdx(toLanguageComboBox.CurrentIndex())
					var mt mt2.MT
					var cfg interface{}

					if currentMTEngineName == "阿里云" {
						cfg = &aliyun.Cfg{AccessKeyId: aLiIdEdit.Text(), AccessKeySecret: aLiKeyEdit.Text(), Location: aLiLocationEdit.Text()}
						mt = new(aliyun.ALiMT)
					} else if currentMTEngineName == "百度" {
						cfg = &bd.Cfg{AppId: bdIdEdit.Text(), AppSecret: bdKeyEdit.Text()}
						mt = new(bd.MT)
					} else {
						walk.MsgBox(aw, "提示", fmt.Sprintf("当前翻译引擎[%s]暂未接入, 尽情期待", currentMTEngineName), walk.MsgBoxIconWarning)
						return
					}
					if err := mt.Init(ctx, cfg); err != nil {
						walk.MsgBox(aw, "错误", fmt.Sprintf("初始化%s服务异常, 错误: %s", mt.GetName(), err.Error()), walk.MsgBoxIconError)
						return
					}
					results, cntError, _ := translate.Translate(ctx, mt, sts, fromLanguage, toLanguage)
					_ = logLabel.SetText(strings.Join(results, "\r\n"))
					_cntWrite, err := translate.Writer(create, sts)
					if err != nil {
						walk.MsgBox(aw, "错误", fmt.Sprintf("写入目标字幕文件内容异常, 错误: %s", err.Error()), walk.MsgBoxIconError)
						return
					}
					walk.MsgBox(aw, "提示", fmt.Sprintf("处理完成, 翻译文件: %s, 总共写入(byte): %d, 错误数量: %d", translatedFilePath, _cntWrite, cntError), walk.MsgBoxOK)

					if !logLabel.Visible() {
						logLabel.SetVisible(true)
					}
					if !stateLabel.Visible() {
						stateLabel.SetVisible(true)
					}
					_ = stateLabel.SetText(fmt.Sprintf(
						"翻译文件: %s\r\n字幕行数: %d\r\n耗时(s): %d", translatedFilePath, len(sts), carbon.Now().DiffAbsInSeconds(timeStart),
					))
					return
				},
			},
			TextEdit{
				AssignTo: &stateLabel, Visible: false, ReadOnly: true, VScroll: true,
				MinSize: Size{Height: 60}, MaxSize: Size{Height: 60},
				Font: Font{PointSize: 8},
			},
			TextEdit{
				Visible: false, ReadOnly: true, VScroll: true,
				AssignTo: &logLabel, MinSize: Size{Height: 120, Width: 100}, MaxSize: Size{Height: 120, Width: 100},
			},

			VSpacer{},
		},
	}.Create()
}

func getCurrentMTEngineName(idx int) string {
	for i, val := range mtEngineValues {
		if idx == i {
			return val
		}
	}
	return ""
}
