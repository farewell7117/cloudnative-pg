package pgbouncer

import "encoding/base64"

func intValOrDefault(val *int32, def int32) int32 {
	if val == nil {
		return def
	}
	return *val
}

func boolYesNo(val *bool, def bool) string {
	if val == nil {
		if def {
			return "yes"
		}
		return "no"
	}
	if *val {
		return "yes"
	}
	return "no"
}

// Helper function to decode base64 strings
func decodeBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
