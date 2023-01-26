package domain

type Dentist struct {
	Id       int    `json:"id"`
	LastName string `json:"lastName" binding:"required"`
	Name     string `json:"name" binding:"required"`
	CRO      string `json:"cro" binding:"required"`
}
