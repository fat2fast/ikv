package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// RawModuleConfig để đọc config dưới dạng strings từ YAML
type RawModuleConfig struct {
	Module struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Enabled     string `yaml:"enabled"`
		Description string `yaml:"description"`
	} `yaml:"module"`
	Database struct {
		Connection struct {
			Driver     string `yaml:"driver"`
			Host       string `yaml:"host"`
			Port       string `yaml:"port"`
			Database   string `yaml:"database"`
			Username   string `yaml:"username"`
			Password   string `yaml:"password"`
			Schema     string `yaml:"schema"`
			AutoCreate string `yaml:"auto_create"`
			SSLMode    string `yaml:"ssl_mode"`
			Timezone   string `yaml:"timezone"`
		} `yaml:"connection"`
		Migration struct {
			Path   string `yaml:"path"`
			Table  string `yaml:"table"`
			Schema string `yaml:"schema"`
		} `yaml:"migration"`
		Performance struct {
			MaxOpenConns    string `yaml:"max_open_conns"`
			MaxIdleConns    string `yaml:"max_idle_conns"`
			ConnMaxLifetime string `yaml:"conn_max_lifetime"`
		} `yaml:"performance"`
	} `yaml:"database"`
}

// ModuleConfig định nghĩa cấu trúc config cho module với đúng kiểu dữ liệu
type ModuleConfig struct {
	Module struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Enabled     bool   `yaml:"enabled"`
		Description string `yaml:"description"`
	} `yaml:"module"`
	Database struct {
		Connection struct {
			Driver     string `yaml:"driver"`
			Host       string `yaml:"host"`
			Port       int    `yaml:"port"`
			Database   string `yaml:"database"`
			Username   string `yaml:"username"`
			Password   string `yaml:"password"`
			Schema     string `yaml:"schema"`
			AutoCreate bool   `yaml:"auto_create"`
			SSLMode    string `yaml:"ssl_mode"`
			Timezone   string `yaml:"timezone"`
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

// expandEnvVar expand một environment variable với format ${VAR:default}
func expandEnvVar(envStr string) string {
	// Nếu không phải format ${VAR:default}, trả về như cũ
	if !strings.HasPrefix(envStr, "${") || !strings.HasSuffix(envStr, "}") {
		return envStr
	}

	// Remove ${ và }
	content := envStr[2 : len(envStr)-1]

	// Split theo :
	parts := strings.SplitN(content, ":", 2)
	varName := parts[0]
	defaultValue := ""

	if len(parts) > 1 {
		defaultValue = parts[1]
	}

	// Lấy giá trị từ environment
	if value := os.Getenv(varName); value != "" {
		return value
	}

	return defaultValue
}

// convertRawToModuleConfig chuyển đổi RawModuleConfig thành ModuleConfig
func convertRawToModuleConfig(raw *RawModuleConfig) (*ModuleConfig, error) {
	config := &ModuleConfig{}

	// Module
	config.Module.Name = raw.Module.Name
	config.Module.Version = raw.Module.Version
	config.Module.Description = raw.Module.Description

	// Parse enabled
	enabledStr := expandEnvVar(raw.Module.Enabled)
	if enabled, err := strconv.ParseBool(enabledStr); err == nil {
		config.Module.Enabled = enabled
	} else {
		config.Module.Enabled = true // default
	}

	// Database Connection
	config.Database.Connection.Driver = expandEnvVar(raw.Database.Connection.Driver)
	config.Database.Connection.Host = expandEnvVar(raw.Database.Connection.Host)
	config.Database.Connection.Database = expandEnvVar(raw.Database.Connection.Database)
	config.Database.Connection.Username = expandEnvVar(raw.Database.Connection.Username)
	config.Database.Connection.Password = expandEnvVar(raw.Database.Connection.Password)
	config.Database.Connection.Schema = expandEnvVar(raw.Database.Connection.Schema)
	config.Database.Connection.SSLMode = expandEnvVar(raw.Database.Connection.SSLMode)
	config.Database.Connection.Timezone = expandEnvVar(raw.Database.Connection.Timezone)

	// Parse port
	portStr := expandEnvVar(raw.Database.Connection.Port)
	if port, err := strconv.Atoi(portStr); err == nil {
		config.Database.Connection.Port = port
	} else {
		config.Database.Connection.Port = 5432 // default
	}

	// Parse auto_create
	autoCreateStr := expandEnvVar(raw.Database.Connection.AutoCreate)
	if autoCreate, err := strconv.ParseBool(autoCreateStr); err == nil {
		config.Database.Connection.AutoCreate = autoCreate
	} else {
		config.Database.Connection.AutoCreate = true // default
	}

	// Database Migration
	config.Database.Migration.Path = expandEnvVar(raw.Database.Migration.Path)
	config.Database.Migration.Table = expandEnvVar(raw.Database.Migration.Table)
	config.Database.Migration.Schema = expandEnvVar(raw.Database.Migration.Schema)

	// Database Performance
	maxOpenConnsStr := expandEnvVar(raw.Database.Performance.MaxOpenConns)
	if maxOpenConns, err := strconv.Atoi(maxOpenConnsStr); err == nil {
		config.Database.Performance.MaxOpenConns = maxOpenConns
	} else {
		config.Database.Performance.MaxOpenConns = 10 // default
	}

	maxIdleConnsStr := expandEnvVar(raw.Database.Performance.MaxIdleConns)
	if maxIdleConns, err := strconv.Atoi(maxIdleConnsStr); err == nil {
		config.Database.Performance.MaxIdleConns = maxIdleConns
	} else {
		config.Database.Performance.MaxIdleConns = 2 // default
	}

	config.Database.Performance.ConnMaxLifetime = expandEnvVar(raw.Database.Performance.ConnMaxLifetime)

	return config, nil
}

// loadModuleConfig load config từ file config.yaml của module
func loadModuleConfig(moduleName string) (*ModuleConfig, error) {
	configPath := filepath.Join("/app/modules", moduleName, "config.yaml")

	// Kiểm tra file có tồn tại không
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Đọc file config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse YAML
	var raw RawModuleConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Convert RawModuleConfig to ModuleConfig
	config, err := convertRawToModuleConfig(&raw)
	if err != nil {
		return nil, fmt.Errorf("failed to convert raw config to module config: %v", err)
	}

	return config, nil
}

// getDatabaseURL tạo database URL từ config
func getDatabaseURL(config *ModuleConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		config.Database.Connection.Username,
		config.Database.Connection.Password,
		config.Database.Connection.Host,
		config.Database.Connection.Port,
		config.Database.Connection.Database,
		config.Database.Connection.SSLMode,
		config.Database.Connection.Schema,
	)
}

// getMigrationPath lấy đường dẫn migration từ config
func getMigrationPath(config *ModuleConfig) string {
	return config.Database.Migration.Path
}

// Migrate command với sub commands
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Quản lý database migrations",
	Long:  "Lệnh này cung cấp các chức năng để quản lý database migrations cho từng module",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Chạy migrations",
	Long:  "Chạy migrations cho module được chỉ định hoặc tất cả modules",
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, _ := cmd.Flags().GetString("module")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			fmt.Println("🚀 Đang chạy migrations cho tất cả modules...")
			// Lấy danh sách tất cả modules từ thư mục modules
			modules, err := getAvailableModules()
			if err != nil {
				log.Fatalf("Lỗi khi lấy danh sách modules: %v", err)
			}

			for _, module := range modules {
				fmt.Printf("📦 Chạy migration cho module: %s\n", module)
				if err := runMigrationUp(module); err != nil {
					log.Printf("❌ Lỗi khi chạy migration cho module %s: %v", module, err)
					continue
				}
				fmt.Printf("✅ Migration hoàn thành cho module: %s\n", module)
			}
			fmt.Println("🎉 Tất cả migrations đã hoàn thành!")
		} else if moduleName != "" {
			fmt.Printf("🚀 Đang chạy migrations cho module: %s\n", moduleName)
			if err := runMigrationUp(moduleName); err != nil {
				log.Fatalf("❌ Lỗi khi chạy migration: %v", err)
			}
			fmt.Println("✅ Migrations đã hoàn thành!")
		} else {
			fmt.Println("❌ Vui lòng chỉ định module với --module hoặc sử dụng --all")
			cmd.Help()
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Long:  "Rollback migrations cho module được chỉ định",
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, _ := cmd.Flags().GetString("module")
		steps, _ := cmd.Flags().GetInt("steps")

		if moduleName == "" {
			fmt.Println("❌ Vui lòng chỉ định module với --module")
			cmd.Help()
			return
		}

		fmt.Printf("🔄 Đang rollback %d migration(s) cho module: %s\n", steps, moduleName)
		if err := runMigrationDown(moduleName, steps); err != nil {
			log.Fatalf("❌ Lỗi khi rollback migration: %v", err)
		}
		fmt.Println("✅ Rollback đã hoàn thành!")
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Tạo migration mới",
	Long:  "Tạo file migration mới cho module được chỉ định",
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, _ := cmd.Flags().GetString("module")
		migrationName, _ := cmd.Flags().GetString("name")

		if moduleName == "" || migrationName == "" {
			fmt.Println("❌ Vui lòng chỉ định module với --module và tên migration với --name")
			cmd.Help()
			return
		}

		fmt.Printf("📝 Đang tạo migration: %s cho module: %s\n", migrationName, moduleName)
		if err := createMigration(moduleName, migrationName); err != nil {
			log.Fatalf("❌ Lỗi khi tạo migration: %v", err)
		}
		fmt.Println("✅ Migration đã được tạo thành công!")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Kiểm tra trạng thái migrations",
	Long:  "Kiểm tra trạng thái migrations cho module được chỉ định",
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, _ := cmd.Flags().GetString("module")

		if moduleName == "" {
			fmt.Println("❌ Vui lòng chỉ định module với --module")
			cmd.Help()
			return
		}

		fmt.Printf("📊 Đang kiểm tra trạng thái migration cho module: %s\n", moduleName)
		if err := checkMigrationStatus(moduleName); err != nil {
			log.Fatalf("❌ Lỗi khi kiểm tra trạng thái migration: %v", err)
		}
	},
}

