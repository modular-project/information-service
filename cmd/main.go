package main

import (
	"fmt"
	"information-service/http/route"
	"information-service/model"
	"information-service/storage"
	"log"
	"net"
	"os"
)

func newDBConnection() storage.DBConnection {
	env := "INFO_DB_USER"
	u, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_DB_PWD"
	pass, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_DB_NAME"
	name, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_DB_HOST"
	host, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_DB_PORT"
	port, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	return storage.DBConnection{
		TypeDB:   storage.POSTGRESQL,
		User:     u,
		Password: pass,
		Port:     port,
		NameDB:   name,
		Host:     host,
	}
}

func main() {
	conn := newDBConnection()
	err := storage.NewDB(&conn)
	if err != nil {
		log.Fatal(err)
	}
	storage.DB().AutoMigrate(&model.Establishment{}, &model.Product{}, &model.Table{})
	env := "INFO_PORT"
	port, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := route.Start()

	log.Printf("Server started at :%s", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to server at :%s, got error: %s", port, err)
	}
}
