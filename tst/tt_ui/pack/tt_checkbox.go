package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTCheckBox(args *TTCheckBoxArgs) CheckBox {
	return CheckBox{
		AssignTo:  args.assignTo,
		Checked:   args.checked,
		MinSize:   args.customSize,
		MaxSize:   args.customSize,
		Text:      args.text,
		OnClicked: args.onClickedFn,
	}
}

func NewTTCheckBoxArgs(assignTo **walk.CheckBox) *TTCheckBoxArgs {
	return &TTCheckBoxArgs{assignTo: assignTo}
}

type TTCheckBoxArgs struct {
	assignTo    **walk.CheckBox
	checked     bool
	customSize  Size
	text        string
	onClickedFn walk.EventHandler
}

func (customT *TTCheckBoxArgs) SetAssignTo(assignTo **walk.CheckBox) *TTCheckBoxArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TTCheckBoxArgs) SetChecked(checked bool) *TTCheckBoxArgs {
	customT.checked = checked
	return customT
}

func (customT *TTCheckBoxArgs) SetCustomSize(customSize Size) *TTCheckBoxArgs {
	customT.customSize = customSize
	return customT
}

func (customT *TTCheckBoxArgs) SetText(text string) *TTCheckBoxArgs {
	customT.text = text
	return customT
}

func (customT *TTCheckBoxArgs) SetOnClickedFn(onClickedFn walk.EventHandler) *TTCheckBoxArgs {
	customT.onClickedFn = onClickedFn
	return customT
}