var migratePendingCmd = &cobra.Command{
	Use:   "pending",
	Short: "Hiển thị migrations chưa được apply",
	Long:  "Hiển thị danh sách migrations chưa được apply cho module được chỉ định hoặc tất cả modules",
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, _ := cmd.Flags().GetString("module")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			fmt.Println("📋 Đang kiểm tra pending migrations cho tất cả modules...")
			modules, err := getAvailableModules()
			if err != nil {
				log.Fatalf("Lỗi khi lấy danh sách modules: %v", err)
			}

			for _, module := range modules {
				fmt.Printf("\n📦 Module: %s\n", module)
				if err := checkPendingMigrations(module); err != nil {
					log.Printf("❌ Lỗi khi kiểm tra pending migrations cho module %s: %v", module, err)
					continue
				}
			}
		} else if moduleName != "" {
			fmt.Printf("📋 Đang kiểm tra pending migrations cho module: %s\n", moduleName)
			if err := checkPendingMigrations(moduleName); err != nil {
				log.Fatalf("❌ Lỗi khi kiểm tra pending migrations: %v", err)
			}
		} else {
			fmt.Println("❌ Vui lòng chỉ định module với --module hoặc sử dụng --all")
			cmd.Help()
		}
	},
}

// Helper functions cho migration
func getAvailableModules() ([]string, error) {
	modulesDir := "modules"
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Kiểm tra xem có file config.yaml không
			configPath := filepath.Join(modulesDir, entry.Name(), "config.yaml")
			if _, err := os.Stat(configPath); err == nil {
				modules = append(modules, entry.Name())
			}
		}
	}
	return modules, nil
}

