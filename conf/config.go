package conf

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Default() *Config {
	return &Config{
		Application: &application{
			Host:   "127.0.0.1",
			Port:   8080,
			Domain: "http://127.0.0.1",
		},
		MySQL: &mySQL{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "mydevice",
			Username: "root",
			Password: "123456",
			Debug:    true,
		},
	}
}

type Config struct {
	Application *application `toml:"app" yaml:"app" json:"app"`
	MySQL       *mySQL       `toml:"mysql" yaml:"mysql" json:"mysql"`
}

func (c *Config) ToYAML() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

type application struct {
	Host   string `toml:"host" yaml:"host" json:"host"`
	Port   int    `toml:"port" yaml:"port" json:"port"`
	Domain string `toml:"domain" yaml:"domain" json:"domain"`

	server *gin.Engine
	lock   sync.Mutex
	root   gin.IRouter
}

func (a *application) GinServer() *gin.Engine {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.server == nil {
		a.server = gin.Default()
		//加载中间件CROS
		a.server.Use(cors.Default())
	}
	return a.server
}

func (a *application) GinRootRouter() gin.IRouter {
	r := a.GinServer()

	if a.root == nil {
		a.root = r.Group("wod").Group("api").Group("v1")
	}
	return a.root
}

func (a *application) Address() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *application) Start() error {
	r := a.GinServer()
	return r.Run(a.Address())
}

type mySQL struct {
	Host     string `json:"host" yaml:"host" toml:"host" env:"DATASOURCE_HOST"`
	Port     int    `json:"port" yaml:"port" toml:"port" env:"DATASOURCE_PORT"`
	DB       string `json:"database" yaml:"database" toml:"database" env:"DATASOURCE_DB"`
	Username string `json:"username" yaml:"username" toml:"username" env:"DATASOURCE_USERNAME"`
	Password string `json:"password" yaml:"password" toml:"password" env:"DATASOURCE_PASSWORD"`
	Debug    bool   `json:"debug" yaml:"debug" toml:"debug" env:"DATASOURCE_DEBUG"`

	// 保证他是单例
	db *gorm.DB
	// 加锁
	lock sync.Mutex
}

func (m *mySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DB,
	)
}

func (m *mySQL) GetDB() *gorm.DB {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.db == nil {
		db, err := gorm.Open(mysql.Open(m.DSN()), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		m.db = db

		if m.Debug {
			m.db = db.Debug()
		}
	}
	return m.db
}
