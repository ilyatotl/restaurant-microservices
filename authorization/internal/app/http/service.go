package http_service

import (
	"authorization/internal/app/core"
	"authorization/internal/app/custom_errors"
	"authorization/internal/app/user"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Service struct {
	r    *gin.Engine
	Core *core.Core
}

func NewHTTPService(r *gin.Engine, core *core.Core) *Service {
	s := &Service{
		r:    r,
		Core: core,
	}
	r.POST("/register/:role", s.RegisterUser)
	r.POST("/authorize", s.AuthorizeUser)
	r.GET("/get_user_info", s.GetUserInfo)
	return s
}

func (s *Service) StartHTTPService(port string) {
	_ = s.r.Run(port)
}

func (s *Service) RegisterUser(c *gin.Context) {
	defer c.Request.Body.Close()

	role := c.Param("role")
	if role != "customer" && role != "chef" && role != "manager" {
		c.JSON(http.StatusNotFound, gin.H{
			"role not found": role,
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var u user.UserModel
	if err := json.Unmarshal(body, &u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	u.Role = role
	id, err := s.Core.RegisterUser(context.Background(), &u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user successfully created with id": id,
	})
}

func (s *Service) AuthorizeUser(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	type data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var d data
	if err := json.Unmarshal(body, &d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := s.Core.AuthorizeUser(context.Background(), d.Email, d.Password)
	if errors.Is(err, custom_errors.ErrEmptyFieldEmail) || errors.Is(err, custom_errors.ErrEmptyFieldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrInvalidEmailOrPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Service) GetUserInfo(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "empty header authorization",
		})
		return
	}

	u, err := s.Core.GetUserInfo(context.Background(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *u)
}
