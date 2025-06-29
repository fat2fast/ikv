package book

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"fat2fast/ikv/shared"
	sharedinfras "fat2fast/ikv/shared/infras"
	"fat2fast/ikv/shared/middleware"

	bookhttpgin "fat2fast/ikv/modules/book/infras/controller/http-gin"
	bookrepository "fat2fast/ikv/modules/book/infras/repository/gorm-pgsql"
	bookservice "fat2fast/ikv/modules/book/service"
	bookurlv1 "fat2fast/ikv/modules/book/urls/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config đại diện cho cấu hình của module Book
type Config struct {
	Module struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Description string `yaml:"description"`
		Enabled     bool   `yaml:"enabled"`
	} `yaml:"module"`

	Database struct {
		Connection struct {
			Driver     string `yaml:"driver"`
			Host       string `yaml:"host"`
			Port       string `yaml:"port"`
			Database   string `yaml:"database"`
			Username   string `yaml:"username"`
			Password   string `yaml:"password"`
			SSLMode    string `yaml:"ssl_mode"`
			Timezone   string `yaml:"timezone"`
			Schema     string `yaml:"schema"`
			AutoCreate bool   `yaml:"auto_create"`
		} `yaml:"connection"`

		Migration struct {
			Path   string `yaml:"path"`
			Table  string `yaml:"table"`
			Schema string `yaml:"schema"`
		} `yaml:"migration"`

		Performance struct {
			MaxOpenConns    int    `yaml:"max_open_conns"`
			MaxIdleConns    int    `yaml:"max_idle_conns"`
			ConnMaxLifetime string `yaml:"conn_max_lifetime"`
		} `yaml:"performance"`
	} `yaml:"database"`
}

// Module đại diện cho module Book
type Module struct {
	config Config
	DB     *gorm.DB
}

// NewModule tạo một instance mới của module Book
func NewModule() (*Module, error) {
	// Lấy đường dẫn của module
	_, filename, _, _ := runtime.Caller(0)
	modulePath := filepath.Dir(filename)

	// Load config từ file YAML
	var config Config
	err := shared.GetModuleConfig(modulePath, &config)
	if err != nil {
		return nil, fmt.Errorf("error loading module config: %v", err)
	}

	// Khởi tạo module
	module := &Module{
		config: config,
	}

	// Kết nối database nếu module được kích hoạt
	if module.IsEnabled() {
		// retry 5 times
		var db *gorm.DB
		var err error
		for i := 0; i < 5; i++ {
			db, err = module.connectDatabase()
			if err == nil {
				break
			}
			log.Printf("Error connecting to database: %v, retrying .. waiting 5 seconds", err)
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			return nil, fmt.Errorf("error connecting to database: %v", err)
		} else {
			module.DB = db
		}
	}

	return module, nil
}

// connectDatabase kết nối đến database dựa trên cấu hình module
func (m *Module) connectDatabase() (*gorm.DB, error) {
	dbConfig := m.config.Database.Connection

	// Xây dựng connection string cho GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.SSLMode, dbConfig.Timezone)

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Thiết lập schema nếu cần
	if dbConfig.AutoCreate {
		log.Printf("Auto create schema %s", dbConfig.Schema)
		db.Exec("CREATE SCHEMA IF NOT EXISTS " + dbConfig.Schema)
	}
	log.Printf("Setting search path to %s", dbConfig.Schema)
	db.Exec("SET search_path TO " + dbConfig.Schema)

	// Thiết lập connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Parse connection max lifetime
	connMaxLifetime, err := time.ParseDuration(m.config.Database.Performance.ConnMaxLifetime)
	if err != nil {
		connMaxLifetime = 5 * time.Minute // Default: 5 minutes
	}

	sqlDB.SetMaxOpenConns(m.config.Database.Performance.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.config.Database.Performance.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	log.Printf("Module %s connected to database %s", m.GetName(), dbConfig.Database)

	return db, nil
}

// RunMigrations chạy migrations cho module
func (m *Module) RunMigrations() error {
	if !m.IsEnabled() {
		return nil
	}

	log.Printf("Running migrations for module %s", m.GetName())

	// TODO: Implement migration logic using golang-migrate or other migration tool
	// Có thể sử dụng golang-migrate để chạy migrations từ thư mục m.config.Database.Migration.Path

	return nil
}

// Register đăng ký module với hệ thống
func (m *Module) Register(router *gin.Engine) error {
	if !m.IsEnabled() {
		log.Printf("Module %s is disabled", m.GetName())
		return nil
	}

	log.Printf("Registering module: %s (v%s)", m.GetName(), m.config.Module.Version)

	// Dependency injection
	controller := m.Initialize()
	routes := bookurlv1.GetRoutes(controller)

	log.Printf("Registering module routes")
	router.Use(middleware.RecoverMiddleware())
	log.Printf("Registering RecoverMiddleware")

	v1 := router.Group("/v1")
	bookV1 := v1.Group("/books")

	for _, route := range routes {
		bookV1.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	return nil
}

// GetName trả về tên của module
func (m *Module) GetName() string {
	return m.config.Module.Name
}

// IsEnabled kiểm tra module có được kích hoạt không
func (m *Module) IsEnabled() bool {
	return m.config.Module.Enabled
}

// GetConfig trả về cấu hình của module
func (m *Module) GetConfig() Config {
	return m.config
}

// GetDB trả về kết nối database của module
func (m *Module) GetDB() *gorm.DB {
	return m.DB
}

// Initialize khởi tạo và dependency injection cho module
func (m *Module) Initialize() *bookhttpgin.BookHTTPController {
	log.Printf("Initializing book module ")
	dbCtx := sharedinfras.NewDbContext(m.DB)

	// Repository
	bookRepository := bookrepository.NewBookRepository(dbCtx)

	// Command handlers
	createCmdHandler := bookservice.NewCreateBookCommandHandler(bookRepository)
	updateCmdHandler := bookservice.NewUpdateBookCommandHandler(bookRepository)
	deleteCmdHandler := bookservice.NewDeleteBookCommandHandler(bookRepository)

	// Query handlers
	getDetailQryHandler := bookservice.NewGetBookDetailQueryHandler(bookRepository)
	listQryHandler := bookservice.NewListBooksQueryHandler(bookRepository)

	// HTTP Controller
	bookHTTPController := bookhttpgin.NewBookHTTPController(
		createCmdHandler,
		updateCmdHandler,
		deleteCmdHandler,
		getDetailQryHandler,
		listQryHandler,
	)

	return bookHTTPController
}
