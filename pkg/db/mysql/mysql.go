package mysqlstorage

import (
	"database/sql"
	"BlankZhu/Go-DNAT-Server/pkg/config"
	"BlankZhu/Go-DNAT-Server/pkg/entity"
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

var once = sync.Once{}
var mysql = MySQLStorage{}

type MySQLStorage struct {
	Config *config.MySQLConfig
	db     *sql.DB
}

func Get() *MySQLStorage {
	return &mysql
}

func Init(config *config.MySQLConfig) {
	once.Do(func() {
		mysql.Config = config
	})
}

func (s *MySQLStorage) Open() error {
	if s.Config == nil {
		return errors.New("Invalid MySQL config")
	}

	connURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		s.Config.Username, s.Config.Password,
		s.Config.IP, s.Config.Port, s.Config.Schema)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *MySQLStorage) Close() error {
	return s.db.Close()
}

func (s *MySQLStorage) Add(rule *entity.Rule) error {
	_, err := s.db.Query("INSERT INTO `t_rule` (`rule_id`, `cidr`, `destination`) VALUES (?, ?, ?)", rule.RuleID, rule.CIDR, rule.Dest)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStorage) Delete(id string) error {
	_, err := s.db.Query("DELETE FROM `t_rule` WHERE `rule_id` = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStorage) Update(rule *entity.Rule) error {
	_, err := s.db.Query("UPDATE `t_rule` SET `cidr` = ?, `destination` = ? WHERE `rule_id` = ?", rule.CIDR, rule.Dest, rule.RuleID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStorage) Get(id string) (*entity.Rule, bool, error) {
	res, err := s.db.Query("SELECT `rule_id`, `cidr`, `destination` FROM `t_rule` WHERE `rule_id` = ?", id)
	if err != nil {
		return nil, false, err
	}

	ret := entity.Rule{}
	found := false
	for res.Next() {
		err = res.Scan(&ret.RuleID, &ret.CIDR, &ret.Dest)
		if err != nil {
			return nil, false, err
		}
		found = true
	}

	return &ret, found, err
}

func (s *MySQLStorage) List() ([]*entity.Rule, error) {
	res, err := s.db.Query("SELECT `rule_id`, `cidr`, `destination` FROM `t_rule`")
	if err != nil {
		return nil, err
	}
	ret := make([]*entity.Rule, 0, 16)
	for res.Next() {
		rule := entity.Rule{}
		err = res.Scan(&rule.RuleID, &rule.CIDR, &rule.Dest)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &rule)
	}
	return ret, nil
}
