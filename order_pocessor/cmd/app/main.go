package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"order_pocessor/internal/app/authentication"
	"order_pocessor/internal/app/cooker"
	"order_pocessor/internal/app/core"
	db "order_pocessor/internal/app/database"
	"order_pocessor/internal/app/dish"
	http_service "order_pocessor/internal/app/http"
	"order_pocessor/internal/app/order"
	"order_pocessor/internal/app/order_dish"
)

const (
	port     = ":9091"
	grpcPort = ":50051"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		return
	}
	defer database.GetPool(ctx).Close()

	orderRepo := order.NewOrderRepo(order.NewRepository(database))
	dishRepo := dish.NewDishRepo(dish.NewRepository(database))
	orderDishRepo := order_dish.NewOrderDishRepo(order_dish.NewRepository(database))
	client, _ := authentication.NewClient(ctx, grpcPort)

	c := core.NewCore(orderRepo, dishRepo, orderDishRepo, client)
	r := gin.Default()
	s := http_service.NewHTTPService(r, c)

	go cooker.StartCooker(ctx, orderRepo)
	s.StartHTTPService(port)
}
