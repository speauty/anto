package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UIScrollView(args *UIScrollViewArgs) Widget {
	return ScrollView{
		MaxSize:         args.customSize,
		MinSize:         args.customSize,
		Visible:         args.visible,
		AssignTo:        args.assignTo,
		Children:        args.children,
		Layout:          args.layout,
		HorizontalFixed: args.horizontalFixed,
		VerticalFixed:   args.verticalFixed,
	}
}

func NewUIScrollViewArgs(assignTo **walk.ScrollView) *UIScrollViewArgs {
	return &UIScrollViewArgs{
		assignTo:        assignTo,
		visible:         true,
		layout:          VBox{MarginsZero: true},
		horizontalFixed: false,
		verticalFixed:   false,
	}
}

type UIScrollViewArgs struct {
	assignTo        **walk.ScrollView
	visible         bool
	customSize      Size
	children        []Widget
	layout          Layout
	horizontalFixed bool
	verticalFixed   bool
}

func (customT *UIScrollViewArgs) SetLayout(layout Layout) *UIScrollViewArgs {
	customT.layout = layout
	return customT
}

func (customT *UIScrollViewArgs) SetHorizontalFixed(horizontalFixed bool) *UIScrollViewArgs {
	customT.horizontalFixed = horizontalFixed
	return customT
}

func (customT *UIScrollViewArgs) HorizontalFixed() *UIScrollViewArgs {
	customT.horizontalFixed = true
	return customT
}

func (customT *UIScrollViewArgs) SetVerticalFixed(verticalFixed bool) *UIScrollViewArgs {
	customT.verticalFixed = verticalFixed
	return customT
}

func (customT *UIScrollViewArgs) VerticalFixed() *UIScrollViewArgs {
	customT.verticalFixed = true
	return customT
}

func (customT *UIScrollViewArgs) SetChildren(children []Widget) *UIScrollViewArgs {
	customT.children = children
	return customT
}

func (customT *UIScrollViewArgs) SetAssignTo(assignTo **walk.ScrollView) *UIScrollViewArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UIScrollViewArgs) SetVisible(visible bool) *UIScrollViewArgs {
	customT.visible = visible
	return customT
}

func (customT *UIScrollViewArgs) SetCustomSize(customSize Size) *UIScrollViewArgs {
	customT.customSize = customSize
	return customT
}
