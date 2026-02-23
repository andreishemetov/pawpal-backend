package data

import "fmt"

type Pet struct {
	ID   int	`json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}

func (p Pet) String() string {
	return fmt.Sprintf("%d - %s (%s)", p.ID, p.Name, p.Type)
}

func (p Pet) Speak() string {
	return "My name is " + p.Name
}

func (p Pet) SpeakInLanguage(language string) string {
	if language == "en" {
		return "My name is " + p.Name
	}
	if language == "es" {
		return "Me llamo " + p.Name
	}
	return p.Name
}
