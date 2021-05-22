package client

// StringPtr returns the pointer of a string value.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns the pointer of an int value.
func IntPtr(i int) *int {
	return &i
}

// TimePtr returns the pointer of a string value.
func BoolPtr(b bool) *bool {
	return &b
}
