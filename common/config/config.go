package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/izquiratops/tango/common/types"
)

const defaultMongoPort = "27017"

func LoadEnvironment(envMap map[EnvironmentType]string) (types.ServerConfig, error) {
	jmdictVersion := os.Getenv("TANGO_VERSION")
	if jmdictVersion == "" {
		return types.ServerConfig{}, fmt.Errorf("TANGO_VERSION environment variable must be set")
	}

	var mongoRunsLocal bool
	var mongoDomain string
	if s := strings.ToLower(os.Getenv("TANGO_MONGO_RUNS_LOCAL")); s == "true" {
		mongoRunsLocal = true
		mongoDomain = envMap[LocalEnv]
	} else {
		mongoRunsLocal = false
		mongoDomain = envMap[ServerEnv]
	}

	mongoURI := ""
	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if mongoUser != "" && mongoPassword != "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%n", mongoUser, mongoPassword, mongoDomain, defaultMongoPort)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%n", mongoDomain, defaultMongoPort)
	}

	return types.ServerConfig{
		JmdictVersion:  jmdictVersion,
		MongoURI:       mongoURI,
		MongoRunsLocal: mongoRunsLocal,
	}, nil
}
