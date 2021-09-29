/**
 * @File: user.go
 * @Author: hsien
 * @Description:
 * @Date: 9/23/21 12:07 PM
 */

package repository

import (
	"context"
	"custom_server/internal/model/postgres"
	"custom_server/internal/model/request"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) AddUser(ctx context.Context, user *postgres.User) error {
	return u.db.WithContext(ctx).Model(user).Create(user).Error
}

func (u *User) QueryUser(ctx context.Context, query *request.UserQuery) (int64, []*postgres.User, error) {
	var (
		count  int64
		result = make([]*postgres.User, 0)
	)
	db := u.db.WithContext(ctx).Model(&postgres.User{})
	if query.Name != "" {
		db = db.Where("user_name like ?", fmt.Sprintf("%s%%", query.Name))
	}
	if query.PhoneNumber != "" {
		db = db.Where("phone_number like ?", fmt.Sprintf("%s%%", query.PhoneNumber))
	}
	if query.Sex != 0 {
		db = db.Where("wechat_sex=?", query.Sex)
	}

	err := db.Offset(query.PageNo - 1).Limit(query.PageSize).Find(&result).
		Offset(-1).Limit(-1).Count(&count).Error
	return count, result, err
}

func (u *User) QueryUserByID(ctx context.Context, userId string) (*postgres.User, error) {
	var result *postgres.User
	err := u.db.WithContext(ctx).Model(&postgres.User{}).
		Where("user_id=?", userId).Limit(1).Find(&result).Error
	return result, err
}

func (u *User) QueryUserByPhone(ctx context.Context, phoneNumber string) (*postgres.User, error) {
	var result *postgres.User
	err := u.db.WithContext(ctx).Model(&postgres.User{}).
		Select("user_id", "shadow").Where("phone_number=?", phoneNumber).Limit(1).Find(&result).Error
	return result, err
}
