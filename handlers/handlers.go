package handlers

import (

	//User-defined packages

	"blog/logs"
	"blog/middleware"
	"blog/models"
	"blog/repository"

	//Inbuild packages
	"fmt"
	"reflect"
	"strconv"

	//Third-party packages
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

// This is for Signup
func (db Database) Signup(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info("signup-API called...")
	defer log.Info("signup-API finished")
	var (
		data models.User
		comp models.User
		role models.Roles
	)
	//Get user details from request body
	if err := c.BodyParser(&data); err != nil {
		log.Error(err)
		return nil
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.SignupReq{})
	for _, field := range fields {
		if reflect.ValueOf(&data).Elem().FieldByName(field).Interface() == reflect.ValueOf(&comp).Elem().FieldByName(field).Interface() {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error(stmt)
			return c.JSON(fiber.Map{
				"message": stmt,
			})
		}
	}

	//To check if the user details already exist or not
	data, err := repository.ReadUserByUserId(db.Db, data)
	if err == nil {
		log.Error("user already exist")
		return c.JSON(fiber.Map{
			"status":  409,
			"message": "user already exist",
		})
	}

	//To change the password into hashedPassword
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil
	}
	data.Password = string(password)

	//Select a role_id for specified role
	role, _ = repository.ReadRoleIdByRole(db.Db, data)
	data.RoleId = role.RoleId

	//Adding a user details into our database
	if err = repository.CreateUser(db.Db, data); err != nil {
		log.Errorf("Error :%s", err)
		return c.JSON(fiber.Map{
			"status":  409,
			"message": "user already exist",
		})
	}
	log.Info("signup successful!!!")
	return c.JSON(fiber.Map{
		"status":    200,
		"message":   "signup successful!!!",
		"user data": data,
	})
}

// This is for Login
func (db Database) Login(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info("login-API called...")
	defer log.Info("login-API finished")
	var (
		data models.User
		auth models.Authentication
		user models.User
	)

	//Get mail-id and password from request body
	if err := c.BodyParser(&data); err != nil {
		log.Error(err)
		return nil
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.LoginReq{})
	for _, field := range fields {
		if reflect.ValueOf(&data).Elem().FieldByName(field).Interface() == reflect.ValueOf(&user).Elem().FieldByName(field).Interface() {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error(stmt)
			return c.JSON(fiber.Map{
				"message": stmt,
			})
		}
	}

	//To verify if the user email is exist or not
	user, err := repository.ReadUserByEmail(db.Db, data)
	if err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err == nil {
			// Fetch a JWT token
			if err := repository.ReadTokenByUserId(db.Db, user); err == nil {
				c.Response().Header.Add("Authorization", auth.Token)
				log.Info("Login Successful!!!")
				return c.JSON(fiber.Map{
					"status":  200,
					"message": "Login Successful!!!",
					"token":   auth.Token,
				})
			}

			//Create a token
			token, err := middleware.CreateToken(db.Db, user, c)
			if err != nil {
				return err
			}
			auth.UserId, auth.Token = user.UserId, token
			if err = repository.AddToken(db.Db, auth); err != nil {
				log.Errorf("Error :%s", err)
				return c.JSON(fiber.Map{
					"status":  409,
					"message": err,
				})
			}
			c.Response().Header.Add("Authorization", token)
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
	log.Error("user not found")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "user not found",
	})
}

// Handler for post a poster
func (db Database) PostPoster(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("poster-API called...")
	defer log.Info("poster-API finished")
	var (
		Post     models.Post
		Catagory models.Catagory
	)
	if err := c.BodyParser(&Post); err != nil {
		log.Error(err)
		return nil
	}

	//To check if any credential is missing or not
	comp := models.Post{}
	fields := structs.Names(&models.PostReq{})
	for _, field := range fields {
		if reflect.ValueOf(&Post).Elem().FieldByName(field).Interface() == reflect.ValueOf(&comp).Elem().FieldByName(field).Interface() {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error(stmt)
			return c.JSON(fiber.Map{
				"message": stmt,
			})
		}
	}
	claims := middleware.GetTokenClaims(c)
	userId, _ := strconv.Atoi(claims.Id)
	Post.UserId = uint(userId)
	Catagory, err := repository.ReadCatagoryIdByCatagory(db.Db, Post)
	if err != nil {
		log.Errorf("Invalid catagory")
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Invalid catagory",
		})
	}
	Post.CatagoryId = Catagory.CatagoryId
	if err = repository.CreatePost(db.Db, Post); err != nil {
		log.Errorf("Error :%s", err)
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Can't add a post",
		})
	}
	log.Info("Post added successfully")
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Post added successfully",
	})
}

