package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTTextEdit(args *TextEditWrapperArgs) Widget {
	return TextEdit{
		MinSize:       args.customSize,
		MaxSize:       args.customSize,
		AssignTo:      args.assignTo,
		VScroll:       args.vScroll,
		HScroll:       args.hScroll,
		OnTextChanged: args.onTextChanged,
		ReadOnly:      args.readOnly,
	}
}

func NewTextEditWrapperArgs(assignTo **walk.TextEdit) *TextEditWrapperArgs {
	return &TextEditWrapperArgs{assignTo: assignTo}
}

type TextEditWrapperArgs struct {
	assignTo      **walk.TextEdit
	customSize    Size
	vScroll       bool
	hScroll       bool
	readOnly      bool
	onTextChanged walk.EventHandler
}

func (customT *TextEditWrapperArgs) SetReadOnly(readOnly bool) *TextEditWrapperArgs {
	customT.readOnly = readOnly
	return customT
}

func (customT *TextEditWrapperArgs) SetOnTextChanged(onTextChanged walk.EventHandler) *TextEditWrapperArgs {
	customT.onTextChanged = onTextChanged
	return customT
}

func (customT *TextEditWrapperArgs) SetAssignTo(assignTo **walk.TextEdit) *TextEditWrapperArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TextEditWrapperArgs) SetCustomSize(customSize Size) *TextEditWrapperArgs {
	customT.customSize = customSize
	return customT
}

func (customT *TextEditWrapperArgs) SetVScroll(vScroll bool) *TextEditWrapperArgs {
	customT.vScroll = vScroll
	return customT
}

func (customT *TextEditWrapperArgs) SetHScroll(hScroll bool) *TextEditWrapperArgs {
	customT.hScroll = hScroll
	return customT
}
