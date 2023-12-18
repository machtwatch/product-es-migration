package product

import "context"

type ProductUC interface {
	GetProducts(ctx context.Context, page int, size int) ([]Product, error)
}
