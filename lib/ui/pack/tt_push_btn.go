package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTPushBtn(args *TTPushBtnArgs) Widget {
	return PushButton{
		AssignTo:  args.assignTo,
		MinSize:   args.customSize,
		MaxSize:   args.customSize,
		Text:      args.text,
		OnClicked: args.onClicked,
	}
}

func NewTTPushBtnArgs(assignTo **walk.PushButton) *TTPushBtnArgs {
	return &TTPushBtnArgs{assignTo: assignTo, text: "btn"}
}

type TTPushBtnArgs struct {
	assignTo   **walk.PushButton
	customSize Size
	text       string
	onClicked  walk.EventHandler
}

func (customT *TTPushBtnArgs) SetAssignTo(assignTo **walk.PushButton) *TTPushBtnArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TTPushBtnArgs) SetCustomSize(customSize Size) *TTPushBtnArgs {
	customT.customSize = customSize
	return customT
}

func (customT *TTPushBtnArgs) SetText(text string) *TTPushBtnArgs {
	customT.text = text
	return customT
}

func (customT *TTPushBtnArgs) SetOnClicked(onClicked walk.EventHandler) *TTPushBtnArgs {
	customT.onClicked = onClicked
	return customT
}
