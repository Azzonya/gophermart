package auth

import (
	"errors"
	"fmt"
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
	userService "github.com/Azzonya/gophermart/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

const (
	sessionCookie              = "jwt"
	defaultJWTCookieExpiration = 24 * time.Hour
)

type Claims struct {
	jwt.RegisteredClaims
	UID string
}

type Auth struct {
	userService *userService.UserServiceI
	JwtSecret   string
}

func New(jwtSecret string) *Auth {
	return &Auth{JwtSecret: jwtSecret}
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/api/user/register" || path == "/api/user/login" {
			c.Next()
			return
		}

		authorizer := New(jwtSecret)

		userID, err := authorizer.GetUserIDFromCookie(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to get cookie",
				"error":   err.Error(),
			})
			return
		}
		if len(userID) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func (a *Auth) GetUserIDFromCookie(c *gin.Context) (string, error) {
	userCookie, err := c.Cookie(sessionCookie)
	if err != nil {
		return "", err
	}

	return a.GetUserIDFromJWT(userCookie)
}

func (a *Auth) GetUserIDFromJWT(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.JwtSecret), nil
		})
	if err != nil {
		return "", fmt.Errorf("token is not valid")
	}

	if !token.Valid {
		return "", errors.New("token is not valid")
	}

	if claims, ok := token.Claims.(*Claims); !ok {
		return "", errors.New("token does not contain user id")
	} else {
		return claims.UID, nil
	}
}

func (a *Auth) CreateJWTCookie(u *userModel.User) (*http.Cookie, error) {
	token, err := a.NewToken(u)
	if err != nil {
		return nil, fmt.Errorf("cannot create auth token: %w", err)
	}
	return &http.Cookie{
		Name:  sessionCookie,
		Value: token,
	}, nil
}

func (a *Auth) NewToken(u *userModel.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultJWTCookieExpiration)),
		},
		UID: u.ID,
	})

	signedToken, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("cannot sign jwt token: %w", err)
	}
	return signedToken, nil
}
