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

func (customPC *pageCtl) PushPages(pages ...IPage) {
	for _, page := range pages {
		if _, isExist := customPC.pages.Load(page.GetId()); isExist {
			continue
		}
		customPC.pages.Store(page.GetId(), page)
	}
}

func (customPC *pageCtl) SetCurrent(pageId int) error {
	currentPage, err := customPC.GetPageById(pageId)
	if err != nil {
		return err
	}
	customPC.current = currentPage
	defer customPC.Render()
	return nil
}

func (customPC *pageCtl) Render() {
	customPC.pages.Range(func(pageId, currentPage any) bool {
		currentPage.(IPage).SetVisible(currentPage.(IPage).GetId() == customPC.current.GetId())
		currentPage.(IPage).Reset()
		return true
	})
}

func (customPC *pageCtl) GetPageById(pageId int) (IPage, error) {
	currentPage, isOk := customPC.pages.Load(pageId)
	if !isOk {
		return nil, fmt.Errorf("当前页面[%d]不存在", pageId)
	}
	return currentPage.(IPage), nil
}
