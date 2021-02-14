package config

import (
	"github.com/labstack/gommon/log"
	"os"

	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/json"
	"github.com/joho/godotenv"
)

type Configuration struct {
	AuthRedirectUrl string

	OauthGoogleClientId          string
	OauthGoogleClientSecret      string
	OauthGoogleClientRedirectUrl string
	OauthGoogleOpenIdConfigUrl   string
	OauthGoogleUserInfoUrl       string
	OauthGoogleScopes            []string

	OauthFacebookClientId          string
	OauthFacebookClientSecret      string
	OauthFacebookClientRedirectUrl string
	OauthFacebookUserInfoUrl       string
	OauthFacebookScopes            []string

	OauthTwitterClientId          string
	OauthTwitterClientSecret      string
	OauthTwitterClientRedirectUrl string
	OauthTwitterClientBearerToken string
	OauthTwitterUserInfoUrl       string
}

func ReadConfig() Configuration {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "development"
	}

	err := godotenv.Load()
	if err != nil {
		log.Info("Error loading .env file")
	}
	configuration := Configuration{}

	//Base configuration values common for all envs
	goconfig.File = "./config/config.base.json"

	err = goconfig.Parse(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	//Load env specific configuration on top
	goconfig.File = "./config/config." + env + ".json"

	err = goconfig.Parse(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	return configuration
}
