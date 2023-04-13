package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UITextLabel(args *UITextLabelArgs) Widget {
	return TextLabel{
		Enabled:  args.enabled,
		MaxSize:  args.customSize,
		MinSize:  args.customSize,
		Visible:  args.visible,
		AssignTo: args.assignTo,
		Text:     args.text,
	}
}

func NewUITextLabelArgs(assignTo **walk.TextLabel) *UITextLabelArgs {
	return &UITextLabelArgs{
		assignTo: assignTo,
		visible:  true,
		enabled:  true,
	}
}

type UITextLabelArgs struct {
	assignTo   **walk.TextLabel
	visible    bool
	enabled    bool
	text       string
	customSize Size
}

func (customT *UITextLabelArgs) SetAssignTo(assignTo **walk.TextLabel) *UITextLabelArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UITextLabelArgs) SetVisible(visible bool) *UITextLabelArgs {
	customT.visible = visible
	return customT
}

func (customT *UITextLabelArgs) SetEnabled(enabled bool) *UITextLabelArgs {
	customT.enabled = enabled
	return customT
}

func (customT *UITextLabelArgs) SetText(text string) *UITextLabelArgs {
	customT.text = text
	return customT
}

func (customT *UITextLabelArgs) SetCustomSize(customSize Size) *UITextLabelArgs {
	customT.customSize = customSize
	return customT
}
