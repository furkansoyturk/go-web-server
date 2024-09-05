package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
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

func CreateJWT(secret []byte, userId string, expiredInSecond int) (string, string, error) {

	ttl := time.Duration(expiredInSecond) * time.Second

	if !(ttl > 0 && ttl < 86000*time.Second) {
		ttl = 86400 * time.Second
	}

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
		return err.Error(), "-", err
	}

	return signedJWT, createRefreshToken(), nil
}

func ReadFrom(token string) (userID string) {
	tkn, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		secretKey := os.Getenv("JWT_SECRET")
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	sbj, _ := tkn.Claims.GetSubject()
	return sbj
}

func createRefreshToken() (refreshToken string) {
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	// The slice should now contain random bytes instead of only zeroes.
	return hex.EncodeToString(b)
}

//TODO: refresh token validation & ttl should be impl.
