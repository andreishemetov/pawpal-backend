package data

type PetQuery struct {
	Page  int
	Limit int
	Type  string // filter
	Q     string // name search
	Sort  string // id|name|age|visits
}
