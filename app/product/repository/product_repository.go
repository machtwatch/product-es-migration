package repository

import (
	"context"
	"product-es-migration/domain/product"

	"gorm.io/gorm"
)

// productRepo represent repository type of product that used as
// a collection oh the product repo.
type productRepo struct {
	readConn *gorm.DB
}

// NewProductRepo instantiate product repository.
// Accepts base repo and some clients provider
func NewProductRepo(dbSlave *gorm.DB) product.ProductRepo {
	return &productRepo{
		readConn: dbSlave,
	}
}

func (r *productRepo) GetProductVariants(ctx context.Context, productIds []int64) (map[int64][]product.ProductVariant, error) {
	query := `
		select 
			mpv.product_id, mpv.id variant_id, mpv.sku, mpv.variant_option1_value, mpv.variant_option2_value, mpv.idx,
			mpvp.our_price, mpvp.retail_price,
			(select sum(qty_available) 
					from ms_product_variant_stock mpvs 
					join ms_office mo on mo.id = mpvs.office_id 
					where mpvs.variant_id = mpv.id
					and mpvs.is_deleted = 0 and mo.is_deleted = 0
					and (mo.event_active = true or mo.event_active is null)) stock
		from ms_product_variant mpv 
		left join ms_product_variant_price mpvp on mpvp.variant_id = mpv.id 
		where mpvp.is_deleted = 0 and mpv.is_deleted = 0 and  mpv.product_id in ?
	`
	var result []product.GetProductVariantsModel

	sql := r.readConn.WithContext(ctx).Raw(query, productIds).Scan(&result)
	if sql.Error != nil {
		return nil, sql.Error
	}
	if sql.RowsAffected == 0 {
		return nil, nil
	}

	productMapper := map[int64][]product.ProductVariant{}
	for _, variant := range result {
		selected, ok := productMapper[variant.ProductId]
		if ok {
			selected = append(selected, product.ProductVariant{
				Id:                  variant.VariantId,
				Sku:                 variant.Sku,
				VariantOption1Value: variant.VariantOption1Value,
				VariantOption2Value: variant.VariantOption2Value,
				OurPrice:            variant.OurPrice,
				RetailPrice:         variant.RetailPrice,
				Stock:               variant.Stock,
				Idx:                 variant.Idx,
				Discount:            int(variant.GetDiscount()),
			})

			productMapper[variant.ProductId] = selected

		} else {
			productMapper[variant.ProductId] = []product.ProductVariant{{
				Id:                  variant.VariantId,
				Sku:                 variant.Sku,
				VariantOption1Value: variant.VariantOption1Value,
				VariantOption2Value: variant.VariantOption2Value,
				OurPrice:            variant.OurPrice,
				RetailPrice:         variant.RetailPrice,
				Stock:               variant.Stock,
				Discount:            int(variant.GetDiscount()),
			}}
		}
	}
	return productMapper, nil
}

func (r *productRepo) GetProducts(ctx context.Context, limit int, offset int) (products []product.Product, productIds []int64, err error) {

	var result []product.GetProductsModel

	query := `
		select 
			mp.id, mp.name product_name, mp.gender, mp.sku,
			mp.tags, mp.description, mp.type, mp.parent_type, mp.handle, mp.slug, mp.published_date,
			mp.product_group_id, mp.is_pre_order, mp.color,
			mb.id brand_id, mb.name brand_name, mb.is_shop_by_whatsapp,
			mpg.id "gender_id", mpt.id "type_id", mppt.id "parent_type_id",
			mp.updated_date, mp.created_date
		from ms_product mp
		left join ms_brand mb on mb.id = mp.brand_id
		left join ms_product_gender mpg on mpg.code = mp.gender
		left join ms_product_type mpt on mpt.code = mp.type
		left join ms_product_parent_type mppt on mppt.code = mp.parent_type
		where mp.is_published = 1
		order by mp.id 
		limit ? offset ?
	`

	sql := r.readConn.WithContext(ctx).Raw(query, limit, offset).Scan(&result)
	if sql.Error != nil {
		return nil, nil, sql.Error
	}
	if sql.RowsAffected == 0 {
		return nil, nil, nil
	}

	productIds = []int64{}
	products = []product.Product{}
	for _, p := range result {
		products = append(products,
			product.Product{
				Id:                 p.Id,
				Name:               p.ProductName,
				Gender:             p.Gender,
				Tags:               p.Tags,
				TagsLowerCaseArr:   p.GetTagCommaDelimeterLowerCase(),
				Sku:                p.Sku,
				Description:        p.Description,
				Type:               p.Type,
				TypeId:             p.TypeId,
				ParentType:         p.ParentType,
				ParentTypeId:       p.ParentTypeId,
				Handle:             p.Handle,
				Slug:               p.Slug,
				PublishedDate:      p.PublishedDate,
				CreatedDate:        p.CreatedDate,
				UpdatedDate:        p.UpdatedDate,
				ProductGroupId:     p.ProductGroupId,
				IsPreOrder:         p.IsPreOrder,
				Color:              p.Color,
				GenderId:           p.GenderId,
				ColorsLowerCaseArr: p.GetColorCommaDelimeterLowerCase(),
				Brand: product.Brand{
					Id:               p.BrandId,
					Name:             p.BrandName,
					IsShopByWhatsapp: p.IsShopByWhatsapp,
				},
			})
		productIds = append(productIds, p.Id)
	}

	return products, productIds, nil
}

func (r *productRepo) GetProductImages(ctx context.Context, productIds []int64) (map[int64][]product.ProductImage, error) {
	query := `
		select 
			id, product_id, idx, src 
		from ms_product_image 
		where is_deleted = 0 and product_id in ?
	`
	var result []product.GetProductImagesModel

	sql := r.readConn.WithContext(ctx).Raw(query, productIds).Scan(&result)
	if sql.Error != nil {
		return nil, sql.Error
	}
	if sql.RowsAffected == 0 {
		return nil, nil
	}

	productMapper := map[int64][]product.ProductImage{}
	for _, image := range result {
		selected, ok := productMapper[image.ProductId]
		if ok {
			selected = append(selected, product.ProductImage{
				Id:  image.Id,
				Src: image.Src,
				Idx: image.Idx,
			})

			productMapper[image.ProductId] = selected

		} else {
			productMapper[image.ProductId] = []product.ProductImage{{
				Id:  image.Id,
				Src: image.Src,
				Idx: image.Idx,
			}}
		}
	}
	return productMapper, nil
}
