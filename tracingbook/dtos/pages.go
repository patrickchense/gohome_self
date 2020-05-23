package dtos

import (
	"gohome_self/tracingbook/models"
	"net/http"
)

func CreateHomeResponse(request *http.Request, tags []models.Book) map[string]interface{} {
	return CreateSuccessDto(map[string]interface{}{
		"books": CreateBookListDto(request, tags),
	})
}