// Handler for get all posters
func (db Database) GetPosters(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("Getposters-API called...")
	defer log.Info("Getposters-API finished")
	claims := middleware.GetTokenClaims(c)
	Posts, err := repository.ReadPostersByUserId(db.Db, claims.Id)
	if err == nil {
		log.Info("posts retrieved Successfully!!!")
		return c.JSON(fiber.Map{
			"status": 200,
			"Posts":  Posts,
		})
	}
	log.Error("post not found")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Post not found",
	})
}

// Handler for get a poster by post-id
func (db Database) GetPosterById(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info("Getposter-API called...")
	defer log.Info("Getposter-API finished")
	Post, err := repository.ReadPostByPostId(db.Db, c.Params("post_id", ""))
	if err == nil {
		log.Info("post retrieved Successfully!!!")
		return c.JSON(fiber.Map{
			"status":    200,
			"Post data": Post,
		})
	}
	log.Error("post not found")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Post not found",
	})
}

// Handler for update a poster by post-id
func (db Database) UpdatePosterById(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Info("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("Updateposter-API called...")
	defer log.Info("Updateposter-API finished")
	Post, err := repository.ReadPostByPostId(db.Db, c.Params("post_id", ""))
	if err == nil {
		if err := c.BodyParser(&Post); err != nil {
			log.Error(err)
			return err
		}
		comp := models.Post{}
		if comp == Post {
			log.Error("no data found")
			return c.JSON(fiber.Map{
				"status":  404,
				"message": "no data found",
			})
		}
		Catagory, err := repository.ReadCatagoryIdByCatagory(db.Db, Post)
		if err != nil {
			log.Errorf("Invalid catagory")
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid catagory",
			})
		}
		Post.CatagoryId = Catagory.CatagoryId
		if err := repository.UpdatePostByPostId(db.Db, c.Params("post_id", ""), Post); err == nil {
			log.Info("post updated Successfully!!!")
			return c.JSON(fiber.Map{
				"status":  200,
				"message": "post updated Successfully!!!",
			})
		}
	}
	log.Error("post not found")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Post not found",
	})
}

// Handler for delete a poster by post-id
func (db Database) DeletePosterById(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("Deleteposter-API called...")
	defer log.Info("Deleteposter-API finished")
	Post, err := repository.ReadPostByPostId(db.Db, c.Params("post_id", ""))
	if err == nil {
		repository.DeletePostByPostId(db.Db, c.Params("post_id", ""), Post)
		log.Info("post deleted Successfully!!!")
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "post deleted Successfully!!!",
		})

	}
	log.Error("post not found")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Post not found",
	})
}

// Handler for add comment to a post
func (db Database) AddComment(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.UserAuth(c); err != nil {
		log.Error("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("AddComment-API called...")
	defer log.Info("AddComment-API finished")
	var commentData models.Comments
	if err := c.BodyParser(&commentData); err != nil {
		log.Error(err)
		return nil
	}
	models.CommentTable = append(models.CommentTable, commentData)
	log.Info("comment added successfully")
	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "comment added successfully",
	})
}

// Handler for get a comment by post-id
func (db Database) GetCommentByPostId(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("GetCommentById-API called...")
	defer log.Info("GetCommentById-API finished")
	post_id, _ := strconv.Atoi(c.Params("post_id", ""))
	var (
		commentData []models.Comments
		check       bool
	)
	for _, value := range models.CommentTable {
		if value.PostId == uint(post_id) {
			check = true
			commentData = append(commentData, value)
		}
	}
	if check {
		log.Info("Comment(s) retrieved Successfully!!!")
		return c.JSON(fiber.Map{
			"status":    200,
			"Post data": commentData,
		})
	}
	log.Info("Comment not found for this post")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Comment not found for this post",
	})
}

func (db Database) DeleteCommentById(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error("unauthorized entry")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized entry",
		})
	}
	log.Info("DeleteCommentById-API called...")
	defer log.Info("DeleteCommentById-API finished")
	comment_id, _ := strconv.Atoi(c.Params("comment_id", ""))
	for index, value := range models.CommentTable {
		if value.CommentId == uint(comment_id) {
			models.CommentTable = append(models.CommentTable[:index], models.CommentTable[index+1:]...)
			log.Info("Comment deleted Successfully!!!")
			return c.JSON(fiber.Map{
				"status":  200,
				"message": "Comment deleted Successfully!!!",
			})
		}
	}
	log.Info("Comment not found")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Comment not found",
	})
}
