package storage

import (
	"fmt"
	"log"
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

type DBConnection struct {
	TypeDB   DRIVER
	User     string
	Password string
	Port     string
	NameDB   string
	Host     string
}

func NewDB(conn *DBConnection) error {
	var err error
	once.Do(func() {
		switch conn.TypeDB {
		case POSTGRESQL:
			err = newPostgresDB(conn)
		default:
			err = fmt.Errorf("invalid database type")
		}
	})
	return err
}

func newTestingDB(conn *DBConnection) error {
	var err error
	fmt.Print(conn)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conn.User, conn.Password, conn.Host, conn.Port, "testing")
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return fmt.Errorf("no se pudo abrir la base de datos: %v", err)
	}

	fmt.Println("conectado a Testing")
	return nil
}

func newPostgresDB(conn *DBConnection) error {
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conn.User, conn.Password, conn.Host, conn.Port, conn.NameDB)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return fmt.Errorf("no se pudo abrir la base de datos: %v", err)
	}

	log.Println("conectado a postgres")
	return nil
}

// DB return a unique instance of db
func DB() *gorm.DB {
	return db
}
