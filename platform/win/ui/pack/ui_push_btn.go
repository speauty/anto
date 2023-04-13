package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UIPushBtn(args *UIPushBtnArgs) Widget {
	return PushButton{
		AssignTo:  args.assignTo,
		MinSize:   args.customSize,
		MaxSize:   args.customSize,
		Text:      args.text,
		OnClicked: args.onClicked,
	}
}

func NewUIPushBtnArgs(assignTo **walk.PushButton) *UIPushBtnArgs {
	return &UIPushBtnArgs{assignTo: assignTo, text: "btn"}
}

type UIPushBtnArgs struct {
	assignTo   **walk.PushButton
	customSize Size
	text       string
	onClicked  walk.EventHandler
}

func (customT *UIPushBtnArgs) SetAssignTo(assignTo **walk.PushButton) *UIPushBtnArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UIPushBtnArgs) SetCustomSize(customSize Size) *UIPushBtnArgs {
	customT.customSize = customSize
	return customT
}

func (customT *UIPushBtnArgs) SetText(text string) *UIPushBtnArgs {
	customT.text = text
	return customT
}

func (customT *UIPushBtnArgs) SetOnClicked(onClicked walk.EventHandler) *UIPushBtnArgs {
	customT.onClicked = onClicked
	return customT
}
