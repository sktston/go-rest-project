package handler

import (
	"github.com/airoasis/go-rest-project/model"
	"github.com/airoasis/go-rest-project/model/entity"
	"github.com/airoasis/go-rest-project/repository"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

//GetBooks ... Get all books
func GetBooks(c *gin.Context) {
	var books []entity.Book
	if err := repository.GetAllBooks(&books); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var bookResponseDTOs []model.BookResponseDTO
		copier.Copy(&bookResponseDTOs, books)
		c.JSON(http.StatusOK, bookResponseDTOs)
	}
}

//CreateBook ... Create Book
func CreateBook(c *gin.Context) {
	var bookDTO model.BookRequestDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var book entity.Book
	copier.Copy(&book, bookDTO)
	if err := repository.CreateBook(&book); err != nil {
		log.Error().Err(err).Msg("ERROR creating book data")
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var bookResponseDTO model.BookResponseDTO
		copier.Copy(&bookResponseDTO, book)
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

//GetBookByID ... Get the book by id
func GetBookByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book entity.Book
	if err := repository.GetBookByID(&book, id); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var bookResponseDTO model.BookResponseDTO
		copier.Copy(&bookResponseDTO, book)
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

//UpdateBook ... Update the book information
func UpdateBook(c *gin.Context) {
	var book entity.Book
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.GetBookByID(&book, id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var bookRequestDTO model.BookRequestDTO
	if err := c.ShouldBindJSON(&bookRequestDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&book, bookRequestDTO)
	if err := repository.UpdateBook(&book); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var bookResponseDTO model.BookResponseDTO
		copier.Copy(&bookResponseDTO, book)
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

//DeleteBook ... Delete the book
func DeleteBook(c *gin.Context) {
	var book entity.Book
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.DeleteBook(&book, id); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.Status(http.StatusOK)
	}
}
