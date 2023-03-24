package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTComboBox(args *TTComboBoxArgs) Widget {
	return ComboBox{
		AssignTo:              args.assignTo,
		Model:                 args.model,
		CurrentIndex:          args.currentIdx,
		OnCurrentIndexChanged: args.onCurrentIdxChangedFn,
	}
}

func NewTTComboBoxArgs(assignTo **walk.ComboBox) *TTComboBoxArgs {
	return &TTComboBoxArgs{assignTo: assignTo}
}

type TTComboBoxArgs struct {
	assignTo              **walk.ComboBox
	model                 interface{}
	currentIdx            int
	onCurrentIdxChangedFn walk.EventHandler
}

func (customT *TTComboBoxArgs) SetAssignTo(assignTo **walk.ComboBox) *TTComboBoxArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TTComboBoxArgs) SetModel(model interface{}) *TTComboBoxArgs {
	customT.model = model
	return customT
}

func (customT *TTComboBoxArgs) SetCurrentIdx(currentIdx int) *TTComboBoxArgs {
	customT.currentIdx = currentIdx
	return customT
}

func (customT *TTComboBoxArgs) SetOnCurrentIdxChangedFn(onCurrentIdxChangedFn walk.EventHandler) *TTComboBoxArgs {
	customT.onCurrentIdxChangedFn = onCurrentIdxChangedFn
	return customT
}
