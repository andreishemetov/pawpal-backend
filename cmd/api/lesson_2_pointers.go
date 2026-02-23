package main

import "fmt"

// Pet with visits (for pointers lesson — different from data.Pet)
type PetWithVisits struct {
	Name   string
	Age    int
	Visits int
}

func (p PetWithVisits) Speak() string {
	return "My name is " + p.Name
}

func (p PetWithVisits) RenameAsCopy(newName string) {
	p.Name = newName
}

func (p *PetWithVisits) RenameAsPointer(newName string) {
	p.Name = newName
}

func (p *PetWithVisits) AddVisit() {
	p.Visits++
}

func (p PetWithVisits) String() string {
	return fmt.Sprintf("%s (%d years) - Visits: %d", p.Name, p.Age, p.Visits)
}

func lesson2() {
	fmt.Println("Lesson 2 starting...")

	pet := PetWithVisits{
		Name:   "Charlie",
		Age:    3,
		Visits: 0,
	}

	pet.RenameAsCopy("Milo")
	fmt.Println(pet.Speak())

	pet.RenameAsPointer("Max")
	fmt.Println(pet.Speak())

	pet.AddVisit()
	fmt.Println(pet.Visits)

	pet.AddVisit()
	fmt.Println(pet.Visits)

	pet.AddVisit()
	fmt.Println(pet.Visits)
}
