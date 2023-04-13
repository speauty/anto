package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UIComboBox(args *UIComboBoxArgs) Widget {
	return ComboBox{
		AssignTo:              args.assignTo,
		Model:                 args.model,
		CurrentIndex:          args.currentIdx,
		OnCurrentIndexChanged: args.onCurrentIdxChangedFn,
		DisplayMember:         args.displayMember,
		BindingMember:         args.bindingMember,
	}
}

func NewUIComboBoxArgs(assignTo **walk.ComboBox) *UIComboBoxArgs {
	return &UIComboBoxArgs{assignTo: assignTo}
}

type UIComboBoxArgs struct {
	assignTo              **walk.ComboBox
	model                 interface{}
	currentIdx            interface{}
	onCurrentIdxChangedFn walk.EventHandler
	displayMember         string
	bindingMember         string
}

func (customT *UIComboBoxArgs) SetDisplayMember(displayMember string) *UIComboBoxArgs {
	customT.displayMember = displayMember
	return customT
}

func (customT *UIComboBoxArgs) SetBindingMember(bindingMember string) *UIComboBoxArgs {
	customT.bindingMember = bindingMember
	return customT
}

func (customT *UIComboBoxArgs) SetAssignTo(assignTo **walk.ComboBox) *UIComboBoxArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UIComboBoxArgs) SetModel(model interface{}) *UIComboBoxArgs {
	customT.model = model
	return customT
}

func (customT *UIComboBoxArgs) SetCurrentIdx(currentIdx interface{}) *UIComboBoxArgs {
	customT.currentIdx = currentIdx
	return customT
}

func (customT *UIComboBoxArgs) SetOnCurrentIdxChangedFn(onCurrentIdxChangedFn walk.EventHandler) *UIComboBoxArgs {
	customT.onCurrentIdxChangedFn = onCurrentIdxChangedFn
	return customT
}
