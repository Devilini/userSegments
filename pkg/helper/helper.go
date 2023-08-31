package helper

func DifferenceSlices(slice1, slice2 []string) []string {
	mb := make(map[string]bool, len(slice2))
	for _, x := range slice2 {
		mb[x] = true
	}
	var diff []string
	for _, x := range slice1 {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

func IntersectionSlices(slice1, slice2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range slice1 {
		hash[e] = true
	}
	for _, e := range slice2 {
		if hash[e] {
			inter = append(inter, e)
		}
	}

	return inter
}
