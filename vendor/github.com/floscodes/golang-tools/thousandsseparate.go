package tools

import (
	"fmt"
	"strings"
)

func ThousandsSeparate(N interface{}, lang string) (string, error) {

	switch N.(type) {
	case int:
	case int16:
	case int32:
	case int64:
	case int8:
	case float32:
	case float64:

	default:
		return "", fmt.Errorf("Not a valid number!")

	}

	n := fmt.Sprintf("%v", N)

	switch lang {

	case "de":

		n = strings.ReplaceAll(n, ",", ".")

		dec := ""

		if strings.Index(n, ".") != -1 {
			dec = n[strings.Index(n, ".")+1 : len(n)]
			n = n[0:strings.Index(n, ".")]

		}

		for i := 0; i <= len(n); i = i + 4 {

			a := n[0 : len(n)-i]
			b := n[len(n)-i : len(n)]
			n = a + "." + b

		}

		if n[0:1] == "." {
			n = n[1:len(n)]
		}

		if n[len(n)-1:len(n)] == "." {
			n = n[0 : len(n)-1]
		}

		if dec != "" {
			n = n + "," + dec
		}

		return n, nil

	case "en":

		n = strings.ReplaceAll(n, ",", "")

		dec := ""

		if strings.Index(n, ".") != -1 {
			dec = n[strings.Index(n, ".")+1 : len(n)]
			n = n[0:strings.Index(n, ".")]

		}

		for i := 0; i <= len(n); i = i + 4 {

			a := n[0 : len(n)-i]
			b := n[len(n)-i : len(n)]
			n = a + "," + b

		}

		if n[0:1] == "," {
			n = n[1:len(n)]
		}

		if n[len(n)-1:len(n)] == "," {
			n = n[0 : len(n)-1]
		}

		if dec != "" {
			n = n + "." + dec
		}

		return n, nil

	}

	return n, nil

}
