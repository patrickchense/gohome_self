package services

import (
	"gohome_self/tracingbook/infrastructure"
	"gohome_self/tracingbook/models"
)

func FetchBooksPage(page int, page_size int) ([]models.Book, int, []int, error) {
	database := infrastructure.GetDb()
	var books []models.Book
	var count int
	tx := database.Begin()
	database.Model(&books).Count(&count)
	database.Offset((page - 1) * page_size).Limit(page_size).Find(&books)
	tx.Model(&books).Order("created_at desc").Offset((page - 1) * page_size).Limit(page_size).Find(&books)
	commentsCount := make([]int, len(books))

	for index, book := range books {
		commentsCount[index] = tx.Model(&book).Association("Comments").Count()
	}
	err := tx.Commit().Error
	return books, count, commentsCount, err
}

func FetchBooks() ([]models.Book, error) {
	database := infrastructure.GetDb()
	var books []models.Book
	var count int
	tx := database.Begin()
	database.Model(&books).Count(&count)
	database.Find(&books)
	tx.Model(&books).Order("created_at desc").Find(&books)
	err := tx.Commit().Error
	return books, err
}

func FetchBookDetails(condition interface{}, optional ...bool) models.Book {
	database := infrastructure.GetDb()
	var book models.Book

	query := database.Where(condition).
		Preload("Tags").Preload("Categories").Preload("Images").Preload("Comments")
	// Unfortunately .Preload("Comments.User") does not work as the doc states ...
	query.First(&book)

	return book
}

func FetchBookId(slug string) (uint, error) {
	bookId := -1
	database := infrastructure.GetDb()
	err := database.Model(&models.Book{}).Where(&models.Book{Slug: slug}).Select("id").Row().Scan(&bookId)
	return uint(bookId), err
}

func Update(book *models.Book, data interface{}) error {
	database := infrastructure.GetDb()
	err := database.Model(book).Update(data).Error
	return err
}

func DeleteBook(condition interface{}) error {
	db := infrastructure.GetDb()
	err := db.Where(condition).Delete(models.Book{}).Error
	return err
}

func FetchBooksIdNameAndPrice(bookIds []uint) (books []models.Book, err error) {
	database := infrastructure.GetDb()
	err = database.Select([]string{"id", "name", "slug", "price"}).Find(&books, bookIds).Error
	return books, err
}
