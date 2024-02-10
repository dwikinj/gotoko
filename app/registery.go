package app

import "github.com/dwikinj/gotoko/app/models"

type Model struct {
	Model any
}

func RegisterModels() []Model {
	return []Model{
		{Model: models.Address{}},
		{Model: models.Category{}},
		{Model: models.Order{}},
		{Model: models.OrderCustomer{}},
		{Model: models.OrderItem{}},
		{Model: models.Payment{}},
		{Model: models.Product{}},
		{Model: models.ProductImage{}},
		{Model: models.Section{}},
		{Model: models.Shipment{}},
		{Model: models.User{}},
	}
}
