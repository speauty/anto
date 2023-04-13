package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTLabel(args *TTLabelArgs) Widget {
	return Label{
		Enabled:  args.enabled,
		MaxSize:  args.customSize,
		MinSize:  args.customSize,
		Visible:  args.visible,
		AssignTo: args.assignTo,
		Text:     args.text,
	}
}

func NewTTLabelArgs(assignTo **walk.Label) *TTLabelArgs {
	return &TTLabelArgs{
		assignTo: assignTo,
		visible:  true,
		enabled:  true,
	}
}

type TTLabelArgs struct {
	assignTo   **walk.Label
	visible    bool
	enabled    bool
	text       string
	customSize Size
}

func (customT *TTLabelArgs) SetAssignTo(assignTo **walk.Label) *TTLabelArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TTLabelArgs) SetVisible(visible bool) *TTLabelArgs {
	customT.visible = visible
	return customT
}

func (customT *TTLabelArgs) SetEnabled(enabled bool) *TTLabelArgs {
	customT.enabled = enabled
	return customT
}

func (customT *TTLabelArgs) SetText(text string) *TTLabelArgs {
	customT.text = text
	return customT
}

func (customT *TTLabelArgs) SetCustomSize(customSize Size) *TTLabelArgs {
	customT.customSize = customSize
	return customT
}
