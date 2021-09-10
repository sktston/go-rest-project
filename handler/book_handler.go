package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/sktston/go-rest-project/model"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/sktston/go-rest-project/repository"
	"net/http"
	"strconv"
)

// GetBooks godoc
// @Summary Get all books
// @Tags books
// @Success 200 {object} []model.BookResponseDTO
// @Router /books [get]
func GetBooks(c *gin.Context) {
	var books []entity.Book
	if err := repository.GetAllBooks(&books); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTOs []model.BookResponseDTO
		copier.Copy(&bookResponseDTOs, books)
		c.JSON(http.StatusOK, bookResponseDTOs)
	}
}

// CreateBook godoc
// @Summary Create Book
// @Tags books
// @Param body body model.BookRequestDTO false "body"
// @Success 200 {object} model.BookResponseDTO
// @Router /books [post]
func CreateBook(c *gin.Context) {
	var bookDTO model.BookRequestDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var book entity.Book
	copier.Copy(&book, bookDTO)
	if err := repository.CreateBook(&book); err != nil {
		log.Error().Err(err).Msg("ERROR creating book data")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTO model.BookResponseDTO
		copier.Copy(&bookResponseDTO, book)
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

// GetBookByID godoc
// @Summary Get the book by id
// @Tags books
// @Param id path string true "id"
// @Success 200 {object} model.BookResponseDTO
// @Router /books/{id} [get]
func GetBookByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book entity.Book
	if err := repository.GetBookByID(&book, id); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTO model.BookResponseDTO
		copier.Copy(&bookResponseDTO, book)
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

// UpdateBook godoc
// @Summary Update the book information
// @Tags books
// @Param id path string true "id"
// @Param body body model.BookRequestDTO false "body"
// @Success 200 {object} model.BookResponseDTO
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	var book entity.Book
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.GetBookByID(&book, id); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var bookRequestDTO model.BookRequestDTO
	if err := c.ShouldBindJSON(&bookRequestDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&book, bookRequestDTO)
	if err := repository.UpdateBook(&book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTO model.BookResponseDTO
		copier.Copy(&bookResponseDTO, book)
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}


// DeleteBook godoc
// @Summary Delete the book
// @Tags books
// @Param id path string true "id"
// @Success 200
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	var book entity.Book
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.DeleteBook(&book, id); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}
