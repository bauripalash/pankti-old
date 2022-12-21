package constants

func GetStdName(n string) (string, bool) {

	SL_BN_EN := map[string]string{
		"গণিত":    "math",
		"তালিকা":  "array",
		"তারিখ":   "date",
		"ফাইল":    "file",
		"সাধারণ":  "std",
		"স্ট্রিং": "string",
		"সিস্টেম": "sys",
	}

	val, ok := SL_BN_EN[n]

	if ok {

		return val, ok
	} else {
		return "", ok
	}

}
