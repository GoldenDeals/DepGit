package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/GoldenDeals/DepGit/internal/git"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/GoldenDeals/DepGit/internal/stroage"
	"github.com/joho/godotenv"
)

var log = logger.New("main")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	printEnv()

	fileStorage, err := stroage.NewFileStorage(os.Getenv("DEPGIT_STORAGE_PATH"))
	if err != nil {
		log.WithError(err).
			Fatalln("Failed to create file storage")
	}

	s, err := git.Init(git.Config{
		Address: os.Getenv("DEPGIT_SSH_GIT_ADDRESS"),
	}, fileStorage)
	if err != nil {
		log.WithError(err).
			Fatalln("Failed to initialize git server")
	}

	log.WithError(s.Serve(context.Background())).
		Fatalln("error serving git server")
}

func printEnv() {
	varf, err := os.Open(".env")
	if err != nil {
		log.WithError(err).
			Fatal("Failed to open .env file")
	}

	defer varf.Close()

	allvars := make(map[string]struct{})
	scanner := bufio.NewScanner(varf)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		arr := strings.Split(line, "=")
		if len(arr) != 2 {
			continue
		}

		key := arr[0]

		allvars[key] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.WithError(err).Warn("Error scanning .env file")
	}

	log.Trace("Environment variables:")
	for _, env := range os.Environ() {
		arr := strings.Split(env, "=")
		if len(arr) != 2 {
			continue
		}

		key := arr[0]
		val := arr[1]

		if _, exists := allvars[key]; !exists {
			continue
		}

		log.WithField("key", key).
			WithField("value", val).
			Trace("environment variable")
	}
}
