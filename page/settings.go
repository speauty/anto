package page

import (
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
	"sync"
	"translator/cfg"
	"translator/tst/tt_log"
	"translator/tst/tt_ui/msg"
	"translator/tst/tt_ui/pack"
	"translator/util"
)

var (
	apiSettings  *Settings
	onceSettings sync.Once
)

func GetSettings() *Settings {
	onceSettings.Do(func() {
		apiSettings = new(Settings)
		apiSettings.id = util.Uid()
		apiSettings.name = "设置"
	})
	return apiSettings
}

type Settings struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite

	ptrEnv *walk.ComboBox

	ptrLingVADataId            *walk.LineEdit
	ptrBaiduAppId              *walk.LineEdit
	ptrBaiduAppKey             *walk.LineEdit
	ptrTencentCloudMTSecretId  *walk.LineEdit
	ptrTencentCloudMTSecretKey *walk.LineEdit
	ptrOpenAPIYouDaoAppKey     *walk.LineEdit
	ptrOpenAPIYouDaoAppSecret  *walk.LineEdit
	ptrHuaweiCloudAKId         *walk.LineEdit
	ptrHuaweiCloudSKKey        *walk.LineEdit
	ptrHuaweiCloudAKProjectId  *walk.LineEdit
	ptrHuaweiCloudAKRegion     *walk.LineEdit
}

func (customPage *Settings) GetId() string {
	return customPage.id
}

func (customPage *Settings) GetName() string {
	return customPage.name
}

func (customPage *Settings) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *Settings) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Settings) GetWidget() Widget {
	defer customPage.Reset()
	return StdRootWidget(&customPage.rootWidget,
		pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutVBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("翻译引擎").SetLayoutVBox(false).SetWidgets(
					pack.NewWidgetGroup().Append(
						pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("LingVA").SetLayoutVBox(false).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
									pack.NewWidgetGroup().Append(
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("数据ID")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrLingVADataId).
											SetText(cfg.GetInstance().LingVA.DataId)),
									).GetWidgets(),
								)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
						pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("百度翻译").SetLayoutVBox(false).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
									pack.NewWidgetGroup().Append(
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("应用ID")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrBaiduAppId).
											SetText(cfg.GetInstance().Baidu.AppId)),
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("应用密钥")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrBaiduAppKey).
											SetText(cfg.GetInstance().Baidu.AppKey)),
									).GetWidgets(),
								)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
						pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("腾讯云翻译").SetLayoutVBox(false).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
									pack.NewWidgetGroup().Append(
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("应用ID")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrTencentCloudMTSecretId).
											SetText(cfg.GetInstance().TencentCloudMT.SecretId)),
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("应用密钥")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrTencentCloudMTSecretKey).
											SetText(cfg.GetInstance().TencentCloudMT.SecretKey)),
									).GetWidgets(),
								)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
						pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("有道智云翻译").SetLayoutVBox(false).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
									pack.NewWidgetGroup().Append(
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("应用ID")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrOpenAPIYouDaoAppKey).
											SetText(cfg.GetInstance().OpenAPIYouDao.AppKey)),
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("应用密钥")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrOpenAPIYouDaoAppSecret).
											SetText(cfg.GetInstance().OpenAPIYouDao.AppSecret)),
									).GetWidgets(),
								)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
						pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("华为云翻译").SetLayoutVBox(false).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
									pack.NewWidgetGroup().Append(
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("访问ID")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrHuaweiCloudAKId).SetText(cfg.GetInstance().HuaweiCloudNlp.AKId)),
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("访问密钥")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrHuaweiCloudSKKey).SetText(cfg.GetInstance().HuaweiCloudNlp.SkKey)),
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("项目ID")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrHuaweiCloudAKProjectId).SetText(cfg.GetInstance().HuaweiCloudNlp.ProjectId)),
										pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("支持区域")),
										pack.TTLineEdit(pack.NewLineEditWrapperArgs(&customPage.ptrHuaweiCloudAKRegion).SetText(cfg.GetInstance().HuaweiCloudNlp.Region)),
									).GetWidgets(),
								)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
						pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTPushBtn(pack.NewTTPushBtnArgs(nil).SetText("同步").SetOnClicked(customPage.eventSync)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
					).AppendZeroHSpacer().AppendZeroVSpacer().GetWidgets(),
				)),
			).AppendZeroVSpacer().GetWidgets(),
		)),
	)
}

