package utils

func PtrToString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
