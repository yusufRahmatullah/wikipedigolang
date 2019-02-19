package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	defaultGinMode  = "debug"
	defaultJobLimit = 10
	defaultPort     = "3000"
	defaultWaitTime = 5
)

var (
	config *Instance
	once   sync.Once
)

// Instance is singleton Configuration instance
// Call GetInstance to get instantiated configuration instance
type Instance struct {
	GinMode        string
	JobLimit       int
	MongoDBURL     string
	Port           string
	SessionSecret  string
	SuperAdminPass string
	WaitTime       int
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Environment variable %v not found", key))
	}
	return val
}

func getEnvOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Environment variable %v not found, use def %v", key, def)
		val = def
	}
	return val
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	ret, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Environment variable %v not found, use def %v", key, def)
		ret = def
	}
	return ret
}

func fromEnv() *Instance {
	ginMode := getEnvOrDefault("GIN_MODE", defaultGinMode)
	mongodbURL := getEnv("MONGODB_URL")
	sapass := getEnv("SUPER_ADMIN")
	port := getEnvOrDefault("PORT", defaultPort)
	jobLimit := getEnvInt("JOB_LIMIT", defaultJobLimit)
	waitTime := getEnvInt("WAIT_TIME", defaultWaitTime)
	sess := getEnv("SESSION_SECRET")
	return &Instance{
		GinMode:        ginMode,
		JobLimit:       jobLimit,
		MongoDBURL:     mongodbURL,
		Port:           port,
		SessionSecret:  sess,
		SuperAdminPass: sapass,
		WaitTime:       waitTime,
	}
}

// GetInstance return instantiated configuration Instance
func GetInstance() *Instance {
	once.Do(func() {
		config = fromEnv()
	})
	return config
}
