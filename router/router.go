package router

import (
	//User-defined packages
	"blog/handlers"
	"blog/logs"
	"blog/middleware"

	//Third-party package
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Router(Db *gorm.DB) {
	log := logs.Log()
	control := handlers.Database{Db: Db}
	app := fiber.New()

	//Public
	app.Post("/signup", control.Signup)
	app.Post("/login", control.Login)
	app.Get("/getPoster/:post_id", control.GetPosterById)

	//Only for user
	app.Post("/user/addComment", middleware.AuthMiddleware(), control.AddComment)

	//Only for admin
	app.Post("/admin/postPoster", middleware.AuthMiddleware(), control.PostPoster)
	app.Get("/admin/getPosters", middleware.AuthMiddleware(), control.GetPosters)
	app.Put("/admin/updatePoster/:post_id", middleware.AuthMiddleware(), control.UpdatePosterById)
	app.Delete("/admin/deletePoster/:post_id", middleware.AuthMiddleware(), control.DeletePosterById)
	app.Get("/admin/getComments/:post_id", middleware.AuthMiddleware(), control.GetCommentByPostId)
	app.Delete("/admin/deleteComment/:comment_id", middleware.AuthMiddleware(), control.DeleteCommentById)

	//start a server
	log.Info("Server starts in port 8000.....")
	app.Listen(":8000")
}
