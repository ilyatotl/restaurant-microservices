package order_dish

type OrderDishModel struct {
	ID       int64 `db:"id"`
	OrderID  int64 `db:"order_id"`
	DishID   int64 `db:"dish_id"`
	Quantity int   `db:"quantity"`
	Price    int   `db:"price"`
}
