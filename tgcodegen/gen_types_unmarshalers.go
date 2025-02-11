// todo: rewrite this garbage. unstar everything?
package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

func (fac *Factory2) buildUnmarshalers() []string {
	result := []string{}
	for _, strct := range fac.Structs {
		fieldVariants := []string{}
		for _, field := range strct.Fields {
			if _, ok := fac.Interfaces[unwrapType(field.Type)]; ok {
				fieldVariants = append(fieldVariants, field.Name)
			}
		}
		if len(fieldVariants) > 0 {
			result = append(result, fac.buildStructUnmarshaler(strct, fieldVariants).build())
		}
	}
	return stringsUnique(result)
}

type structUnmarshalerVariant struct {
	Name          string
	Type          string
	Tag           string
	Options       []*goTypeStruct
	InterfaceType string
	IsArray       bool
}

func (variant *structUnmarshalerVariant) buildStructDeclaration() string {
	result := []string{
		fmt.Sprintf("type %s struct {", variant.Type),
	}

	for _, option := range variant.Options {
		for _, field := range option.Fields {
			result = append(result, fmt.Sprintf("%s *%s %s",
				field.Name, field.Type, removeDefaultFromTag(strings.ReplaceAll(field.Tag, ",omitempty", ""))),
			)
		}
	}
	result = stringsUnique(result)

	result = append(result, "}")
	return strings.Join(result, "\n")
}

func (variant *structUnmarshalerVariant) buildParsingCode() string {
	if discr, ok := discriminators[variant.InterfaceType]; ok {
		if variant.IsArray {
			return variant.buildParsingCodeSliceWithDiscriminators(discr)
		}
		return variant.buildParsingCodeWithDiscriminators(discr)
	}
	slog.Info("unmarshalers.buildParsingCode#no_discriminators", "interface", variant.InterfaceType)

	if variant.IsArray {
		return variant.buildParsingCodeSlice()
	}
	return variant.buildParsingCodeSingle()
}

type structUnmarshaler struct {
	OriginalStruct *goTypeStruct
	MainFields     []*goTypeStructField
	Variants       []*structUnmarshalerVariant
}

func (unmarshaler *structUnmarshaler) build() string {
	result := []string{
		fmt.Sprintf("func (impl *%s) UnmarshalJSON(data []byte) error {", unmarshaler.OriginalStruct.Name),
	}

	for _, variant := range unmarshaler.Variants {
		result = append(result, strings.ReplaceAll(variant.buildStructDeclaration(), ",omitempty", ""))
	}

	baseInstanceMainFields := []string{}
	for _, field := range unmarshaler.MainFields {
		baseInstanceMainFields = append(baseInstanceMainFields, strings.ReplaceAll(field.buildDeclarationWithoutDefault(), ",omitempty", ""))
	}
	baseInstanceVariantFields := []string{}
	for _, variant := range unmarshaler.Variants {
		slice := ""
		if variant.IsArray {
			slice = "[]"
		}
		baseInstanceVariantFields = append(baseInstanceVariantFields,
			"// Joint of structs, used for parsing variant interfaces.",
			fmt.Sprintf("%s %s*%s %s", variant.Name, slice, variant.Type, variant.Tag),
		)
	}

	result = append(result,
		fmt.Sprintf("type BaseInstance struct {\n%s\n%s\n}",
			strings.Join(stringsUnique(baseInstanceMainFields), "\n"),
			strings.Join(stringsUnique(baseInstanceVariantFields), "\n"),
		),
		"var inst BaseInstance",
		"if err := json.Unmarshal(data, &inst); err != nil {",
		"return err",
		"}",
	)
	for _, field := range unmarshaler.MainFields {
		result = append(result,
			fmt.Sprintf("impl.%s = inst.%s", field.Name, field.Name),
		)
	}

	for _, variant := range unmarshaler.Variants {
		result = append(result, variant.buildParsingCode())
	}

	result = append(result,
		"return nil",
		"}",
	)
	return strings.Join(result, "\n")
}

