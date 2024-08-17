package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateJWT(secret []byte, userId string, expiredInSecond int) (string, error) {

	ttl := time.Duration(expiredInSecond) * time.Second

	if !(ttl > 0 && ttl < 86000*time.Second) {
		ttl = 86400 * time.Second
	}

	fmt.Printf("input ttl %v\n", ttl)

	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ttl)),
			Subject:   userId,
		},
	)
	signedJWT, err := JWTToken.SignedString(secret)
	if err != nil {
		return err.Error(), err
	}

	return signedJWT, nil
}
