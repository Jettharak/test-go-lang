package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type serverConfig struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
}

func loadServerConfig(path string) (serverConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return serverConfig{}, fmt.Errorf("open server config: %w", err)
	}
	defer file.Close()

	var config serverConfig
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return serverConfig{}, fmt.Errorf("decode server config: %w", err)
	}

	var extra yaml.Node
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return serverConfig{}, errors.New("decode server config: multiple YAML documents are not supported")
		}
		return serverConfig{}, fmt.Errorf("decode server config: %w", err)
	}
	if config.Server.Address == "" {
		return serverConfig{}, errors.New("server.address must not be empty")
	}

	return config, nil
}

func startServer() error {
	config, err := loadServerConfig("config.yml")
	if err != nil {
		return err
	}
	log.Printf("Starting server on %s", config.Server.Address)
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addHandler)
	return http.ListenAndServe(config.Server.Address, mux)
}
