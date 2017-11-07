package main

func Contains(arr []string, str string) bool {
	for _, item := range arr {
		// fmt.Println("item == str", item, str, (item == str))
		if item == str {
			return true
		}
	}
	return false
}
