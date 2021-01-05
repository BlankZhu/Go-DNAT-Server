package memstorage

import (
	"BlankZhu/Go-DNAT-Server/pkg/entity"
)

var cache = make(map[string]*entity.Rule)
var mem = MemStorage{}

type MemStorage struct{}

func Get() *MemStorage {
	return &mem
}

func (s MemStorage) Open() {
	// do nothing
}

func (s MemStorage) Close() {
	// clear all
	cache = make(map[string]*entity.Rule)
}

func (s MemStorage) Add(rule *entity.Rule) {
	cache[rule.RuleID] = rule
}

func (s MemStorage) Delete(id string) {
	delete(cache, id)
}

func (s MemStorage) Update(rule *entity.Rule) {
	cache[rule.RuleID] = rule
}

func (s MemStorage) Get(id string) (*entity.Rule, bool) {
	r, ok := cache[id]
	return r, ok
}

func (s MemStorage) List() []*entity.Rule {
	ret := make([]*entity.Rule, 0, len(cache))
	for _, rule := range cache {
		ret = append(ret, rule)
	}
	return ret
}
