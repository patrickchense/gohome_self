package dtos

import (
	"gohome_self/tracingbook/models"
	"net/http"
)

type CreateSiteRequestDto struct {
	Name string `form:"name" json:"name" xml:"name"`
	Url  string `form:"url" json:"url" xml:"url"`
}

func CreateSitePagedResponse(request *http.Request, sites []models.Site, page, page_size, totalSitesCount int, includes ...bool) map[string]interface{} {
	var resources = make([]interface{}, len(sites))
	for index, site := range sites {

		includeUser := getIncludeFlags(includes...)

		resources[index] = CreateSiteDto(&site, includeUser)
	}
	return CreatePagedResponse(request, resources, "sites", page, page_size, totalSitesCount)
}

func CreateSiteDto(site *models.Site, includes ...bool) map[string]interface{} {

	includeUser := getIncludeFlags(includes...)

	result := map[string]interface{}{
		"id":          site.ID,
		"url":         site.Url,
		"site_status": site.GetSiteStatusAsString(),
	}

	if includeUser {
		result["user"] = map[string]interface{}{
			"id": site.UserId,
		}
	}

	return CreateSuccessDto(result)
}

func CreateSiteDetailsDto(site *models.Site) map[string]interface{} {
	// includeUser -> false
	// includeSiteItems -> true
	// includeUser -> false
	return CreateSuccessDto(CreateSiteDto(site, true, true, false))
}

func getIncludeFlags(includes ...bool) bool {
	includeUser := false
	if len(includes) > 1 {
		includeUser = includes[1]
	}
	return includeUser
}

func CreateSiteCreatedDto(site *models.Site) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateSiteDetailsDto(site), "Site created successfully")
}
