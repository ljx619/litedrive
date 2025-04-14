package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Ceph     CephConfig     `mapstructure:"ceph"`
	Cos      CosConfig      `mapstructure:"cos"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret  string `mapstructure:"secret"`
	Expires int    `mapstructure:"expires"`
}

type StorageConfig struct {
	Root             string `mapstructure:"root"`
	TempLocalRoot    string `mapstructure:"temp_local_root"`
	TempPartRoot     string `mapstructure:"temp_part_root"`
	CephRootDir      string `mapstructure:"ceph_root_dir"`
	CosRootDir       string `mapstructure:"cos_root_dir"`
	CurrentStoreType string `mapstructure:"current_store_type"`
}

type CephConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

type CosConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
}

const defaultConfigPath = "./configs"

func LoadConfig() (*Config, error) {
	// 读取环境变量 APP_ENV（默认为 "prod"）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod"
	}

	viper.SetConfigName("config") // 先加载默认配置
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultConfigPath)

	// 读取基础配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 读取环境配置（如果存在）
	viper.SetConfigName("config." + env)
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("No specific config for %s environment, using default\n", env)
	} else {
		// 打印使用的配置文件路径
		//fmt.Printf("Loaded config for %s environment: %s\n", env, viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
