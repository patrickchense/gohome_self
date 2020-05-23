package dtos

import (
	"gohome_self/tracingbook/models"
	"net/http"
	"time"
)

type ManagedModel models.Book

type CreateBook struct {
	Name        string `form:"name" json:"name" xml:"name" binding:"required"`
	Description string `form:"description" json:"description" xml:"description" binding:"required"`
	Author      string `form:"author" json:"price" xml:"season" binding:"required"`
	Latest      string `form:"latest" json:"stock" xml:"episode" binding:"required"`
}

func CreatedBookPagedResponse(request *http.Request, books []models.Book, page, page_size, count int, commentsCount []int) interface{} {
	var resources = make([]interface{}, len(books))
	for index, book := range books {
		resources[index] = CreateBookDto(&book, commentsCount[index])
	}
	return CreatePagedResponse(request, resources, "books", page, page_size, count)
}

func CreateBookListDto(request *http.Request, books []models.Book) map[string]interface{} {
	var resources = make([]interface{}, len(books))
	for index, book := range books {
		resources[index] = CreateBookDto(&book, 0)
	}
	return CreateResponse(request, resources, "books")
}

func CreateBookDto(book *models.Book, commentCount int) map[string]interface{} {

	result := map[string]interface{}{
		"id":         book.ID,
		"name":       book.Name,
		"slug":       book.Slug,
		"author":     book.Author,
		"latest":     book.Latest,
		"created_at": book.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		"updated_at": book.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}

	if commentCount >= 0 {
		// "comments_count": book.CommentsCount,
		result["comments_count"] = commentCount
	}
	return result
}

func CreateBookDetailsDto(book models.Book) map[string]interface{} {
	result := CreateBookDto(&book, -1)
	result["description"] = book.Description
	return result
}
func CreateBookCreatedDto(book models.Book) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateBookDetailsDto(book), "Book crated successfully")
}
