package models

// StringPtr returns a pointer to the string value s.
// If s is empty, it returns nil.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntPtr returns a pointer to the int value i.
// This can be useful for optional integer fields in structs.
func IntPtr(i int) *int {
    return &i
}

// Int32Ptr returns a pointer to the int32 value i.
func Int32Ptr(i int32) *int32 {
    return &i
}

// Int64Ptr returns a pointer to the int64 value i.
func Int64Ptr(i int64) *int64 {
    return &i
}
