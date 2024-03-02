package products_repository

import (
	"errors"
	"trab02/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts() ([]models.Product, error) {
	var products []models.Product
	result := r.db.Find(&products)
	return products, result.Error
}

func (r *ProductRepository) GetProduct(id uint64) (models.Product, error) {
	var product models.Product
	result := r.db.First(&product, id)
	return product, result.Error
}

func (r *ProductRepository) PostProduct(newProduct models.Product) (uint64, error) {
	result := r.db.Create(&newProduct)
	if result.Error != nil {
		return 0, result.Error
	}
	return newProduct.Id, nil
}

func (r *ProductRepository) DeleteProduct(id uint64) error {
	product, err := r.GetProduct(id)
	if err != nil || product.Id == 0 {
		return errors.New("produto não está cadastrado no banco")
	}

	result := r.db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Needed to update the database as products come and go
func (r *ProductRepository) DecrementProductQuantity(id uint64) error {
	product, err := r.GetProduct(id)
	if err != nil || product.Id == 0 {
		return errors.New("produto não está cadastrado no banco")
	}

	if product.Quantity == 0 {
		return errors.New("não há mais produtos em estoque")
	}

	product.Quantity--
	if product.Quantity == 0 {
		r.DeleteProduct(id)
		return nil
	}
	result := r.db.Save(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProductRepository) IncrementProductQuantity(id uint64) error {
	product, err := r.GetProduct(id)
	if err != nil || product.Id == 0 {
		return errors.New("produto não está cadastrado no banco")
	}

	product.Quantity++
	result := r.db.Save(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
