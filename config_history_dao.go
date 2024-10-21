package main

import (
	"database/sql"
	"log"
)

type ConfigHistory struct {
	Id         int    `json:"id" json:"id,omitempty"`
	ConfigId   int    `json:"config_id"`
	OldValue   string `json:"old_value"`
	NewValue   string `json:"new_value"`
	Nickname   string `json:"nickname"`
	Enable     bool   `json:"enable"`
	Message    string `json:"message"`
	CreateTime int64  `json:"create_time"`
}

type ConfigHistoryDao struct{}

var configHistoryDao = &ConfigHistoryDao{}

func (dao *ConfigHistoryDao) create(configHistory ConfigHistory) error {
	_, err := db.Exec("INSERT INTO config_history(config_id, old_value, new_value, nickname, enable, message, create_time) VALUES(?, ?, ?, ?, ?, ?, ?)", configHistory.ConfigId, configHistory.OldValue, configHistory.NewValue, configHistory.Nickname, configHistory.Enable, configHistory.Message, configHistory.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ConfigHistoryDao) list(configId int, pageSize int, pageIndex int) ([]*ConfigHistory, error) {
	rows, err := db.Query("SELECT * FROM config_history WHERE config_id = ? ORDER BY create_time DESC LIMIT ? OFFSET ?", configId, pageSize, (pageIndex-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var configHistoryList []*ConfigHistory
	for rows.Next() {
		var configHistory ConfigHistory
		err := rows.Scan(&configHistory.Id, &configHistory.ConfigId, &configHistory.OldValue, &configHistory.NewValue, &configHistory.Nickname, &configHistory.Enable, &configHistory.Message, &configHistory.CreateTime)
		if err != nil {
			return nil, err
		}
		configHistoryList = append(configHistoryList, &configHistory)
	}
	return configHistoryList, nil
}

func (dao *ConfigHistoryDao) getCount(configId int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM config_history WHERE config_id = ?", configId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}
