package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type ConfigServer struct {
	Hostname string
	Port     int
	Database string
	Username string
	Password string
}

func main() {
	config := mustConfig()
	http.HandleFunc("/", config.httpHandleFunc)
	err := http.ListenAndServe(":80", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func mustConfig() *ConfigServer {
	yamlFile, err := os.ReadFile("/config/postgresql.yaml")
	if err != nil {
		fmt.Printf("config read %s error: #%v ", "/config/postgresql.yaml", err)
		return nil
	}

	config := &ConfigServer{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		fmt.Printf("config unmarshal error: #%v", err)
		return nil
	}
	return config
}

func (config *ConfigServer) httpHandleFunc(w http.ResponseWriter, r *http.Request) {
	if config == nil {
		fmt.Printf("%s - config is nil", r.URL)
		http.Error(w, "config is nil", 500)
		return
	}
	if config.Hostname == "" {
		fmt.Printf("%s - hostname not set", r.URL)
		http.Error(w, "hostname not set", 500)
	}
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Hostname, config.Port, config.Username, config.Password, config.Database)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Printf("%s %v", r.URL, err)
		http.Error(w, err.Error(), 500)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Printf("%s %v", r.URL, err)
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("postgresql - ok"))
	fmt.Printf("%s - ok", r.URL)
}
