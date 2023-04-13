package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UIComposite(args *UICompositeArgs) Widget {
	return Composite{
		MinSize:   args.customSize,
		MaxSize:   args.customSize,
		Layout:    args.layout,
		Alignment: args.alignment,
		Children:  args.widgets,
		AssignTo:  args.assignTo,
		Visible:   args.visible,
	}
}

func NewUICompositeArgs(assignTo **walk.Composite) *UICompositeArgs {
	return &UICompositeArgs{
		assignTo:  assignTo,
		alignment: AlignHCenterVCenter,
		layout:    VBox{},
		visible:   true,
	}
}

type UICompositeArgs struct {
	alignment  Alignment2D
	assignTo   **walk.Composite
	customSize Size
	layout     Layout
	widgets    []Widget
	visible    Property
}

func (customT *UICompositeArgs) SetVisible(visible Property) *UICompositeArgs {
	customT.visible = visible
	return customT
}
func (customT *UICompositeArgs) SetAssignTo(assignTo **walk.Composite) *UICompositeArgs {
	customT.assignTo = assignTo
	return customT
}
func (customT *UICompositeArgs) SetAlignment(alignment Alignment2D) *UICompositeArgs {
	customT.alignment = alignment
	return customT
}
func (customT *UICompositeArgs) SetCustomSize(customSize Size) *UICompositeArgs {
	customT.customSize = customSize
	return customT
}
func (customT *UICompositeArgs) SetLayout(layout Layout) *UICompositeArgs {
	customT.layout = layout
	return customT
}
func (customT *UICompositeArgs) SetLayoutHBox(flagMarginsZero bool) *UICompositeArgs {
	customT.layout = HBox{MarginsZero: flagMarginsZero}
	return customT
}
func (customT *UICompositeArgs) SetLayoutVBox(flagMarginsZero bool) *UICompositeArgs {
	customT.layout = VBox{MarginsZero: flagMarginsZero}
	return customT
}

func (customT *UICompositeArgs) SetWidgets(widgets []Widget) *UICompositeArgs {
	customT.widgets = widgets
	return customT
}
