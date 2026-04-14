package data

type PetsListMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PetsListResponse struct {
	Items []Pet        `json:"items"`
	Meta  PetsListMeta `json:"meta"`
}
