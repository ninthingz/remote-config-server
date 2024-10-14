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
	Enable     bool   `json:"enable"`
	CreateTime int    `json:"create_time"`
}

type ConfigHistoryDao struct{}

var configHistoryDao = &ConfigHistoryDao{}

func (dao *ConfigHistoryDao) create(configHistory ConfigHistory) error {
	_, err := db.Exec("INSERT INTO config_history(config_id, old_value, new_value, enable, create_time) VALUES(?, ?, ?, ?, ?)", configHistory.ConfigId, configHistory.OldValue, configHistory.NewValue, configHistory.Enable, configHistory.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ConfigHistoryDao) list(configId int) ([]*ConfigHistory, error) {
	rows, err := db.Query("SELECT * FROM config_history WHERE config_id = ? ORDER BY create_time DESC", configId)
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
		err := rows.Scan(&configHistory.Id, &configHistory.ConfigId, &configHistory.OldValue, &configHistory.NewValue, &configHistory.Enable, &configHistory.CreateTime)
		if err != nil {
			return nil, err
		}
		configHistoryList = append(configHistoryList, &configHistory)
	}
	return configHistoryList, nil
}
