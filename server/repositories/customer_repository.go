package repositories

import (
	"server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db}
}

func (cr *CustomerRepository) Create(customer *models.Customer) error {
	return cr.db.Create(customer).Error
}

func (cr *CustomerRepository) List(userID uuid.UUID) ([]models.Customer, error) {
	var customers []models.Customer

	if err := cr.db.Where("user_id = ?", userID).Find(&customers).Error; err != nil {
		return nil, err
	}

	return customers, nil
}

func (cr *CustomerRepository) GetCustomerById(id uint) (*models.Customer, error) {
	var customer models.Customer

	if err := cr.db.Where("id = ?", id).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (cr *CustomerRepository) Update(customer *models.UpdateCustomerRequest) (*models.Customer, error) {
	updateData := make(map[string]interface{})

	if customer.CPF != nil {
		updateData["cpf"] = *customer.CPF
	}

	if customer.Name != nil {
		updateData["name"] = *customer.Name
	}

	if customer.Email != nil {
		updateData["email"] = *customer.Email
	}

	if customer.Phone != nil {
		updateData["phone"] = *customer.Phone
	}

	err := cr.db.Model(&models.Customer{}).Where("id = ?", customer.ID).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	updatedCustomer := &models.Customer{}
	err = cr.db.Where("id = ?", customer.ID).First(updatedCustomer).Error
	if err != nil {
		return nil, err
	}

	return updatedCustomer, nil
}

func (cr *CustomerRepository) Delete(customerID uint) error {
	customer := models.Customer{ID: customerID}
	return cr.db.Delete(&customer).Error
}