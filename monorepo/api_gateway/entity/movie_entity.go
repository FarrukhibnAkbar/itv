package entity

type Movie struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
	Year        int32  `json:"year"`
	Plot        string `json:"plot"`
}

type MovieList struct {
	Movies []Movie `json:"movies"`
	Count  int32   `json:"count"`
}
