package services

import (
	"GORushShoping/datamodels"
	"GORushShoping/repositories"
)

//业务逻辑的处理

// IProductService 定义接口功能
type IProductService interface {
	GetProductById(id int64) (product *datamodels.Product, err error)
	GetAllProduct() (products []*datamodels.Product, err error)
	DeleteProductById(id int64) (state bool)
	InsertProduct(product *datamodels.Product) (affected int64, err error)
	UpdateProduct(product *datamodels.Product) (err error)
	SubNumberOne(productID int64) error
}

type ProductService struct {
	productRepository repositories.IProduct
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}

func NewProductService(rp repositories.IProduct) IProductService {
	return &ProductService{productRepository: rp}
}

func (p *ProductService) GetProductById(id int64) (product *datamodels.Product, err error) {
	return p.productRepository.SelectById(id)
}

func (p *ProductService) GetAllProduct() (products []*datamodels.Product, err error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductById(id int64) (state bool) {
	return p.productRepository.Delete(id)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (affected int64, err error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) (err error) {
	return p.productRepository.Update(product)
}
