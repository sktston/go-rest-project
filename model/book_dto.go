package model

type BookRequestDTO struct {
	Title    	string	`json:"title" binding:"required,min=5""`
	Author   	string	`json:"author" binding:"required,min=5""`
	Publisher 	string	`json:"publisher" binding:"required,min=5""`
}

type BookResponseDTO struct {
	ID			uint	`json:"id"`
	Title    	string	`json:"title"`
	Author   	string	`json:"author"`
	Publisher 	string	`json:"publisher"`
}
