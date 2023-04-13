package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UICheckBox(args *UICheckBoxArgs) CheckBox {
	return CheckBox{
		AssignTo:  args.assignTo,
		Checked:   args.checked,
		MinSize:   args.customSize,
		MaxSize:   args.customSize,
		Text:      args.text,
		OnClicked: args.onClickedFn,
	}
}

func NewUICheckBoxArgs(assignTo **walk.CheckBox) *UICheckBoxArgs {
	return &UICheckBoxArgs{assignTo: assignTo}
}

type UICheckBoxArgs struct {
	assignTo    **walk.CheckBox
	checked     bool
	customSize  Size
	text        string
	onClickedFn walk.EventHandler
}

func (customT *UICheckBoxArgs) SetAssignTo(assignTo **walk.CheckBox) *UICheckBoxArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UICheckBoxArgs) SetChecked(checked bool) *UICheckBoxArgs {
	customT.checked = checked
	return customT
}

func (customT *UICheckBoxArgs) SetCustomSize(customSize Size) *UICheckBoxArgs {
	customT.customSize = customSize
	return customT
}

func (customT *UICheckBoxArgs) SetText(text string) *UICheckBoxArgs {
	customT.text = text
	return customT
}

func (customT *UICheckBoxArgs) SetOnClickedFn(onClickedFn walk.EventHandler) *UICheckBoxArgs {
	customT.onClickedFn = onClickedFn
	return customT
}
