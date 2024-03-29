package db

import (
	"fmt"
	"strings"
)

// Returning return for helping build query
// input : id, name, category
// output : "RETURNING id, name, category"
func Returning(columns ...string) string {
	sb := strings.Builder{}
	sb.WriteString("RETURNING ")
	for _, key := range columns {
		sb.WriteString(key + ", ")
	}
	return strings.TrimSuffix(sb.String(), ", ")
}

// A return A.text for helping join query
// input : updated_at
// output : "A.updated_at"
func A(text string) string {
	return fmt.Sprintf("A.%s", text)
}

// B return B.text for helping join query
// input : updated_at
// output : "B.updated_at"
func B(text string) string {
	return fmt.Sprintf("B.%s", text)
}

// C return C.text for helping join query
// input : updated_at
// output : "C.updated_at"
func C(text string) string {
	return fmt.Sprintf("C.%s", text)
}

// D return D.text for helping join query
// input : updated_at
// output : "D.updated_at"
func D(text string) string {
	return fmt.Sprintf("D.%s", text)
}

// E return E.text for helping join query
// input : updated_at
// output : "E.updated_at"
func E(text string) string {
	return fmt.Sprintf("E.%s", text)
}

// Dot return table.column for helping join query
// input : (user , updated_at)
// output : "user.updated_at"
func Dot(table, column string) string {
	return fmt.Sprintf("%s.%s", table, column)
}

// CoalesceInt Coalesce(null,default)
func CoalesceInt(text string, def int) string {
	return fmt.Sprintf("Coalesce(%s,%d)", text, def)
}

// CoalesceInt Coalesce(null,default)
func CoalesceFloat64(text string, def float64) string {
	return fmt.Sprintf("Coalesce(%s,%f)", text, def)
}

// CoalesceString Coalesce(null,default)
func CoalesceString(text string, def string) string {
	return fmt.Sprintf("Coalesce(%s,'%s')", text, def)
}
