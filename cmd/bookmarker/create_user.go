package main

import (
	"bookmarker/internal/dbutil"
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"fmt"
	"log"
)

func createUserCommand(username, password string) {
	db, err := dbutil.OpenSqliteDB("../../internal/db/bookmarker_db1.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)
	user, err := service.CreateUser(username, password)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("User created: ID=%d, Username=%s\n", user.ID, user.Username)
}
