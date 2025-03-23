package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var users = []Credentials{
	{Username: "bob", Password: "123"},
	{Username: "dog", Password: "223"},
	{Username: "gob", Password: "113"},
}

var admins = []Credentials{
	{Username: "admin1", Password: "1234567"},
	{Username: "admin2", Password: "000000"},
}

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func generateToken(username string, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func login(c *gin.Context) {
	var creds Credentials
	var flag bool = false
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "in...."})
	}

	for _, name := range users {
		if creds.Username == name.Username && creds.Password == name.Password {
			flag = true
			role := "user"
			token, err := generateToken(creds.Username, role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
			return

		}

	}

	for _, name := range admins {
		if creds.Username == name.Username && creds.Password == name.Password {
			flag = true
			role := "admin"
			token, err := generateToken(creds.Username, role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
			return

		}

	}
	if !flag {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

}

func authMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		if claims.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
