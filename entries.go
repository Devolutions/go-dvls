package dvls

import (
	"strconv"
	"strings"
)

const (
	entryEndpoint            string = "/api/connections/partial"
	entryConnectionsEndpoint string = "/api/connections"
)

type Entries struct {
	UserCredential *EntryUserCredentialService
	Certificate    *EntryCertificateService
}

func keywordsToSlice(kw string) []string {
	var spacedTag bool
	tags := strings.FieldsFunc(string(kw), func(r rune) bool {
		if r == '"' {
			spacedTag = !spacedTag
		}
		return !spacedTag && r == ' '
	})
	for i, v := range tags {
		unquotedTag, err := strconv.Unquote(v)
		if err != nil {
			continue
		}

		tags[i] = unquotedTag
	}

	return tags
}

func sliceToKeywords(kw []string) string {
	keywords := []string(kw)
	for i, v := range keywords {
		if strings.Contains(v, " ") {
			kw[i] = "\"" + v + "\""
		}
	}

	kString := strings.Join(keywords, " ")

	return kString
}
