package storage

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DRIVER string

const (
	POSTGRESQL DRIVER = "POSTGRES"
	TESTING    DRIVER = "TESTING"
)

var (
	db   *gorm.DB
	once sync.Once
)

type dbUser struct {
	User     string
	Password string
	Port     string
	NameDB   string
	Host     string
}

func New(driver DRIVER) {
	once.Do(func() {
		u := loadData()
		switch driver {
		case POSTGRESQL:
			newPostgresDB(&u)
		case TESTING:
			newTestingDB(&u)
		}
	})
}

func newTestingDB(u *dbUser) {
	var err error
	fmt.Print(u)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", u.User, u.Password, u.Host, u.Port, "testing")
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}

	fmt.Println("conectado a Testing")
}

func newPostgresDB(u *dbUser) {
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", u.User, u.Password, u.Host, u.Port, u.NameDB)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}

	fmt.Println("conectado a postgres")
}

// DB return a unique instance of db
func DB() *gorm.DB {
	return db
}

func getEnv(env string) (string, error) {
	s, f := os.LookupEnv(env)
	if !f {
		return "", fmt.Errorf("environment variable (%s) not found", env)
	}
	return s, nil
}

func loadData() dbUser {
	user, err := getEnv("DB_USER")
	if err != nil {
		log.Fatalf(err.Error())
	}
	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		log.Fatalf(err.Error())
	}
	port, err := getEnv("INFO_DB_PORT")
	if err != nil {
		log.Fatalf(err.Error())
	}
	name, err := getEnv("INFO_DB_NAME")
	if err != nil {
		log.Fatalf(err.Error())
	}
	host, err := getEnv("INFO_DB_HOST")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return dbUser{user, password, port, name, host}
}
