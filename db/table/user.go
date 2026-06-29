package table

import (
	"errors"

	"gorm.io/gorm"

	"ride/db"
)

// User 用户表，当前只有业务ID和手机号。
type User struct {
	Id    int64  `gorm:"primaryKey;autoIncrement;comment:业务ID" json:"id"`
	Phone string `gorm:"size:20;uniqueIndex;not null;comment:手机号" json:"phone"`
}

func (User) TableName() string {
	return "user"
}

// CreateUser 创建一个用户（仅手机号）。手机号已存在时返回错误。
func CreateUser(phone string) (*User, error) {
	if phone == "" {
		return nil, errors.New("phone required")
	}
	u := &User{Phone: phone}
	if err := db.GetClient().Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindUserByPhone 按手机号查询用户。
func FindUserByPhone(phone string) (*User, error) {
	var u User
	if err := db.GetClient().Where("phone = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func init() {
	if err := db.GetClient().AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}

// keep gorm import referenced
var _ = gorm.ErrRecordNotFound
