package main

import (
	"fmt"

	"github.com/andreishemetov/pawpal/internal/data"
)

func lesson1() {

	fmt.Println("Lesson 1 starting...")

	pet := data.Pet{
		ID:   1,
		Name: "Charlie",
		Type: "Dog",
	}

	fmt.Println("PawPal backend starting...", pet)
	fmt.Println(pet.Speak())
	fmt.Println(pet.SpeakInLanguage("es"))
	fmt.Println(SpeakInLanguage2(pet, "en"))

	pets := []data.Pet{
		{ID: 1, Name: "Charlie", Type: "Dog"},
		{ID: 2, Name: "Milo", Type: "Cat"},
		{ID: 3, Name: "Bella", Type: "Dog"},
		{ID: 4, Name: "Max", Type: "Cat"},
		{ID: 5, Name: "Lucy", Type: "Dog"},
		{ID: 6, Name: "Daisy", Type: "Cat"},
		{ID: 7, Name: "Rocky", Type: "Dog"},
		{ID: 8, Name: "Mia", Type: "Cat"},
		{ID: 9, Name: "Buddy", Type: "Dog"},
		{ID: 10, Name: "Luna", Type: "Cat"},
	}

	printPets(pets)
	fmt.Println("Total pets:", countPets(pets))
}

func SpeakInLanguage2(p data.Pet, language string) string {
	if language == "en" {
		return "My name is " + p.Name
	}
	if language == "es" {
		return "Me llamo " + p.Name
	}
	return p.Name
}

func printPets(pets []data.Pet) {
	for _, p := range pets {
		fmt.Println(p.Name)
	}
}

func countPets(pets []data.Pet) int {
	return len(pets)
}
