package glom

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

// Based on sliceToInterface
// Converts a map[something]something to map[string]interface
func mapToInterface(data interface{}) (map[string]interface{}, error) {
	mapV := reflect.ValueOf(data)
	if mapV.Kind() != reflect.Map {
		return nil, fmt.Errorf("failed to convert %v, given %v type to map[string]interface{}", mapV, reflect.TypeOf(data))
	}
	if mapV.IsNil() || !mapV.IsValid() {
		return nil, fmt.Errorf("given nil or empty map")
	}

	result := make(map[string]interface{})
	keys := mapV.MapKeys()
	for k := range keys {
		//fmt.Printf("%d/%d = %v", k, len(keys), mapV.MapIndex(keys[k]))
		result[keys[k].String()] = mapV.MapIndex(keys[k]).Interface()
	}

	return result, nil
}

// Converts the slice of something into a slice of interface
// https://gist.github.com/heri16/077282d46ae95d48d430a90fb6accdff
// I only need the length
func sliceToInterface(data interface{}) ([]interface{}, error) {
	sliceV := reflect.ValueOf(data)
	if sliceV.Kind() == reflect.Slice { // Prevent us from converting an interface to interface
		switch data.(type) {
		case []interface{}:
			return data.([]interface{}), nil
		}
	}
	if sliceV.Kind() != reflect.Slice && sliceV.Kind() != reflect.Array {
		return nil, fmt.Errorf("failed to convert %v, given %v type to []interface{}", sliceV, reflect.TypeOf(data))
	}
	if sliceV.IsNil() || !sliceV.IsValid() {
		return nil, fmt.Errorf("given nil or empty slice")
	}

	length := sliceV.Len()
	result := make([]interface{}, length)

	for i := 0; i < length; i++ {
		//fmt.Printf("%d/%d = %v\r\n", i, length-1, sliceV.Index(i))
		result[i] = sliceV.Index(i).Interface()
	}

	return result, nil
}

// Returns a array/slice of possible choices (e.g. key names, indexs, struct fields made Public)
func GetPossible(data interface{}) []string {
	var result []string
	//fmt.Printf("%v (%v)\r\n", reflect.TypeOf(data).Kind(), reflect.TypeOf(data))
	//fmt.Println(data)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		mapV := reflect.ValueOf(data)
		keysV := mapV.MapKeys()
		for key := range keysV {
			result = append(result, keysV[key].String())
		}
	case reflect.Array, reflect.Slice:
		sliceV := reflect.ValueOf(data)
		for idx := 0; idx < sliceV.Len(); idx++ {
			result = append(result, fmt.Sprintf("%d", idx))
		}
	case reflect.Struct:
		result = structs.Names(data)
	}
	return result
}

// Is what we are asking for a possible option
func inside(possible []string, target string) bool {
	for _, val := range possible {
		if target == val {
			return true
		}
	}
	return false
}

// Returns the next level of the interface
// Supporting nested data
func next_level(current_level interface{}, go_to string) (interface{}, error) {
	if inside(GetPossible(current_level), go_to) {
		//fmt.Printf("%v (%v)\r\n", reflect.TypeOf(current_level).Kind(), reflect.TypeOf(current_level))
		switch reflect.TypeOf(current_level).Kind() {
		case reflect.Map:
			CL, err := mapToInterface(current_level)
			if err != nil {
				return nil, err
			}
			return CL[go_to], nil
		case reflect.Array, reflect.Slice:
			val, err := strconv.Atoi(go_to)
			if err == nil {
				CL, err := sliceToInterface(current_level)
				if err != nil {
					return nil, err
				}
				return CL[val], nil
			} else {
				return nil, err
			}
		case reflect.Struct:
			structV := reflect.ValueOf(current_level)
			return structV.FieldByName(go_to).Interface(), nil
		}
	}
	return nil, fmt.Errorf("failed moving to '%s' from '%s' (%v)", go_to, current_level, reflect.TypeOf(current_level))
}

func list_possible(possible []string) []string {
	var result []string
	for _, val := range possible {
		result = append(result, fmt.Sprintf("'%s'", val))
	}
	return result
}

// Once you've got a single return type maybe you want it a particular type
// Returns string
func String(data interface{}) (string, error) {
	at := GetPossible(data)
	if len(at) != 0 {
		return "", fmt.Errorf("can't convert multiple values to string")
	}
	return fmt.Sprintf("%v", data), nil
}

// Once you've got a single return type maybe you want it a particular type
// Returns int
func Int(data interface{}) (int, error) {
	at := GetPossible(data)
	if len(at) != 0 {
		return 0, fmt.Errorf("can't convert multiple values to string")
	}
	return data.(int), nil
}

// Once you've got a single return type maybe you want it a particular type
// Returns float64
func Float64(data interface{}) (float64, error) {
	at := GetPossible(data)
	if len(at) != 0 {
		return 0.0, fmt.Errorf("can't convert multiple values to string")
	}
	return data.(float64), nil
}

// The main function, call this to walk your data
func Glom(data interface{}, path string) (interface{}, error) {
	complete_path := strings.Split(path, ".")
	//fmt.Printf("Seeking '%s' will take %d steps\r\n", path, len(complete_path))
	var path_taken []string
	var currently interface{}
	currently = data
	for _, hop := range complete_path {
		//fmt.Printf("current: %v\r\n", currently)
		//fmt.Printf("Path: '%v'\r\n", strings.Join(path_taken, "."))
		if hop != "*" && !inside(GetPossible(currently), hop) {
			return nil, fmt.Errorf("failed moving to '%s' from path of '%s', options are %s (%d)", hop, strings.Join(path_taken, "."), strings.Join(list_possible(GetPossible(currently)), ", "), len(GetPossible(currently)))
		} else {
			if hop != "*" {
				next, err := next_level(currently, hop)
				if err != nil {
					//return nil, fmt.Errorf("Failed moving to '%s' from path of '%s', options are %s (%d)", hop, strings.Join(path_taken, "."), strings.Join(list_possible(getPossible(next)), ", "), len(getPossible(next)))
					return nil, err
				} else {
					path_taken = append(path_taken, hop)
					currently = next
				}
			} else {
				return currently, nil
			}
		}
	}
	return currently, nil
}
