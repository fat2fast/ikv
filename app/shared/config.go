package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

// LoadModuleConfig đọc file config.yaml của module và trả về map[string]interface{}
func LoadModuleConfig(modulePath string) (map[string]interface{}, error) {
	configPath := filepath.Join(modulePath, "config.yaml")

	// Kiểm tra file tồn tại
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Đọc file config
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Xử lý các biến môi trường trong config
	yamlContent := string(yamlFile)
	yamlContent = replaceEnvVars(yamlContent)

	// Parse YAML thành map
	var config map[string]interface{}
	err = yaml.Unmarshal([]byte(yamlContent), &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return config, nil
}

// replaceEnvVars thay thế các biến môi trường trong chuỗi
// Format: ${ENV_VAR:default_value}
func replaceEnvVars(content string) string {
	re := regexp.MustCompile(`\${([^:}]+)(?::([^}]+))?}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// Trích xuất tên biến và giá trị mặc định
		parts := re.FindStringSubmatch(match)
		envVar := parts[1]
		defaultValue := ""
		if len(parts) > 2 {
			defaultValue = parts[2]
		}

		// Lấy giá trị từ biến môi trường hoặc sử dụng giá trị mặc định
		value := os.Getenv(envVar)
		if value == "" {
			value = defaultValue
		}

		return value
	})
}

// GetModuleConfig là hàm wrapper để lấy config của module dưới dạng struct
func GetModuleConfig(modulePath string, config interface{}) error {
	configMap, err := LoadModuleConfig(modulePath)
	if err != nil {
		return err
	}

	// Convert map thành YAML bytes
	yamlBytes, err := yaml.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("error converting config map to YAML: %v", err)
	}

	// Unmarshal YAML bytes vào struct config
	err = yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		return fmt.Errorf("error parsing config into struct: %v", err)
	}

	return nil
}
