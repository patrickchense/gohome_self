package services

import (
	"gohome_self/tracingbook/infrastructure"
	"gohome_self/tracingbook/models"
)

func FetchSitesPage(userId uint32, page, pageSize int) (sites []models.Site, totalSitesCount int, err error) {
	database := infrastructure.GetDb()

	totalSitesCount = 0

	query := database.Model(&models.Site{}).Where(&models.Site{UserId: userId})
	query.Count(&totalSitesCount)

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&sites).Error
	if err != nil {
		return
	}

	var siteIds = make([]uint, len(sites))
	for i := 0; i < len(sites); i++ {
		siteIds[i] = sites[i].ID
	}

	return sites, totalSitesCount, err
}

func FetchSiteById(siteId uint) (site models.Site, err error) {
	database := infrastructure.GetDb()
	err = database.Model(models.Site{}).First(&site, siteId).Error
	return site, err
}

func FetchSiteDetails(siteId uint) (updateItems []models.UpdateItem, err error) {
	database := infrastructure.GetDb()
	err = database.Model(models.UpdateItem{}).Find(&updateItems, siteId).Error
	return updateItems, err
}
