package models

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type AppConfig struct {
	Host string
	Port string
}

type AuthConfig struct {
	PrivateKeyPath            string
	PublicKeyPath             string
	TokenExpirySeconds        int
	RefreshTokenExpirySeconds int
	Issuer                    string
	Audience                  string
	Kid                       string
}

type ApplicationConfig struct {
	App  AppConfig
	DB   DBConfig
	Auth AuthConfig
}
