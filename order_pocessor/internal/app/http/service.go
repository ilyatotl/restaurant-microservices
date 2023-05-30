package http_service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"order_pocessor/internal/app/core"
	"order_pocessor/internal/app/custom_errors"
	"order_pocessor/internal/app/dish"
	"strconv"
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
	s.r.POST("/order", s.CreateOrder)
	s.r.GET("/get_order_info", s.GetOrder)
	s.r.GET("/menu", s.ShowMenu)

	s.r.POST("/dish/add", s.AddDish)
	s.r.GET("/dish/get", s.GetDish)
	s.r.PUT("/dish/update", s.UpdateDish)
	s.r.DELETE("/dish/delete", s.DeleteDish)
	return s
}

func (s *Service) StartHTTPService(port string) {
	_ = s.r.Run(port)
}

func (s *Service) CreateOrder(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
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
	dishes := make([]dish.DishRequest, 0)
	if err := json.Unmarshal(body, &dishes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := s.Core.CreateOrder(token, dishes)
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrNotEnoughDishes) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrDishNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order created successfully with id": id,
	})
}

func (s *Service) GetOrder(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
		})
		return
	}

	q, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "specify the id of the order",
		})
		return
	}
	id, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "specify the correct id",
		})
		return
	}

	status, err := s.Core.GetOrder(token, int64(id))
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrOrderNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order status": status,
	})
}

func (s *Service) AddDish(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
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
	var d dish.DishModel
	if err := json.Unmarshal(body, &d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := s.Core.AddDish(token, &d)
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, gin.H{
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
		"added product with id": id,
	})
}

func (s *Service) GetDish(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
		})
		return
	}

	q, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "specify id in query",
		})
		return
	}
	id, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "specify correct id",
		})
		return
	}

	d, err := s.Core.GetDish(token, int64(id))
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, gin.H{
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
	c.JSON(http.StatusOK, d)
}

func (s *Service) UpdateDish(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
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
	var d dish.DishModel
	if err := json.Unmarshal(body, &d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.Core.UpdateDish(token, &d)
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, gin.H{
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
		"dish updated correctly": "ok",
	})
}

func (s *Service) DeleteDish(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
		})
		return
	}

	q, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "specify id in query",
		})
		return
	}
	id, err := strconv.Atoi(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "specify correct id",
		})
		return
	}

	err = s.Core.DeleteDish(token, int64(id))
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, custom_errors.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, gin.H{
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
		"dish was deleted correctly": "ok",
	})
}

func (s *Service) ShowMenu(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authorized",
		})
		return
	}

	menu, err := s.Core.ShowMenu(token)
	if errors.Is(err, custom_errors.ErrUserNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, menu)
}
