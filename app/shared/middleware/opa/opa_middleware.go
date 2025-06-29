package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
)

// OPAMiddleware struct để lưu trữ cấu hình
type OPAMiddleware struct {
	config *OPAConfig
}

// NewOPAMiddleware tạo instance mới của OPA middleware
func NewOPAMiddleware(config *OPAConfig) (*OPAMiddleware, error) {
	if config == nil {
		return nil, fmt.Errorf("config không được nil")
	}

	// Thiết lập giá trị mặc định
	config.SetDefaults()

	// Validate cấu hình
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	middleware := &OPAMiddleware{
		config: config,
	}

	return middleware, nil
}

// validateConfig kiểm tra tính hợp lệ của cấu hình
func validateConfig(config *OPAConfig) error {
	if config.Query == "" {
		return fmt.Errorf("query không được rỗng")
	}

	if config.URL == "" {
		return fmt.Errorf("URL của OPA server không được rỗng")
	}

	if config.InputCreationMethod == nil {
		return fmt.Errorf("InputCreationMethod không được nil")
	}

	return nil
}

// Use trả về gin.HandlerFunc để sử dụng làm middleware
func (m *OPAMiddleware) Use() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tạo input data từ gin.Context
		input, err := m.config.InputCreationMethod(c)
		if err != nil {
			m.handleError(c, fmt.Errorf("không thể tạo input data: %w", err))
			return
		}

		// Đánh giá policy từ remote OPA server
		result, err := m.evaluateRemotePolicy(input)
		if err != nil {
			m.handleError(c, fmt.Errorf("không thể đánh giá policy: %w", err))
			return
		}

		// Kiểm tra kết quả với expected result
		if !m.compareResults(result, m.config.ExceptedResult) {
			m.handleDenied(c, result)
			return
		}

		// Policy cho phép, tiếp tục xử lý request
		c.Next()
	}
}

// evaluateRemotePolicy đánh giá policy sử dụng remote OPA server
func (m *OPAMiddleware) evaluateRemotePolicy(input map[string]interface{}) (interface{}, error) {
	// Chuẩn bị request body
	requestBody := map[string]interface{}{
		"input": input,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("không thể marshal request body: %w", err)
	}

	// Tạo URL cho OPA API
	url := strings.TrimSuffix(m.config.URL, "/") + "/v1/data/" + strings.TrimPrefix(m.config.Query, "data.")

	// Tạo HTTP request
	req, err := http.NewRequestWithContext(m.config.Context, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("không thể tạo HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Gửi request
	resp, err := m.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("không thể gửi request đến OPA server: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OPA server trả về status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("không thể parse JSON response: %w", err)
	}

	// Lấy kết quả từ response
	result, exists := response["result"]
	if !exists {
		return false, nil
	}

	return result, nil
}

// compareResults so sánh kết quả với expected result
func (m *OPAMiddleware) compareResults(actual, expected interface{}) bool {
	return actual == expected
}

// handleError xử lý lỗi và trả về response
func (m *OPAMiddleware) handleError(c *gin.Context, err error) {
	appError := datatype.ErrInternalServerError.WithDebug(err.Error())
	c.JSON(appError.StatusCode(), appError)
	c.Abort()
}

// handleDenied xử lý trường hợp request bị từ chối
func (m *OPAMiddleware) handleDenied(c *gin.Context, result interface{}) {
	appError := datatype.ErrForbidden.WithReason(m.config.DeniedMessage).WithDebug(fmt.Sprintf("Policy result: %v", result))
	c.JSON(appError.StatusCode(), appError)
	c.Abort()
}
