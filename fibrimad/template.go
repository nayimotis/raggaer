package main

import (
	"html/template"
	"strings"
	"time"

	"github.com/raggaer/fibrimad/app/role"
)

func setTemplateFuncMap() template.FuncMap {
	m := template.FuncMap{
		"eq":                        eq,
		"getDetailedRole":           detailedRole,
		"date":                      date,
		"longString":                longString,
		"splitString":               splitStringArray,
		"getLastElementStringArray": getLastElementStringArray,
		"validViewExtension":        validViewExtension,
	}

	return m
}

func validViewExtension(ext string) bool {
	valid := []string{"pdf", "txt", "css", "xml", "png", "jpeg", "jpg", "gif"}
	for _, e := range valid {
		if e == ext {
			return true
		}
	}
	return false
}

func eq(a, b interface{}) bool {
	return a == b
}

func detailedRole(rolename string) *role.DetailedRole {
	r, ok := role.List[rolename]
	if !ok {
		return nil
	}
	return &r
}

func date(d time.Time, layout string) string {
	return d.Format(layout)
}

func longString(s string) string {
	t := 15
	if len(s) <= t {
		return s
	}

	return s[0:t] + "..."
}

func splitStringArray(s, d string) []string {
	return strings.Split(s, d)
}

func getLastElementStringArray(s []string) string {
	return s[len(s)-1]
}
