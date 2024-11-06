package server

func keysSliceFromMap[K comparable, V any](mymap map[K]V) []K {
	keys := make([]K, 0, len(mymap))
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}
