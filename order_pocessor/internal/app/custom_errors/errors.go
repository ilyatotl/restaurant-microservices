package custom_errors

import "errors"

var (
	ErrOrderNotFound     = errors.New("order with this id not found")
	ErrDishNotFound      = errors.New("dish with this id not found")
	ErrOrderDishNotFound = errors.New("order-dish with this id not found")
)

var (
	ErrUserNotAuthorized = errors.New("user not authorized")
	ErrNotEnoughDishes   = errors.New("not enough dishes available")
	ErrPermissionDenied  = errors.New("permission denied")
)
