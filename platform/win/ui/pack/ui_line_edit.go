package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UILineEdit(args *UILineEditArgs) Widget {
	return LineEdit{
		MinSize:       args.customSize,
		MaxSize:       args.customSize,
		AssignTo:      args.assignTo,
		ReadOnly:      args.readOnly,
		OnTextChanged: args.onTextChanged,
		Text:          args.text,
	}
}

func NewUILineEditArgs(assignTo **walk.LineEdit) *UILineEditArgs {
	return &UILineEditArgs{assignTo: assignTo}
}

type UILineEditArgs struct {
	assignTo      **walk.LineEdit
	customSize    Size
	readOnly      bool
	onTextChanged walk.EventHandler
	text          string
}

func (customT *UILineEditArgs) SetText(text string) *UILineEditArgs {
	customT.text = text
	return customT
}

func (customT *UILineEditArgs) SetReadOnly(readOnly bool) *UILineEditArgs {
	customT.readOnly = readOnly
	return customT
}

func (customT *UILineEditArgs) SetOnTextChanged(onTextChanged walk.EventHandler) *UILineEditArgs {
	customT.onTextChanged = onTextChanged
	return customT
}

func (customT *UILineEditArgs) SetAssignTo(assignTo **walk.LineEdit) *UILineEditArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UILineEditArgs) SetCustomSize(customSize Size) *UILineEditArgs {
	customT.customSize = customSize
	return customT
}
