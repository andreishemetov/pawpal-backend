package data

type PetsListResponse struct {
	Items []Pet `json:"items"`
	Meta  struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	} `json:"meta"`
}
