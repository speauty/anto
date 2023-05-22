package e_fyne

import (
	"fmt"
	"github.com/golang-module/carbon"
)

func (ag *AppGui) Pages() map[string]ImplPage {
	return ag.pages
}

func (ag *AppGui) RegisterPages(morePages ...ImplPage) {
	for _, pageItem := range morePages {
		if _, isExisted := ag.pages[pageItem.GetID()]; !isExisted {
			ag.pages[pageItem.GetID()] = pageItem
			if pageItem.IsDefault() {
				ag.chanNextPageId <- pageItem.GetID()
			}
			if !ag.config.IsRelease() {
				ag.PushToConsole(fmt.Sprintf("时间: %s, 注册页面: %s(%s)", carbon.Now(), pageItem.GetName(), pageItem.GetID()))
			}
			continue
		}
		ag.PushToConsole(fmt.Sprintf("时间: %s, 注册页面: %s(%s), 警告: 已注册", carbon.Now(), pageItem.GetName(), pageItem.GetID()))
	}
}

func (ag *AppGui) toPage(pageId string) {
	if ag.currentPageId == pageId {
		return
	}
	ag.chanNextPageId <- pageId
}

func (ag *AppGui) initPages() {
	ag.pages = make(map[string]ImplPage)
	ag.chanNextPageId = make(chan string, 1)
	ag.dispatcherPage()
}

func (ag *AppGui) dispatcherPage() {
	go func() {
		for {
			currentPageId := <-ag.chanNextPageId
			ag.currentPageId = currentPageId
			currentPage, isExisted := ag.pages[currentPageId]
			if !isExisted {
				ag.PushToConsole(fmt.Sprintf("时间: %s, 错误: 页面[%s]不存在, 疑似未注册", carbon.Now(), currentPageId))
				continue
			}
			ag.mainWindow.SetContent(currentPage.OnRender())
		}
	}()
}
