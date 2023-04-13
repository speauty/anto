package domain

import (
	serviceTranslator "anto/dependency/service/translator"
	"anto/lib/log"
	_type "anto/type"
	"fmt"
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
	names []*_type.StdComboBoxModel
}

func (customT *Translators) Register(translators ...serviceTranslator.InterfaceTranslator) {
	for _, translator := range translators {
		if _, isExisted := customT.list.Load(translator.GetId()); isExisted {
			continue
		}
		customT.list.Store(translator.GetId(), translator)
	}
	customT.genNames2ComboBox()
}

func (customT *Translators) GetById(id string) serviceTranslator.InterfaceTranslator {
	obj, isExisted := customT.list.Load(id)
	if !isExisted {
		return nil
	}
	return obj.(serviceTranslator.InterfaceTranslator)
}

func (customT *Translators) GetNames() []*_type.StdComboBoxModel {
	return customT.names
}

func (customT *Translators) genNames2ComboBox() {
	customT.names = []*_type.StdComboBoxModel{}
	customT.list.Range(func(idx, translator any) bool {
		if translator.(serviceTranslator.InterfaceTranslator).IsValid() {
			customT.names = append(customT.names, &_type.StdComboBoxModel{
				Key:  translator.(serviceTranslator.InterfaceTranslator).GetId(),
				Name: translator.(serviceTranslator.InterfaceTranslator).GetName(),
			})
		} else {
			log.Singleton().Warn(fmt.Sprintf("当前翻译引擎无效: %s", translator.(serviceTranslator.InterfaceTranslator).GetName()))
		}
		return true
	})

	if len(customT.names) > 1 {
		sort.Slice(customT.names, func(i, j int) bool {
			return customT.names[i].Key < customT.names[j].Key
		})
	}
}
