package page

import (
	pack2 "anto/lib/ui/pack"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func StdRootWidget(rootWidget **walk.Composite, widgets ...Widget) Widget {
	return pack2.TTComposite(
		pack2.NewTTCompositeArgs(rootWidget).
			SetVisible(false).
			SetLayoutVBox(false).
			SetWidgets(widgets),
	)
}

func StdBrowserSelectorWidget(title string, btOnClickFn walk.EventHandler, echoTarget **walk.Label) Widget {
	return pack2.TTComposite(
		pack2.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack2.NewWidgetGroup().Append(
				pack2.TTLabel(pack2.NewTTLabelArgs(nil).SetText(title)),
				pack2.TTPushBtn(pack2.NewTTPushBtnArgs(nil).SetText("选择").SetOnClicked(btOnClickFn)),
				pack2.TTLabel(pack2.NewTTLabelArgs(echoTarget).SetEnabled(false)),
			).AppendZeroHSpacer().GetWidgets()))
}
