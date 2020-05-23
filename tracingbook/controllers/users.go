package controllers

import (
	"errors"
	"gohome_self/tracingbook/dtos"
	"gohome_self/tracingbook/models"
	"gohome_self/tracingbook/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UsersRegistration(c *gin.Context) {

	var json dtos.RegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err := services.CreateOne(&models.User{
		Username:  json.Username,
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
	}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"full_messages": []string{"User created successfully"}})
}

func UsersLogin(c *gin.Context) {

	var json dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	user := models.User{}
	userGoten, err := user.FindOneUser(&models.User{Username: json.Username})

	if err != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login_error", err))
		return
	}

	if userGoten.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login", errors.New("invalid credentials")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateLoginSuccessful(userGoten))

}
