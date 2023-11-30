package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joaoribeirodasilva/mqtt-course/api/configuration"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID
	Name     string
	Surename string
	Admin    bool
}

type Token struct {
	conf        *configuration.Configuration
	User        *User
	token       *jwt.Token
	TokenString string
}

const (
	iss = "api.mqtt-course.io"
	aud = "mqtt-course"
)

func New(conf *configuration.Configuration) *Token {

	t := &Token{}

	t.conf = conf

	return t
}

func (t *Token) Create(user *User) error {

	now := time.Now()
	expires := now.Add(time.Duration(time.Hour * 24 * 30))

	sub := make(map[string]interface{})
	sub["id"] = user.ID.Hex()
	sub["name"] = user.Name
	sub["surename"] = user.Surename
	sub["admin"] = user.Admin

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": iss,
		"sub": sub,
		"aud": aud,
		"iat": now.Unix(),
		"exp": expires.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(t.conf.Server.JwtKey))
	if err != nil {
		return fmt.Errorf("ERROR: [JWT TOKEN] failed to encrypt token")
	}

	t.TokenString = tokenStr

	return nil
}

func (t *Token) IsValid(header string) bool {

	header = strings.TrimSpace(header)
	if header == "" {
		return false
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return false
	}

	if jwtToken[0] != "Bearer" {
		return false
	}

	if err := t.parseToken(jwtToken[1]); err != nil {
		return false
	}

	claims, Ok := t.token.Claims.(jwt.MapClaims)
	if !Ok || claims.Valid() != nil || !claims.VerifyAudience(aud, true) || !claims.VerifyIssuer(iss, true) {
		return false
	}

	defer func() {
		recover()
	}()

	sub := claims["sub"].(map[string]interface{})

	iid, ok := sub["id"]
	if !ok {
		return false
	}
	id := iid.(string)
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}

	iname, ok := sub["name"]
	if !ok {
		return false
	}
	name := iname.(string)

	isurename, ok := sub["surename"]
	if !ok {
		return false
	}
	surename := isurename.(string)

	iadmin, ok := sub["admin"]
	if !ok {
		return false
	}
	admin := iadmin.(bool)

	t.User = &User{
		ID:       mid,
		Name:     name,
		Surename: surename,
		Admin:    admin,
	}

	return true
}

func (t *Token) parseToken(jwtToken string) error {

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, fmt.Errorf("ERROR: [JWT TOKEN] invalid token")
		}
		return []byte(t.conf.Server.JwtKey), nil
	})

	if err != nil {
		return fmt.Errorf("ERROR: [JWT TOKEN] invalid token")
	}

	t.token = token

	return nil
}
