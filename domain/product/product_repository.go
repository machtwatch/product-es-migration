package product

import (
	"context"
)

type ProductRepo interface {
	GetProducts(ctx context.Context, limit int, offset int) (products []Product, productIds []int64, err error)
	GetProductVariants(ctx context.Context, productIds []int64) (map[int64][]ProductVariant, error)
	GetProductImages(ctx context.Context, productIds []int64) (map[int64][]ProductImage, error)
}
