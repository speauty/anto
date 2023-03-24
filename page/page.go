package page

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"translator/tst/tt_ui/pack"
)

func StdRootWidget(rootWidget **walk.Composite, widgets ...Widget) Widget {
	return pack.TTComposite(
		pack.NewTTCompositeArgs(rootWidget).
			SetVisible(false).
			SetLayoutVBox(false).
			SetWidgets(widgets),
	)
}
