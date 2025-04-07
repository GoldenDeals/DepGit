package main

import (
	"fmt"

	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var log = logger.New("main")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	err = fmt.Errorf("some error")
	log.
		WithField("user_id", uuid.New().String()).
		WithError(err).
		Error("error creating user")

	// s, err := git.Init(git.Config{os.Getenv("DEPGIT_SSH_GIT_ADDRESS")})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Fatalln(s.Serve(context.Background()))
}
