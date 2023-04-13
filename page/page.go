package page

import (
	"anto/tst/tt_ui/pack"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func StdRootWidget(rootWidget **walk.Composite, widgets ...Widget) Widget {
	return pack.TTComposite(
		pack.NewTTCompositeArgs(rootWidget).
			SetVisible(false).
			SetLayoutVBox(false).
			SetWidgets(widgets),
	)
}

func StdBrowserSelectorWidget(title string, btOnClickFn walk.EventHandler, echoTarget **walk.Label) Widget {
	return pack.TTComposite(
		pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText(title)),
				pack.TTPushBtn(pack.NewTTPushBtnArgs(nil).SetText("选择").SetOnClicked(btOnClickFn)),
				pack.TTLabel(pack.NewTTLabelArgs(echoTarget).SetEnabled(false)),
			).AppendZeroHSpacer().GetWidgets()))
}
