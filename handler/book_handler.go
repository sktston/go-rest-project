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


// CreateBook godoc
// @Summary Create Book
// @Tags books
// @Accept json
// @Produce json
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
	if err := copier.Copy(&book, bookDTO); err != nil {
		log.Error().Err(err).Msg("")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if err := repository.CreateBook(&book); err != nil {
		log.Error().Err(err).Msg("ERROR creating book data")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTO model.BookResponseDTO
		if err := copier.Copy(&bookResponseDTO, book); err != nil {
			log.Error().Err(err).Msg("")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

// GetBooks godoc
// @Summary Get all books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {object} []model.BookResponseDTO
// @Router /books [get]
func GetBooks(c *gin.Context) {
	var books []entity.Book
	if err := repository.GetAllBooks(&books); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTOs []model.BookResponseDTO
		if err := copier.Copy(&bookResponseDTOs, books); err != nil {
			log.Error().Err(err).Msg("")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, bookResponseDTOs)
	}
}

// GetBookByID godoc
// @Summary Get the book by id
// @Tags books
// @Accept json
// @Produce json
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
		if err := copier.Copy(&bookResponseDTO, book); err != nil {
			log.Error().Err(err).Msg("")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}

// UpdateBook godoc
// @Summary Update the book information
// @Tags books
// @Accept json
// @Produce json
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

	if err := copier.Copy(&book, bookRequestDTO); err != nil {
		log.Error().Err(err).Msg("")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if err := repository.UpdateBook(&book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		var bookResponseDTO model.BookResponseDTO
		if err := copier.Copy(&bookResponseDTO, book); err != nil {
			log.Error().Err(err).Msg("")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, bookResponseDTO)
	}
}


// DeleteBook godoc
// @Summary Delete the book
// @Tags books
// @Accept json
// @Produce json
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
