package db

import "BlankZhu/Go-DNAT-Server/pkg/entity"

type DataStorage interface {
	Open() error
	Close() error
	Add(rule *entity.Rule) error
	Delete(id string) error
	Update(rule *entity.Rule) error
	Get(id string) (*entity.Rule, bool, error)
	List() ([]*entity.Rule, error)
}
