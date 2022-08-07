package models

import (
	"fmt"
	"log"
	"time"

	"errors"
	"html"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Shop struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Address   string    `gorm:"size:255" json:"address"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Shop) BeforeSave() error {
	hashedPassword, err := Hash(s.Password)
	if err != nil {
		return err
	}
	s.Password = string(hashedPassword)
	return nil
}

func (s *Shop) Prepare() {
	s.ID = 0
	s.Title = html.EscapeString(strings.TrimSpace(s.Title))
	s.Email = html.EscapeString(strings.TrimSpace(s.Email))
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *Shop) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":

		if err := checkmail.ValidateFormat(s.Email); err != nil && s.Email != "" {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if s.Password == "" {
			return errors.New("Required Password")
		}
		if s.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(s.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "create":
		if s.Title == "" {
			return errors.New("Required Title")
		}
		if s.Password == "" {
			return errors.New("Required Password")
		}
		if s.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(s.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	default:
		return nil
	}
}

func (s *Shop) SaveShop(db *gorm.DB) (*Shop, error) {

	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return &Shop{}, err
	}
	return s, nil
}

func (u *Shop) FindAllShops(db *gorm.DB) (*[]Shop, error) {
	var err error
	shops := []Shop{}
	err = db.Debug().Model(&Shop{}).Select([]string{"email", "title", "id"}).Limit(100).Find(&shops).Error

	if err != nil {
		return &[]Shop{}, err
	}
	return &shops, err
}

func (s *Shop) FindShopByID(db *gorm.DB, uid uint32) (*Shop, error) {
	var err error
	err = db.Debug().Model(Shop{}).Select([]string{"email", "title", "id", "created_at", "updated_at"}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return &Shop{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Shop{}, errors.New("Shop Not Found")
	}

	return s, err
}

func (s *Shop) UpdateShop(db *gorm.DB, uid uint32) (*Shop, error) {

	// To hash the password
	err := s.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&Shop{}).Where("id = ?", uid).Take(&Shop{}).UpdateColumns(
		map[string]interface{}{
			"password":  s.Password,
			"title":     s.Title,
			"email":     s.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Shop{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Shop{}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return &Shop{}, err
	}
	return s, nil
}

func (s *Shop) DeleteShop(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Shop{}).Where("id = ?", uid).Take(&Shop{}).Delete(&Shop{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (shop *Shop) PartialUpdateShop(db *gorm.DB, uid uint32, password string, title string, email string, address string) (*Shop, error) {
	var err error
	err = db.Debug().Model(&Shop{}).Where("id = ?", uid).Updates(Shop{Title: title, Password: password, Email: email, Address: address, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Shop{}, err
	}
	if db.Error != nil {
		return &Shop{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Shop{}).Where("id = ?", uid).Take(&shop).Error
	fmt.Print(err)
	if err != nil {
		return &Shop{}, err
	}
	return shop, nil
}
