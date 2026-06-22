package util

import (
	"fmt"
	"strconv"
)

func StringValue(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil

	case int:
		return strconv.Itoa(v), nil

	case int32:
		return strconv.FormatInt(int64(v), 10), nil

	case int64:
		return strconv.FormatInt(v, 10), nil

	case bool:
		if v {
			return "true", nil
		}

		return "false", nil
	}

	return "", fmt.Errorf("Unsupported string conversion")
}

func NumberValue(value any) (float64, error) {
	switch v := value.(type) {
	case string:
		number, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return 0, fmt.Errorf("Unsupported number conversion")
		}

		return number, nil

	case float64:
		return v, nil

	case int:
		return float64(v), nil

	case int32:
		return float64(v), nil

	case int64:
		return float64(v), nil

	case bool:
		if v {
			return 1, nil
		}

		return 0, nil
	}

	return 0, fmt.Errorf("Unsupported number conversion")
}
