package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UILabel(args *UILabelArgs) Widget {
	return Label{
		Enabled:  args.enabled,
		MaxSize:  args.customSize,
		MinSize:  args.customSize,
		Visible:  args.visible,
		AssignTo: args.assignTo,
		Text:     args.text,
	}
}

func NewUILabelArgs(assignTo **walk.Label) *UILabelArgs {
	return &UILabelArgs{
		assignTo: assignTo,
		visible:  true,
		enabled:  true,
	}
}

type UILabelArgs struct {
	assignTo   **walk.Label
	visible    bool
	enabled    bool
	text       string
	customSize Size
}

func (customT *UILabelArgs) SetAssignTo(assignTo **walk.Label) *UILabelArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UILabelArgs) SetVisible(visible bool) *UILabelArgs {
	customT.visible = visible
	return customT
}

func (customT *UILabelArgs) SetEnabled(enabled bool) *UILabelArgs {
	customT.enabled = enabled
	return customT
}

func (customT *UILabelArgs) SetText(text string) *UILabelArgs {
	customT.text = text
	return customT
}

func (customT *UILabelArgs) SetCustomSize(customSize Size) *UILabelArgs {
	customT.customSize = customSize
	return customT
}
