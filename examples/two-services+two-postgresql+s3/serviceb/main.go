package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// const (
// 	host            = "serviceb-db"
// 	port            = 5432
// 	user            = "postgres"
// 	password        = "HOrFk14CyX"
// 	dbname          = "postgres"
// 	endpoint        = "serviceb-s3"
// 	accessKeyID     = "Lh25XnG11NuoGIPbCTHc"
// 	secretAccessKey = "Yjl2T2N2yKda6g7ebGfljUnCy6CPwv33L2rDkOMc"
// 	useSSL          = false
// )

type Config struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	BucketName      string `mapstructure:"bucketName"`
}

func (c *Config) any(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

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

	ctx := context.Background()
	// Initialize minio client object.
	minioClient, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient)

	objectName := "test"
	filePath := "/tmp/testdata"
	contentType := "application/octet-stream"

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println(f.Name())

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, c.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	minioClient.RemoveObject(ctx, c.BucketName, objectName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully delete %s ", objectName)

	io.WriteString(w, "This is my website!\n")
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
	// config.ParseEnv: will parse env var in string value. eg: shell: ${SHELL}

	// can also
	// fmt.Printf("config data: \n %#v\n", config.Data())
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
		fmt.Printf("error starting server: %s\n", err)
		panic(err)
	}
}
