package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user
// @ID sign-up
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userRegisterRequest true "User info for registration"
// @Success 201 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users [post]
func (handler *Handler) SignUp(context echo.Context) error {
	var u model.User
	req := &userRegisterRequest{}
	if err := req.bind(context, &u); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := handler.userStore.Create(&u); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return utils.ResponseByContentType(context, http.StatusCreated, newUserResponse(&u))
}

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userLoginRequest true "Credentials to use"
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users/login [post]
func (handler *Handler) Login(context echo.Context) error {
	req := &userLoginRequest{}
	if err := req.bind(context); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u, err := handler.userStore.GetByEmail(req.User.Email)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return utils.ResponseByContentType(context, http.StatusForbidden, utils.AccessForbidden())
	}
	if !u.CheckPassword(req.User.Password) {
		return utils.ResponseByContentType(context, http.StatusForbidden, utils.AccessForbidden())
	}
	return utils.ResponseByContentType(context, http.StatusOK, newUserResponse(u))
}

func userIDFromToken(context echo.Context) uint {
	id, ok := context.Get("user").(uint)
	if !ok {
		return 0
	}
	return id
}

// CurrentUser godoc
// @Summary Get the current user
// @Description Gets the currently logged-in user
// @ID current-user
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /user [get]
func (handler *Handler) CurrentUser(context echo.Context) error {
	u, err := handler.userStore.GetByID(userIDFromToken(context))
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}
	return utils.ResponseByContentType(context, http.StatusOK, newUserResponse(u))
}

func (handler *Handler) OauthLoginOrSignUp(context echo.Context, email string) error {
	u, _ := handler.userStore.GetByEmail(email)
	if u == nil {
		u = new(model.User)
		u.Email = email
		if err := handler.userStore.Create(u); err != nil {
			return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
		}
	}
	return utils.ResponseByContentType(context, http.StatusOK, newUserResponse(u))
}

func getUserInfoFromGoogle(conf *oauth2.Config, tok *oauth2.Token) map[string]interface{} {
	client := conf.Client(oauth2.NoContext, tok)
	response, _ := client.Get(`https://accounts.google.com/.well-known/openid-configuration`)
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	json := utils.JsonToMap(body)
	response, _ = client.Get(json[`userinfo_endpoint`].(string))
	body, _ = ioutil.ReadAll(response.Body)
	response.Body.Close()

	json = utils.JsonToMap(body)
	fmt.Println(json[`email`].(string))

	return json
}

func (handler *Handler) OauthGoogle(context echo.Context) error {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &oauth2.Config{
		ClientID:     handler.config.Oauth_Google_Client_Id,
		ClientSecret: handler.config.Oauth_Google_Client_Secret,
		RedirectURL:  handler.config.Oauth_Google_Client_Redirect_Url,
		Scopes: []string{
			"https://www.googleapis.com/auth/plus.profile.emails.read",
			"https://www.googleapis.com/auth/plus.login",
			"https://www.googleapis.com/auth/plus.me",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	if context.QueryParam("code") != "" {
		// Handle the exchange code to initiate a transport.
		tok, err := conf.Exchange(oauth2.NoContext, context.QueryParam("code"))
		if err != nil {
			log.Fatal(err)
		}
		json := getUserInfoFromGoogle(conf, tok)
		handler.OauthLoginOrSignUp(context, json[`email`].(string))
		return context.Redirect(http.StatusSeeOther, "/api/posts")
	} else {
		// Redirect user to Google's consent page to ask for permission
		// for the scopes specified above.
		url := conf.AuthCodeURL("state")
		fmt.Printf("Visit the URL for the auth dialog: %v", url)
		return context.Redirect(http.StatusSeeOther, url)
	}

	return utils.ResponseByContentType(context, http.StatusOK, map[string]interface{}{"result": "ok"})
}

func getUserInfoFromFacebook(conf *oauth2.Config, tok *oauth2.Token) map[string]interface{} {
	client := conf.Client(oauth2.NoContext, tok)

	response, _ := client.Get(`https://graph.facebook.com/me?fields=email,name&access_token=` + tok.AccessToken)
	body, _ := ioutil.ReadAll(response.Body)

	json := utils.JsonToMap(body)
	response.Body.Close()

	return json
}

func (handler *Handler) OauthFacebook(context echo.Context) error {
	conf := &oauth2.Config{
		ClientID:     handler.config.Oauth_Facebook_Client_Id,
		ClientSecret: handler.config.Oauth_Facebook_Client_Secret,
		RedirectURL:  handler.config.Oauth_Facebook_Client_Redirect_Url,
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}

	if context.QueryParam("code") != "" {
		// Handle the exchange code to initiate a transport.
		tok, err := conf.Exchange(oauth2.NoContext, context.QueryParam("code"))
		if err != nil {
			log.Fatal(err)
		}
		json := getUserInfoFromFacebook(conf, tok)
		handler.OauthLoginOrSignUp(context, json[`email`].(string))
		return context.Redirect(http.StatusSeeOther, "/api/posts")
	} else {
		// Redirect user to Google's consent page to ask for permission
		// for the scopes specified above.
		url := conf.AuthCodeURL("state")
		fmt.Printf("Visit the URL for the auth dialog: %v", url)
		return context.Redirect(http.StatusSeeOther, url)
	}

	return utils.ResponseByContentType(context, http.StatusOK, map[string]interface{}{"result": "ok"})
}
