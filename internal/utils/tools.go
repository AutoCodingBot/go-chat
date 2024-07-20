package utils

import (
	"sort"
)

func GenerateConversationId(useNameAlpha string, userNameBeta string) string {
	names := []string{useNameAlpha, userNameBeta}
	sort.Strings(names)
	return names[0] + "_" + names[1]
}
