package oath2

import (
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Auth return objects about auth
type Auth struct {
	SigningKey string
	TokenName  string
}

// GetJwtMiddleware for http
func (auth *Auth) GetJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	var middleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return auth.SigningKey, nil
		},
		UserProperty:  "user",
		Debug:         false,
		SigningMethod: jwt.SigningMethodHS256,
		Extractor: jwtmiddleware.FromFirst(auth.FromCookie(auth.TokenName), jwtmiddleware.FromAuthHeader,
			jwtmiddleware.FromParameter("auth_code")),
	})
	return middleware
}

// FromCookie Extractor
func (auth *Auth) FromCookie(param string) jwtmiddleware.TokenExtractor {
	return func(r *http.Request) (string, error) {
		cookie, err := r.Cookie(param)
		return cookie.Value, err
	}
}

// GetToken create a jwt token with user claims
func (auth *Auth) GetToken(user map[string]interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["username"] = user["username"]
	claims["roles"] = user["roles"]
	claims["org"] = user["org"]
	signedToken, _ := token.SignedString(auth.SigningKey)
	return signedToken
}

// GetJSONToken create a JSON token string
// func (auth *Auth) GetJSONToken(user *models.User) string {
// 	token := GetToken(user)
// 	jsontoken := "{\"id_token\": \"" + token + "\"}"
// 	return jsontoken
// }

// GetUserClaimsFromContext return "user" claims as a map from request
func (auth *Auth) GetUserClaimsFromContext(req *http.Request) map[string]interface{} {
	//claims := context.Get(req, "user").(*jwt.Token).Claims.(jwt.MapClaims)
	claims := req.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return claims
}

// GetHashPassword get the password
func (auth *Auth) GetHashPassword(username string, password string) (string, error) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(username+password), bcrypt.DefaultCost)
	return string(hpass), err
}