func (customPage *Settings) Reset() {

}

func (customPage *Settings) eventSync() {
	lingVADataId := customPage.ptrLingVADataId.Text()
	baiduAppId := customPage.ptrBaiduAppId.Text()
	baiduAppKey := customPage.ptrBaiduAppKey.Text()
	tencentCloudMTSecretId := customPage.ptrTencentCloudMTSecretId.Text()
	tencentCloudMTSecretKey := customPage.ptrTencentCloudMTSecretKey.Text()
	openAPIYouDaoAppKey := customPage.ptrOpenAPIYouDaoAppKey.Text()
	openAPIYouDaoAppSecret := customPage.ptrOpenAPIYouDaoAppSecret.Text()
	huaweiCloudAKId := customPage.ptrHuaweiCloudAKId.Text()
	huaweiCloudSKKey := customPage.ptrHuaweiCloudSKKey.Text()
	huaweiCloudAKProjectId := customPage.ptrHuaweiCloudAKProjectId.Text()
	huaweiCloudAKRegion := customPage.ptrHuaweiCloudAKRegion.Text()
	cntModified := 0
	if lingVADataId != cfg.GetInstance().LingVA.DataId {
		cfg.GetInstance().LingVA.DataId = lingVADataId
		cntModified++
	}
	if baiduAppId != cfg.GetInstance().Baidu.AppId {
		cfg.GetInstance().Baidu.AppId = baiduAppId
		cntModified++
	}
	if baiduAppKey != cfg.GetInstance().Baidu.AppKey {
		cfg.GetInstance().Baidu.AppKey = baiduAppKey
		cntModified++
	}
	if tencentCloudMTSecretId != cfg.GetInstance().TencentCloudMT.SecretId {
		cfg.GetInstance().TencentCloudMT.SecretId = tencentCloudMTSecretId
		cntModified++
	}
	if tencentCloudMTSecretKey != cfg.GetInstance().TencentCloudMT.SecretKey {
		cfg.GetInstance().TencentCloudMT.SecretKey = tencentCloudMTSecretKey
		cntModified++
	}

	if openAPIYouDaoAppKey != cfg.GetInstance().OpenAPIYouDao.AppKey {
		cfg.GetInstance().OpenAPIYouDao.AppKey = openAPIYouDaoAppKey
		cntModified++
	}
	if openAPIYouDaoAppSecret != cfg.GetInstance().OpenAPIYouDao.AppSecret {
		cfg.GetInstance().OpenAPIYouDao.AppSecret = openAPIYouDaoAppSecret
		cntModified++
	}

	if huaweiCloudAKId != cfg.GetInstance().HuaweiCloudNlp.AKId {
		cfg.GetInstance().HuaweiCloudNlp.AKId = huaweiCloudAKId
		cntModified++
	}
	if huaweiCloudSKKey != cfg.GetInstance().HuaweiCloudNlp.SkKey {
		cfg.GetInstance().HuaweiCloudNlp.SkKey = huaweiCloudSKKey
		cntModified++
	}
	if huaweiCloudAKProjectId != cfg.GetInstance().HuaweiCloudNlp.ProjectId {
		cfg.GetInstance().HuaweiCloudNlp.ProjectId = huaweiCloudAKProjectId
		cntModified++
	}
	if huaweiCloudAKRegion != cfg.GetInstance().HuaweiCloudNlp.Region {
		cfg.GetInstance().HuaweiCloudNlp.Region = huaweiCloudAKRegion
		cntModified++
	}
	if cntModified == 0 {
		msg.Info(customPage.mainWindow, "暂无配置需要同步")
		return
	}
	if err := cfg.GetInstance().Sync(); err != nil {
		tt_log.GetInstance().Error("同步配置到文件失败", zap.Error(err))
		msg.Err(customPage.mainWindow, errors.New("同步配置到文件失败"))
		return
	}
	msg.Info(customPage.mainWindow, "同步配置成功, 请重启当前应用")
}
