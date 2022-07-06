package handlers

import (
	"net/http"

	"github.com/darolpz/students/internal/auth"
	"github.com/darolpz/students/internal/model"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Login
// @Description  get authorization token
// @Tags         auth
// @Param        student body model.Authentication true "Authentication"
// @Accept       json
// @Success      200 {string} token
// @Failure      400 {string} string
// @Failure      401 {string} string
// @Failure      500 {string} string
// @Router       /auth/login [post]
func Login(repo repository.IUsersRepository, authService auth.IAuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var auth model.Authentication
		// Bind the JSON request body to the auth struct.
		if err := c.BindJSON(&auth); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		// Retrieve user from repository
		user, err := repo.FindUserByEmail(auth.Email)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Check if password is correct
		if isPasswordValid := authService.CheckPasswordHash(auth.Password, user.Password); !isPasswordValid {
			c.String(http.StatusUnauthorized, "invalid password")
			return
		}

		// Generate JWT token
		token, err := authService.GenerateJWT(user.Email, user.Role)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, token)
	}
}

// Register godoc
// @Summary      Register user
// @Description  create new user
// @Tags         auth
// @Param        student body model.User true "user"
// @Accept       json
// @Success      200 {object} model.Student
// @Failure      400 {string} string
// @Failure      500 {string} string
// @Router       /auth/register [post]
func Register(repo repository.IUsersRepository, authService auth.IAuthService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var newUser model.User
		// Bind the JSON request body to the newUser struct.
		if err := c.BindJSON(&newUser); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		// Hash password
		password, err := authService.GenerateHashPassword(newUser.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Override password with hashed password
		newUser.Password = password

		// Persist user to repository
		user, err := repo.CreateUser(newUser)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
