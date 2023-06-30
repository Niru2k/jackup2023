package handlers

import (
	//User-defined packages
	"blog/helper"
	"blog/logs"
	"blog/models"

	//Inbuild package
	"time"

	//Third-party packages
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// This is for Signup
func Signup(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info("signup-API called...")
	defer log.Info("signup-API finished")
	var data models.User
	check := false

	//Get user details from request body
	if err := c.BodyParser(&data); err != nil {
		log.Error(err)
		return nil
	}
	if data.Email == "" {
		log.Error("missing Email")
		return c.JSON(fiber.Map{
			"message": "missing Email",
		})
	} else if data.Id <= 0 {
		log.Error("invalid id")
		return c.JSON(fiber.Map{
			"message": "invalid id",
		})
	} else if data.Password == "" {
		log.Error("missing password")
		return c.JSON(fiber.Map{
			"message": "missing password",
		})
	} else if data.Username == "" {
		log.Error("missing username")
		return c.JSON(fiber.Map{
			"message": "missing username",
		})
	} else if data.Role == "" {
		log.Error("missing role")
		return c.JSON(fiber.Map{
			"message": "missing role",
		})
	}

	for _, val := range models.Database {
		if val.Email == data.Email {
			check = true
		}
	}
	if !check {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err)
			return nil
		}
		data.Password = string(password)
		models.Database = append(models.Database, data)
		return c.JSON(fiber.Map{
			"status":    200,
			"message":   "signup successful!!!",
			"user data": data,
		})
	}
	log.Error("user already exist")
	return c.JSON(fiber.Map{
		"status":  409,
		"message": "user already exist",
	})
}

// This is for Login
func Login(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info("login-API called...")
	defer log.Info("login-API finished")
	var data models.User
	var auth models.Authentication
	check := false

	//Get mail-id and password from request body
	if err := c.BodyParser(&data); err != nil {
		log.Error(err)
		return nil
	}
	for _, value := range models.Database {
		if data.Email == value.Email {
			check = true
		}
	}
	if !check {
		log.Error("user not found")
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"status":  404,
			"message": "user not found",
		})
	} else {
		for _, value := range models.Database {
			if value.Email == data.Email {
				if err := bcrypt.CompareHashAndPassword([]byte(value.Password), []byte(data.Password)); err == nil {
					//create a JWT token
					claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
						Issuer:    "JWT",
						ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
					})

					token, err := claims.SignedString([]byte(helper.SecretKey))
					if err != nil {
						log.Error(err)
						return nil
					}
					auth.Id, auth.Token = data.Id, token
					models.Auth = append(models.Auth, auth)
					log.Info("Login Successful!!!")
					return c.JSON(fiber.Map{
						"status":  200,
						"message": "Login Successful!!!",
						"token":   token,
					})
				}

				log.Error("incorrect password")
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"status":  400,
					"message": "incorrect password",
				})
			}
		}
	}
	return nil
}
