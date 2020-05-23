package controllers

// import "C"
import (
	"errors"
	"gohome_self/tracingbook/dtos"
	"gohome_self/tracingbook/middlewares"
	"gohome_self/tracingbook/models"
	"gohome_self/tracingbook/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.RouterGroup) {
	router.GET("/", BookList)
	router.GET("/:slug", GetBookDetailsBySlug)

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("/", CreateBook)
		router.DELETE("/:slug", BookDelete)
	}
}

func BookList(c *gin.Context) {

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	productModels, modelCount, commentsCount, err := services.FetchBooksPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid param")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreatedBookPagedResponse(c.Request, productModels, page, pageSize, modelCount, commentsCount))
}

func GetBookDetailsBySlug(c *gin.Context) {
	productSlug := c.Param("slug")

	product := services.FetchBookDetails(&models.Book{Slug: productSlug}, true)
	if product.ID == 0 {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid slug")))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateBookDetailsDto(product))
}

func CreateBook(c *gin.Context) {
	// Only admin users can create products
	user := c.Keys["currentUser"].(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you must be admin"))
		return
	}

	var formDto dtos.CreateBook
	if err := c.ShouldBind(&formDto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	name := formDto.Name
	description := formDto.Description

	author := formDto.Author
	latest := formDto.Latest
	form, err := c.MultipartForm()

	tagCount := 0
	catCount := 0
	for key := range form.Value {
		if strings.HasPrefix(key, "tags[") {
			tagCount++
		}
		if strings.HasPrefix(key, "category[") {
			catCount++
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateDetailedErrorDto("form_error", err))
		return
	}

	product := models.Book{
		Name:        name,
		Description: description,
		Author:      author,
		Latest:      latest,
	}

	if err := services.CreateOne(&product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateBookCreatedDto(product))

}

func BookDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := services.DeleteBook(&models.Book{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid slug")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": "Delete success"})
}
