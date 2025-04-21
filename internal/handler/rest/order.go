package rest

import (
	"eventdrivensystem/internal/models/order"
	"eventdrivensystem/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *RouterHandler) CreateOrder(c echo.Context) error {
	var param order.CreateOrderParam

	if err := c.Bind(&param); err != nil {
		return errors.ErrBadRequest
	}

	if err := r.validator.Struct(param); err != nil {
		return errors.ErrBadRequest
	}

	err := r.uc.Order.CreateOrder(c.Request().Context(), &param)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Order created successfully",
	})
}

func (r *RouterHandler) RegisterOrderRoutes(base *echo.Group) {
	orders := base.Group("/orders")
	orders.POST("", r.CreateOrder)
}
