package middleware

import (
	//User-defined packages
	"blog/helper"
	"blog/models"

	//Inbuild packages
	"errors"
	"strconv"
	"time"

	//Third-party packages
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateToken(db *gorm.DB, user models.User, c *fiber.Ctx) (string, error) {
	exp := time.Now().Add(time.Hour * 24).Unix()
	userId := strconv.Itoa(int(user.UserId))
	roleId := strconv.Itoa(int(user.RoleId))
	claims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: exp,
		Id:        userId,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "JWT",
		NotBefore: 0,
		Subject:   roleId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(helper.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		for index, char := range tokenString {
			if char == ' ' {
				tokenString = tokenString[index+1:]
			}
		}
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(helper.SecretKey), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid token signature",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Expired token",
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}
		// Check the user's role
		if claims.Subject == "1" {
			c.Locals("role", "admin")
		} else if claims.Subject == "2" {
			c.Locals("role", "user")
		}

		return c.Next()
	}
}

func GetTokenClaims(c *fiber.Ctx) jwt.StandardClaims {
	tokenString := c.Get("Authorization")
	for index, char := range tokenString {
		if char == ' ' {
			tokenString = tokenString[index+1:]
		}
	}
	claims := jwt.StandardClaims{}
	jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(helper.SecretKey), nil
	})
	return claims
}

func AdminAuth(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return errors.New("unauthorized entry")
	}
	return nil
}

func UserAuth(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "user" {
		return errors.New("unauthorized entry")
	}
	return nil
}
