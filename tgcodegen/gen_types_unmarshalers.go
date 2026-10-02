// todo: rewrite this garbage. unstar everything?
package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"sort"
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
	funcs := stringsUnique(result)
	sort.Strings(funcs)
	return funcs
}

type structUnmarshalerVariant struct {
	Name          string
	Type          string
	Tag           string
	Options       []*goTypeStruct
	InterfaceType string
	IsArray       bool
}

// ambiguousFields returns the set of field names that appear, with different
// Go types, on more than one option of this variant (e.g. OwnedGiftRegular.Gift
// is *Gift, OwnedGiftUnique.Gift is *UniqueGift). A single merged field can't
// be declared for these under their normal pointer-to-type shape, since Go
// doesn't allow two struct fields with the same name.
func (variant *structUnmarshalerVariant) ambiguousFields() map[string]bool {
	seenType := map[string]string{}
	ambiguous := map[string]bool{}
	for _, option := range variant.Options {
		for _, field := range option.Fields {
			if prevType, ok := seenType[field.Name]; ok {
				if prevType != field.Type {
					ambiguous[field.Name] = true
				}
				continue
			}
			seenType[field.Name] = field.Type
		}
	}
	return ambiguous
}

func (variant *structUnmarshalerVariant) buildStructDeclaration() string {
	ambiguous := variant.ambiguousFields()
	result := []string{
		fmt.Sprintf("type %s struct {", variant.Type),
	}

	for _, option := range variant.Options {
		for _, field := range option.Fields {
			tag := removeDefaultFromTag(strings.ReplaceAll(field.Tag, ",omitempty", ""))
			if ambiguous[field.Name] {
				// Deferred: decoded into its real type only once the
				// discriminant has picked the concrete variant.
				result = append(result, fmt.Sprintf("%s json.RawMessage %s", field.Name, tag))
				continue
			}
			result = append(result, fmt.Sprintf("%s *%s %s", field.Name, field.Type, tag))
		}
	}
	result = stringsUnique(result)

	result = append(result, "}")
	return strings.Join(result, "\n")
}

// buildFieldValue renders the expression that reads optionField off of the
// joined instance variable named inst, choosing the raw-JSON decode path for
// fields whose name is ambiguous across this variant's options.
func (variant *structUnmarshalerVariant) buildFieldValue(inst string, optionField *goTypeStructField) string {
	if variant.ambiguousFields()[optionField.Name] {
		return fmt.Sprintf("unmarshalRawOrZero[%s](%s.%s)", optionField.Type, inst, optionField.Name)
	}
	return fmt.Sprintf("deref(%s.%s)", inst, optionField.Name)
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
			variant.Options = append(variant.Options, fac.findStructByGoName(option))
		}
		result.Variants = append(result.Variants, variant)
	}

	return result
}

