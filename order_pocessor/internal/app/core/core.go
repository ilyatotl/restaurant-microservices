package core

import (
	"context"
	"order_pocessor/internal/app/authentication"
	"order_pocessor/internal/app/custom_errors"
	"order_pocessor/internal/app/dish"
	"order_pocessor/internal/app/order"
	"order_pocessor/internal/app/order_dish"
	"sync"
)

const (
	Pending    = "Pending"
	InProgress = "In progress"
	Completed  = "Completed"
)

type Core struct {
	sync.Mutex

	OrderRepo     *order.OrderRepository
	DishRepo      *dish.DishRepository
	OrderDishRepo *order_dish.OrderDishRepository
	AuthClient    *authentication.Client
}

func NewCore(orderRepo *order.OrderRepository, dishRepo *dish.DishRepository,
	odRepo *order_dish.OrderDishRepository, c *authentication.Client) *Core {
	return &Core{
		OrderRepo:     orderRepo,
		DishRepo:      dishRepo,
		OrderDishRepo: odRepo,
		AuthClient:    c,
	}
}

func (c *Core) CreateOrder(token string, dishes []dish.DishRequest) (int64, error) {
	u, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return 0, custom_errors.ErrUserNotAuthorized
	}

	c.Lock()
	defer c.Unlock()

	dishList, _ := c.DishRepo.List(context.Background())
	for _, d := range dishes {
		found := false
		for _, dishElem := range dishList {
			if dishElem.ID == d.ID {
				found = true
				if dishElem.Quantity < d.Quantity {
					return 0, custom_errors.ErrNotEnoughDishes
				}
				break
			}
		}
		if !found {
			return 0, custom_errors.ErrDishNotFound
		}
	}

	id, err := c.OrderRepo.Add(context.Background(), &order.OrderModel{UserID: u.Id, Status: Pending})
	if err != nil {
		return 0, err
	}
	for _, d := range dishes {
		for _, dishElem := range dishList {
			if d.ID == dishElem.ID {
				_, err := c.DishRepo.Update(context.Background(), &dish.DishModel{
					ID:          dishElem.ID,
					Name:        dishElem.Name,
					Description: dishElem.Description,
					Price:       dishElem.Price,
					Quantity:    dishElem.Quantity - d.Quantity})

				if err != nil {
					return 0, err
				}

				_, err = c.OrderDishRepo.Add(context.Background(), &order_dish.OrderDishModel{
					OrderID:  id,
					DishID:   d.ID,
					Quantity: d.Quantity,
					Price:    dishElem.Price})
				if err != nil {
					return 0, err
				}
			}
		}
	}
	return id, nil
}

func (c *Core) GetOrder(token string, id int64) (string, error) {
	u, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return "", custom_errors.ErrUserNotAuthorized
	}

	o, err := c.OrderRepo.GetById(context.Background(), id)
	if err != nil {
		return "", custom_errors.ErrOrderNotFound
	}

	if u.Role == "customer" && u.Id != o.UserID {
		return "", custom_errors.ErrPermissionDenied
	}
	return o.Status, err
}

func (c *Core) AddDish(token string, d *dish.DishModel) (int64, error) {
	u, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return 0, custom_errors.ErrUserNotAuthorized
	}
	if u.Role != "manager" {
		return 0, custom_errors.ErrPermissionDenied
	}
	return c.DishRepo.Add(context.Background(), d)
}

func (c *Core) GetDish(token string, id int64) (*dish.DishModel, error) {
	u, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthorized
	}
	if u.Role != "manager" {
		return nil, custom_errors.ErrPermissionDenied
	}
	return c.DishRepo.GetById(context.Background(), id)
}

func (c *Core) UpdateDish(token string, d *dish.DishModel) error {
	u, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return custom_errors.ErrUserNotAuthorized
	}
	if u.Role != "manager" {
		return custom_errors.ErrPermissionDenied
	}

	c.Lock()
	defer c.Unlock()
	_, err = c.DishRepo.Update(context.Background(), d)
	return err
}

func (c *Core) DeleteDish(token string, id int64) error {
	u, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return custom_errors.ErrUserNotAuthorized
	}
	if u.Role != "manager" {
		return custom_errors.ErrPermissionDenied
	}

	c.Lock()
	defer c.Unlock()
	_, err = c.DishRepo.Delete(context.Background(), id)
	return err
}

func (c *Core) ShowMenu(token string) ([]dish.DishModel, error) {
	_, err := c.AuthClient.Authorize(context.Background(), token)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthorized
	}

	dishList, _ := c.DishRepo.List(context.Background())
	menu := make([]dish.DishModel, 0)
	for _, d := range dishList {
		if d.Quantity > 0 {
			menu = append(menu, *d)
		}
	}
	return menu, nil
}
