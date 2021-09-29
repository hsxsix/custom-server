/**
 * @File: orm_test.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 11:22 AM
 */

package database

import (
	"custom_server/pkg/config"
	"testing"
)

func TestInitDB(t *testing.T) {
	type args struct {
		cfg *config.DataBase
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "mysql_test",
			args: struct {
				cfg *config.DataBase
			}{
				cfg: &config.DataBase{
					Type:    "mysql",
					DSN:     "root:mariadb@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
					MaxConn: 3,
				}},
			wantErr: false,
		},
		{
			name: "postgres_test",
			args: struct {
				cfg *config.DataBase
			}{
				cfg: &config.DataBase{
					Type:    "postgres",
					DSN:     "host=localhost user=postgres password=postgres dbname=server port=5432 sslmode=disable TimeZone=Asia/Shanghai",
					MaxConn: 3,
				}},
			wantErr: false,
		},
		{
			name: "sqlite_test",
			args: struct {
				cfg *config.DataBase
			}{
				cfg: &config.DataBase{
					Type:    "sqlite",
					DSN:     "gorm.db",
					MaxConn: 1,
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := InitDB(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
