package types

type ServerConfig struct {
	IsLocalEnvironment bool
	ShouldRebuild      bool
	JmdictVersion      string
	MongoURI           string
}
