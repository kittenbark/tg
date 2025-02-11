package main

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

/*
note: these words are blessed by the god
the key was: splitting parsing and stringify.
*/

type Type struct {
	Name        string   `json:"name"`
	Href        string   `json:"href"`
	Description []string `json:"description"`
	Fields      []*Field `json:"fields"`
	SubtypeOf   []string `json:"subtype_of"`
	Subtypes    []string `json:"subtypes"`
}

func NewFactory() *Factory2 {
	basicsList := []*goTypeBasic{
		{Name: "Integer", Association: "int64"},
		{Name: "Float", Association: "float64"},
		{Name: "Boolean", Association: "bool"},
		{Name: "String", Association: "string"},
	}
	basics := map[string]*goTypeBasic{}
	for _, basic := range basicsList {
		basics[basic.Name] = basic
	}
	return &Factory2{
		Basics:  basics,
		Structs: map[string]*goTypeStruct{},
		Interfaces: map[string]*goTypeInterfaceVariant{
			"Album": {
				Name:    "Album",
				Comment: "",
				Options: []string{
					"InputMediaAudio", "InputMediaDocument", "InputMediaPhoto", "InputMediaVideo",
				},
				Prefix: "",
			},
		},
		Specials: map[string]string{
			"Message":                  "*Message",  // Message depends on itself, so we have to give the codegen a hint.
			"InaccessibleMessage":      "*Message",  // InaccessibleMessage is a subset of Message (field-wise).
			"MaybeInaccessibleMessage": "*Message",  // MaybeInaccessibleMessage = Message | InaccessibleMessage.
			"InputFile":                "InputFile", // This is how Telegram represents mutlipart file uploading.
		},
		Renames: map[string]string{
			"InputMediaVideo":     "Video",
			"InputMediaPhoto":     "Photo",
			"InputMediaAudio":     "Audio",
			"InputMediaDocument":  "Document",
			"InputMediaAnimation": "Animation",
			"Photo":               "TelegramPhoto",
			"Video":               "TelegramVideo",
			"Audio":               "TelegramAudio",
			"Document":            "TelegramDocument",
			"Animation":           "TelegramAnimation",
		},
		Skip: []string{"InputFile", "MaybeInaccessibleMessage", "InaccessibleMessage"},
	}
}

func (fac *Factory2) BuildTypesFromSchema(schema map[string]*Type) (string, string) {
	slog.Debug("BuildTypesFromSchema", "schema_len", len(schema))

	result := []string{
		"package tg",
		`import (
"context"
)`,
	}

	resultUnmarshalers := []string{
		"package tg",
		`import "encoding/json"`,
	}

	typesBuilt := []string{}
	builtTotal := len(fac.Skip)
	for builtTotal < len(schema) {
		builtIteration := 0
		for _, typ := range schema {
			if fac.contains(typ.Name, true) || !fac.buildable(typ) {
				continue
			}
			typesBuilt = append(typesBuilt, fac.build(typ))
			builtIteration++
		}
		builtTotal += builtIteration
		slog.Info("BuildTypesFromSchema#build_loop", "total", builtTotal, "iter", builtIteration)
		if builtIteration == 0 {
			panic("BuildTypesFromSchema: no types built in iteration")
		}
	}
	slices.Sort(typesBuilt)
	result = append(result, typesBuilt...)
	resultUnmarshalers = append(resultUnmarshalers, fac.buildUnmarshalers()...)

	return strings.Join(result, "\n\n"), strings.Join(resultUnmarshalers, "\n\n")
}

type goTypeBasic struct {
	Name        string
	Association string
}

type goTypeStructField struct {
	Name    string
	Comment string
	Type    string
	Tag     string
}

func (field *goTypeStructField) buildDeclaration() string {
	return strings.Trim(fmt.Sprintf("%s\n%s %s %s", field.Comment, field.Name, field.Type, field.Tag), "\n")
}

func (field *goTypeStructField) buildDeclarationWithoutDefault() string {
	return strings.Trim(fmt.Sprintf("%s\n%s %s %s", field.Comment, field.Name, field.Type, removeDefaultFromTag(field.Tag)), "\n")
}

type goTypeStruct struct {
	Name    string
	Comment string
	Fields  []*goTypeStructField
	Extra   []string
}

