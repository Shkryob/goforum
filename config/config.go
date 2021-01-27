package config

import (
	"fmt"
	"os"

	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/json"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Oauth_Google_Client_Id           string
	Oauth_Google_Client_Secret       string
	Oauth_Google_Client_Redirect_Url string

	Oauth_Facebook_Client_Id           string
	Oauth_Facebook_Client_Secret       string
	Oauth_Facebook_Client_Redirect_Url string

	Oauth_Twitter_Client_Id           string
	Oauth_Twitter_Client_Secret       string
	Oauth_Twitter_Client_Redirect_Url string
	Oauth_Twitter_Client_Bearer_Token string
}

func ReadConfig() Configuration {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "development"
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	goconfig.File = "./config/config." + env + ".json"
	configuration := Configuration{}

	err = goconfig.Parse(&configuration)
	if err != nil {
		fmt.Println(err)
	}

	return configuration
}
