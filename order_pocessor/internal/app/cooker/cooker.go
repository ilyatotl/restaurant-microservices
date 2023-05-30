package cooker

import (
	"context"
	"order_pocessor/internal/app/order"
	"time"
)

const (
	Pending    = "Pending"
	InProgress = "In progress"
	Completed  = "Completed"
)

type cooker struct {
	OrderRepo *order.OrderRepository
}

func StartCooker(ctx context.Context, repository *order.OrderRepository) {
	cooker := &cooker{OrderRepo: repository}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			o, err := cooker.OrderRepo.GetPendingOrder(ctx)
			if err != nil {
				continue
			}
			cooker.OrderRepo.ChangeStatus(ctx, o.ID, InProgress)
			time.Sleep(time.Second * 30)
			cooker.OrderRepo.ChangeStatus(ctx, o.ID, Completed)
		}
	}
}
