package storage

import (
	_ "embed" // embed SQL sqhema
	"strings"
)

// DBSchema contains the storage scheme for PostgreSQL
//
//go:embed schema.sql
var DBSchema string

func removeComments(schema string) string {
	ret := ""
	lines := strings.Split(schema, "\n")
	for _, v := range lines {
		line := v
		pos := strings.Index(line, "--")
		if pos == 0 {
			continue
		}
		if pos > 0 {
			line = line[0:pos]
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		ret += line + "\n"
	}
	return ret
}

func schemaCommandList() []string {
	var ret []string
	commands := strings.Split(removeComments(DBSchema), ";")
	for _, v := range commands {
		if len(strings.TrimSpace(v)) > 0 {
			ret = append(ret, strings.TrimSpace(v))
		}
	}
	return ret
}
