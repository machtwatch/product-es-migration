package usecase

import (
	"context"
	"product-es-migration/domain"
	"product-es-migration/domain/product"
)

// Usecase type of product usecase.
//
// It contains the repositories and messaging event producer.
type productUC struct {
	productRepo product.ProductRepo
}

// NewProductUC instantiate product usecase.
//
// Accept domain repo collections, event producer, and xms integration. It will return the instance of usecase.
func NewProductUC(repoCollection domain.RepositoryCollection) product.ProductUC {
	return &productUC{
		productRepo: repoCollection.ProductRepo,
	}
}

func (r *productUC) GetProducts(ctx context.Context, page int, size int) ([]product.Product, error) {
	offset := (page - 1) * size
	limit := size
	products, productIds, err := r.productRepo.GetProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	productVariants, err := r.productRepo.GetProductVariants(ctx, productIds)
	if err != nil {
		return nil, err
	}
	productImages, err := r.productRepo.GetProductImages(ctx, productIds)
	if err != nil {
		return nil, err
	}
	for i := range products {

		variants, any := productVariants[products[i].Id]
		if any {
			selectedVariant := variants[0]
			if len(variants) > 1 {
				for j := 1; j < len(variants); j++ {
					if selectedVariant.Stock > 0 {
						if variants[j].Stock > 0 && variants[j].OurPrice < selectedVariant.OurPrice {
							selectedVariant = variants[j]
						} else if variants[j].Stock > 0 && variants[j].OurPrice == selectedVariant.OurPrice &&
							variants[j].Idx < selectedVariant.Idx {
							selectedVariant = variants[j]
						}
					} else {
						if variants[j].Stock > 0 {
							selectedVariant = variants[j]
						}
					}

				}
			}
			products[i].SelectedVariant = selectedVariant
			products[i].Variants = variants
			prodStock := 0
			for _, variant := range variants {
				prodStock += variant.Stock
			}

			if prodStock == 0 {
				products[i].IsOutOfStock = 1
			} else {
				products[i].IsOutOfStock = 0
			}
		}

		images, any := productImages[products[i].Id]
		if any {
			products[i].Images = images
		}
	}

	return products, err
}
