package config

// EnvironmentType represents the type of environment the application is running in
type EnvironmentType string

const (
	LocalEnv  EnvironmentType = "local"
	ServerEnv EnvironmentType = "server"
)
