package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Db              *DbConf `yaml:"db"`
	DeadLetterQueue string  `yaml:"dead_letter_queue,omitempty"`

	Gotify   []GotifyConf   `yaml:"gotify" validate:"dive"`
	Awtrix   []AwtrixConf   `yaml:"awtrix" validate:"dive"`
	Telegram []TelegramConf `yaml:"telegram" validate:"dive"`
	Email    []EmailConf    `yaml:"email" validate:"dive"`

	EventSourceImpl []string   `yaml:"events_impl" validate:"required,dive,oneof=kafka rabbitmq http"`
	Kafka           *KafkaConf `yaml:"kafka"`
	RabbitMq        *RabbitMq  `yaml:"rabbitmq"`
	Http            *HttpConf  `yaml:"http"`
	Smtp            *SmtpConf  `yaml:"smtp"`

	MetricsAddr string `yaml:"metrics_addr" validate:"omitempty,tcp_addr"`
}

type DbConf struct {
	Type string `yaml:"type" validate:"oneof=memory sqlite"`

	Name string `yaml:"name"`
}

type GotifyConf struct {
	ServiceUri string `yaml:"uri" validate:"required"`
	GotifyAddr string `yaml:"addr" validate:"required,url"`
	Token      string `yaml:"token" validate:"required_without=TokenFile"`
	TokenFile  string `yaml:"token_file" validate:"required_without=Token,omitempty,file"`
}

type AwtrixConf struct {
	ServiceUri string `yaml:"uri" validate:"required"`
	Addr       string `yaml:"addr" validate:"required,url"`
}

type TelegramConf struct {
	ServiceUri string  `yaml:"uri" validate:"required"`
	Token      string  `yaml:"token" validate:"required_without=TokenFile"`
	TokenFile  string  `yaml:"token_file" validate:"required_without=Token,omitempty,file"`
	Receivers  []int64 `yaml:"receivers" validate:"required"`
}

type EmailConf struct {
	ServiceUri string   `yaml:"uri" validate:"required"`
	Sender     string   `yaml:"token" validate:"required"`
	Receivers  []string `yaml:"receivers" validate:"required"`
	Host       string   `yaml:"receiver" validate:"required"`

	UserName     string `yaml:"user_name" validate:"required_without=UserNameFile"`
	UserNameFile string `yaml:"user_name_file" validate:"required_without=UserName,omitempty,file"`

	Password     string `yaml:"password" validate:"required_without=PasswordFile"`
	PasswordFile string `yaml:"password_file" validate:"required_without=Password,omitempty,file"`
}

type KafkaConf struct {
	// Mandatory options
	Enabled bool     `yaml:"enabled"`
	Brokers []string `yaml:"brokers" validate:"dive,required"`
	Topic   string   `yaml:"topic" validate:"required"`
	GroupId string   `yaml:"group_id" validate:"required"`

	// Advanced options
	Partition   int    `yaml:"partition" validate:"gte=0"`
	TlsCertFile string `yaml:"tls_cert_file" validate:"omitempty,file"`
	TlsKeyFile  string `yaml:"tls_key_file" validate:"omitempty,file"`
}

type RabbitMq struct {
	// Mandatory options
	Broker       string `yaml:"broker" validate:"required"`
	Port         int    `yaml:"port" validate:"omitempty,gte=80,lt=65535"`
	QueueName    string `yaml:"queue" validate:"required"`
	ConsumerName string `yaml:"consumer"`
	Vhost        string `yaml:"vhost" validate:"required,startswith=/"`
	Username     string `yaml:"username" validate:"required"`
	Password     string `yaml:"password" validate:"required"`
	UseSsl       bool   `yaml:"use_ssl"`

	// Advanced options
	TlsCertFile string `yaml:"tls_cert_file" validate:"omitempty,file"`
	TlsKeyFile  string `yaml:"tls_key_file" validate:"omitempty,file"`
}

type HttpConf struct {
	// Mandatory options
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address" validate:"required"`

	// Advanced options
	TlsCertFile string `yaml:"tls_cert_file" validate:"required_with=TlsKeyFile TlsClientCa,omitempty,filepath"`
	TlsKeyFile  string `yaml:"tls_key_file" validate:"required_with=TlsCertFile TlsClientCa,omitempty,filepath"`
	TlsClientCa string `yaml:"tls_client_ca_file" validate:"omitempty,filepath"`
}

type SmtpConf struct {
	// Mandatory options
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address" validate:"required_if=EventSourceImpl smtp"`
	Domain  string `yaml:"domain" validate:"required_with=Address"`

	// Advanced options
	TlsCertFile string `yaml:"tls_cert_file" validate:"required_unless=TlsKeyFile '',omitempty,filepath"`
	TlsKeyFile  string `yaml:"tls_key_file" validate:"required_unless=TlsCertFile '',omitempty,filepath"`
}

func getDefaultConfig() *Config {
	return &Config{
		Http: &HttpConf{
			Enabled: true,
			Address: "0.0.0.0:8080",
		},
		MetricsAddr: "0.0.0.0:9223",
	}
}

func Read(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf := getDefaultConfig()
	err = yaml.Unmarshal(content, conf)
	return conf, err
}
