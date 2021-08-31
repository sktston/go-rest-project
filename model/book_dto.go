package model

type BookRequestDTO struct {
	Title    	string	`json:"title" binding:"required,min=5" example:"The Three Pigs"` // book title
	Author   	string	`json:"author" binding:"required,min=5" example:"David Wiesner"` // book author
	Publisher 	string	`json:"publisher" binding:"required,min=5" example:"Clarion Books"` // book publisher
}

type BookResponseDTO struct {
	ID			uint	`json:"id" example:"1234"` // book id
	Title    	string	`json:"title" example:"The Three Pigs"` // book title
	Author   	string	`json:"author" example:"David Wiesner"` // book author
	Publisher 	string	`json:"publisher" example:"Clarion Books"` // book publisher
}
