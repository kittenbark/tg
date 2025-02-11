package main

import (
	"slices"
	"strings"
)

var telegramCoreTypes = []string{
	"Integer",
	"Float",
	"Boolean",
	"String",
}

func unwrapArrayOf(typ string) string {
	if !strings.HasPrefix(typ, "Array of ") {
		return typ
	}
	return unwrapArrayOf(typ[len("Array of "):])
}

func convertArrayOfToBrackets(typ string) string {
	return strings.Repeat("[]", strings.Count(typ, "Array of "))
}

func snakeCaseToCamelCase(str string) string {
	return strings.Join(stringsMap(strings.Split(str, "_"), func(s string) string { return strings.Title(s) }), "")
}

func makeComment(name string, description ...string) string {
	lines := description
	if len(description) == 1 {
		lines = strings.Split(strings.TrimSuffix(strings.ReplaceAll(description[0], ". ", ".\n"), "\n"), "\n")
		i, j := 0, 1
		for j < len(lines) {
			if i == j {
				j++
				continue
			}
			if lines[i] == "" || len(lines[i]) >= 80 {
				i++
				continue
			}
			if len(lines[i])+len(lines[j]) < 117 {
				lines[i] += " " + lines[j]
				lines[j] = ""
			}
			j++
		}
		newLines := []string{}
		for _, line := range lines {
			if line != "" {
				newLines = append(newLines, line)
			}
		}
		lines = newLines
	}

	if len(lines) == 0 {
		return ""
	}

	comments := strings.Join(stringsMap(lines, func(s string) string { return "// " + s }), "\n")
	if name == "" {
		return comments
	}
	return strings.Replace(comments, "//", "// "+name, 1)
}

func stringsMap(list []string, lambda func(string) string) []string {
	result := make([]string, len(list))
	for i := range list {
		result[i] = lambda(list[i])
	}
	return result
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func firstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func slicesContainsAny[T comparable](arr []T, items ...T) bool {
	for _, item := range items {
		if slices.Contains(arr, item) {
			return true
		}
	}
	return false
}

func stringsUnique(list []string) []string {
	result := []string{}
	for _, str := range list {
		if !slices.Contains(result, str) {
			result = append(result, str)
		}
	}
	return result
}

type set[T comparable] struct {
	data map[T]struct{}
}

func newSet[T comparable]() *set[T] {
	return &set[T]{data: make(map[T]struct{})}
}

func (s *set[T]) insert(item T) {
	s.data[item] = struct{}{}
}

func (s *set[T]) erase(item T) {
	delete(s.data, item)
}

func (s *set[T]) contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

func (s *set[T]) intersection(other *set[T]) *set[T] {
	if other == nil || other.data == nil {
		return newSet[T]()
	}

	result := newSet[T]()
	for item := range s.data {
		if other.contains(item) {
			result.insert(item)
		}
	}
	return result
}

func (s *set[T]) union(other *set[T]) *set[T] {
	if other == nil || other.data == nil {
		return newSet[T]()
	}

	result := newSet[T]()
	for item := range s.data {
		result.insert(item)
	}
	for item := range other.data {
		result.insert(item)
	}
	return result
}

func (s *set[T]) diff(other *set[T]) *set[T] {
	if other == nil || other.data == nil {
		return newSet[T]()
	}

	result := newSet[T]()
	for item := range s.data {
		if !other.contains(item) {
			result.insert(item)
		}
	}
	for item := range other.data {
		if !s.contains(item) {
			result.insert(item)
		}
	}
	return result
}

func stringsTrim(list []string, cutset string) []string {
	result := []string{}
	for _, str := range list {
		result = append(result, strings.Trim(str, cutset))
	}
	return result
}

func stringsTrimPrefix(list []string, prefix string) []string {
	result := []string{}
	for _, str := range list {
		result = append(result, strings.TrimPrefix(str, prefix))
	}
	return result
}

func findDiscriminators(sets []*set[string]) []string {
	result := []string{}
	for i := range sets {
		union := newSet[string]()
		for j := range sets {
			if i != j {
				union = union.union(sets[j])
			}
		}

		discriminator := ""
		for el := range sets[i].data {
			if !union.contains(el) {
				discriminator = el
				break
			}
		}
		result = append(result, discriminator)
	}
	return result
}

func getOrDefault[T any](list []T, i int, def T) T {
	if len(list) <= i {
		return def
	}
	return list[i]
}

func removeDefaultFromTag(tag string) string {
	if !strings.Contains(tag, "default:") {
		return tag
	}
	tag, _, _ = strings.Cut(tag, "default:")
	return strings.TrimSpace(tag) + "`"
}
