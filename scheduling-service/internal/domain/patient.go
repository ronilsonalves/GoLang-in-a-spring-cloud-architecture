package domain

type Patient struct {
	Id        int    `json:"id"`
	LastName  string `json:"lastName" binding:"required"`
	Name      string `json:"name" binding:"required"`
	RG        string `json:"rg" binding:"required"`
	CreatedAt string `json:"createdAt" binding:"required"`
}
