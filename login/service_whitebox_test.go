package login

import (
	"fmt"
	"testing"

	"github.com/almighty/almighty-core/account"
	"github.com/almighty/almighty-core/configuration"
	"github.com/almighty/almighty-core/resource"
	"github.com/almighty/almighty-core/token"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"golang.org/x/oauth2/github"
)

var loginService *gitHubOAuth

func setup() {

	var err error
	if err = configuration.Setup(""); err != nil {
		panic(fmt.Errorf("Failed to setup the configuration: %s", err.Error()))
	}

	oauth := &oauth2.Config{
		ClientID:     configuration.GetGithubClientID(),
		ClientSecret: configuration.GetGithubSecret(),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	publicKey, err := token.ParsePublicKey([]byte(configuration.GetTokenPublicKey()))
	if err != nil {
		panic(err)
	}

	privateKey, err := token.ParsePrivateKey([]byte(configuration.GetTokenPrivateKey()))
	if err != nil {
		panic(err)
	}

	tokenManager := token.NewManager(publicKey, privateKey)
	userRepository := account.NewUserRepository(nil)
	identityRepository := account.NewIdentityRepository(nil)
	loginService = &gitHubOAuth{
		config:       oauth,
		identities:   identityRepository,
		users:        userRepository,
		tokenManager: tokenManager,
	}
}

func tearDown() {
	loginService = nil
}

func TestValidOAuthAccessToken(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	setup()
	defer tearDown()

	accessToken := &oauth2.Token{
		AccessToken: configuration.GetGithubAuthToken(),
		TokenType:   "Bearer",
	}

	emails, err := loginService.getUserEmails(context.Background(), accessToken)
	assert.Nil(t, err)
	assert.NotEmpty(t, emails)
}

func TestInvalidOAuthAccessToken(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	setup()
	defer tearDown()

	if loginService == nil {
		setup()
	}

	invalidAccessToken := "7423742yuuiy-INVALID-73842342389h"

	accessToken := &oauth2.Token{
		AccessToken: invalidAccessToken,
		TokenType:   "Bearer",
	}

	emails, err := loginService.getUserEmails(context.Background(), accessToken)
	assert.Nil(t, err)
	assert.Empty(t, emails)
}

func TestGetUserEmails(t *testing.T) {
	resource.Require(t, resource.UnitTest)

	setup()
	defer tearDown()

	accessToken := &oauth2.Token{
		AccessToken: configuration.GetGithubAuthToken(),
		TokenType:   "Bearer",
	}

	githubEmails, err := loginService.getUserEmails(context.Background(), accessToken)
	t.Log(githubEmails)
	assert.Nil(t, err)
	assert.NotNil(t, githubEmails)
	assert.NotEmpty(t, githubEmails)
}

func TestGetUser(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	setup()
	defer tearDown()

	accessToken := &oauth2.Token{
		AccessToken: configuration.GetGithubAuthToken(),
		TokenType:   "Bearer",
	}

	githubUser, err := loginService.getUser(context.Background(), accessToken)
	assert.Nil(t, err)
	assert.NotNil(t, githubUser)
	t.Log(githubUser)
}

func TestFilterPrimaryEmail(t *testing.T) {
	resource.Require(t, resource.UnitTest)

	t.Skip("Not implemented")
}