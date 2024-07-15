package util

import (
	"context"
	"log"
	"main/prisma/db"
)

var database = db.NewClient()
var ctx = context.Background()

func init() {
	if err := database.Connect(); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}

func ReturnDatabase() *db.PrismaClient {
	return database
}

func ReturnContext() *context.Context {
	return &ctx
}
