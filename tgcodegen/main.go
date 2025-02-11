package main

import (
	"log/slog"
	"os"
	"os/exec"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	schema, err := Read[DynamicSchema]("tgcodegen/data/schema.json")
	if err != nil {
		panic(err)
	}

	typesFile := "api_types.go"
	typesUnmarshalersFile := "api_types_unmarshalers.go"
	methodsFile := "api_methods.go"

	fac := NewFactory()
	typesString, unmarshalersString := fac.BuildTypesFromSchema(schema.Types)
	if err := os.WriteFile(typesFile, []byte(typesString), 0644); err != nil {
		panic(err)
	}
	_, _ = exec.Command("go", "fix", "tg").CombinedOutput()
	_, _ = exec.Command("go", "fmt", typesFile).CombinedOutput()

	if err := os.WriteFile(typesUnmarshalersFile, []byte(unmarshalersString), 0644); err != nil {
		panic(err)
	}
	_, _ = exec.Command("go", "fix", "tg").CombinedOutput()
	_, _ = exec.Command("go", "fmt", typesUnmarshalersFile).CombinedOutput()

	if err := os.WriteFile(methodsFile, []byte(fac.BuildMethods(schema.Methods)), 0644); err != nil {
		panic(err)
	}
	_, _ = exec.Command("go", "fix", "tg").CombinedOutput()
	_, _ = exec.Command("go", "fmt", methodsFile).CombinedOutput()
}
