package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"gondos/internal/app"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil, jwt.WithAcceptableSkew(30*time.Second))
}

// Register a new account
// (POST /auth/register)
func (si *serverImpl) AuthRegister(w http.ResponseWriter, r *http.Request) {
	var request AuthRegisterRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	user, err := app.NewUser(request.Name, request.Email, request.Password)
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	if err := si.app.CreateUser(r.Context(), user); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login
// (POST /auth/login)
func (si *serverImpl) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var request AuthLoginRequest
	if err := parseJSON(r, &request); err != nil {
		si.deliverErr(w, r, err)
		return
	}

	user, err := si.app.AuthEmail(r.Context(), request.Email, request.Password)
	if err != nil {
		// whatever the problem just send wrong credential
		// or return the real error
		si.deliverErr(w, r, app.ErrCredentialsIncorrect)
		return
	}

	_, tk, err := tokenAuth.Encode(map[string]interface{}{
		jwt.SubjectKey:    strconv.Itoa(int(user.ID())),
		jwt.IssuedAtKey:   time.Now().Unix(),
		jwt.ExpirationKey: time.Now().Add(10 * time.Minute).Unix(),
	})
	if err != nil {
		si.deliverErr(w, r, err)
		return
	}

	sendJSON(w, http.StatusOK, AuthLoginResponse{
		AccessToken: tk,
	})
}

// authenticate checks if request authenticated or not
func authenticate(r *http.Request) error {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return err
	}

	// get subject and put to context
	sub, ok := claims[jwt.SubjectKey].(string)
	if !ok {
		return jwtauth.ErrUnauthorized
	}
	subInt, err := strconv.Atoi(sub)
	if err != nil {
		return err
	}
	*r = *r.WithContext(app.SetUserIDCtx(r.Context(), int64(subInt)))
	return nil
}
