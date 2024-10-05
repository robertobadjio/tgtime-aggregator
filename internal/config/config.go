package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
)

type Config struct {
	HttpPort          string
	GrpcPort          string
	DataBaseHost      string
	DataBasePort      string
	DataBaseName      string
	DataBaseUser      string
	DataBasePassword  string
	DataBaseSslMode   string
	ApiURL            string
	ApiMasterEmail    string
	ApiMasterPassword string
	KafkaHost         string
	KafkaPort         string
}

const projectDirName = "tgtime-aggregator"

func init() {
	loadEnv()
}

func New() *Config {
	return &Config{
		HttpPort:          getEnv("HTTP_PORT", "8081"), // TODO: const
		GrpcPort:          getEnv("GRPC_PORT", "8082"), // TODO: const
		DataBaseHost:      getEnv("DATABASE_HOST", ""),
		DataBasePort:      getEnv("DATABASE_PORT", "5432"), // TODO: const
		DataBaseName:      getEnv("DATABASE_NAME", ""),
		DataBaseUser:      getEnv("DATABASE_USER", ""),
		DataBasePassword:  getEnv("DATABASE_PASSWORD", ""),
		DataBaseSslMode:   getEnv("DATABASE_SSL_MODE", ""),
		ApiURL:            getEnv("API_URL", ""),
		ApiMasterEmail:    getEnv("API_MASTER_EMAIL", ""),
		ApiMasterPassword: getEnv("API_MASTER_PASSWORD", ""),
		KafkaHost:         getEnv("KAFKA_HOST", ""),
		KafkaPort:         getEnv("KAFKA_PORT", ""),
	}
}

func loadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Problem loading .env file")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
