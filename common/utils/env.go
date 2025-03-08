package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/izquiratops/tango/common/types"
)

func resolveBooleanFromEnv(envName string) bool {
	value := strings.ToLower(os.Getenv(envName))
	return value == "true" || value == "1" || value == "yes"
}

func LoadEnvironmentConfig(mongoDomainMap map[bool]string) (types.ServerConfig, error) {
	jmdictVersion := os.Getenv("TANGO_VERSION")
	if jmdictVersion == "" {
		return types.ServerConfig{}, fmt.Errorf("TANGO_VERSION environment variable must be set")
	}

	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoDomain := mongoDomainMap[resolveBooleanFromEnv("TANGO_LOCAL")]

	mongoURI := ""
	if mongoUser != "" && mongoPassword != "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:27017", mongoUser, mongoPassword, mongoDomain)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:27017", mongoDomain)
	}

	return types.ServerConfig{
		JmdictVersion: jmdictVersion,
		MongoURI:      mongoURI,
	}, nil
}
