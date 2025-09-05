package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type DB struct {
	Name     string `yaml:"POSTGRES_DB"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Port     string `yaml:"PORT"`
	Host     string `yaml:"HOST"`

	Features struct {
		Quote bool `yaml:"FEATURES_QUOTE"`
	}
}

func PostgresConnection() (*sql.DB, error) {
	config := getDBConfig()

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getDBConfig() DB {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения config.yaml: %v", err)
	}

	var config DB
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Ошибка разбора config.yaml: %v", err)
	}

	return config
}
