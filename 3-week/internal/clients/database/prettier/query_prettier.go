package prettier

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	PlaceholderDollar   = "$"
	PlaceholderQuestion = "?"
)

func Pretty(query string, placeholder string, args ...any) string {
	if len(args) == 0 {
		return clean(query)
	}

	replacements := make([]string, 0, len(args)*2)
	for i, param := range args {
		key := placeholder + strconv.Itoa(i+1)
		replacements = append(replacements, key, formatValue(param))
	}

	query = strings.NewReplacer(replacements...).Replace(query)
	return clean(query)
}

func clean(query string) string {
	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")
	return strings.TrimSpace(query)
}

func formatValue(param any) string {
	switch v := param.(type) {
	case string, time.Time:
		return fmt.Sprintf("%q", v)
	case []byte:
		return fmt.Sprintf("%q", string(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}
