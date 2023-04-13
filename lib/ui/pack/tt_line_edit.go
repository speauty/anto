package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTLineEdit(args *LineEditWrapperArgs) Widget {
	return LineEdit{
		MinSize:       args.customSize,
		MaxSize:       args.customSize,
		AssignTo:      args.assignTo,
		ReadOnly:      args.readOnly,
		OnTextChanged: args.onTextChanged,
		Text:          args.text,
	}
}

func NewLineEditWrapperArgs(assignTo **walk.LineEdit) *LineEditWrapperArgs {
	return &LineEditWrapperArgs{assignTo: assignTo}
}

type LineEditWrapperArgs struct {
	assignTo      **walk.LineEdit
	customSize    Size
	readOnly      bool
	onTextChanged walk.EventHandler
	text          string
}

func (customT *LineEditWrapperArgs) SetText(text string) *LineEditWrapperArgs {
	customT.text = text
	return customT
}

func (customT *LineEditWrapperArgs) SetReadOnly(readOnly bool) *LineEditWrapperArgs {
	customT.readOnly = readOnly
	return customT
}

func (customT *LineEditWrapperArgs) SetOnTextChanged(onTextChanged walk.EventHandler) *LineEditWrapperArgs {
	customT.onTextChanged = onTextChanged
	return customT
}

func (customT *LineEditWrapperArgs) SetAssignTo(assignTo **walk.LineEdit) *LineEditWrapperArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *LineEditWrapperArgs) SetCustomSize(customSize Size) *LineEditWrapperArgs {
	customT.customSize = customSize
	return customT
}
