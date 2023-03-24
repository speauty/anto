package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTTextLabel(args *TTTextLabelArgs) Widget {
	return TextLabel{
		Enabled:  args.enabled,
		MaxSize:  args.customSize,
		MinSize:  args.customSize,
		Visible:  args.visible,
		AssignTo: args.assignTo,
		Text:     args.text,
	}
}

func NewTTTextLabelArgs(assignTo **walk.TextLabel) *TTTextLabelArgs {
	return &TTTextLabelArgs{
		assignTo: assignTo,
		visible:  true,
		enabled:  true,
	}
}

type TTTextLabelArgs struct {
	assignTo   **walk.TextLabel
	visible    bool
	enabled    bool
	text       string
	customSize Size
}

func (customT *TTTextLabelArgs) SetAssignTo(assignTo **walk.TextLabel) *TTTextLabelArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TTTextLabelArgs) SetVisible(visible bool) *TTTextLabelArgs {
	customT.visible = visible
	return customT
}

func (customT *TTTextLabelArgs) SetEnabled(enabled bool) *TTTextLabelArgs {
	customT.enabled = enabled
	return customT
}

func (customT *TTTextLabelArgs) SetText(text string) *TTTextLabelArgs {
	customT.text = text
	return customT
}

func (customT *TTTextLabelArgs) SetCustomSize(customSize Size) *TTTextLabelArgs {
	customT.customSize = customSize
	return customT
}
