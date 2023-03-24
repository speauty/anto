package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTGroupBox(args *TTGroupBoxArgs) Widget {
	return GroupBox{
		Title:     args.title,
		MinSize:   args.customSize,
		MaxSize:   args.customSize,
		Layout:    args.layout,
		Alignment: args.alignment,
		Children:  args.widgets,
		AssignTo:  args.assignTo,
		Visible:   args.visible,
	}
}

func NewTTGroupBoxArgs(assignTo **walk.GroupBox) *TTGroupBoxArgs {
	return &TTGroupBoxArgs{
		assignTo:  assignTo,
		alignment: AlignHCenterVCenter,
		layout:    VBox{},
	}
}

type TTGroupBoxArgs struct {
	title      string
	alignment  Alignment2D
	assignTo   **walk.GroupBox
	customSize Size
	layout     Layout
	widgets    []Widget
	visible    bool
}

func (customT *TTGroupBoxArgs) SetVisible(visible bool) *TTGroupBoxArgs {
	customT.visible = visible
	return customT
}

func (customT *TTGroupBoxArgs) SetTitle(title string) *TTGroupBoxArgs {
	customT.title = title
	return customT
}

func (customT *TTGroupBoxArgs) SetAssignTo(assignTo **walk.GroupBox) *TTGroupBoxArgs {
	customT.assignTo = assignTo
	return customT
}
func (customT *TTGroupBoxArgs) SetAlignment(alignment Alignment2D) *TTGroupBoxArgs {
	customT.alignment = alignment
	return customT
}
func (customT *TTGroupBoxArgs) SetCustomSize(customSize Size) *TTGroupBoxArgs {
	customT.customSize = customSize
	return customT
}
func (customT *TTGroupBoxArgs) SetLayout(layout Layout) *TTGroupBoxArgs {
	customT.layout = layout
	return customT
}
func (customT *TTGroupBoxArgs) SetLayoutHBox(flagMarginsZero bool) *TTGroupBoxArgs {
	return customT.SetLayout(HBox{MarginsZero: flagMarginsZero})
}
func (customT *TTGroupBoxArgs) SetLayoutVBox(flagMarginsZero bool) *TTGroupBoxArgs {
	return customT.SetLayout(VBox{MarginsZero: flagMarginsZero})
}
func (customT *TTGroupBoxArgs) SetWidgets(widgets []Widget) *TTGroupBoxArgs {
	customT.widgets = widgets
	return customT
}
