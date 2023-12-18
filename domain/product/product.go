package product

import (
	"math"
	"strings"
)

type GetProductsModel struct {
	Id               int64   `json:"id" gorm:"column:id"`
	Sku              string  `json:"sku" gorm:"column:sku"`
	ProductName      string  `json:"product_name" gorm:"column:product_name"`
	Gender           string  `json:"gender" gorm:"column:gender"`
	GenderId         *int    `json:"gender_id" gorm:"column:gender_id"`
	Tags             string  `json:"tags" gorm:"column:tags"`
	Description      string  `json:"description" gorm:"column:description"`
	Type             string  `json:"type" gorm:"column:type"`
	TypeId           *int    `json:"type_id" gorm:"column:type_id"`
	ParentType       string  `json:"parent_type" gorm:"column:parent_type"`
	ParentTypeId     *int    `json:"parent_type_id" gorm:"column:parent_type_id"`
	Handle           string  `json:"handle" gorm:"column:handle"`
	Slug             string  `json:"slug" gorm:"column:slug"`
	PublishedDate    string  `json:"published_date" gorm:"column:published_date"`
	CreatedDate      string  `json:"created_date" gorm:"column:created_date"`
	UpdatedDate      *string `json:"updated_date" gorm:"column:updated_date"`
	ProductGroupId   int64   `json:"product_group_id" gorm:"column:product_group_id"`
	IsPreOrder       bool    `json:"is_pre_order" gorm:"column:is_pre_order"`
	Color            string  `json:"color" gorm:"column:color"`
	BrandId          int64   `json:"brand_id" gorm:"column:brand_id"`
	BrandName        string  `json:"brand_name" gorm:"column:brand_name"`
	IsShopByWhatsapp bool    `json:"is_shop_by_whatsapp" gorm:"column:is_shop_by_whatsapp"`
}

type GetProductVariantsModel struct {
	ProductId           int64   `json:"product_id" gorm:"column:product_id"`
	VariantId           int64   `json:"variant_id" gorm:"column:variant_id"`
	Sku                 string  `json:"sku" gorm:"column:sku"`
	VariantOption1Value string  `json:"variant_option1_value" gorm:"column:variant_option1_value"`
	VariantOption2Value string  `json:"variant_option2_value" gorm:"column:variant_option1_value"`
	OurPrice            float64 `json:"our_price" gorm:"column:our_price"`
	RetailPrice         float64 `json:"retail_price" gorm:"column:retail_price"`
	Stock               int     `json:"stock" gorm:"column:stock"`
	Idx                 float64 `json:"idx" gorm:"column:idx"`
}

type GetProductImagesModel struct {
	Id        int64  `json:"id" gorm:"column:id"`
	ProductId int64  `json:"product_id" gorm:"column:product_id"`
	Src       string `json:"src" gorm:"column:src"`
	Idx       int    `json:"idx" gorm:"column:idx"`
}

type Product struct {
	Id                 int64            `json:"id"`
	Name               string           `json:"name"`
	Sku                string           `json:"sku"`
	Gender             string           `json:"gender"`
	GenderId           *int             `json:"gender_id"`
	Tags               string           `json:"tags"`
	TagsLowerCaseArr   []string         `json:"tags_lowercase_arr"`
	Description        string           `json:"description"`
	Type               string           `json:"type"`
	TypeId             *int             `json:"type_id"`
	ParentType         string           `json:"parent_type"`
	ParentTypeId       *int             `json:"parent_type_id"`
	Handle             string           `json:"handle"`
	Slug               string           `json:"slug"`
	PublishedDate      string           `json:"published_date"`
	UpdatedDate        *string          `json:"updated_date"`
	CreatedDate        string           `json:"created_date"`
	ProductGroupId     int64            `json:"product_group_id"`
	IsPreOrder         bool             `json:"is_pre_order"`
	Color              string           `json:"color"`
	ColorsLowerCaseArr []string         `json:"colors_lowercase_arr"`
	Brand              Brand            `json:"brand"`
	Variants           []ProductVariant `json:"variants"`
	SelectedVariant    ProductVariant   `json:"selected_variant"`
	Images             []ProductImage   `json:"images"`
	IsOutOfStock       int              `json:"is_out_of_stock"`
}

type Brand struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	IsShopByWhatsapp bool   `json:"is_shop_by_whatsapp"`
}

type ProductVariant struct {
	Id                  int64   `json:"id"`
	Sku                 string  `json:"sku"`
	VariantOption1Value string  `json:"variant_option1_value"`
	VariantOption2Value string  `json:"variant_option2_value"`
	OurPrice            float64 `json:"our_price"`
	RetailPrice         float64 `json:"retail_price"`
	Stock               int     `json:"stock"`
	Discount            int     `json:"discount"`
	Idx                 float64 `json:"idx"`
}

type ProductImage struct {
	Id  int64  `json:"id"`
	Src string `json:"src"`
	Idx int    `json:"idx"`
}

func (v *GetProductVariantsModel) GetDiscount() int64 {
	if v.OurPrice == 0 || v.RetailPrice == 0 || v.OurPrice >= v.RetailPrice {
		return 0
	}
	return 100 - int64(math.Round(v.OurPrice/v.RetailPrice*100))
}

func (p *GetProductsModel) GetTagCommaDelimeterLowerCase() []string {
	tags := strings.Split(p.Tags, ",")

	for i := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tags[i]))
	}
	return tags
}

func (p *GetProductsModel) GetColorCommaDelimeterLowerCase() []string {
	colors := strings.Split(p.Color, ",")

	for i := range colors {
		colors[i] = strings.ToLower(strings.TrimSpace(colors[i]))
	}
	return colors
}
