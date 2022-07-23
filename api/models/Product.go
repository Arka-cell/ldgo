package models

import (
	"errors"
	"html"
	"math"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:512" json:"description"`
	Price       float64   `gorm:"default:0" json:"price"`
	Shop        Shop      `json:"shop"`
	ShopID      uint32    `gorm:"not null" json:"shop_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	adjustedPrice := math.Floor(p.Price*100) / 100
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Price = adjustedPrice
	p.Shop = Shop{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Shop{}).Where("id = ?", p.ShopID).Take(&p.Shop).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i, _ := range products {
			err := db.Debug().Model(&Shop{}).Where("id = ?", products[i].ShopID).Take(&products[i].Shop).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Shop{}).Where("id = ?", p.ShopID).Take(&p.Shop).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) UpdateAProduct(db *gorm.DB) (*Product, error) {

	var err error

	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{Name: p.Name, Description: p.Description, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Shop{}).Where("id = ?", p.ShopID).Take(&p.Shop).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) DeleteAProduct(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ? and Shop_id = ?", pid, uid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
