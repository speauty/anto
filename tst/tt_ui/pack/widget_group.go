package pack

import . "github.com/lxn/walk/declarative"

func NewWidgetGroup() *WidgetGroup {
	return new(WidgetGroup)
}

type WidgetGroup struct {
	widgets []Widget
}

func (customT *WidgetGroup) Append(widgets ...Widget) *WidgetGroup {
	if len(widgets) > 0 {
		customT.widgets = append(customT.widgets, widgets...)
	}
	return customT
}

func (customT *WidgetGroup) AppendZeroHSpacer() *WidgetGroup {
	customT.widgets = append(customT.widgets, HSpacer{})
	return customT
}

func (customT *WidgetGroup) AppendZeroVSpacer() *WidgetGroup {
	customT.widgets = append(customT.widgets, VSpacer{})
	return customT
}

func (customT *WidgetGroup) GetWidgets() []Widget {
	return customT.widgets
}
