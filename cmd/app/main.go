package main

import (
	"context"
	"os"

	"github.com/GoldenDeals/DepGit/internal/git"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/joho/godotenv"
)

var log = logger.New("main")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	s, err := git.Init(git.Config{os.Getenv("DEPGIT_SSH_GIT_ADDRESS")})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(s.Serve(context.Background()))
}
