package token

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminClaims struct {
	jwt.RegisteredClaims
	Email string
	ID    primitive.ObjectID
}

var secretKey = os.Getenv("TRAVAS_KEY")

func Generate(email string, id primitive.ObjectID) (string, string, error) {
	travasClaims := AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "travasAdmin",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * time.Hour),
			},
		},
		Email: email,
		ID:    id,
	}
	refTravasClaims := &jwt.RegisteredClaims{
		Issuer:   "travasAdmin",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(48 * time.Hour),
		},
	}
	travasToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, travasClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	refTravasToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refTravasClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	return travasToken, refTravasToken, nil
}

func Parse(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Fatalf("error while parsing token with it claims %v", err)
	}
	claims, ok := token.Claims.(*AdminClaims)
	if !ok {
		log.Fatalf("error %v controller not authorized access", http.StatusUnauthorized)
	}
	if err := claims.Valid(); err != nil {
		log.Fatalf("error %v %s", http.StatusUnauthorized, err)
	}
	return claims, nil
}
