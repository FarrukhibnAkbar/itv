package entity

type Movie struct {
	Id          string `json:"id" gorm:"column:id"`
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	Director    string `json:"director" gorm:"column:director"`
	Year        int32  `json:"year" gorm:"column:year"`
	Plot        string `json:"plot" gorm:"column:plot"`
}

type MovieList struct {
	Movies []Movie `json:"movies"`
	Count  int32   `json:"count"`
}
