package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var cfg Config
var doOnce sync.Once

type Config struct {
	Application struct {
		Transport struct {
			Grpc struct {
				PORT int `mapstructure:"PORT"`
			} `mapstructure: "GRPC"`
			Http struct {
				PORT int `mapstructure:"PORT"`
			} `mapstructure:"HTTP"`
		} `mapstructure:"TRANSPORT"`
		Graceful struct {
			MaxSecond time.Duration `mapstructure:"MAX_SECOND"`
		} `mapstructure:"GRACEFUL"`
		EnablePprof bool `mapstructure:"ENABLE_PPROF"`
	} `mapstructure:"APPLICATION"`

	Auth struct {
		JwtToken struct {
			PublicKey       string        `mapstructure:"PUBLIC_KEY"`
			PrivateKey      string        `mapstructure:"PRIVATE_KEY"`
			Duration        time.Duration `mapstructure:"DURATION"`
			RefreshDuration time.Duration `mapstructure:"REFRESH_DURATION"`
		} `mapstructure:"JWT_TOKEN"`
	} `mapstructure:"AUTH"`

	DB struct {
		MySQL struct {
			Host        string `mapstructure:"HOST"`
			Port        int    `mapstructure:"PORT"`
			Name        string `mapstructure:"NAME"`
			User        string `mapstructure:"USER"`
			Pass        string `mapstructure:"PASS"`
			MaxPoolSize int    `mapstructure:"MAX_POOL_SIZE"`
		} `mapstructure:"MYSQL"`
		MongoDB struct {
			Host        string `mapstructure:"HOST"`
			Port        int    `mapstructure:"PORT"`
			Name        string `mapstructure:"NAME"`
			User        string `mapstructure:"USER"`
			Pass        string `mapstructure:"PASS"`
			MaxPoolSize int    `mapstructure:"MAX_POOL_SIZE"`
		} `mapstructure:"MONGODB"`
	} `mapstructure:"DB"`

	Cache struct {
		Redis struct {
			Host string `mapstructure:"HOST"`
			Port int    `mapstructure:"PORT"`
			DB   int    `mapstructure:"DB"`
			Pass string `mapstructure:"PASS"`
		}
	}
}

// ReadPublicKey will return public key
func (cfg Config) JwtPublicKey() (*rsa.PublicKey, error) {
	data, _ := pem.Decode([]byte(cfg.Auth.JwtToken.PublicKey))

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		return nil, errors.New("cannot reflect the interface")
	}

	return publicKey, nil
}

func (cfg Config) JwtPrivateKey() (*rsa.PrivateKey, error) {
	data, _ := pem.Decode([]byte(cfg.Auth.JwtToken.PrivateKey))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		return nil, err
	}
	return privateKeyImported, nil
}

func Get() Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("cannot read .env file: %v\n", err)
	}

	doOnce.Do(func() {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatalln("cannot unmarshaling config", err)
		}
	})

	return cfg
}