// sortedDiscriminatorValues returns discr.mapping's keys in a stable order.
// discr.mapping is a Go map, so ranging over it directly (as all three
// discriminator-switch builders used to) makes the generated switch's case
// order - and therefore the generated file's exact bytes - vary from one
// generator run to the next, even with no schema or logic change.
func sortedDiscriminatorValues(discr *discriminator) []string {
	values := make([]string, 0, len(discr.mapping))
	for value := range discr.mapping {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
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
		fmt.Sprintf("if %s != nil && %s.%s != nil {", inst, inst, snakeCaseToCamelCase(discr.property)),
		fmt.Sprintf("switch *%s.%s {", inst, snakeCaseToCamelCase(discr.property)),
	}

	findOption := func(name string) *goTypeStruct {
		i := slices.IndexFunc(variant.Options, func(el *goTypeStruct) bool { return unwrapType(el.Name) == name })
		return variant.Options[i]
	}

	for _, value := range sortedDiscriminatorValues(discr) {
		option := findOption(discr.mapping[value])
		result = append(result,
			fmt.Sprintf("case \"%s\":", value),
			fmt.Sprintf("impl.%s = &%s{", variant.Name, option.Name),
		)

		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: %s,", optionField.Name, variant.buildFieldValue(inst, optionField)),
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

	for _, value := range sortedDiscriminatorValues(discr) {
		option := findOption(discr.mapping[value])
		result = append(result,
			fmt.Sprintf("case \"%s\":", value),
			fmt.Sprintf("%s = append(%s, &%s{", target, target, option.Name),
		)

		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: %s,", optionField.Name, variant.buildFieldValue("item", optionField)),
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

// buildBareUnmarshalerFunc generates a standalone decoder for an interface
// type used as a bare value (e.g. an API method's return type) rather than as
// a struct field. Unlike buildParsingCodeWithDiscriminators, there is no
// containing impl/inst to assign into, so this returns the concrete value
// directly instead of assigning to a field.
func (variant *structUnmarshalerVariant) buildBareUnmarshalerFunc(discr *discriminator) string {
	slices.SortFunc(variant.Options, func(a, b *goTypeStruct) int {
		return cmp.Compare(len(a.Fields), len(b.Fields))
	})

	discriminatorField := snakeCaseToCamelCase(discr.property)
	result := []string{
		fmt.Sprintf("func Unmarshal%s(data []byte) (%s, error) {", variant.InterfaceType, variant.InterfaceType),
		strings.ReplaceAll(variant.buildStructDeclaration(), ",omitempty", ""),
		fmt.Sprintf("var inst %s", variant.Type),
		"if err := json.Unmarshal(data, &inst); err != nil {",
		"return nil, err",
		"}",
		fmt.Sprintf("if inst.%s == nil {", discriminatorField),
		fmt.Sprintf("return nil, fmt.Errorf(\"tg: unmarshal %s: missing %%q\", %q)", variant.InterfaceType, discr.property),
		"}",
		fmt.Sprintf("switch *inst.%s {", discriminatorField),
	}

	findOption := func(name string) *goTypeStruct {
		i := slices.IndexFunc(variant.Options, func(el *goTypeStruct) bool { return unwrapType(el.Name) == name })
		return variant.Options[i]
	}

	for _, value := range sortedDiscriminatorValues(discr) {
		option := findOption(discr.mapping[value])
		result = append(result,
			fmt.Sprintf("case %q:", value),
			fmt.Sprintf("return &%s{", option.Name),
		)
		for _, optionField := range option.Fields {
			result = append(result,
				fmt.Sprintf("%s: %s,", optionField.Name, variant.buildFieldValue("inst", optionField)),
			)
		}
		result = append(result, "}, nil")
	}

	result = append(result,
		"}",
		fmt.Sprintf("return nil, fmt.Errorf(\"tg: unmarshal %s: unknown %%q %%q\", %q, *inst.%s)", variant.InterfaceType, discr.property, discriminatorField),
		"}",
	)
	return strings.Join(result, "\n")
}

// buildBareInterfaceUnmarshalers builds a standalone UnmarshalXxx(data []byte)
// (Xxx, error) function for each interface name in names. It exists for
// interfaces used as a bare API method return type (or the element type of a
// bare []Interface return) — unlike a struct field, those never pass through
// buildUnmarshalers, so encoding/json has no way to pick a concrete type for
// them on its own.
func (fac *Factory2) buildBareInterfaceUnmarshalers(names map[string]bool) []string {
	result := []string{}
	for name := range names {
		variant, ok := fac.Interfaces[name]
		if !ok {
			panic(fmt.Sprintf("tg: method returns bare interface %q that was never registered as an interface type", name))
		}
		discr, ok := discriminators[name]
		if !ok {
			panic(fmt.Sprintf("tg: method returns bare interface %q with no registered discriminator - add one in schema.go", name))
		}

		options := []*goTypeStruct{}
		for _, option := range variant.Options {
			options = append(options, fac.findStructByGoName(option))
		}

		unmarshaler := &structUnmarshalerVariant{
			Type:          name + "UnmarshalJoined",
			Options:       options,
			InterfaceType: name,
		}
		result = append(result, unmarshaler.buildBareUnmarshalerFunc(discr))
	}
	sort.Strings(result)
	return result
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
				fmt.Sprintf("%s: %s,", optionField.Name, variant.buildFieldValue(inst, optionField)),
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
				fmt.Sprintf("%s: %s,", optionField.Name, variant.buildFieldValue("item", optionField)),
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
