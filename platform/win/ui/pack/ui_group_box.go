package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UIGroupBox(args *UIGroupBoxArgs) Widget {
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

func NewUIGroupBoxArgs(assignTo **walk.GroupBox) *UIGroupBoxArgs {
	return &UIGroupBoxArgs{
		assignTo:  assignTo,
		alignment: AlignHCenterVCenter,
		layout:    VBox{},
		visible:   true,
	}
}

type UIGroupBoxArgs struct {
	title      string
	alignment  Alignment2D
	assignTo   **walk.GroupBox
	customSize Size
	layout     Layout
	widgets    []Widget
	visible    bool
}

func (customT *UIGroupBoxArgs) SetVisible(visible bool) *UIGroupBoxArgs {
	customT.visible = visible
	return customT
}

func (customT *UIGroupBoxArgs) SetTitle(title string) *UIGroupBoxArgs {
	customT.title = title
	return customT
}

func (customT *UIGroupBoxArgs) SetAssignTo(assignTo **walk.GroupBox) *UIGroupBoxArgs {
	customT.assignTo = assignTo
	return customT
}
func (customT *UIGroupBoxArgs) SetAlignment(alignment Alignment2D) *UIGroupBoxArgs {
	customT.alignment = alignment
	return customT
}
func (customT *UIGroupBoxArgs) SetCustomSize(customSize Size) *UIGroupBoxArgs {
	customT.customSize = customSize
	return customT
}
func (customT *UIGroupBoxArgs) SetLayout(layout Layout) *UIGroupBoxArgs {
	customT.layout = layout
	return customT
}
func (customT *UIGroupBoxArgs) SetLayoutHBox(flagMarginsZero bool) *UIGroupBoxArgs {
	return customT.SetLayout(HBox{MarginsZero: flagMarginsZero})
}
func (customT *UIGroupBoxArgs) SetLayoutVBox(flagMarginsZero bool) *UIGroupBoxArgs {
	return customT.SetLayout(VBox{MarginsZero: flagMarginsZero})
}
func (customT *UIGroupBoxArgs) SetWidgets(widgets []Widget) *UIGroupBoxArgs {
	customT.widgets = widgets
	return customT
}
