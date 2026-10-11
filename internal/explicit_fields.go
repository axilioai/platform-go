package internal

// Repo-owned: regen.sh keeps this file out of every regen. Fern's version
// drops the embedded model whenever a MarshalJSON wrapper has shadow fields
// beside it (every model with a date field), so a setter call sent only the
// shadowed dates. See generated_runtime_test.go.

import (
	"math/big"
	"reflect"
	"strings"
)

// HandleExplicitFields processes a struct to remove `omitempty` from
// fields that have been explicitly set (as indicated by their corresponding bit in explicitFields).
// Note that `marshaler` should be an embedded struct to avoid infinite recursion.
// Returns an interface{} that can be passed to json.Marshal.
//
// A generated MarshalJSON wraps the model as
//
//	struct {
//		embed
//		StartAt *internal.DateTime `json:"start_at,omitempty"`
//	}
//
// where each field after the embed shadows the model field of the same name
// with a date-aware encoding. The result has the model's exported fields in
// order, with each shadowed one replaced by its shadow; a field's explicit bit
// is its index in the model.
func HandleExplicitFields(marshaler interface{}, explicitFields *big.Int) interface{} {
	val := reflect.ValueOf(marshaler)
	typ := reflect.TypeOf(marshaler)

	// Handle pointer types
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// Only handle struct types
	if val.Kind() != reflect.Struct {
		return marshaler
	}

	// If no explicit fields set, use standard marshaling
	if explicitFields == nil || explicitFields.Sign() == 0 {
		return marshaler
	}

	// The wrapper pattern: the model embedded first, shadows after it.
	sourceVal, sourceType := val, typ
	shadows := map[string]int{}
	if typ.NumField() > 0 && typ.Field(0).Anonymous && typ.Field(0).Type.Kind() == reflect.Struct {
		sourceVal = val.Field(0)
		sourceType = sourceVal.Type()
		for i := 1; i < typ.NumField(); i++ {
			if typ.Field(i).IsExported() {
				shadows[typ.Field(i).Name] = i
			}
		}
	}

	fields := make([]reflect.StructField, 0, sourceType.NumField()+len(shadows))
	values := make([]reflect.Value, 0, cap(fields))

	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)

		// Skip unexported fields and the explicitFields field itself
		if !field.IsExported() || field.Name == "explicitFields" {
			continue
		}

		value := sourceVal.Field(i)
		if shadowIndex, ok := shadows[field.Name]; ok {
			field = typ.Field(shadowIndex)
			value = val.Field(shadowIndex)
			delete(shadows, field.Name)
		}

		// Check if this field has been explicitly set
		fieldBit := big.NewInt(1)
		fieldBit.Lsh(fieldBit, uint(i))
		if big.NewInt(0).And(explicitFields, fieldBit).Sign() != 0 {
			field.Tag = withoutOmitempty(field.Tag)
		}

		fields = append(fields, field)
		values = append(values, value)
	}

	// Shadows with no model field of their name are never explicit.
	for i := 1; i < typ.NumField(); i++ {
		if _, ok := shadows[typ.Field(i).Name]; ok {
			fields = append(fields, typ.Field(i))
			values = append(values, val.Field(i))
		}
	}

	// Create new struct type with modified tags
	newType := reflect.StructOf(fields)
	newVal := reflect.New(newType).Elem()
	for i, value := range values {
		newVal.Field(i).Set(value)
	}

	return newVal.Interface()
}

// withoutOmitempty returns tag with `omitempty` removed from its json options,
// keeping its url tag.
func withoutOmitempty(tag reflect.StructTag) reflect.StructTag {
	jsonTag := tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return tag
	}

	// Parse the json tag, remove omitempty from options
	parts := strings.Split(jsonTag, ",")
	if len(parts) > 1 {
		newParts := []string{parts[0]} // Keep the field name
		for _, part := range parts[1:] {
			if strings.TrimSpace(part) != "omitempty" {
				newParts = append(newParts, part)
			}
		}
		jsonTag = strings.Join(newParts, ",")
	}

	// Reconstruct the struct tag
	newTag := `json:"` + jsonTag + `"`
	if urlTag := tag.Get("url"); urlTag != "" {
		newTag += ` url:"` + urlTag + `"`
	}
	return reflect.StructTag(newTag)
}
