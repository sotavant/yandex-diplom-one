package workers

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sotavant/yandex-diplom-one/accrual"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"net/http"
	"time"
)

const (
	requestInterval = time.Second * 2
	requestTimeout  = time.Millisecond * 100
	accrualURI      = "/api/orders/"
)

type AccrualWorker struct {
	Service        *accrual.Service
	AccrualAddress string
}

func NewAccrualWorker(s *accrual.Service, address string) *AccrualWorker {
	return &AccrualWorker{
		Service:        s,
		AccrualAddress: address,
	}
}

func (a *AccrualWorker) Run(ctx context.Context) {
	request := make(chan bool)

	go func(ctx context.Context) {
		for {
			select {
			case <-request:
				return
			default:
				<-time.After(requestInterval)
				a.UpdateOrders(ctx)
			}
		}
	}(ctx)
}

func (a *AccrualWorker) UpdateOrders(ctx context.Context) {
	orders, err := a.Service.GetNotProcessedOrders(ctx)
	if err != nil {
		panic(err)
	}

	for _, order := range orders {
		processedOrder, err := a.GetOrderInfo(order)
		if err != nil {
			continue
		}

		err = a.Service.UpdateOrderState(ctx, order, processedOrder)
		if err != nil {
			continue
		}
	}

	time.Sleep(requestTimeout)
}

func (a *AccrualWorker) GetOrderInfo(order domain.Order) (domain.OrderAccrual, error) {
	var resOrder domain.OrderAccrual
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&resOrder).
		Get(a.AccrualAddress + accrualURI + order.Number)

	if err != nil {
		//internal.Logger.Infow("get accrual info error", "err", err)
		return resOrder, err
	}

	code := resp.StatusCode()
	if code != http.StatusOK {
		internal.Logger.Infow("bad accrual response code", "errcode", code, "response", string(resp.Body()))
		return resOrder, fmt.Errorf("bad status %d", code)
	}

	return resOrder, nil
}
