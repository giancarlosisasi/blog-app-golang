package security

import (
	"blog-app/internal/config"
	"blog-app/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTRefreshClaims struct {
	UserID string `json:"user_id"`
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

func GenerateJWT(jwtConfig *JWTConfig, userID string, email string, username string) (jwtTokenString string, jwtRefreshTokenString string, err error) {
	if jwtConfig.SecretKey == "" {
		return "", "", errors.New("JWT secret key is required")
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
		return "", "", err
	}

	// refresh token
	refreshClaims := JWTRefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.RefreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.AppName,
			Subject:   email,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtConfig.SecretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func SetJWTCookie(c *fiber.Ctx, jwtConfig *JWTConfig, tokenString string, refreshTokenString string) {
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

	refreshCookie := &fiber.Cookie{
		// TODO: create a new const for the refresh cookie name
		Name:     fmt.Sprintf("%s_refresh", jwtConfig.CookieName),
		Value:    refreshTokenString,
		Path:     jwtConfig.CookiePath,
		Domain:   jwtConfig.CookieDomain,
		MaxAge:   int(jwtConfig.RefreshDuration.Seconds()),
		Secure:   jwtConfig.IsProduction,
		HTTPOnly: true,
		SameSite: jwtConfig.SameSite,
	}

	c.Cookie(cookie)
	c.Cookie(refreshCookie)
}

func ValidateJWT(tokenString string, config *JWTConfig) (*JWTClaims, error) {
	token, err := validateJwtToken(tokenString, config, &JWTClaims{})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			log.Error().Err(err).Msg("jwt token claims: token is expired")
			return nil, utils.ErrAuthTokenExpired
		}
		return claims, nil
	}

	return nil, utils.ErrAuthTokenInvalid
}

func ValidateRefreshJWT(tokenString string, config *JWTConfig) (*JWTRefreshClaims, error) {
	token, err := validateJwtToken(tokenString, config, &JWTRefreshClaims{})

	if claims, ok := token.Claims.(*JWTRefreshClaims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			log.Error().Err(err).Msg("jwt token claims: token is expired")
			return nil, utils.ErrAuthRefreshTokenExpired
		}
		return claims, nil
	}

	return nil, utils.ErrAuthRefreshTokenInvalid
}

func validateJwtToken(tokenString string, config *JWTConfig, claims jwt.Claims) (*jwt.Token, error) {
	if config.SecretKey == "" {
		return nil, errors.New("JWT secret key is required")
	}

	if tokenString == "" {
		return nil, utils.ErrAuthTokenMissing
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(config.SecretKey), nil
	})

	if err != nil {
		log.Error().Err(err).Msg("error to parse with claims")
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Error().Err(err).Msg("jwt token parseWithClaims: token is expired")
			return nil, utils.ErrAuthTokenExpired
		}

		if errors.Is(err, jwt.ErrTokenMalformed) ||
			errors.Is(err, jwt.ErrTokenSignatureInvalid) ||
			errors.Is(err, jwt.ErrTokenInvalidClaims) {
			log.Error().Err(err).Msg("jwt token parseWithClaims: token is invalid (malformed or invalid signature or invalid claims)")
			return nil, utils.ErrAuthTokenInvalid
		}

		return nil, utils.ErrAuthTokenInvalid
	}

	return token, nil
}

func GetJWTFromCookie(c *fiber.Ctx, config *JWTConfig) (string, error) {
	token := c.Cookies(config.CookieName)
	if token == "" {
		return "", utils.ErrAuthTokenMissing
	}

	return token, nil
}

func GetRefreshJWTFromCookie(c *fiber.Ctx, config *JWTConfig) (string, error) {
	// TODO: create a new const for the refresh cookie name
	refreshToken := c.Cookies(fmt.Sprintf("%s_refresh", config.CookieName))
	if refreshToken == "" {
		return "", utils.ErrAuthRefreshTokenMissing
	}

	return refreshToken, nil
}

// for logout
func ClearJWTCookie(c *fiber.Ctx, config *JWTConfig) {
	cookie := &fiber.Cookie{
		Name:     config.CookieName,
		Value:    "",
		Path:     config.CookieDomain,
		MaxAge:   -1, // expire immediately
		Secure:   config.SecureOnly,
		HTTPOnly: true,
		SameSite: config.SameSite,
	}

	c.Cookie(cookie)
}
