package tt_ui

import (
	"fmt"
	"github.com/lxn/walk"
	"sync"
)

type IPage interface {
	GetId() int
	GetName() string
	BindWindow(win walk.Form)
	SetVisible(isVisible bool)
	GetWidget() *walk.Widget
	Reset()
}

type pageCtl struct {
	current IPage
	pages   sync.Map
}

func (customPC *pageCtl) setCurrent(pageId int) error {
	currentPage, err := customPC.getPageById(pageId)
	if err != nil {
		return err
	}
	customPC.current = currentPage
	defer customPC.render()
	return nil
}

func (customPC *pageCtl) render() {
	customPC.pages.Range(func(pageId, currentPage any) bool {
		currentPage.(IPage).SetVisible(currentPage.(IPage).GetId() == customPC.current.GetId())
		currentPage.(IPage).Reset()
		return true
	})
}

func (customPC *pageCtl) getPageById(pageId int) (IPage, error) {
	currentPage, isOk := customPC.pages.Load(pageId)
	if !isOk {
		return nil, fmt.Errorf("当前页面[%d]不存在", pageId)
	}
	return currentPage.(IPage), nil
}
