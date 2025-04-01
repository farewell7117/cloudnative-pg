package pgbouncer

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
