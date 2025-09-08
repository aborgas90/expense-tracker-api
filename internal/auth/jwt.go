package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"log"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

type JWTCLAIM struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (accessToken string, refreshToken string, err error) {
	// Access Token
	log.Println("➡️ GenerateToken dipanggil untuk user:", userID)
	accessClaims := &JWTCLAIM{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(accessSecret)
	if err != nil {
		log.Println("❌ gagal buat access token:", err)
		return "", "", err
	}

	log.Println("✅ Access Token berhasil dibuat")
	// Refresh Token
	refreshClaim := &JWTCLAIM{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshToken, err = rt.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken
func ValidateAccessToken(tokenStr string) (*JWTCLAIM, error) {
	return validateToken(tokenStr, accessSecret)
}

// ValidateRefreshToken
func ValidateRefreshToken(tokenStr string) (*JWTCLAIM, error) {
	return validateToken(tokenStr, refreshSecret)
}

func validateToken(tokenStr string, secret []byte) (*JWTCLAIM, error) {
	claims := &JWTCLAIM{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}


// func CSRFToken() string {

// }