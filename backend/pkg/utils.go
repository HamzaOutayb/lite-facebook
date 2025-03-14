package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func WriteJson(w http.ResponseWriter, statuscode int, Data any) error {
	w.WriteHeader(statuscode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Data)
	if err != nil {
		return err
	}
	return nil
}

func ParseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// GetScanFields returns a slice of pointers to struct fields for scanning SQL results.
// Pass a pointer to the struct.
// Example: GetScanFields(&user) => []*interface{}{&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Image}
func GetScanFields(s interface{}) []interface{} {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		log.Fatal("Input must be a pointer to a struct")
	}
	val = val.Elem() // like &user => user

	fields := make([]interface{}, val.NumField()) // like user => user.ID, user.FirstName, ...
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}
	return fields
}

// GetExecFields returns a slice of struct field values, excluding specified fields.
// Example: GetExecFields(user, "ID", "CreatedAt") => []interface{}{user.FirstName, user.LastName, user.Nickname, user.Image}
func GetExecFields(s interface{}, excludeFields ...string) []interface{} {
	val := reflect.ValueOf(s) // like user => user.FirstName, user.LastName, ...
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		log.Fatal("Input must be a struct or pointer to a struct")
	}

	// Convert excluded fields to a map for fast lookup
	excluded := make(map[string]bool)
	for _, field := range excludeFields {
		excluded[field] = true
	}

	var fields []interface{}
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		if !excluded[fieldName] {
			fields = append(fields, val.Field(i).Interface())
		}
	}
	return fields
}

func SetSessionCookie(w http.ResponseWriter, uuid string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  uuid,
		Path:   "/",
		HttpOnly: true,
		MaxAge: 31536000,
	})
}

func DeleteSessionCookie(w http.ResponseWriter, uid string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  uid,
		Path:   "/",
		MaxAge: -1,
	})
}

func Length(a, b int, e string) bool {
	return len(e) < a || len(e) > b
}

func Placeholders(n int) string {
	return strings.Repeat("?,", n)[:2*n-1]
}
