package bookhttpgin

import (
	"net/http"
	"strconv"
	"time"

	bookmodel "fat2fast/ikv/modules/book/model"
	bookservice "fat2fast/ikv/modules/book/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
)

// ActionListBooks lấy danh sách books - GET /
func (c *BookHTTPController) ActionListBooks(ctx *gin.Context) {
	// Parse query parameters
	filter, err := c.parseListQueryParams(ctx)
	if err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug("Invalid query parameters"))
	}

	// Tạo query
	query := &bookservice.ListBooksQuery{Filter: filter}

	// Thực thi query
	response, err := c.listQryHdl.Execute(ctx.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	// Trả về response
	ctx.JSON(http.StatusOK, datatype.ResponseSuccess(response))
}

// parseListQueryParams parse các query parameters cho list API
func (c *BookHTTPController) parseListQueryParams(ctx *gin.Context) (*bookmodel.ListBookFilter, error) {
	// Default values
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Parse dates
	var createdFrom, createdTo time.Time
	if dateStr := ctx.Query("created_from"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			createdFrom = parsed
		}
	}
	if dateStr := ctx.Query("created_to"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			createdTo = parsed
		}
	}

	// Parse price range
	priceMin, _ := strconv.ParseFloat(ctx.Query("price_min"), 64)
	priceMax, _ := strconv.ParseFloat(ctx.Query("price_max"), 64)

	return &bookmodel.ListBookFilter{
		Page:        page,
		PerPage:     perPage,
		Status:      ctx.Query("status"),
		Search:      ctx.Query("search"),
		Author:      ctx.Query("author"),
		SortBy:      ctx.DefaultQuery("sort_by", "created_at"),
		SortOrder:   ctx.DefaultQuery("sort_order", "DESC"),
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		PriceMin:    priceMin,
		PriceMax:    priceMax,
	}, nil
}
