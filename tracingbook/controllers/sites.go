package controllers

import (
	"gohome_self/tracingbook/dtos"
	"gohome_self/tracingbook/middlewares"
	"gohome_self/tracingbook/models"
	"gohome_self/tracingbook/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterSiteRoutes(router *gin.RouterGroup) {
	router.POST("", CreateSite)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("", ListSites)
		router.GET("/:id", ShowSite)
	}
}

func ListSites(c *gin.Context) {
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
	userId := c.MustGet("currentUserId").(uint32)

	sites, totalCommentCount, err := services.FetchSitesPage(userId, page, pageSize)

	c.JSON(http.StatusOK, dtos.CreateSitePagedResponse(c.Request, sites, page, pageSize, totalCommentCount))
}

func ShowSite(c *gin.Context) {
	siteId, err := strconv.Atoi(c.Param("id"))
	user := c.MustGet("currentUser").(models.User)
	site, err := services.FetchSiteById(uint(siteId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if site.UserId == user.ID || user.IsAdmin() {
		c.JSON(http.StatusOK, dtos.CreateSiteDetailsDto(&site))
	} else {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you can not view this site"))
		return
	}
}

func CreateSite(c *gin.Context) {
	var siteRequest dtos.CreateSiteRequestDto
	if err := c.ShouldBind(&siteRequest); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	userObj, userLoggedIn := c.Get("currentUser")
	var user models.User
	if userLoggedIn {
		user = (userObj).(models.User)
	}

	site := models.Site{
		Url:        siteRequest.Url,
		Name:       siteRequest.Name,
		SiteStatus: 0,
	}

	if userLoggedIn {
		site.UserId = user.ID
	}

	c.JSON(http.StatusOK, dtos.CreateSiteCreatedDto(&site))

}
