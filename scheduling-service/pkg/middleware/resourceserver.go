package middleware

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/web"
	"net/http"
	"os"
	"strings"
	"time"
)

type Claims struct {
	RealmAccess roles  `json:"realm_access,omitempty"`
	JTI         string `json:"jti,omitempty"`
}

type roles struct {
	Roles []string `json:"roles,omitempty"`
}

var RealmConfigURL = os.Getenv("REALM_CONFIG_URL")
var clientID = os.Getenv("CLIENT_ID")
var authorizedRole = "ADMIN"

func IsAuthorizedJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawAccessToken := strings.Replace(c.GetHeader("Authorization"), "Bearer", "", 1)
		fmt.Println(rawAccessToken)
		trans := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: trans,
		}

		ctx := oidc.ClientContext(context.Background(), client)
		provider, err := oidc.NewProvider(ctx, RealmConfigURL)
		if err != nil {
			authorizationFailed("an authorization error occurred while getting the provider: "+err.Error(), c)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: clientID,
		}

		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken)
		if err != nil {
			authorizationFailed("an authorization error occurred while verifying the token: "+err.Error(), c)
			return
		}

		var IDTokenClaims Claims
		if err := idToken.Claims(&IDTokenClaims); err != nil {
			authorizationFailed("An error occurred while extracting claims: "+err.Error(), c)
			return
		}

		userAccessRoles := IDTokenClaims.RealmAccess.Roles
		for _, userRole := range userAccessRoles {
			if userRole == authorizedRole {
				c.Next()
				return
			}
		}

		authorizationFailed("The user has no permission to access this API", c)
	}

}

func authorizationFailed(message string, c *gin.Context) {
	web.BadResponse(c, http.StatusUnauthorized, "ERROR", message)
	return
}
