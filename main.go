package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/dgrijalva/jwt-go"
	jwtware "github.com/gofiber/jwt"
)

const jwtSecret = "asecret"

func authRequired() func(ctx *fiber.Ctx) {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) {
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		SigningKey:     []byte(jwtSecret),
	})
}

func main() {
	app := fiber.New()

	app.Use(middleware.Logger())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Hello world")
	})

	app.Post("/login", login)

	app.Get("/hello", authRequired(), func(ctx *fiber.Ctx) {
		user := ctx.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["sub"].(string)
		ctx.Send(fmt.Sprintf("Hello user with id: %s", id))
	})

	err := app.Listen(3000)
	if err != nil {
		panic(err)
	}
}

func login(ctx *fiber.Ctx) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	if body.Email != "bob@gmail.com" || body.Password != "password123" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // a week

	s, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": s,
		"user": struct {
			Id int `json:"id"`
			Email string `json:"email"`
		}{
			Id: 1,
			Email: "bob@gmail.com",
		},
	})
}
