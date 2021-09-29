/**
 * @File: struct.go
 * @Author: hsien
 * @Description:
 * @Date: 9/17/21 6:15 PM
 */

package config

type Config struct {
	Server   *Server              `toml:"server"`
	Log      *Log                 `toml:"log"`
	DataBase map[string]*DataBase `toml:"database"`
}

type Server struct {
	Addr  string `toml:"addr"`
	Debug bool   `toml:"debug"`
}

type Log struct {
	LogLevel string `toml:"log_level"`
	LogFile  string `toml:"log_file"`
}

type DataBase struct {
	Type     string `toml:"type"`
	DSN      string `toml:"dsn"`
	MaxConn  int    `toml:"max_conn"`
	LogLevel string `toml:"-"`
}
