package main

import (
	"encoding/json"
	"fmt"

	"github.com/andreishemetov/pawpal/internal/data"
)

func lesson3() {

	fmt.Println("Lesson 3 starting...")

	pet := data.Pet{
		ID:   1,
		Name: "Charlie",
		Type: "Dog",
		Age:  3,
	}

	jsonBytes, err := json.Marshal(pet)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))

	jsonInput := `{"id":2,"name":"Milo","age":2,"type":"pekinese"}`

	var newPet data.Pet

	err = json.Unmarshal([]byte(jsonInput), &newPet)
	if err != nil {
		panic(err)
	}

	fmt.Println(newPet)
}
