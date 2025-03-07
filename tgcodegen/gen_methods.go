// this is good
// TODO:
// 1. support mutlipart
// 2. files only and maybe ids split
// 3. whats with IFile in api_types.go? is it parsable? should we use some dull struct? -- OK
package main

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

type Method struct {
	Name        string   `json:"name"`
	Href        string   `json:"href"`
	Description []string `json:"description"`
	Returns     []string `json:"returns"`
	Fields      []*Field `json:"fields"`
}

func (fac *Factory2) BuildMethods(methods map[string]*Method) string {
	result := []string{
		"package tg",
		"import (",
		`"context"`,
		")",
	}

	methodsBuilt := []string{}
	for methodName, method := range methods {
		methodsBuilt = append(methodsBuilt, fac.buildMethod(methodName, method))
	}
	slices.Sort(methodsBuilt)
	result = append(result, methodsBuilt...)

	return strings.Join(result, "\n\n")
}

func (fac *Factory2) buildMethod(name string, method *Method) string {
	slog.Info("buildMethod", "name", name)
	result := &goFunc{
		name:    firstUpper(name),
		argsReq: []*goFuncReqArgument{},
		argsOpt: nil,
		requestStruct: &goTypeStruct{
			Name:    "Request",
			Comment: "",
			Fields:  []*goTypeStructField{},
		},
		genericFunc: "GenericRequest",
		returns:     fac.goTypeSlice(method.Returns),
		comment:     method.Description,
		extra:       []string{},
	}

	for _, field := range method.Fields {
		fieldType, fieldExtra := fac.parseMethodArgumentType(field.Types)
		result.extra = append(result.extra, fieldExtra...)
		fieldName := snakeCaseToCamelCase(field.Name)

		if slicesContainsAny(field.Types, "InputFile", "InputMedia", "Array of InputMediaPhoto") {
			result.genericFunc = "GenericRequestMultipart"
		}

		result.requestStruct.Fields = append(result.requestStruct.Fields, &goTypeStructField{
			Name:    fieldName,
			Comment: "",
			Type:    fieldType,
			Tag: fmt.Sprintf("`json:\"%s%s\"`", field.Name, func() string {
				if field.Required {
					return ""
				}
				return ",omitempty"
			}()),
		})

		if field.Required {
			result.argsReq = append(result.argsReq, &goFuncReqArgument{
				Name: firstLower(fieldName),
				Type: fieldType,
			})
			continue
		}

		if result.argsOpt == nil {
			result.argsOpt = &goTypeStruct{
				Name:    "Opt" + result.name,
				Comment: "",
				Fields:  []*goTypeStructField{},
			}
		}
		result.argsOpt.Fields = append(result.argsOpt.Fields, &goTypeStructField{
			Name:    fieldName,
			Comment: "",
			Type:    fieldType,
			Tag:     "",
		})
	}

	return result.build()
}

func (fac *Factory2) parseMethodArgumentType(types []string) (argType string, built []string) {
	slog.Debug("parseMethodArgumentType", "types", types)
	switch {
	case slices.Contains(types, "InputFile") && slices.Contains(types, "String"):
		return "InputFile", nil
	case slices.Contains(types, "InputFile"):
		return "*LocalFile", nil
	case len(types) > 1 && !slicesContainsAny(types, telegramCoreTypes...):
		return fac.findOrBuildVariantByFields(types)
	default:
		return fac.goType(types[0]), nil
	}
}

type goFuncReqArgument struct {
	Name string
	Type string
}

type goFunc struct {
	name          string
	argsReq       []*goFuncReqArgument
	argsOpt       *goTypeStruct
	requestStruct *goTypeStruct
	genericFunc   string
	returns       []string
	comment       []string
	extra         []string
}

func (fn *goFunc) build() string {
	result := []string{
		makeComment(fn.name, fn.comment...),
	}

	arguments := []string{"ctx context.Context"}
	reqArgumentsFill := []string{}
	schedulerChatId := "0"
	schedulerWeight := "1"
	if fn.name == "SendMediaGroup" || fn.name == "SendPaidMedia" {
		schedulerWeight = "len(media)"
	}
	for _, arg := range fn.argsReq {
		arguments = append(arguments, fmt.Sprintf("%s %s", arg.Name, arg.Type))
		reqArgumentsFill = append(reqArgumentsFill, fmt.Sprintf("%s: %s,", firstUpper(arg.Name), arg.Name))
		if arg.Name == "chatId" && arg.Type == "int64" {
			schedulerChatId = "chatId"
		}
	}
	if fn.argsOpt != nil {
		arguments = append(arguments, fmt.Sprintf("opts ...*%s", fn.argsOpt.Name))
	}

	funcReturns := fmt.Sprintf("(%s, error)", fn.returns[0])
	if len(fn.returns) > 1 {
		funcReturns += fmt.Sprintf("/* >> either: %v */", fn.returns[1:])
	}

	result = append(result,
		fmt.Sprintf("func %s(%s) %s {", fn.name, strings.Join(arguments, ", "), funcReturns),
		fmt.Sprintf("schedule(ctx, %s, %s)", schedulerChatId, schedulerWeight),
		fmt.Sprintf("defer scheduleDone(ctx, %s, %s)", schedulerChatId, schedulerWeight),
		strings.TrimSpace(fn.requestStruct.build()),
		fmt.Sprintf("request := &Request{\n%s\n}", strings.Join(reqArgumentsFill, "\n")),
	)

	if fn.argsOpt != nil {
		if len(fn.argsOpt.Fields) == 0 {
			panic("need at least one field " + fn.argsOpt.Name)
		}

		result = append(result, "for _, opt := range opts {")

		for _, field := range fn.argsOpt.Fields {
			checkForZero := "if %s != nil {"
			switch field.Type {
			case "int64":
				checkForZero = "if %s != 0 {"
			case "float64":
				checkForZero = "if %s != 0.0 {"
			case "bool":
				checkForZero = "if %s {"
			case "string":
				checkForZero = "if %s != \"\" {"
			}

			result = append(result,
				fmt.Sprintf(checkForZero, "opt."+field.Name),
				fmt.Sprintf("request.%s = opt.%s", field.Name, field.Name),
				"}",
			)
		}

		result = append(result, "}")
	}

	result = append(result,
		fmt.Sprintf("return %s[Request, %s](ctx, \"%s\", request)", fn.genericFunc, fn.returns[0], firstLower(fn.name)),
		"}",
	)
	if fn.argsOpt != nil {
		result = append(result, fn.argsOpt.build())
	}

	result = append(result, fn.extra...)
	return strings.Join(result, "\n")
}
