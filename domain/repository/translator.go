package repository

import (
	serviceTranslator "anto/domain/service/translator"
	"anto/lib/restrictor"
	"sort"
	"sync"
)

var (
	apiTranslators  *Translators
	onceTranslators sync.Once
)

func GetTranslators() *Translators {
	onceTranslators.Do(func() {
		apiTranslators = new(Translators)
	})
	return apiTranslators
}

type Translators struct {
	list  sync.Map
	names []string
}

func (customT *Translators) Register(translators ...serviceTranslator.ImplTranslator) {
	tmpRestrictor := restrictor.Singleton()
	for _, translator := range translators {
		if _, isExisted := customT.list.Load(translator.GetId()); isExisted {
			continue
		}
		customT.list.Store(translator.GetId(), translator)
		tmpLimiter := tmpRestrictor.Get(translator.GetId())
		limited := translator.GetCfg().GetQPS() / 4 * 3 // 缓冲
		if limited < 1 {
			limited = 1
		}
		tmpLimiter.SetLimit(1)
		tmpLimiter.SetBurst(limited)

		tmpRestrictor.Set(translator.GetId(), tmpLimiter)
	}
	customT.genNames()
}

func (customT *Translators) GetById(id string) serviceTranslator.ImplTranslator {
	obj, isExisted := customT.list.Load(id)
	if !isExisted {
		return nil
	}
	return obj.(serviceTranslator.ImplTranslator)
}

func (customT *Translators) GetByName(name string) (currentTranslator serviceTranslator.ImplTranslator) {
	customT.list.Range(func(id, translatorItem any) bool {
		if translatorItem.(serviceTranslator.ImplTranslator).GetName() == name {
			currentTranslator = translatorItem.(serviceTranslator.ImplTranslator)
			return false
		}
		return true
	})
	return
}

func (customT *Translators) GetNames() []string {
	return customT.names
}

func (customT *Translators) GetAllNames() []string {
	var names []string
	customT.list.Range(func(idx, translator any) bool {
		names = append(names, translator.(serviceTranslator.ImplTranslator).GetName())
		return true
	})

	if len(names) > 1 {
		sort.Slice(names, func(i, j int) bool {
			return names[i] < names[j]
		})
	}
	return names
}

func (customT *Translators) genNames() {
	customT.list.Range(func(idx, translator any) bool {
		if translator.(serviceTranslator.ImplTranslator).IsValid() {
			customT.names = append(customT.names, translator.(serviceTranslator.ImplTranslator).GetName())
		}
		return true
	})

	if len(customT.names) > 1 {
		sort.Slice(customT.names, func(i, j int) bool {
			return customT.names[i] < customT.names[j]
		})
	}
}
