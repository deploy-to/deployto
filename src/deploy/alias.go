package deploy

import "strings"

func BuildAlias(names []string) string {
	return strings.Join(names, "-")
}
