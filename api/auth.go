package api

import (
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil, jwt.WithAcceptableSkew(30*time.Second))
}
