package main

import "os"

func Contains(arr []string, str string) bool {
	for _, item := range arr {
		// fmt.Println("item == str", item, str, (item == str))
		if item == str {
			return true
		}
	}
	return false
}

func Exist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
