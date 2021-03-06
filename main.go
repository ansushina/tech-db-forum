/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"log"
	"net/http"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
	sw "github.com/ansushina/tech-db-forum/app/http"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func main() {
	log.Printf("Server started")
	err := database.CreateConn()
	if err != nil {
		log.Fatal(err)
	}

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":5000", router))
}
