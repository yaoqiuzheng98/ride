package table

import (
	"crypto/rand"
	"errors"
	"strings"

	"gorm.io/gorm"

	"ride/db"
)

// User 用户表。
type User struct {
	gorm.Model
	BizID string `gorm:"size:8;uniqueIndex;not null;comment:业务ID" json:"biz_id"` // 对外暴露的业务ID，大写字母+数字，保证唯一
	Phone string `gorm:"size:20;uniqueIndex;not null;comment:手机号" json:"phone"`
}

func (User) TableName() string {
	return "user"
}

const bizIDCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bizIDLength = 8

// generateBizID 生成一个 8 位大写字母+数字的随机串。
func generateBizID() (string, error) {
	buf := make([]byte, bizIDLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i, b := range buf {
		buf[i] = bizIDCharset[int(b)%len(bizIDCharset)]
	}
	return string(buf), nil
}

// ErrPhoneAlreadyExists 手机号已注册。
var ErrPhoneAlreadyExists = errors.New("phone already exists")

// CreateUser 创建一个用户（仅手机号）。手机号已存在时返回 ErrPhoneAlreadyExists。
// BizID 自动生成大写字母+数字的随机串，冲突时自动重试保证唯一。
func CreateUser(phone string) (*User, error) {
	if phone == "" {
		return nil, errors.New("phone required")
	}
	for i := 0; i < 5; i++ {
		bizID, err := generateBizID()
		if err != nil {
			return nil, err
		}
		u := &User{BizID: bizID, Phone: phone}
		if err := db.GetClient().Create(u).Error; err != nil {
			// 区分 phone 冲突和 biz_id 冲突：MySQL 错误信息含冲突的键名
			if isDuplicateKey(err, "phone") {
				return nil, ErrPhoneAlreadyExists
			}
			if isDuplicateKey(err, "biz_id") {
				continue // 极小概率冲突，重试
			}
			return nil, err
		}
		return u, nil
	}
	return nil, errors.New("failed to generate unique biz_id")
}

// isDuplicateKey 判断是否为指定键的唯一索引冲突。
func isDuplicateKey(err error, key string) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if !strings.Contains(msg, "Duplicate entry") {
		return false
	}
	// MySQL 错误格式: Error 1062: Duplicate entry 'xxx' for key 'uk_xxx'
	return strings.Contains(msg, key)
}

// FindUserByPhone 按手机号查询用户。
func FindUserByPhone(phone string) (*User, error) {
	var u User
	if err := db.GetClient().Where("phone = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindUserByBizID 按业务ID查询用户。
func FindUserByBizID(bizID string) (*User, error) {
	var u User
	if err := db.GetClient().Where("biz_id = ?", bizID).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func init() {
	if err := db.GetClient().AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}
