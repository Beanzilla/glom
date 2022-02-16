package glom

import (
	"fmt"
	"strings"
	"testing"
)

func TestGlomArray(t *testing.T) {
	var data []interface{}
	data = append(data, "Goose")

	test1 := make(map[string]interface{})

	test1a := make(map[string]interface{})
	test1a["name"] = "Ducky"
	test1a["age"] = 62
	test1a["race"] = "Duck"

	test1b := make(map[string]interface{})
	test1b["name"] = "Sir Meow"
	test1b["age"] = 12
	test1b["race"] = "Cat"

	var animals []interface{}
	animals = append(animals, test1a)
	animals = append(animals, test1b)

	test1["animals"] = animals
	data = append(data, test1)

	/*
		HJSON Representation of data
		data = [
			"Goose"
			{
				"animals": [
					{"name": "Ducky", "age": 62, "race": "Duck"}
					{"name": "Sir Meow", "age": 12, "race": "Cat"}
				]
			}
		]
	*/
	test, err := Glom(data, "1.animals.1.name")
	if err != nil {
		t.Errorf("Unexpected Error: \"%v\"", err)
	} else if test != "Sir Meow" {
		t.Errorf("Failed getting 'Sir Meow' got \"%v\"", test)
	}
}

func TestGetPossible(t *testing.T) {
	var data []string
	data = append(data, "One")
	data = append(data, "Two")
	data = append(data, "Three")
	data = append(data, "Four")

	result := GetPossible(data)
	if len(result) != len(data) {
		t.Errorf("Expected even size, %d != %d", len(result), len(data))
	}
}

func TestStruct(t *testing.T) {
	type Animal struct {
		Name     string
		Lifespan int
	}

	cat := Animal{"Cat", 12}
	dog := Animal{"Dog", 13}

	var data []Animal
	data = append(data, cat)
	data = append(data, dog)

	test, err := Glom(data, "1.*")
	if err != nil {
		t.Errorf("TestStruct 1/3: Unexpected Error: \"%v\"", err)
	} else {
		if fmt.Sprintf("%v", test) != fmt.Sprintf("%v", dog) {
			t.Errorf("TestStruct 1/3: Failed getting '%v' got '%v'", dog, test)
		}
	}

	test2, err2 := Glom(cat, "Lifespan")
	if err2 != nil {
		t.Errorf("TestStruct 2/3: Unexpected Error: \"%v\"", err2)
	} else {
		if test2 != cat.Lifespan {
			t.Errorf("TestStruct 2/3: Failed getting '%v' got '%v'", cat.Lifespan, test2)
		}
	}

	data2 := make(map[string]Animal)

	data2["Squirrel"] = Animal{"Squirrel", 999}
	data2["Hamster"] = Animal{"Hamster", 4}

	test3, err3 := Glom(data2, "Squirrel.Name")
	if err3 != nil {
		t.Errorf("TestStruct 3/3: Unexpected Error: \"%v\"", err3)
	} else {
		if test3 != "Squirrel" {
			t.Errorf("TestStruct 3/3: Failed getting 'Squirrel' got '%v'", test3)
		}
	}
}

func TestListPossible(t *testing.T) {
	var list []string
	list = append(list, "One")
	list = append(list, "Two")
	list = append(list, "Three")

	result := list_possible(list)

	if strings.Join(result, ", ") != "'One', 'Two', 'Three'" {
		t.Errorf("Failed getting \"%s\" got \"%v\"", "'One', 'Two', 'Three'", strings.Join(result, ", "))
	}
}

func TestFail(t *testing.T) {
	data := make(map[string]interface{})
	data["Duck"] = "Quack"
	data["Cheese"] = 3
	data["Mouse"] = true

	test, err := Glom(data, "Moose")
	if err == nil {
		t.Errorf("Expected Error, got '%v'", test)
	}
}

func TestMapToInter(t *testing.T) {
	m := make(map[string]string)
	m["Duck"] = "Quack"
	m["Cheese"] = "Yes Please!"
	m["Mouse"] = "true"
	var s []string
	s = append(s, "Duck")
	s = append(s, "Cheese")
	s = append(s, "Mouse")
	var m2 map[string]int

	_, err1 := mapToInterface(m)
	if err1 != nil {
		t.Errorf("Unexpected Error given map: %v", err1)
	}

	test2, err2 := mapToInterface(s)
	if err2 == nil {
		t.Errorf("Expected Error given slice, got '%v'", test2)
	}

	test3, err3 := mapToInterface(m2)
	if err3 == nil {
		t.Errorf("Expected Error given invalid/empty map, got '%v'", test3)
	}
}

func TestSliceToInter(t *testing.T) {
	m := make(map[string]string)
	m["Duck"] = "Quack"
	m["Cheese"] = "Yes Please!"
	m["Mouse"] = "true"
	var s []string
	s = append(s, "Duck")
	s = append(s, "Cheese")
	s = append(s, "Mouse")
	var m2 map[string]int
	var s2 []int

	test1, err1 := sliceToInterface(m)
	if err1 == nil {
		t.Errorf("Expected Error given map, got '%v'", test1)
	}

	_, err2 := sliceToInterface(s)
	if err2 != nil {
		t.Errorf("Unexpected Error given slice: %v", err2)
	}

	test3, err3 := sliceToInterface(m2)
	if err3 == nil {
		t.Errorf("Expected Error given invalid/empty map, got '%v'", test3)
	}

	test4, err4 := sliceToInterface(s2)
	if err4 == nil {
		t.Errorf("Expected Error given invalid/empty slice, got '%v'", test4)
	}
}

