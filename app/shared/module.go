package shared

import (
	"github.com/gin-gonic/gin"
)

// Module là interface mà tất cả các module phải implement
type Module interface {
	// Register đăng ký module với hệ thống (routes, middleware, etc.)
	Register(router *gin.Engine) error

	// GetName trả về tên của module
	GetName() string

	// IsEnabled kiểm tra module có được kích hoạt không
	IsEnabled() bool
}

// ModuleRegistry quản lý tất cả các module trong hệ thống
type ModuleRegistry struct {
	modules []Module
}

// NewModuleRegistry tạo một instance mới của ModuleRegistry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make([]Module, 0),
	}
}

// RegisterModule đăng ký một module vào registry
func (r *ModuleRegistry) RegisterModule(module Module) {
	r.modules = append(r.modules, module)
}

// RegisterAllModules đăng ký tất cả các module với router
func (r *ModuleRegistry) RegisterAllModules(router *gin.Engine) error {
	for _, module := range r.modules {
		if module.IsEnabled() {
			if err := module.Register(router); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetModules trả về danh sách tất cả các module
func (r *ModuleRegistry) GetModules() []Module {
	return r.modules
}

// GetEnabledModules trả về danh sách các module được kích hoạt
func (r *ModuleRegistry) GetEnabledModules() []Module {
	enabledModules := make([]Module, 0)
	for _, module := range r.modules {
		if module.IsEnabled() {
			enabledModules = append(enabledModules, module)
		}
	}
	return enabledModules
}

// GetModuleByName tìm module theo tên
func (r *ModuleRegistry) GetModuleByName(name string) Module {
	for _, module := range r.modules {
		if module.GetName() == name {
			return module
		}
	}
	return nil
}
