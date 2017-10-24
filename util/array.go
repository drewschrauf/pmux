package util

func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ArrayContainsAll(set []string, subset []string) bool {
	var found []string
	for _, e := range subset {
		if ArrayContains(set, e) {
			found = append(found, e)
		}
	}
	return len(found) == len(subset)
}
