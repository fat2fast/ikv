package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"fat2fast/ikv/modules/book"
	"fat2fast/ikv/modules/user"
	"fat2fast/ikv/shared"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start service",
	Run: func(cmd *cobra.Command, args []string) {
		// Set timezone default
		loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
		if err != nil {
			panic(err)
		}
		time.Local = loc
		// Setup Logger
		shared.SetupLogger()
		// Start Gin
		r := gin.Default()
		err = r.SetTrustedProxies([]string{"192.168.1.1", "10.1.0.0/24"})
		if err != nil {
			panic(err)
		}
		// // Setup Database
		// db := shared.ConnectDb()
		// // init db to app
		// r.Use(shared.InitDb(db))

		var routes = shared.GetUrl()
		for _, route := range routes {
			r.Handle(route.Method, route.Path, route.HandlerFunc)
		}

		// Khởi tạo module registry
		registry := shared.NewModuleRegistry()

		// Khởi tạo và đăng ký module Book
		bookModule, err := book.NewModule()
		if err != nil {
			log.Fatalf("Failed to initialize Book module: %v", err)
		}
		registry.RegisterModule(bookModule)
		// Khởi tạo và đăng ký module User
		userModule, err := user.NewModule()
		if err != nil {
			log.Fatalf("Failed to initialize User module: %v", err)
		}
		registry.RegisterModule(userModule)

		// Đăng ký tất cả các module với router
		if err := registry.RegisterAllModules(r); err != nil {
			log.Fatalf("Failed to register modules: %v", err)
		}

		// In ra danh sách các module đã đăng ký
		log.Println("Registered modules:")
		for _, module := range registry.GetEnabledModules() {
			log.Printf("- %s", module.GetName())
		}

		// Khởi chạy server
		r.Run(":" + os.Getenv("SERVICE_PORT"))
	},
}

// Sub command ví dụ 1: Version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "In ra phiên bản của ứng dụng",
	Long:  "Lệnh này sẽ in ra thông tin phiên bản chi tiết của ứng dụng IKV",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("IKV Application v1.0.0")
		fmt.Println("Build Date: 2024-01-01")
		fmt.Println("Go Version: go1.21+")
	},
}

// Sub command ví dụ 2: Database command với flags
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Quản lý database",
	Long:  "Lệnh này cung cấp các chức năng để quản lý database",
}

var dbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Kiểm tra trạng thái kết nối database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Đang kiểm tra kết nối database...")

		// Setup Logger
		shared.SetupLogger()

		// Thử kết nối database
		db := shared.ConnectDb()
		if db != nil {
			fmt.Println("✅ Kết nối database thành công!")
		} else {
			fmt.Println("❌ Không thể kết nối database!")
		}
	},
}

// Sub command ví dụ 3: Module command với flags
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Quản lý modules",
	Long:  "Lệnh này cung cấp các chức năng để quản lý các modules trong hệ thống",
}

var moduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "Liệt kê tất cả modules",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📋 Danh sách modules hiện có:")

		// Khởi tạo module registry
		registry := shared.NewModuleRegistry()

		// Khởi tạo và đăng ký module Book
		bookModule, err := book.NewModule()
		if err != nil {
			log.Fatalf("Failed to initialize Book module: %v", err)
		}
		registry.RegisterModule(bookModule)

		// In ra danh sách modules
		for _, module := range registry.GetEnabledModules() {
			fmt.Printf("  - %s\n", module.GetName())
		}
	},
}

// Khởi tạo và thêm tất cả sub commands
func init() {
	// Thêm version command
	rootCmd.AddCommand(versionCmd)

	// Thêm migrate command từ migrate.go
	rootCmd.AddCommand(migrateCmd)

	// Thêm db command và sub commands
	dbCmd.AddCommand(dbStatusCmd)
	rootCmd.AddCommand(dbCmd)

	// Thêm module command và sub commands
	moduleCmd.AddCommand(moduleListCmd)
	rootCmd.AddCommand(moduleCmd)

	// Có thể thêm flags cho các commands
	// Ví dụ: thêm flag --verbose cho version command
	versionCmd.Flags().BoolP("verbose", "v", false, "In thông tin chi tiết")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("failed to execute command", err)
	}
}
