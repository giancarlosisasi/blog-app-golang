package security

import (
	"blog-app/internal/config"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type JWTClaims struct {
	UserID   pgtype.UUID `json:"user_id"`
	Email    string      `json:"email"`
	Username string      `json:"username"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	AppName         string
	SecretKey       string
	ExpirationTime  time.Duration
	RefreshDuration time.Duration
	CookieName      string
	CookieDomain    string
	CookiePath      string
	SecureOnly      bool
	SameSite        string
	IsProduction    bool
}

func JWTDefaultConfig(appConf *config.Config) *JWTConfig {
	isProduction := appConf.AppEnv == "production"

	return &JWTConfig{
		AppName:         appConf.AppDomain,
		SecretKey:       appConf.JWTSecret,
		ExpirationTime:  time.Duration(appConf.JWTAccessExpiry) * time.Minute,
		RefreshDuration: time.Duration(appConf.JWTRefreshExpiry) * time.Hour,
		CookieName:      "auth_token",
		CookieDomain:    appConf.AppDomain, // set with our production domain later
		CookiePath:      "/",
		SecureOnly:      isProduction,
		SameSite:        "Strict",
		IsProduction:    isProduction,
	}
}

func GenerateJWT(jwtConfig *JWTConfig, userID pgtype.UUID, email string, username string) (string, error) {
	if jwtConfig.SecretKey == "" {
		return "", errors.New("JWT secret key is required")
	}

	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.AppName,
			Subject:   email,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signin token with secret
	tokenString, err := token.SignedString([]byte(jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetJWTCookie(c *fiber.Ctx, jwtConfig *JWTConfig, tokenString string) {
	cookie := &fiber.Cookie{
		Name:     jwtConfig.CookieName,
		Value:    tokenString,
		Path:     jwtConfig.CookiePath,
		Domain:   jwtConfig.CookieDomain,
		MaxAge:   int(jwtConfig.ExpirationTime.Seconds()),
		Secure:   jwtConfig.IsProduction,
		HTTPOnly: true, // Prevents XSS attacks
		SameSite: jwtConfig.SameSite,
	}

	c.Cookie(cookie)
}
