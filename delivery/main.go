package main

import (
	"github.com/InkForge/Blog_Website/infrastructures/db/mongo"
	"github.com/InkForge/Blog_Website/repositories"
)

func main() {
	client := mongo.NewMongoClient()
	db := client.Database("blogapp")

	userRepo := repositories.NewUserRepository(db)
	
	// other repo defined here as the above one

}
