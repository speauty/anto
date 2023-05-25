package repository

import (
	_type "anto/common"
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
	list     sync.Map
	names    []*_type.StdComboBoxModel
	namesAll []*_type.StdComboBoxModel
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
	customT.genNames2ComboBox()
}

func (customT *Translators) GetById(id string) serviceTranslator.ImplTranslator {
	obj, isExisted := customT.list.Load(id)
	if !isExisted {
		return nil
	}
	return obj.(serviceTranslator.ImplTranslator)
}

func (customT *Translators) GetNames() []*_type.StdComboBoxModel {
	return customT.names
}

func (customT *Translators) GetNamesAll() []*_type.StdComboBoxModel {
	return customT.namesAll
}

func (customT *Translators) genNames2ComboBox() {
	customT.names = []*_type.StdComboBoxModel{}
	customT.namesAll = []*_type.StdComboBoxModel{}
	customT.list.Range(func(idx, translator any) bool {
		customT.namesAll = append(customT.namesAll, &_type.StdComboBoxModel{
			Key:  translator.(serviceTranslator.ImplTranslator).GetId(),
			Name: translator.(serviceTranslator.ImplTranslator).GetName(),
		})
		if translator.(serviceTranslator.ImplTranslator).IsValid() {
			customT.names = append(customT.names, &_type.StdComboBoxModel{
				Key:  translator.(serviceTranslator.ImplTranslator).GetId(),
				Name: translator.(serviceTranslator.ImplTranslator).GetName(),
			})
		}
		return true
	})
	sort.Slice(customT.namesAll, func(i, j int) bool {
		return customT.namesAll[i].Key < customT.namesAll[j].Key
	})
	if len(customT.names) > 1 {
		sort.Slice(customT.names, func(i, j int) bool {
			return customT.names[i].Key < customT.names[j].Key
		})
	}
}
