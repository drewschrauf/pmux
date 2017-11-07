package util

// ArrayContains : Check if slice contains element
func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ArrayContainsAll : Check if slace contins all elements
func ArrayContainsAll(set []string, subset []string) bool {
	var found []string
	for _, e := range subset {
		if ArrayContains(set, e) {
			found = append(found, e)
		}
	}
	return len(found) == len(subset)
}
