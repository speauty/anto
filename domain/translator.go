package domain

import (
	"fmt"
	"sync"
	"translator/tst/tt_log"
	"translator/tst/tt_translator"
	_type "translator/type"
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

func (customT *Translators) Register(translators ...tt_translator.ITranslator) {
	for _, translator := range translators {
		if _, isExisted := customT.list.Load(translator.GetId()); isExisted {
			continue
		}
		customT.list.Store(translator.GetId(), translator)
	}
	customT.genNames2ComboBox()
}

func (customT *Translators) GetById(id string) tt_translator.ITranslator {
	obj, isExisted := customT.list.Load(id)
	if !isExisted {
		return nil
	}
	return obj.(tt_translator.ITranslator)
}

func (customT *Translators) GetNames() []*_type.StdComboBoxModel {
	return customT.names
}

func (customT *Translators) genNames2ComboBox() {
	customT.names = []*_type.StdComboBoxModel{}
	customT.list.Range(func(idx, translator any) bool {
		if translator.(tt_translator.ITranslator).IsValid() {
			customT.names = append(customT.names, &_type.StdComboBoxModel{
				Key:  translator.(tt_translator.ITranslator).GetId(),
				Name: translator.(tt_translator.ITranslator).GetName(),
			})
		} else {
			tt_log.GetInstance().Warn(fmt.Sprintf("当前翻译引擎无效: %s", translator.(tt_translator.ITranslator).GetName()))
		}
		return true
	})
}
