/**
 * @File: user.go
 * @Author: hsien
 * @Description:
 * @Date: 9/22/21 5:20 PM
 */

package postgres

import "time"

type User struct {
	ID          uint   `gorm:"primarykey"`
	UserID      string `gorm:"column:user_id"`
	UserName    string `gorm:"column:user_name"`
	PhoneNumber string `gorm:"column:phone_number"`
	UserAvatar  string `gorm:"column:user_avatar"`
	UserSex     int    `gorm:"column:user_sex"`
	UserDesc    string `gorm:"column:user_desc"`
	Shadow      string `gorm:"column:shadow"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u User) TableName() string {
	return "app.user"
}