func (fac *Factory2) buildStructUnmarshaler(strct *goTypeStruct, variants []string) *structUnmarshaler {
	slog.Debug("buildStructUnmarshaler", "type", strct.Name, "variants", variants)
	result := &structUnmarshaler{
		OriginalStruct: strct,
		MainFields:     []*goTypeStructField{},
		Variants:       []*structUnmarshalerVariant{},
	}

	for _, field := range strct.Fields {
		if !slices.Contains(variants, field.Name) {
			slog.Debug("buildUnmarshaler#normal", "type", strct.Name, "field", field.Name, "type", field.Type)
			result.MainFields = append(result.MainFields, field)
			continue
		}

		fieldOptions := fac.Interfaces[unwrapType(field.Type)].Options
		slog.Debug("buildUnmarshaler#variant",
			"type", strct.Name,
			"field", field.Name,
			"type", field.Type,
			"options", fieldOptions,
		)
		variant := &structUnmarshalerVariant{
			Name:          field.Name,
			Type:          strings.TrimLeft(field.Type, "[]") + "UnmarshalJoined" + field.Name,
			InterfaceType: unwrapType(field.Type),
			Tag:           field.Tag,
			IsArray:       strings.HasPrefix(field.Type, "[]"),
			Options:       []*goTypeStruct{},
		}
		for _, option := range fieldOptions {
			variant.Options = append(variant.Options, fac.Structs[option])
		}
		result.Variants = append(result.Variants, variant)
	}

	return result
}

func (variant *structUnmarshalerVariant) buildParsingCodeWithDiscriminators(discr *discriminator) string {
	slices.SortFunc(variant.Options, func(a, b *goTypeStruct) int {
		return cmp.Compare(len(a.Fields), len(b.Fields))
	})

	if len(variant.Options) != len(discr.mapping) {
		slog.Error("unmarshalers.buildParsingCodeWithDiscriminators#bad_disciminators",
			"type", variant.Type,
			"variant.Options", variant.Options,
			"discriminators", discr.mapping,
		)
	}

	inst := fmt.Sprintf("inst.%s", variant.Name)
	result := []string{
		fmt.Sprintf("if %s != nil && %s.%s == nil {", inst, inst, snakeCaseToCamelCase(discr.property)),
		fmt.Sprintf("switch *%s.%s {", inst, snakeCaseToCamelCase(discr.property)),
	}

	findOption := func(name string) *goTypeStruct {
		i := slices.IndexFunc(variant.Options, func(el *goTypeStruct) bool { return unwrapType(el.Name) == name })
		return variant.Options[i]
	}

	for value, typeName := range discr.mapping {
		option := findOption(typeName)
		result = append(result,
			fmt.Sprintf("case \"%s\":", value),
			fmt.Sprintf("impl.%s = &%s{", variant.Name, option.Name),
		)

		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: deref(%s.%s),", optionField.Name, inst, optionField.Name),
			)
		}
		result = append(result, "}")
	}

	result = append(result,
		"}", // switch
		"}", // if
	)
	return strings.Join(result, "\n")
}

func (variant *structUnmarshalerVariant) buildParsingCodeSliceWithDiscriminators(discr *discriminator) string {
	slices.SortFunc(variant.Options, func(a, b *goTypeStruct) int {
		return cmp.Compare(len(a.Fields), len(b.Fields))
	})

	if len(variant.Options) != len(discr.mapping) {
		slog.Error("unmarshalers.buildParsingCodeWithDiscriminators#bad_disciminators",
			"type", variant.Type,
			"variant.Options", variant.Options,
			"discriminators", discr.mapping,
		)
	}

	inst := fmt.Sprintf("inst.%s", variant.Name)
	target := fmt.Sprintf("impl.%s", variant.Name)
	result := []string{
		fmt.Sprintf("if len(%s) != 0 {", inst),
		fmt.Sprintf("%s = []%s{}", target, variant.InterfaceType),
		fmt.Sprintf("for _, item := range %s {", inst),
		fmt.Sprintf("if item == nil || item.%s == nil { continue }", snakeCaseToCamelCase(discr.property)),
		fmt.Sprintf("switch *item.%s {", snakeCaseToCamelCase(discr.property)),
	}

	findOption := func(name string) *goTypeStruct {
		i := slices.IndexFunc(variant.Options, func(el *goTypeStruct) bool { return unwrapType(el.Name) == name })
		return variant.Options[i]
	}

	for value, typeName := range discr.mapping {
		option := findOption(typeName)
		result = append(result,
			fmt.Sprintf("case \"%s\":", value),
			fmt.Sprintf("%s = append(%s, &%s{", target, target, option.Name),
		)

		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: deref(item.%s),", optionField.Name, optionField.Name),
			)
		}
		result = append(result, "})")
	}

	result = append(result,
		"}", // switch
		"}", // for
		"}", // if
	)
	return strings.Join(result, "\n")
}

