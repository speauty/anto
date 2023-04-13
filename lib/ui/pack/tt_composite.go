package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTComposite(args *TTCompositeArgs) Widget {
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

func NewTTCompositeArgs(assignTo **walk.Composite) *TTCompositeArgs {
	return &TTCompositeArgs{
		assignTo:  assignTo,
		alignment: AlignHCenterVCenter,
		layout:    VBox{},
		visible:   true,
	}
}

type TTCompositeArgs struct {
	alignment  Alignment2D
	assignTo   **walk.Composite
	customSize Size
	layout     Layout
	widgets    []Widget
	visible    Property
}

func (customT *TTCompositeArgs) SetVisible(visible Property) *TTCompositeArgs {
	customT.visible = visible
	return customT
}
func (customT *TTCompositeArgs) SetAssignTo(assignTo **walk.Composite) *TTCompositeArgs {
	customT.assignTo = assignTo
	return customT
}
func (customT *TTCompositeArgs) SetAlignment(alignment Alignment2D) *TTCompositeArgs {
	customT.alignment = alignment
	return customT
}
func (customT *TTCompositeArgs) SetCustomSize(customSize Size) *TTCompositeArgs {
	customT.customSize = customSize
	return customT
}
func (customT *TTCompositeArgs) SetLayout(layout Layout) *TTCompositeArgs {
	customT.layout = layout
	return customT
}
func (customT *TTCompositeArgs) SetLayoutHBox(flagMarginsZero bool) *TTCompositeArgs {
	customT.layout = HBox{MarginsZero: flagMarginsZero}
	return customT
}
func (customT *TTCompositeArgs) SetLayoutVBox(flagMarginsZero bool) *TTCompositeArgs {
	customT.layout = VBox{MarginsZero: flagMarginsZero}
	return customT
}

func (customT *TTCompositeArgs) SetWidgets(widgets []Widget) *TTCompositeArgs {
	customT.widgets = widgets
	return customT
}
