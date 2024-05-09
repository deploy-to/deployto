package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	_ "github.com/lib/pq"
)

type Config struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Dbname    string `mapstructure:"dbname"`
	EndpointB string `mapstructure:"endpointb"`
}

func (c *Config) any(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	requestURL := "http://" + c.EndpointB
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
}

func getConfig() (c Config, err error) {
	config.AddDriver(yaml.Driver)
	fmt.Printf("config file /config/config.yaml")
	// load more files
	err = config.LoadFiles("/config/config.yaml")
	// can also load multi at once
	// err := config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")
	if err != nil {
		return c, err
	}
	c = Config{}
	fmt.Printf("config data: \n %#v\n", config.Data())
	config.BindStruct("", &c)
	err = config.Decode(&c)
	if err != nil {
		return c, err
	}
	fmt.Printf("%+v\n", c)

	return c, err
}

func main() {
	c, _ := getConfig()
	http.HandleFunc("/", c.any)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
