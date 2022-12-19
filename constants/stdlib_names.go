package constants

func GetStdName(n string) (string, bool) {

	SL_BN_EN := map[string]string{
		"গণিত": "math",
	}

	val, ok := SL_BN_EN[n]

	if ok {

		return val, ok
	} else {
		return "", ok
	}

}