func runMigrationUp(moduleName string) error {
	config, err := loadModuleConfig(moduleName)
	if err != nil {
		return err
	}

	databaseURL := getDatabaseURL(config)
	migrationPath := getMigrationPath(config)
	fmt.Printf("📁 Migration Path: %s\n", migrationPath)

	// Tạo thư mục migrations nếu chưa có
	if err := os.MkdirAll(migrationPath, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %v", err)
	}

	// Tạo migration instance từ golang-migrate
	sourceURL := fmt.Sprintf("file://%s", migrationPath)
	fmt.Printf("📂 Source URL: %s\n", sourceURL)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Lấy version trước khi chạy migration
	oldVersion, _, err := m.Version()
	var beforeVersion uint = 0
	if err == nil {
		beforeVersion = oldVersion
	}

	// Chạy migrations lên version mới nhất
	fmt.Printf("🚀 Đang chạy migrations cho module: %s\n", moduleName)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Lấy version hiện tại
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %v", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Printf("✅ No migrations found, database is empty\n")
	} else if dirty {
		fmt.Printf("⚠️ Migration version %d is dirty (incomplete)\n", version)
	} else {
		fmt.Printf("✅ Migration completed, current version: %d\n", version)

		// Log các migration files đã được apply
		if version > beforeVersion {
			upFiles, err := filepath.Glob(filepath.Join(migrationPath, "*.up.sql"))
			if err == nil {
				fmt.Printf("📄 Files đã được apply:\n")
				for _, file := range upFiles {
					fileName := filepath.Base(file)
					migrationName := strings.TrimSuffix(fileName, ".up.sql")
					parts := strings.SplitN(migrationName, "_", 2)
					if len(parts) >= 1 {
						if fileVersion, parseErr := strconv.ParseUint(parts[0], 10, 64); parseErr == nil {
							if fileVersion > uint64(beforeVersion) && fileVersion <= uint64(version) {
								fmt.Printf("   ✅ %s\n", migrationName)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func runMigrationDown(moduleName string, steps int) error {
	config, err := loadModuleConfig(moduleName)
	if err != nil {
		return err
	}

	databaseURL := getDatabaseURL(config)
	migrationPath := getMigrationPath(config)
	fmt.Printf("📁 Migration Path: %s\n", migrationPath)
	fmt.Printf("🔢 Steps to rollback: %d\n", steps)

	// Tạo migration instance từ golang-migrate
	sourceURL := fmt.Sprintf("file://%s", migrationPath)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Lấy version trước khi rollback
	oldVersion, _, err := m.Version()
	var beforeVersion uint = 0
	if err == nil {
		beforeVersion = oldVersion
	}

	// Rollback migrations
	fmt.Printf("🔄 Đang rollback %d migration(s) cho module: %s\n", steps, moduleName)
	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %v", err)
	}

	// Lấy version hiện tại
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %v", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Printf("✅ Rollback completed, database is empty\n")
	} else if dirty {
		fmt.Printf("⚠️ Migration version %d is dirty (incomplete)\n", version)
	} else {
		fmt.Printf("✅ Rollback completed, current version: %d\n", version)
	}

	// Log các migration files đã được rollback
	if beforeVersion > 0 {
		downFiles, err := filepath.Glob(filepath.Join(migrationPath, "*.down.sql"))
		if err == nil {
			fmt.Printf("📄 Files đã được rollback:\n")
			for _, file := range downFiles {
				fileName := filepath.Base(file)
				migrationName := strings.TrimSuffix(fileName, ".down.sql")
				parts := strings.SplitN(migrationName, "_", 2)
				if len(parts) >= 1 {
					if fileVersion, parseErr := strconv.ParseUint(parts[0], 10, 64); parseErr == nil {
						var currentVersion uint64 = 0
						if err != migrate.ErrNilVersion {
							currentVersion = uint64(version)
						}
						if fileVersion <= uint64(beforeVersion) && fileVersion > currentVersion {
							fmt.Printf("   🔄 %s\n", migrationName)
						}
					}
				}
			}
		}
	}

	return nil
}

func createMigration(moduleName, migrationName string) error {
	config, err := loadModuleConfig(moduleName)
	if err != nil {
		return err
	}
	migrationPath := getMigrationPath(config)
	log.Printf("📁 Migration Path: %s", migrationPath)

	// Tạo thư mục migrations nếu chưa có
	if err := os.MkdirAll(migrationPath, 0755); err != nil {
		return err
	}

	// Tạo timestamp cho migration
	timestamp := time.Now().Format("20060102150405")

	upFile := filepath.Join(migrationPath, fmt.Sprintf("%s_%s.up.sql", timestamp, migrationName))
	downFile := filepath.Join(migrationPath, fmt.Sprintf("%s_%s.down.sql", timestamp, migrationName))

	// Tạo file up
	upContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Write your up migration here\n", migrationName, time.Now().Format("2006-01-02 15:04:05"))
	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return err
	}

	// Tạo file down
	downContent := fmt.Sprintf("-- Rollback: %s\n-- Created at: %s\n\n-- Write your down migration here\n", migrationName, time.Now().Format("2006-01-02 15:04:05"))
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return err
	}

	fmt.Printf("📄 Đã tạo file: %s\n", upFile)
	fmt.Printf("📄 Đã tạo file: %s\n", downFile)

	return nil
}

func checkMigrationStatus(moduleName string) error {
	config, err := loadModuleConfig(moduleName)
	if err != nil {
		return err
	}

	databaseURL := getDatabaseURL(config)
	migrationPath := getMigrationPath(config)

	fmt.Printf("📁 Migration Path: %s\n", migrationPath)

	// Kiểm tra số lượng file migration
	files, err := filepath.Glob(filepath.Join(migrationPath, "*.sql"))
	if err != nil {
		return err
	}

	fmt.Printf("📊 Tổng số file migration: %d\n", len(files))

	// Tạo migration instance từ golang-migrate
	sourceURL := fmt.Sprintf("file://%s", migrationPath)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Lấy version hiện tại
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %v", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Printf("📊 Database status: No migrations applied yet\n")
		fmt.Printf("🆕 Database is empty, ready for first migration\n")
	} else if dirty {
		fmt.Printf("⚠️ Database status: Migration version %d is dirty (incomplete)\n", version)
		fmt.Printf("🔧 You may need to fix this migration manually\n")
	} else {
		fmt.Printf("✅ Database status: Current migration version %d\n", version)
		fmt.Printf("🎯 Database is up to date\n")
	}

	return nil
}

func checkPendingMigrations(moduleName string) error {
	config, err := loadModuleConfig(moduleName)
	if err != nil {
		return err
	}

	databaseURL := getDatabaseURL(config)
	migrationPath := getMigrationPath(config)

	// Lấy tất cả file migration
	upFiles, err := filepath.Glob(filepath.Join(migrationPath, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	if len(upFiles) == 0 {
		fmt.Printf("📭 Không có migration nào trong module: %s\n", moduleName)
		return nil
	}

	// Tạo migration instance từ golang-migrate
	sourceURL := fmt.Sprintf("file://%s", migrationPath)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Lấy version hiện tại
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %v", err)
	}

	fmt.Printf("📊 Tổng số migration files: %d\n", len(upFiles))

	if err == migrate.ErrNilVersion {
		// Chưa có migration nào được apply
		fmt.Printf("🆕 Tất cả %d migrations đang pending (chưa apply):\n", len(upFiles))
		for _, file := range upFiles {
			fileName := filepath.Base(file)
			migrationName := strings.TrimSuffix(fileName, ".up.sql")
			fmt.Printf("   📄 %s\n", migrationName)
		}
		return nil
	}

	if dirty {
		fmt.Printf("⚠️ Database có dirty migration version %d\n", currentVersion)
		return nil
	}

	// So sánh với current version để tìm pending migrations
	pendingCount := 0
	fmt.Printf("✅ Current applied version: %d\n", currentVersion)
	fmt.Printf("📋 Pending migrations:\n")

	for _, file := range upFiles {
		fileName := filepath.Base(file)
		migrationName := strings.TrimSuffix(fileName, ".up.sql")

		// Extract version từ tên file (format: timestamp_name.up.sql)
		parts := strings.SplitN(migrationName, "_", 2)
		if len(parts) >= 1 {
			if versionStr := parts[0]; len(versionStr) >= 1 {
				// Parse version number
				if version, parseErr := strconv.ParseUint(versionStr, 10, 64); parseErr == nil {
					if version > uint64(currentVersion) {
						fmt.Printf("   📄 %s (version: %d)\n", migrationName, version)
						pendingCount++
					}
				}
			}
		}
	}

	if pendingCount == 0 {
		fmt.Printf("✅ Không có pending migrations. Database đã up-to-date!\n")
	} else {
		fmt.Printf("📊 Tổng số pending migrations: %d\n", pendingCount)
	}

	return nil
}

// init function để setup migrate commands
func init() {
	// Cấu hình flags cho migrate commands
	migrateUpCmd.Flags().StringP("module", "m", "", "Tên module để chạy migration")
	migrateUpCmd.Flags().BoolP("all", "a", false, "Chạy migration cho tất cả modules")

	migrateDownCmd.Flags().StringP("module", "m", "", "Tên module để rollback migration")
	migrateDownCmd.Flags().IntP("steps", "s", 1, "Số bước rollback")

	migrateCreateCmd.Flags().StringP("module", "m", "", "Tên module để tạo migration")
	migrateCreateCmd.Flags().StringP("name", "n", "", "Tên migration")

	migrateStatusCmd.Flags().StringP("module", "m", "", "Tên module để kiểm tra trạng thái")

	migratePendingCmd.Flags().StringP("module", "m", "", "Tên module để kiểm tra pending migrations")
	migratePendingCmd.Flags().BoolP("all", "a", false, "Kiểm tra pending migrations cho tất cả modules")

	// Thêm sub commands vào migrate command
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migratePendingCmd)
}
