package tt_ui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
)

type IPage interface {
	GetId() string
	GetName() string
	BindWindow(win *walk.MainWindow)
	SetVisible(isVisible bool)
	GetWidget() Widget
	Reset()
}

type PageCtl struct {
	current IPage
	pages   sync.Map
	menus   []MenuItem
}

func (customPC *PageCtl) PushPages(pages ...IPage) {
	for _, page := range pages {
		if _, isExist := customPC.pages.Load(page.GetId()); isExist {
			continue
		}
		customPC.pages.Store(page.GetId(), page)
	}
}

func (customPC *PageCtl) SetCurrent(pageId string) error {
	currentPage, err := customPC.GetPageById(pageId)
	if err != nil {
		return err
	}
	customPC.current = currentPage
	defer customPC.Render()
	return nil
}

func (customPC *PageCtl) Render() {
	customPC.pages.Range(func(pageId, currentPage any) bool {
		currentPage.(IPage).SetVisible(currentPage.(IPage).GetId() == customPC.current.GetId())
		currentPage.(IPage).Reset()
		return true
	})
}

func (customPC *PageCtl) GetPageById(pageId string) (IPage, error) {
	currentPage, isOk := customPC.pages.Load(pageId)
	if !isOk {
		return nil, fmt.Errorf("当前页面[%s]不存在", pageId)
	}
	return currentPage.(IPage), nil
}

func (customPC *PageCtl) Bind(win *walk.MainWindow) {
	customPC.pages.Range(func(pageId, currentPage any) bool {
		currentPage.(IPage).BindWindow(win)
		return true
	})
	return
}

func (customPC *PageCtl) GetWidgets() []Widget {
	var widgets []Widget
	customPC.pages.Range(func(pageId, currentPage any) bool {
		widgets = append(widgets, currentPage.(IPage).GetWidget())
		return true
	})
	return widgets
}
