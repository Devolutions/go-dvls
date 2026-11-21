package dvls

import (
	"encoding/json"
	"fmt"
)

func jsonBody(v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	return body, nil
}

func unmarshalBody(body []byte, v any) error {
	err := json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}

// singleByName returns the only item of items, which the server filtered on name ignoring case. When several items
// differ only by case, the one whose name equals name exactly is returned. notFound or multiple is returned when no
// item or more than one item matches.
func singleByName[T any](items []T, name string, nameOf func(T) string, notFound error, multiple error) (T, error) {
	var zero T

	switch len(items) {
	case 0:
		return zero, notFound
	case 1:
		return items[0], nil
	}

	var match T
	count := 0
	for _, item := range items {
		if nameOf(item) == name {
			match = item
			count++
		}
	}

	if count != 1 {
		return zero, multiple
	}

	return match, nil
}
