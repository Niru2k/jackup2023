package handlers

import (
	"blog/helper"
	"blog/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// This is for Signup
func Signup(c *fiber.Ctx) error {
	var data models.User
	check := false

	//Get user details from request body
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if data.Email == "" || data.Id <= 0 || data.Password == "" || data.Username == "" || data.Role == "" {
		return c.JSON(fiber.Map{
			"message": "missing credentials",
		})
	}

	for _, val := range models.Database {
		if val.Id == data.Id || val.Email == data.Email {
			check = true
		}
	}
	if !check {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(password)
		models.Database = append(models.Database, data)

		//create a JWT token
		var auth models.Authentication
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(data.Id)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		})

		token, err := claims.SignedString([]byte(helper.SecretKey))
		if err != nil {
			return err
		}
		auth.Id, auth.Token = data.Id, token
		models.Auth = append(models.Auth, auth)
		logrus.Println("Database :", models.Database)
		logrus.Println("Token :", models.Auth)
		return c.JSON(data)
	}

	return c.JSON(fiber.Map{
		"status":  409,
		"message": "user already exist",
	})
}

// This is for Login
func Login(c *fiber.Ctx) error {
	var data models.User
	check := false

	//Get mail-id and password from request body
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	for _, value := range models.Database {
		if data.Email == value.Email {
			check = true
		}
	}
	if !check {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"status":  404,
			"message": "user not found",
		})
	} else {
		for _, value := range models.Database {
			if value.Email == data.Email {
				if err := bcrypt.CompareHashAndPassword([]byte(value.Password), []byte(data.Password)); err != nil {
					c.Status(fiber.StatusBadRequest)
					return c.JSON(fiber.Map{
						"message": "incorrect password",
					})
				} else {
					return c.JSON(fiber.Map{
						"message": "Login Successful!!!",
					})
				}
			}
		}
	}
	return nil
}