func TestEdgeCasesMapNextLvl(t *testing.T) {
	// This doesn't work, I thought it would but it does not
	var m map[string]int
	m2 := make(map[string]int)
	m2["Cheese"] = 6
	m2["C"] = 1
	m2["h"] = 1
	m2["e"] = 3
	m2["s"] = 1

	test1, err1 := next_level(m, "failwhale")
	if err1 == nil {
		t.Errorf("Expected Error given invalid/empty map, got '%v'", test1)
	}

	test2, err2 := next_level(m2, "n")
	if err2 == nil {
		t.Errorf("Expected Error given map but invalid key, got '%v'", test2)
	}
}

func TestEdgeCasesGlom(t *testing.T) {
	// This is just a generic test, nothing fancy
	data := make(map[string]interface{})

	lvl2 := make(map[string]interface{})
	lvl2["Duck"] = "Quack"
	lvl2["Cheese"] = 6
	lvl2["Mouse"] = true
	data["part1"] = lvl2

	var lvl1 []interface{}
	lvl1 = append(lvl1, "Pig")
	lvl1 = append(lvl1, "Chicken")
	lvl1 = append(lvl1, "Cow")
	lvl1 = append(lvl1, "Dog")
	lvl1 = append(lvl1, "Cat")
	lvl1 = append(lvl1, "Horse")
	lvl1 = append(lvl1, true)
	lvl1 = append(lvl1, 42)
	data["part2"] = lvl1

	_, err1 := Glom(data, "part1.Mouse")
	if err1 != nil {
		t.Errorf("Unexpected Error (part1.Mouse = true): %v", err1)
	}

	_, err2 := Glom(data, "part2.3")
	if err2 != nil {
		t.Errorf("Unexpected Error (part2.3 = 'Dog'): %v", err2)
	}
}

func TestTypeConversions(t *testing.T) {
	// This is just a generic test, nothing fancy
	data := make(map[string]interface{})

	lvl2 := make(map[string]interface{})
	lvl2["Duck"] = "Quack"
	lvl2["Cheese"] = 6
	lvl2["Mouse"] = true
	lvl2["Gravity"] = 9.81
	data["part1"] = lvl2

	var lvl1 []interface{}
	lvl1 = append(lvl1, "Pig")
	lvl1 = append(lvl1, "Chicken")
	lvl1 = append(lvl1, "Cow")
	lvl1 = append(lvl1, "Dog")
	lvl1 = append(lvl1, "Cat")
	lvl1 = append(lvl1, "Horse")
	lvl1 = append(lvl1, true)
	lvl1 = append(lvl1, 42)
	data["part2"] = lvl1

	// Part 1/6: Attempt to convert part2 a slice into a string (invalid)
	p1, err := Glom(data, "part2")
	if err != nil {
		t.Errorf("Unexpected error, expected part2 []interface{}, got %v", err)
	} else {
		// Now the real test
		_, err := String(p1)
		if err == nil {
			t.Errorf("Expected error, got a value")
		}
	}

	// Part 2/6: Converting part2.1 == "Chicken" to string
	p2, err := Glom(data, "part2.1")
	if err != nil {
		t.Errorf("Unexpected error, expected part2.1 interface{}, got %v", err)
	} else {
		// Note, p2 is still an interface, let's test converting it to a string
		d1, err := String(p2)
		if err != nil {
			t.Errorf("Unexpected error, expected to convert interface{} to string, got %v", err)
		} else {
			// Compare
			if d1 != "Chicken" {
				t.Errorf("Failed to convert interface{} to string, expected 'Chicken' got '%s'", d1)
			}
		}
	}

	// Part 3/6: Attempt to convert part1 a map[string]interface{} into a int (invalid)
	p3, err := Glom(data, "part1")
	if err != nil {
		t.Errorf("Unexpected error, expected part1 map[string]interface{}, got %v", err)
	} else {
		// Now the real test, int convert
		_, err := Int(p3)
		if err == nil {
			t.Errorf("Expected error, got a value")
		}
	}

	// Part 4/6: Convert part1.Cheese == 6 to int
	p4, err := Glom(data, "part1.Cheese")
	if err != nil {
		t.Errorf("Unexpected error, expected part1.Cheese interface{}, got %v", err)
	} else {
		// Note, p4 is still an interface, let's test converting it to a int
		d2, err := Int(p4)
		if err != nil {
			t.Errorf("Unexpected error, expected to convert interface{} to int, got %v", err)
		} else {
			// Compare
			if d2 != 6 {
				t.Errorf("Failed to convert interface{} to int, expected 6 got %d", d2)
			}
		}
	}

	// Part 5/6: Attempt to convert part2 a []interface{} into a float64 (invalid)
	p5, err := Glom(data, "part2")
	if err != nil {
		t.Errorf("Unexpected error, expected part2 []interface{}, got %v", err)
	} else {
		// Now the test, float64 convert
		_, err := Float64(p5)
		if err == nil {
			t.Errorf("Expected error, got a value")
		}
	}

	// Part 6/6: Convert part1.Gravity == 9.81 to float64
	p6, err := Glom(data, "part1.Gravity")
	if err != nil {
		t.Errorf("Unexpected error, expected part1.Gravity interface{}, got %v", err)
	} else {
		// Note, p6 is still an interface, let's test converting it to a float64
		d3, err := Float64(p6)
		if err != nil {
			t.Errorf("Unexpected error, expected to convert interface{} to float64, got %v", err)
		} else {
			// Compare
			if d3 != 9.81 {
				t.Errorf("Failed to convert interface{} to float64, expected 9.81 got %f", d3)
				// Something is wrong with the Earth's gravitational field!
				// ...
				// Or this test is wrong. :)
			}
		}
	}
}
