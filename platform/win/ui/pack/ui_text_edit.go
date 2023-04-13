package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UITextEdit(args *UITextEditArgs) Widget {
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

func NewUITextEditArgs(assignTo **walk.TextEdit) *UITextEditArgs {
	return &UITextEditArgs{assignTo: assignTo}
}

type UITextEditArgs struct {
	assignTo      **walk.TextEdit
	customSize    Size
	vScroll       bool
	hScroll       bool
	readOnly      bool
	onTextChanged walk.EventHandler
}

func (customT *UITextEditArgs) SetReadOnly(readOnly bool) *UITextEditArgs {
	customT.readOnly = readOnly
	return customT
}

func (customT *UITextEditArgs) SetOnTextChanged(onTextChanged walk.EventHandler) *UITextEditArgs {
	customT.onTextChanged = onTextChanged
	return customT
}

func (customT *UITextEditArgs) SetAssignTo(assignTo **walk.TextEdit) *UITextEditArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UITextEditArgs) SetCustomSize(customSize Size) *UITextEditArgs {
	customT.customSize = customSize
	return customT
}

func (customT *UITextEditArgs) SetVScroll(vScroll bool) *UITextEditArgs {
	customT.vScroll = vScroll
	return customT
}

func (customT *UITextEditArgs) SetHScroll(hScroll bool) *UITextEditArgs {
	customT.hScroll = hScroll
	return customT
}
