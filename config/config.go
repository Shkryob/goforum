package config

import (
	"fmt"
	"os"

	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/json"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Auth_Redirect_Url string

	Oauth_Google_Client_Id           string
	Oauth_Google_Client_Secret       string
	Oauth_Google_Client_Redirect_Url string
	Oauth_Google_Open_Id_Config_Url  string
	Oauth_Google_User_Info_Url       string
	Oauth_Google_Scopes              []string

	Oauth_Facebook_Client_Id           string
	Oauth_Facebook_Client_Secret       string
	Oauth_Facebook_Client_Redirect_Url string
	Oauth_Facebook_User_Info_Url       string
	Oauth_Facebook_Scopes              []string

	Oauth_Twitter_Client_Id           string
	Oauth_Twitter_Client_Secret       string
	Oauth_Twitter_Client_Redirect_Url string
	Oauth_Twitter_Client_Bearer_Token string
	Oauth_Twitter_User_Info_Url       string
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
	configuration := Configuration{}

	//Base configuration values common for all envs
	goconfig.File = "./config/config.base.json"

	err = goconfig.Parse(&configuration)
	if err != nil {
		fmt.Println(err)
	}

	//Load env specific configuration on top
	goconfig.File = "./config/config." + env + ".json"

	err = goconfig.Parse(&configuration)
	if err != nil {
		fmt.Println(err)
	}

	return configuration
}
