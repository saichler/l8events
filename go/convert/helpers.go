package convert

import (
	"fmt"
	evt "github.com/saichler/l8events/go/types/l8events"
	"strconv"
	"strings"
)

// str returns the attribute value or empty string if missing.
func str(attrs map[string]string, key string) string {
	return attrs[key]
}

// i32 parses an attribute as int32. Returns 0 for missing keys, error for malformed values.
func i32(attrs map[string]string, key string) (int32, error) {
	v, ok := attrs[key]
	if !ok || v == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("attribute %q: %w", key, err)
	}
	return int32(n), nil
}

// i64 parses an attribute as int64. Returns 0 for missing keys, error for malformed values.
func i64(attrs map[string]string, key string) (int64, error) {
	v, ok := attrs[key]
	if !ok || v == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("attribute %q: %w", key, err)
	}
	return n, nil
}

// f64 parses an attribute as float64. Returns 0 for missing keys, error for malformed values.
func f64(attrs map[string]string, key string) (float64, error) {
	v, ok := attrs[key]
	if !ok || v == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("attribute %q: %w", key, err)
	}
	return f, nil
}

// boolean parses an attribute as bool. Returns false for missing keys, error for malformed values.
func boolean(attrs map[string]string, key string) (bool, error) {
	v, ok := attrs[key]
	if !ok || v == "" {
		return false, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("attribute %q: %w", key, err)
	}
	return b, nil
}

// subCategory parses the "subCategory" attribute as int32 for use as an enum value.
func subCategory(attrs map[string]string) (int32, error) {
	return i32(attrs, "subCategory")
}

// setCommon copies the common fields from an EventRecord into the target fields.
func setCommon(r *evt.EventRecord, eventId *string, propertyId *string, sourceId *string, sourceType *string, message *string) {
	*eventId = r.EventId
	*propertyId = str(r.Attributes, "propertyId")
	*sourceId = r.SourceId
	*sourceType = r.SourceType
	*message = r.Message
}

// extractMap collects all attributes with the given prefix into a map.
// Keys have the prefix stripped. For example, prefix "varbinds." with key
// "varbinds.1.3.6.1" yields map entry "1.3.6.1".
func extractMap(attrs map[string]string, prefix string) map[string]string {
	result := make(map[string]string)
	for k, v := range attrs {
		if strings.HasPrefix(k, prefix) {
			result[strings.TrimPrefix(k, prefix)] = v
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
