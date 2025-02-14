package controllers

import (
	db "go-fiber-auth-jwt/database"
	h "go-fiber-auth-jwt/helpers"
	v "go-fiber-auth-jwt/helpers/validations"
	m "go-fiber-auth-jwt/models"
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	dbConn := db.DB

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Receive user's information
	name := strings.TrimSpace(data["name"])
	email := strings.TrimSpace(data["email"])
	username := strings.TrimSpace(data["username"])
	password := strings.TrimSpace(data["password"])
	confirmPassword := strings.TrimSpace(data["confirmPassword"])

	var checkUser m.User

	// Validate

	if err := v.ValidateRegisterName(name); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	dbConn.Where("email = ?", email).First(&checkUser)
	if err := v.ValidateRegisterEmail(email, checkUser.Email); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	dbConn.Where("username = ?", username).First(&checkUser)
	if err := v.ValidateRegisterUsername(username, checkUser.Username); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	if err := v.ValidateRegisterPassword(password); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	if err := v.ValidateRegisterConfirmPassword(password, confirmPassword); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	// Encrypt password
	encryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	// Create user

	user := m.User{
		Name:     name,
		Email:    email,
		Username: username,
		Password: encryptPassword,
	}

	db.DB.Create(&user)

	return h.ResultJSON(c, fiber.StatusOK, false, "User has successfully registered", user)
}

func Login(c *fiber.Ctx) error {
	dbConn := db.DB

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Receive user's information
	username := strings.TrimSpace(data["username"])
	password := strings.TrimSpace(data["password"])

	// Validate

	if err := v.ValidateLoginUsername(username); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	if err := v.ValidateLoginPassword(password); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, err.Error(), "")
	}

	// Check data

	var user m.User

	dbConn.Where("username = ?", data["username"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return h.ResultJSON(c, fiber.StatusNotFound, true, "User not found", "")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return h.ResultJSON(c, fiber.StatusBadRequest, true, "Incorrect password", "")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_email": user.Email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return h.ResultJSON(c, fiber.StatusInternalServerError, true, "Could not login", "")
	}

	// Set cookie

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	datas := map[string]interface{}{}

	datas["id"] = user.Id
	datas["name"] = user.Name
	datas["email"] = user.Email
	datas["username"] = user.Username
	datas["token"] = token

	return h.ResultJSON(c, fiber.StatusOK, false, "User has successfully logged in", datas)
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return h.ResultJSON(c, fiber.StatusUnauthorized, true, "Unauthenticated", "")
	}

	claims := token.Claims.(jwt.MapClaims)

	var user m.User

	db.DB.Where("email = ?", claims["user_email"]).First(&user)

	return h.ResultJSON(c, fiber.StatusOK, false, "User fetched successfully", user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return h.ResultJSON(c, fiber.StatusOK, false, "User has successfully logged out", "")
}