func (variant *structUnmarshalerVariant) buildParsingCodeSingle() string {
	slices.SortFunc(variant.Options, func(a, b *goTypeStruct) int {
		return cmp.Compare(len(a.Fields), len(b.Fields))
	})

	inst := fmt.Sprintf("inst.%s", variant.Name)
	result := []string{
		fmt.Sprintf("if %s != nil {", inst),
		"nonEmptyFields := []string{}",
	}

	visited := []string{}
	for _, option := range variant.Options {
		for _, optionField := range option.Fields {
			if slices.Contains(visited, optionField.Name) {
				continue
			}
			result = append(result,
				fmt.Sprintf("if %s.%s != nil {", inst, optionField.Name),
				fmt.Sprintf("nonEmptyFields = append(nonEmptyFields, \"%s\")", optionField.Name),
				"}",
			)
			visited = append(visited, optionField.Name)
		}
	}

	result = append(result, "switch {")
	for _, option := range variant.Options {
		fieldsNames := []string{}
		for _, field := range option.Fields {
			fieldsNames = append(fieldsNames, `"`+field.Name+`"`)
		}
		result = append(result,
			fmt.Sprintf("case containsAll([]string{%s}, nonEmptyFields):", strings.Join(fieldsNames, ", ")),
			fmt.Sprintf("impl.%s = &%s{", variant.Name, option.Name),
		)
		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: deref(%s.%s),", optionField.Name, inst, optionField.Name),
			)
		}
		result = append(result, "}")
	}
	result = append(result, "}")

	result = append(result, "}")
	return strings.Join(result, "\n")
}

func (variant *structUnmarshalerVariant) buildParsingCodeSlice() string {
	slices.SortFunc(variant.Options, func(a, b *goTypeStruct) int {
		return cmp.Compare(len(a.Fields), len(b.Fields))
	})

	inst := fmt.Sprintf("inst.%s", variant.Name)
	target := fmt.Sprintf("impl.%s", variant.Name)
	result := []string{
		fmt.Sprintf("if len(%s) != 0 {", inst),
		fmt.Sprintf("%s = []%s{}", target, variant.InterfaceType),
		fmt.Sprintf("for _, item := range %s {", inst),
		"if item == nil { continue }",
		"nonEmptyFields := []string{}",
	}

	visited := []string{}
	for _, option := range variant.Options {
		for _, optionField := range option.Fields {
			if slices.Contains(visited, optionField.Name) {
				continue
			}
			result = append(result,
				fmt.Sprintf("if item.%s != nil {", optionField.Name),
				fmt.Sprintf("nonEmptyFields = append(nonEmptyFields, \"%s\")", optionField.Name),
				"}",
			)
			visited = append(visited, optionField.Name)
		}
	}

	result = append(result, "switch {")
	for _, option := range variant.Options {
		fieldsNames := []string{}
		for _, field := range option.Fields {
			fieldsNames = append(fieldsNames, `"`+field.Name+`"`)
		}
		result = append(result,
			fmt.Sprintf("case containsAll([]string{%s}, nonEmptyFields):", strings.Join(fieldsNames, ", ")),
			fmt.Sprintf("%s = append(%s, &%s{", target, target, option.Name),
		)
		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: deref(item.%s),", optionField.Name, optionField.Name),
			)
		}
		result = append(result, "})")
	}
	result = append(result, "}")

	result = append(result,
		"}", // for
		"}", // if
	)
	return strings.Join(result, "\n")
}

type goTypeStructFieldJoined struct {
	Name   string
	Type   string
	Tag    string
	UsedIn *set[string]
}

type goTypeStructJoined struct {
	Name   string
	Type   string
	Fields []*goTypeStructFieldJoined
}

func (goStruct *goTypeStructJoined) hasField(name string) bool {
	for _, field := range goStruct.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func (goStruct *goTypeStructJoined) build() string {
	fieldLines := []string{}
	for _, field := range goStruct.Fields {
		fieldLines = append(fieldLines,
			strings.Trim(fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag), "\n"),
		)
	}
	return fmt.Sprintf("type %s struct {\n%s\n}", goStruct.Name, strings.Join(fieldLines, "\n"))
}
