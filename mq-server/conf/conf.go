package conf

import (
	"log"
	"mq-server/model"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	RabbitMQ         string
	RabbitMQUser     string
	RabbitMQPassWord string
	RabbitMQHost     string
	RabbitMQPort     string
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")     // 设置配置文件格式为YAML
	viper.AutomaticEnv()            // 读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") // 读取环境变量的前缀为APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	dbUser := viper.GetString("db.user")
	dbPassword := viper.GetString("db.password")
	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbName := viper.GetString("db.name")
	conn := strings.Join([]string{dbUser, ":", dbPassword, "@tcp(", dbHost, ":", dbPort, ")/", dbName, "?charset=utf8mb4&parseTime=true"}, "")
	model.Database(conn)

	RabbitMQName := viper.GetString("rabbitmq.name")
	RabbitMQUser := viper.GetString("rabbitmq.User")
	RabbitMQPassWord := viper.GetString("rabbitmq.PassWord")
	RabbitMQHost := viper.GetString("rabbitmq.host")
	RabbitMQPort := viper.GetString("rabbitmq.port")

	// 连接RabbitMQ
	pathRabbitMQ := strings.Join([]string{RabbitMQName, "://", RabbitMQUser, ":", RabbitMQPassWord, "@", RabbitMQHost, ":", RabbitMQPort, "/"}, "")
	model.RabbitMQ(pathRabbitMQ)
	return nil
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}