func (goStruct *goTypeStruct) hasField(name string) bool {
	for _, field := range goStruct.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func (goStruct *goTypeStruct) build() string {
	fieldLines := []string{}
	for _, field := range goStruct.Fields {
		fieldLines = append(fieldLines, field.buildDeclaration())
	}
	return fmt.Sprintf("%s\ntype %s struct {\n%s\n}\n%s",
		goStruct.Comment,
		goStruct.Name,
		strings.Join(fieldLines, "\n"),
		strings.Join(goStruct.Extra, "\n"),
	)
}

type goTypeInterfaceVariant struct {
	Name    string
	Comment string
	Options []string
	Prefix  string
}

func (variant *goTypeInterfaceVariant) build() string {
	if len(variant.Options) == 0 {
		return fmt.Sprintf("%s\ntype %s any", variant.Comment, variant.Name)
	}

	optionsDeclVarLines := []string{}
	for _, option := range variant.Options {
		optionsDeclVarLines = append(optionsDeclVarLines, fmt.Sprintf("_ %s = &%s{}", variant.Name, option))
	}
	optionsVarDecl := fmt.Sprintf("var (\n%s\n)", strings.Join(optionsDeclVarLines, "\n"))

	optionsImplementationLines := []string{}
	optionsDeclarationLines := []string{}
	for _, optionImplementing := range variant.Options {
		for _, option := range variant.Options {
			returnStatement := "return nil"
			if optionImplementing == option {
				returnStatement = "return impl"
			}
			funcName := variant.Prefix + strings.TrimPrefix(option, variant.Name)
			optionsImplementationLines = append(optionsImplementationLines,
				fmt.Sprintf("func (impl *%s) %s() *%s { %s }", optionImplementing, funcName, option, returnStatement),
			)
			optionsDeclarationLines = append(optionsDeclarationLines,
				fmt.Sprintf("%s() *%s", funcName, option),
			)
		}
		optionsImplementationLines = append(optionsImplementationLines, "")
	}
	optionsImplementation := strings.Join(optionsImplementationLines, "\n")
	optionsDeclaration := strings.Join(stringsUnique(optionsDeclarationLines), "\n")

	return fmt.Sprintf("%s\ntype %s interface{\n%s\n}\n\n%s\n\n%s",
		variant.Comment, variant.Name, optionsDeclaration,
		optionsVarDecl,
		optionsImplementation,
	)
}

type discriminator struct {
	property string
	mapping  map[string]string
}

type Factory2 struct {
	Basics     map[string]*goTypeBasic
	Structs    map[string]*goTypeStruct
	Interfaces map[string]*goTypeInterfaceVariant
	Renames    map[string]string
	Specials   map[string]string
	Skip       []string
}

func (fac *Factory2) build(typ *Type) string {
	switch {
	case slices.Contains(fac.Skip, typ.Name):
		return ""
	case len(typ.Subtypes) > 0:
		return fac.buildInterfaceVariant(typ)
	default:
		return fac.buildStruct(typ)
	}
}

// todo: this variantsPrefix is a temporary measure. pls fix one day.
func (fac *Factory2) buildInterfaceVariant(typ *Type, variantsPrefix ...string) string {
	fac.Interfaces[typ.Name] = &goTypeInterfaceVariant{
		Name:    typ.Name,
		Comment: makeComment(typ.Name, typ.Description...),
		Options: typ.Subtypes,
		Prefix:  getOrDefault(variantsPrefix, 0, "Opt"),
	}
	defer slog.Info("buildInterfaceVariant#done", "type", typ.Name)

	for i := range fac.Interfaces[typ.Name].Options {
		fac.Interfaces[typ.Name].Options[i] = fac.renames(fac.Interfaces[typ.Name].Options[i])
	}

	return fac.Interfaces[typ.Name].build()
}

func (fac *Factory2) buildStruct(typ *Type) string {
	fac.Structs[typ.Name] = &goTypeStruct{
		Name:    fac.renames(typ.Name),
		Comment: makeComment(fac.renames(typ.Name), typ.Description...),
		Fields:  []*goTypeStructField{},
		Extra:   []string{},
	}
	defer slog.Info("buildStruct#done", "type", typ.Name)

	for _, field := range typ.Fields {
		fieldTag := fmt.Sprintf("`json:\"%s,omitempty\"`", field.Name)
		if field.Required {
			fieldTag = fmt.Sprintf("`json:\"%s\"`", field.Name)
		}
		fieldTag = fac.withDefaultByDiscriminators(typ.Name, field, fieldTag)
		fieldType := fac.goType(field.Types[0])
		fieldName := firstUpper(snakeCaseToCamelCase(field.Name))

		guiltyComment := ""
		if len(field.Types) > 1 {
			guiltyComment = "\n// >> either: " + strings.Join(field.Types[1:], ", ")
		}

		if fieldName != "SmallFileId" && strings.HasSuffix(fieldName, "FileId") && fieldType == "string" {
			fac.Structs[typ.Name].Extra = append(fac.Structs[typ.Name].Extra,
				fmt.Sprintf("func (impl *%s) Download(ctx context.Context, path string) error {\n"+
					"return GenericDownload(ctx, path, impl.%s)\n"+
					"}", fac.Structs[typ.Name].Name, fieldName,
				),
				fmt.Sprintf("func (impl *%s) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {\n"+
					"return GenericDownloadTemp(ctx, impl.%s, dirAndPattern...)\n"+
					"}", fac.Structs[typ.Name].Name, fieldName,
				),
			)
		}
		if strings.HasPrefix(typ.Name, "InputMedia") {
			if fieldName == "Media" {
				fieldType = "InputFile"
			}
		}

		fac.Structs[typ.Name].Fields = append(fac.Structs[typ.Name].Fields, &goTypeStructField{
			Name:    fieldName,
			Comment: makeComment("", field.Description) + guiltyComment,
			Type:    fieldType,
			Tag:     fieldTag,
		})
	}

	if strings.HasPrefix(typ.Name, "InputMedia") {
		fac.Structs[typ.Name].Fields = append(fac.Structs[typ.Name].Fields, &goTypeStructField{
			Name:    fac.Specials["InputFile"],
			Comment: "// Used for uploading media.",
			Type:    fac.Specials["InputFile"],
			Tag:     "`json:\"-\"`",
		})
	}

	return fac.Structs[typ.Name].build()
}

func (fac *Factory2) contains(name string, strict ...bool) bool {
	name = unwrapArrayOf(name)
	if slices.Contains(fac.Skip, name) {
		return true
	}
	if _, ok := fac.Specials[name]; ok && (len(strict) == 0 || !strict[0]) {
		return true
	}
	if _, ok := fac.Interfaces[name]; ok {
		return true
	}
	if _, ok := fac.Structs[name]; ok {
		return true
	}
	if _, ok := fac.Basics[name]; ok {
		return true
	}
	slog.Debug("fac.contains#no", "name", name)
	return false
}

func (fac *Factory2) buildable(typ *Type) bool {
	if len(typ.Subtypes) > 0 {
		for _, subtype := range typ.Subtypes {
			if !fac.contains(subtype) {
				slog.Debug("buildable#not", "type", typ.Name, "missing_subtype", subtype)
				return false
			}
			return true
		}
	}

	for _, field := range typ.Fields {
		for _, fieldType := range field.Types {
			if !fac.contains(fieldType) {
				slog.Debug("buildable#not", "type", typ.Name, "missing_type", fieldType)
				return false
			}
		}
	}
	return true
}

func (fac *Factory2) renames(name string) string {
	if rename, ok := fac.Renames[name]; ok {
		return rename
	}
	return name
}

func (fac *Factory2) withDefaultByDiscriminators(typ string, field *Field, tag string) string {
	fieldDiscr := ""
	for _, discr := range discriminators {
		if discr.property != field.Name {
			continue
		}
		for discrValue, discrTyp := range discr.mapping {
			if discrTyp != typ {
				continue
			}
			fieldDiscr = discrValue
		}

	}
	if fieldDiscr == "" {
		return tag
	}
	return fmt.Sprintf("`%s default:\"%s\"`",
		strings.Trim(tag, "`"), fieldDiscr,
	)
}

func (fac *Factory2) goType(name string) (result string) {
	defer func() { result = strings.ReplaceAll(result, "[]*PhotoSize", "TelegramPhoto") }()

	array := convertArrayOfToBrackets(name)
	name = unwrapArrayOf(name)
	if variant, ok := fac.Interfaces[name]; ok {
		return array + variant.Name
	}
	if srct, ok := fac.Structs[name]; ok {
		return array + "*" + srct.Name
	}
	if basic, ok := fac.Basics[name]; ok {
		return array + basic.Association
	}
	if special, ok := fac.Specials[name]; ok {
		return array + special
	}
	return name + "/* TODO: this type not found, fix it. */"
}

func (fac *Factory2) goTypeSlice(slice []string) []string {
	slice = slices.Clone(slice)
	for i := range slice {
		slice[i] = fac.goType(slice[i])
	}
	slices.Sort(slice)
	return slice
}

func (fac *Factory2) goTypeSliceWithoutArrayOf(slice []string) []string {
	slice = slices.Clone(slice)
	for i := range slice {
		slice[i] = strings.TrimLeft(fac.goType(slice[i]), "[]*")
	}
	slices.Sort(slice)
	return slice
}

func (fac *Factory2) tryFindVariantByFields(types []string) (*goTypeInterfaceVariant, bool) {
	if len(types) <= 1 || slicesContainsAny(types, telegramCoreTypes...) {
		return nil, false
	}
	types = fac.goTypeSliceWithoutArrayOf(types)

	slog.Debug("fac.tryFindVariantByFields#start", "types", types)
	for _, variant := range fac.Interfaces {
		slog.Debug("fac.tryFindVariantByFields#compare", "variant", variant.Name, "types", types, "options", variant.Options)
		if slices.Equal(fac.goTypeSliceWithoutArrayOf(variant.Options), types) {
			slog.Debug("fac.tryFindVariantByFields#equal", "types", types, "variant", variant.Options)
			return variant, true
		}
	}
	return nil, false
}

func (fac *Factory2) findOrBuildVariantByFields(types []string) (string, []string) {
	result, ok := fac.tryFindVariantByFields(types)
	if ok {
		return result.Name, nil
	}

	typesWithoutArrays := stringsMap(types, func(s string) string { return strings.TrimSpace(strings.ReplaceAll(s, "Array of", " ")) })
	name := "Variant" + strings.Join(typesWithoutArrays, "")
	return strings.Repeat("[]", strings.Count(types[0], "Array of")) + name, []string{fac.buildInterfaceVariant(&Type{
		Name:     name,
		Subtypes: typesWithoutArrays,
	}, strings.ToLower(name))}
}

func unwrapType(s string) string {
	return strings.NewReplacer(
		"[]", "",
		"*", "",
	).Replace(s)
}
