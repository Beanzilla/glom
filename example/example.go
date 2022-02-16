package main

import (
	"fmt"
	"github.com/beanzilla/glom"
)

func main() {

	// An example structure (Shown using all maps, but slice/array and even structures are supported)
	critters := make(map[string]interface{})
	cat := make(map[string]interface{})
	cat["name"] = "Cat"
	cat["sounds"] = "Meow"
	cat["food"] = "Fish"
	critters["Cat"] = cat
	dog := make(map[string]interface{})
	dog["name"] = "Dog"
	dog["sounds"] = "Woof"
	dog["food"] = "Anything"
	critters["Dog"] = dog

	test := make(map[string]interface{})
	test["Animals"] = critters

	/* In JSON it would be represented as:
	{
		"Animals": {
			"Cat": {
				"name": "Cat",
				"sounds": "Meow",
				"food": "Fish"
			},
			"Dog": {
				"name": "Dog",
				"sounds": "Woof",
				"food": "Anything"
			}
		}
	}
	Where accessing name would be
	data["Animals"]["Cat"]["name"] or data["Animals"]["Dog"]["name"]

	But with glom it's
	"Animals.Cat.name" or "Animals.Dog.name"
	*/

	// An example of accessing something that doesn't exist
	_, err := glom.Glom(test, "Animals.Dog.hates")
	if err != nil {
		fmt.Println(err) // Failed moving to 'hates' from path of 'Animals.Dog', options are 'name', 'sounds', 'food' (3)
	}

	// An example of successfully geting something
	value, err := glom.Glom(test, "Animals.Cat.sounds")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Cat's make '%v' sounds.\r\n", value)
		// Note, value is of interface type
	}
}
