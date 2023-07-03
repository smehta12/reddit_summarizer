package utils

func GetMapValues(m map[string]string) []string {
	values := make([]string, len(m))

	for _, v := range m {
		values = append(values, v)
	}

	return values
}
