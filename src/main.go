package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aune/prismic-drone-trigger/v2/drone"
	"github.com/aune/prismic-drone-trigger/v2/prismic"

	"github.com/gin-gonic/gin"
)

const (
	envDroneHost     = "DRONE_HOST"
	envDroneToken    = "DRONE_TOKEN"
	envRepoOwner     = "REPO_OWNER"
	envRepoName      = "REPO_NAME"
	envRepoBranch    = "REPO_BRANCH"
	envPrismicSecret = "PRISMIC_SECRET"
	envHttpAddress   = "HTTP_ADDRESS"
	envHttpPort      = "HTTP_PORT"
	envHttpRoute     = "HTTP_ROUTE"
	defaultHttpPort  = "80"
	defaultHttpRoute = "/handle"
)

// config is used to keep all the configuration from env vars in a single object after validation
type config struct {
	drone         drone.Client
	repoOwner     string
	repoName      string
	repoBranch    string
	httpAddress   string
	httpPort      string
	httpRoute     string
	prismicSecret string
}

// initConfig validates the configuration and returns it as an object
func initConfig() *config {
	// Formal parameters validation
	droneHost := os.Getenv(envDroneHost)
	droneToken := os.Getenv(envDroneToken)
	if len(droneHost) == 0 || len(droneToken) == 0 {
		log.Fatal("Drone configuration missing.")
	}

	repoOwner := os.Getenv(envRepoOwner)
	repoName := os.Getenv(envRepoName)
	repoBranch := os.Getenv(envRepoBranch)
	if len(repoOwner) == 0 || len(repoName) == 0 || len(repoBranch) == 0 {
		log.Fatal("Repo configuration missing.")
	}

	prismicSecret := os.Getenv(envPrismicSecret)
	if len(prismicSecret) == 0 {
		log.Fatal("Prismic configuration missing.")
	}

	httpAddress := os.Getenv(envHttpAddress)
	httpPort := os.Getenv(envHttpPort)
	if len(httpPort) == 0 {
		httpPort = defaultHttpPort
	}

	httpRoute := os.Getenv(envHttpRoute)
	if len(httpRoute) == 0 {
		httpRoute = defaultHttpRoute
	}

	// Check the connection to Drone
	drone := drone.NewClient(droneHost, droneToken)

	log.Printf(
		"Connecting to Drone host %s and checking last build for %s/%s on branch %s",
		droneHost,
		repoOwner,
		repoName,
		repoBranch,
	)

	buildNumber, err := drone.GetLastBuildNumber(repoOwner, repoName, repoBranch)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found last build with ID %d", buildNumber)

	return &config{
		*drone,
		repoOwner,
		repoName,
		repoBranch,
		httpAddress,
		httpPort,
		httpRoute,
		prismicSecret,
	}
}

func main() {
	config := initConfig()

	log.Printf("Starting server on port %s", config.httpPort)

	r := gin.Default()
	r.POST(config.httpRoute, func(c *gin.Context) {
		var webhook prismic.Webhook
		if err := c.ShouldBindJSON(&webhook); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !webhook.VerifySecret(config.prismicSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong webhook secret"})
			return
		}

		log.Printf(
			"Triggering new build for %s/%s on branch %s",
			config.repoOwner,
			config.repoName,
			config.repoBranch,
		)

		buildID, err := config.drone.TriggerBuild(
			config.repoOwner,
			config.repoName,
			config.repoBranch,
		)

		if err != nil {
			log.Printf("Error triggering build: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		}
		log.Printf("New build %d started", buildID)

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Started build %d", buildID)})
	})
	r.Run(fmt.Sprintf("%s:%s", config.httpAddress, config.httpPort))
}
