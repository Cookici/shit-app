package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 应用配置
type Config struct {
	Server ServerConfig
	DB     DBConfig
	Wechat WechatConfig
	Aliyun AliyunConfig // 新增阿里云配置
	JWT    JWTConfig    // 新增JWT配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int
}

// DBConfig 数据库配置
type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// GetDSN 获取数据库连接字符串
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// WechatConfig 微信配置
type WechatConfig struct {
	AppID     string
	AppSecret string
}

// AliyunConfig 阿里云配置
type AliyunConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	URLPrefix       string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 168), // 默认7天
		},
		DB: DBConfig{
			Host:         getEnv("DB_HOST", ""),
			Port:         getEnvAsInt("DB_PORT", 63950),
			User:         getEnv("DB_USER", ""),
			Password:     getEnv("DB_PASSWORD", ""),
			DBName:       getEnv("DB_NAME", ""),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxLifetime:  time.Duration(getEnvAsInt("DB_MAX_LIFETIME", 3600)) * time.Second,
		},
		Wechat: WechatConfig{
			AppID:     getEnv("WECHAT_APPID", ""),
			AppSecret: getEnv("WECHAT_SECRET", ""),
		},
		Aliyun: AliyunConfig{
			Endpoint:        getEnv("ALIYUN_ENDPOINT", ""),
			AccessKeyID:     getEnv("ALIYUN_ACCESS_KEY_ID", ""),
			AccessKeySecret: getEnv("ALIYUN_ACCESS_KEY_SECRET", ""),
			BucketName:      getEnv("ALIYUN_BUCKET_NAME", ""),
			URLPrefix:       getEnv("ALIYUN_URL_PREFIX", ""),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取环境变量并转换为整数，如果不存在则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
