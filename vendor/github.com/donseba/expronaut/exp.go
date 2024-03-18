package expronaut

import (
	"context"
	"errors"
	"fmt"
)

func ToGoTemplate(comparison string) string {
	lexer := NewLexer(comparison)
	p := NewParser(lexer)
	tree := p.Parse()

	return tree.GoTemplate()
}

func SetVariables(ctx context.Context, variables map[string]any) context.Context {
	return context.WithValue(ctx, ContextKey, variables)
}

func Evaluate(ctx context.Context, comparison string) (any, error) {
	lexer := NewLexer(comparison)
	p := NewParser(lexer)
	tree := p.Parse()

	return tree.Evaluate(ctx)
}

func EvaluateBool(ctx context.Context, comparison string) (bool, error) {
	ev, err := Evaluate(ctx, comparison)
	if err != nil {
		return false, err
	}

	b, ok := ev.(bool)
	if !ok {
		return false, fmt.Errorf("non-boolean result")
	}

	return b, nil
}

// Exp evaluates a comparison expression with the given variables.
// ideally this is used in a template engine
// Example:
//
//	{{ if exp "foo == 5", "foo", 5 }} Hello {{ end }}
//	{{ if exp "foo == 5 && bar == 10", "foo", 5, "bar", 10 }} Hello {{ end }}
//	{{ if exp "foo == 5 && bar == 10", "foo", 5, "bar", 10 }} Hello {{ end }}
func Exp(comparison string, params ...any) (bool, error) {
	var dict map[string]any
	if len(params) > 0 {
		if len(params)%2 != 0 {
			return false, errors.New("invalid dict call")
		}

		dict = make(map[string]any, len(params)/2)
		for i := 0; i < len(params); i += 2 {
			key, ok := params[i].(string)
			if !ok {
				return false, errors.New("dict keys must be strings")
			}
			dict[key] = params[i+1]
		}
	}

	ctx := context.TODO()
	ctx = SetVariables(ctx, dict)

	return EvaluateBool(ctx, comparison)
}
