package page

import (
	"anto/platform/win/ui/pack"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func StdRootWidget(rootWidget **walk.Composite, widgets ...Widget) Widget {
	return pack.UIComposite(
		pack.NewUICompositeArgs(rootWidget).
			SetVisible(false).
			SetLayoutVBox(false).
			SetWidgets(widgets),
	)
}

func StdBrowserSelectorWidget(title string, btOnClickFn walk.EventHandler, echoTarget **walk.Label) Widget {
	return pack.UIComposite(
		pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UILabel(pack.NewUILabelArgs(nil).SetText(title)),
				pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("选择").SetOnClicked(btOnClickFn)),
				pack.UILabel(pack.NewUILabelArgs(echoTarget).SetEnabled(false)),
			).AppendZeroHSpacer().GetWidgets()))
}
