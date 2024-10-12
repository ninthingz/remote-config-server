package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Id      int    `json:"id" json:"id,omitempty"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Secret  string `json:"secret"`
}

type ConfigDao struct {
}

var configDao = &ConfigDao{}

func (dao *ConfigDao) create(config Config) error {
	_, err := db.Exec("INSERT INTO config(name, value, message, secret) VALUES(?, ?, ?, ?)", config.Name, config.Value, config.Message, config.Secret)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ConfigDao) update(config Config) error {
	_, err := db.Exec("UPDATE config SET name = ?, value = ?, message = ?, secret = ? WHERE id = ?", config.Name, config.Value, config.Message, config.Secret, config.Id)
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
	err := db.QueryRow("SELECT * FROM config WHERE id = ?", id).Scan(&config.Id, &config.Name, &config.Value, &config.Message, &config.Secret)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (dao *ConfigDao) list() ([]*Config, error) {
	rows, err := db.Query("SELECT * FROM config")
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
		err := rows.Scan(&config.Id, &config.Name, &config.Value, &config.Message, &config.Secret)
		if err != nil {
			return nil, err
		}
		configList = append(configList, &config)
	}
	return configList, nil
}

func (dao *ConfigDao) getByName(name string) (*Config, error) {
	var config Config
	err := db.QueryRow("SELECT * FROM config WHERE name = ?", name).Scan(&config.Id, &config.Name, &config.Value, &config.Message, &config.Secret)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
