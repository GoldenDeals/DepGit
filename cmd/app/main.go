// Package main is the entry point for the DepGit application.
// It initializes and runs the git server and web server with file storage capabilities.
package main

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/GoldenDeals/DepGit/internal/git"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/GoldenDeals/DepGit/internal/stroage"
	"github.com/GoldenDeals/DepGit/internal/web"
	"github.com/joho/godotenv"
)

var log = logger.New("main")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	printEnv()

	// Create a cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalCh
		log.WithField("signal", sig.String()).Info("Received shutdown signal")
		cancel()
	}()

	// Initialize storage
	fileStorage, err := stroage.NewFileStorage(os.Getenv("DEPGIT_STORAGE_PATH"))
	if err != nil {
		log.WithError(err).Fatalln("Failed to create file storage")
	}

	// Initialize git server
	gitServer, err := git.Init(git.Config{
		Address: os.Getenv("DEPGIT_SSH_GIT_ADDRESS"),
	}, fileStorage)
	if err != nil {
		log.WithError(err).Fatalln("Failed to initialize git server")
	}

	// Initialize web server
	webServer := web.New(web.Config{
		Address:     os.Getenv("DEPGIT_WEB_ADDRESS"),
		StaticDir:   os.Getenv("DEPGIT_WEB_STATIC_DIR"),
		APIBasePath: "/api/v1",
	}, web.NewAPIHandler())

	// Start servers in goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Start git server
	go func() {
		defer wg.Done()
		if err := gitServer.Serve(ctx); err != nil {
			log.WithError(err).Error("Git server error")
		}
	}()

	// Start web server
	go func() {
		defer wg.Done()
		if err := webServer.Start(ctx); err != nil {
			log.WithError(err).Error("Web server error")
		}
	}()

	// Wait for all servers to shut down
	wg.Wait()
	log.Info("All servers shut down, exiting")
}

func printEnv() {
	varf, err := os.Open(".env")
	if err != nil {
		log.WithError(err).Fatal("Failed to open .env file")
	}

	defer func() {
		if err := varf.Close(); err != nil {
			log.WithError(err).Error("Failed to close .env file")
		}
	}()

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
			Trace("Environment variable")
	}
}
