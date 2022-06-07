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

var (
	host string
	port string
)

func init() {
	env := "DB_TYPE"
	driver, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_HOST"
	host, f = os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_PORT"
	port, f = os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	storage.New(storage.DRIVER(driver))
	storage.DB().AutoMigrate(&model.Establishment{}, &model.Product{}, &model.Table{})
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := route.Start()

	log.Printf("Server started at %s:%s", host, port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to server at %s:%s, got error: %s", host, port, err)
	}
}
