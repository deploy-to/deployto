package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gopkg.in/yaml.v2"
)

type server struct {
	connConfig *pgx.ConnConfig
}

func main() {
	server := newServer()
	http.HandleFunc("/", server.httpHandleFunc)
	err := http.ListenAndServe(getFromEnv("PORT", ":80"), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func newServer() *server {
	configPath := getFromEnv("CONFIG_PATH", "/config")

	configFile, err := os.ReadFile(filepath.Join(configPath, "postgresql.yaml"))
	if err != nil {
		fmt.Printf("config read %s error: #%v ", "postgresql.yaml", err)
		os.Exit(1)
	}
	config := &pgconn.Config{}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		fmt.Printf("Unable to unmarshal config: #%v", err)
		os.Exit(1)
	}
	connstring := fmt.Sprintf(
		"host=%s     port=%d      dbname=%s        user=%s      password=%s       target_session_attrs=read-write",
		config.Host, config.Port, config.Database, config.User, config.Password)
	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	pem, err := os.ReadFile(filepath.Join(configPath, "ssl_ca.crt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read ssl_ca.crt: %v\nContinue without SSL\n", err)
	} else {
		rootCertPool := x509.NewCertPool()
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			panic("Failed to append PEM.")
		}
		connConfig.TLSConfig = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}

	return &server{connConfig: connConfig}
}

func (s *server) httpHandleFunc(w http.ResponseWriter, r *http.Request) {
	db, err := pgx.ConnectConfig(r.Context(), s.connConfig)
	if err != nil {
		fmt.Printf("%s %v", r.URL, err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer db.Close(r.Context())

	err = db.Ping(r.Context())
	if err != nil {
		fmt.Printf("%s %v", r.URL, err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("postgresql - ok"))
	fmt.Printf("%s - ok", r.URL)
}

func getFromEnv(key, def string) string {
	value, valueExists := os.LookupEnv(key)
	if !valueExists {
		return def
	}
	return value
}
