package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Config struct {
	Id          int    `json:"id" json:"id,omitempty"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Message     string `json:"message"`
	Secret      string `json:"secret"`
	LastGetTime int64  `json:"last_get_time" json:"last_get_time,omitempty"`
	Enable      bool   `json:"enable"`
}

type ConfigDao struct {
}

var configDao = &ConfigDao{}

func (dao *ConfigDao) create(config Config, nickname string) error {
	config.LastGetTime = time.Now().Unix()
	_, err := db.Exec("INSERT INTO config(name, value, message, secret, last_get_time, enable) VALUES(?, ?, ?, ?, ?, ?)", config.Name, config.Value, config.Message, config.Secret, config.LastGetTime, config.Enable)
	if err != nil {
		return err
	}
	selectConfig, err := configDao.getByName(config.Name)
	if err != nil {
		return err
	}

	configHistory := ConfigHistory{
		ConfigId:   selectConfig.Id,
		OldValue:   "",
		NewValue:   selectConfig.Value,
		Nickname:   nickname,
		Enable:     selectConfig.Enable,
		Message:    selectConfig.Message,
		CreateTime: time.Now().Unix(),
	}

	err = configHistoryDao.create(configHistory)
	if err != nil {
		return err
	}

	return nil
}

func (dao *ConfigDao) update(config Config, nickname string) error {

	oldConfig, err := configDao.get(config.Id)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE config SET name = ?, value = ?, message = ?, secret = ?, enable = ? WHERE id = ?", config.Name, config.Value, config.Message, config.Secret, config.Enable, config.Id)
	if err != nil {
		return err
	}

	if oldConfig.Value == config.Value && oldConfig.Enable == oldConfig.Enable {
		return nil
	}

	configHistory := ConfigHistory{
		ConfigId:   config.Id,
		OldValue:   oldConfig.Value,
		NewValue:   config.Value,
		Nickname:   nickname,
		Enable:     config.Enable,
		Message:    config.Message,
		CreateTime: time.Now().Unix(),
	}

	err = configHistoryDao.create(configHistory)
	if err != nil {
		return err
	}

	return nil
}

func (dao *ConfigDao) delete(id int) error {
	_, err := db.Exec("DELETE FROM config WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ConfigDao) get(id int) (*Config, error) {
	var config Config
	err := db.QueryRow("SELECT * FROM config WHERE id = ?", id).Scan(&config.Id, &config.Name, &config.Value, &config.Message, &config.Secret, &config.LastGetTime, &config.Enable)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (dao *ConfigDao) list(keyword string, pageSize int, pageIndex int, orderBy string, sortType int) ([]*Config, error) {
	sortTypeStr := ""
	if sortType == 1 {
		sortTypeStr = "ASC"
	} else {
		sortTypeStr = "DESC"
	}
	if orderBy != "id" && orderBy != "name" && orderBy != "last_get_time" && orderBy != "enable" {
		orderBy = "id"
	}
	rows, err := db.Query("SELECT * FROM config WHERE name LIKE ? ORDER BY "+orderBy+" "+sortTypeStr+" LIMIT ? OFFSET ?", "%"+keyword+"%", pageSize, pageSize*(pageIndex-1))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			println("close rows error")
			return
		}
	}(rows)

	var configList []*Config
	for rows.Next() {
		var config Config
		err := rows.Scan(&config.Id, &config.Name, &config.Value, &config.Message, &config.Secret, &config.LastGetTime, &config.Enable)
		if err != nil {
			return nil, err
		}
		configList = append(configList, &config)
	}
	return configList, nil
}

func (dao *ConfigDao) getCount(keyword string) (int, error) {
	var count int
	err := db.QueryRow("SELECT count(*) FROM config WHERE name LIKE ?", "%"+keyword+"%").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dao *ConfigDao) getByName(name string) (*Config, error) {
	var config Config
	err := db.QueryRow("SELECT * FROM config WHERE name = ?", name).Scan(&config.Id, &config.Name, &config.Value, &config.Message, &config.Secret, &config.LastGetTime, &config.Enable)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (dao *ConfigDao) updateLastGetTime(config *Config) error {
	_, err := db.Exec("UPDATE config SET last_get_time = ? WHERE id = ?", time.Now().Unix(), config.Id)
	if err != nil {
		return err
	}
	return nil
}
