package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"galvin/golang-gin/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
}

type GetProductsBySlugV1Param struct {
	Slug string `uri:"slug" binding:"slug,min=3,max=5"`
}

type GetProductsV1Param struct {
	Search string `form:"search" binding:"required,min=3,max=50,search"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Email  string `form:"email" binding:"omitempty,email"`
	Date   string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

type ProductImage struct {
	ImageName string `json:"image_name" binding:"required"`
	ImageLink string `json:"image_link" binding:"required,file_ext=jpg png gif"`
}

type ProductAttribute struct {
	AttributeName  string `json:"attribute_name" binding:"required"`
	AttributeValue string `json:"attribute_value" binding:"required"`
}

type ProductInfo struct {
	InfoKey   string `json:"info_key" binding:"required"`
	InfoValue string `json:"info_value" binding:"required"`
}

type PostProductsV1Param struct {
	Name             string                 `json:"name" binding:"required,min=3,max=100"`
	Price            int                    `json:"price" binding:"required,min_int=100000"`
	Display          *bool                  `json:"display" binding:"omitempty"`
	ProductImage     ProductImage           `json:"product_image" binding:"required"`
	Tags             []string               `json:"tags" binding:"required,gt=3,lt=5"`
	ProductAttribute []ProductAttribute     `json:"product_attribute" binding:"required,gt=0,dive"`
	ProductInfo      map[string]ProductInfo `json:"product_info" binding:"required,gt=0,dive"`
	ProductMetadata  map[string]any         `json:"product_metadata" binding:"omitempty"`
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (u *ProductHandler) GetProductsV1(ctx *gin.Context) {
	// var params GetProductsV1Param
	// if err := ctx.ShouldBindQuery(&params); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
	// 	return
	// }

	// if params.Limit == 0 {
	// 	params.Limit = 1
	// }

	// if params.Email == "" {
	// 	params.Email = "No Email"
	// }

	// if params.Date == "" {
	// 	params.Date = time.Now().Format("2006-01-02")
	// }

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "List all products (V1)",
	// 	"search":  params.Search,
	// 	"limit":   params.Limit,
	// 	"email":   params.Email,
	// 	"date":    params.Date,
	// })
	searchRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	search := ctx.Query("search")

	if err := utils.ValidationRequired("Search", search); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationStringLength("Search", search, 3, 50); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationRegex("Search", search, searchRegex, "Invalid search format"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be between 1 and 100"})
		return
	}
	
	email := ctx.DefaultQuery("email", "No Email")
	if email != "No Email" && !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	date := ctx.DefaultQuery("date", time.Now().Format("2006-01-02"))
	if _, err := time.Parse("2006-01-02", date); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
 
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List all products (V1)",
		"search":  search,
		"limit":   limit,
		"email":   email,
		"date":    date,
	})
}

func (u *ProductHandler) GetProductsBySlugV1(ctx *gin.Context) {
	// var params GetProductsBySlugV1Param
	// if err := ctx.ShouldBindUri(&params); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
	// 	return
	// }

	// ctx.JSON(http.StatusCreated, gin.H{
	// 	"message": "Get product by Slug (V1)",
	// 	"slug":    params.Slug,
	// })
	slug := ctx.Param("slug")
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:[-.] [a-z0-9]+)*$`)

	if err := utils.ValidationRequired("Slug", slug); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationStringLength("Slug", slug, 3, 5); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationRegex("Slug", slug, slugRegex, "Invalid slug format"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get product by Slug (V1)",
		"slug":    slug,
	})
}

func (u *ProductHandler) PostProductsV1(ctx *gin.Context) {
	var params PostProductsV1Param
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	for key := range params.ProductInfo {
		if _, err := uuid.Parse(key); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": gin.H{
					"product_info": fmt.Sprintf("Key '%s' trong product_info không phải là UUUID hợp lệ", key),
				},
			})

			return
		}
	}

	if params.Display == nil {
		defaultDisplay := true
		params.Display = &defaultDisplay
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":           "Create product (V1)",
		"name":              params.Name,
		"price":             params.Price,
		"display":           params.Display,
		"product_image":     params.ProductImage,
		"tags":              params.Tags,
		"product_attribute": params.ProductAttribute,
		"product_info":      params.ProductInfo,
		"product_metadata":  params.ProductMetadata,
	})
}

func (u *ProductHandler) PutProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Update product (V1)"})
}

func (u *ProductHandler) DeleteProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Delete product (V1)"})
}