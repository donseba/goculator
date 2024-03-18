package expronaut

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	// ContextKey is used to store the variables in the context.
	ContextKey = "_exp"
)

// ASTNode is the interface for all nodes in the AST.
type ASTNode interface {
	Evaluate(ctx context.Context) (any, error) // Evaluate computes the value of the node.
	GoTemplate() string
	String() string
}

// IntLiteralNode represents an int literal in the AST.
type IntLiteralNode struct {
	Value int
}

// Evaluate computes the value of the number literal.
func (n *IntLiteralNode) Evaluate(ctx context.Context) (any, error) {
	return n.Value, nil // For a literal, Evaluate simply returns the literal's value.
}

func (n *IntLiteralNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// GoTemplate returns the Go template representation of the number literal.
func (n *IntLiteralNode) GoTemplate() string { return fmt.Sprintf("%d", n.Value) }

// FloatLiteralNode represents a numeric literal in the AST.
type FloatLiteralNode struct {
	Value float64
}

// Evaluate computes the value of the number literal.
func (n *FloatLiteralNode) Evaluate(ctx context.Context) (any, error) {
	return n.Value, nil // For a literal, Evaluate simply returns the literal's value.
}

func (n *FloatLiteralNode) String() string {
	return fmt.Sprintf("%f", n.Value)
}

// GoTemplate returns the Go template representation of the number literal.
func (n *FloatLiteralNode) GoTemplate() string { return fmt.Sprintf("%f", n.Value) }

// BinaryOperationNode represents a binary operation (e.g., addition, subtraction) in the AST.
type BinaryOperationNode struct {
	Left     ASTNode   // The left operand
	Operator TokenType // The operator
	Right    ASTNode   // The right operand
}

// Evaluate computes the value of the binary operation.
func (n *BinaryOperationNode) Evaluate(ctx context.Context) (any, error) {
	leftEval, err := n.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	rightEval, err := n.Right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	switch n.Operator {
	case TokenTypePlus:
		return BuiltinFunctions.Add(ctx, leftEval, rightEval)
	case TokenTypeMinus:
		return BuiltinFunctions.Sub(ctx, leftEval, rightEval)
	case TokenTypeMultiply:
		return BuiltinFunctions.Mul(ctx, leftEval, rightEval)
	case TokenTypeDivide:
		return BuiltinFunctions.Div(ctx, leftEval, rightEval)
	case TokenTypeDivideInteger:
		return BuiltinFunctions.DivInt(ctx, leftEval, rightEval)
	case TokenTypeModulo:
		return BuiltinFunctions.Mod(ctx, leftEval, rightEval)
	case TokenTypeExponent:
		return BuiltinFunctions.Exp(ctx, leftEval, rightEval)
	case TokenTypeEqual, TokenTypeNotEqual,
		TokenTypeLessThan, TokenTypeLessThanOrEqual,
		TokenTypeGreaterThan, TokenTypeGreaterThanOrEqual:

		if left, ok := leftEval.(string); ok {
			if right, ok := rightEval.(string); ok {
				return applyStringComparison(left, right, n.Operator), nil
			}
		} else if left, ok := leftEval.(float64); ok {
			if right, ok := rightEval.(float64); ok {
				return applyFloatComparison(left, right, n.Operator), nil
			} else if right, ok := rightEval.(int); ok {
				return applyFloatComparison(left, float64(right), n.Operator), nil
			}
		} else if left, ok := leftEval.(int); ok {
			if right, ok := rightEval.(int); ok {
				return applyIntComparison(left, right, n.Operator), nil
			} else if right, ok := rightEval.(float64); ok {
				return applyFloatComparison(float64(left), right, n.Operator), nil
			}
		} else if left, ok := leftEval.(time.Time); ok {
			if right, ok := rightEval.(time.Time); ok {
				return applyTimeComparison(left, right, n.Operator), nil
			}
		}
	case TokenTypeLeftShift:
		if left, ok := leftEval.(int); ok {
			if right, ok := rightEval.(int); ok {
				return left << right, nil
			}
		}
	case TokenTypeRightShift:
		if left, ok := leftEval.(int); ok {
			if right, ok := rightEval.(int); ok {
				return left >> right, nil
			}
		}
	default:
		return nil, fmt.Errorf("unknown or unsupported operator: %v", n.Operator)
	}

	return nil, fmt.Errorf(fmt.Sprintf("type mismatch or operation not applicable (%T(%v), %s, %T(%v))", leftEval, leftEval, n.Operator, rightEval, rightEval))
}

func (n *BinaryOperationNode) String() string {
	// Recursively call String() on the left and right operands to build the representation
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Operator, n.Right.String())
}

// GoTemplate returns the Go template representation of the binary operation.
func (n *BinaryOperationNode) GoTemplate() string {
	var (
		er = n.Right.GoTemplate()
		el = n.Left.GoTemplate()
	)

	if _, ok := n.Right.(*BinaryOperationNode); ok {
		er = fmt.Sprintf("(%s)", er)
	}

	if _, ok := n.Left.(*BinaryOperationNode); ok {
		el = fmt.Sprintf("(%s)", el)
	}

	return fmt.Sprintf("%s %s %s", TokenGoTemplate(n.Operator), el, er)
}

// StringLiteralNode represents a string literal in the AST.
type StringLiteralNode struct {
	Value string
}

// Evaluate computes the value of the string literal.
func (n *StringLiteralNode) Evaluate(ctx context.Context) (any, error) {
	return n.Value, nil // For a string literal, Evaluate simply returns the string's value.
}

// GoTemplate returns the Go template representation of the string literal.
func (n *StringLiteralNode) GoTemplate() string { return fmt.Sprintf(`"%s"`, n.Value) }

func (n *StringLiteralNode) String() string {
	return n.Value
}

// VariableNode represents a variable in the AST.
type VariableNode struct {
	Name string
}

// Evaluate computes the value of the variable.
func (n *VariableNode) Evaluate(ctx context.Context) (any, error) {
	// Assuming env is a map passed to Evaluate containing variable values.
	value, exists := n.lookup(ctx, n.Name)
	if !exists {
		return nil, fmt.Errorf("variable %s not defined", n.Name)
	}
	return value, nil
}

// GoTemplate returns the Go template representation of the variable.
func (n *VariableNode) lookup(ctx context.Context, name string) (any, bool) {
	vars, ok := ctx.Value(ContextKey).(map[string]any)
	if !ok {
		return nil, false
	}

	// Split the name into parts
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		return vars[parts[0]], true
	}

	// Traverse the parts to find the value
	value, exists := vars[parts[0]]
	if !exists {
		return nil, false
	}

	// do a recursive lookup
	for _, part := range parts[1:] {
		switch v := value.(type) {
		case map[string]any:
			value, exists = v[part]
			if !exists {
				return nil, false
			}
		default:
			return nil, false
		}
	}

	return value, true
}

func (n *VariableNode) String() string {
	return n.Name
}

// GoTemplate returns the Go template representation of the variable.
func (n *VariableNode) GoTemplate() string {
	return fmt.Sprintf(`.%s`, n.Name)
}

// LogicalOperationNode represents a logical operation (e.g., AND, OR) in the AST.
type LogicalOperationNode struct {
	Left     ASTNode
	Operator TokenType
	Right    ASTNode
}

// Evaluate computes the value of the logical operation.
func (n *LogicalOperationNode) Evaluate(ctx context.Context) (any, error) {
	leftEval, err := n.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	rightEval, err := n.Right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	leftBool, okLeft := leftEval.(bool)
	rightBool, okRight := rightEval.(bool)
	if !okLeft || !okRight {
		return nil, fmt.Errorf("operands for logical operation must be boolean")
	}

	switch n.Operator {
	case TokenTypeAnd:
		return leftBool && rightBool, nil
	case TokenTypeOr:
		return leftBool || rightBool, nil
	default:
		return nil, fmt.Errorf("unknown logical operator: %v", n.Operator)
	}
}

func (n *LogicalOperationNode) String() string {
	// Recursively call String() on the left and right operands to build the representation
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Operator, n.Right.String())
}

// GoTemplate returns the Go template representation of the logical operation.
func (n *LogicalOperationNode) GoTemplate() string {
	return fmt.Sprintf("%s ( %s ) ( %s )", TokenGoTemplate(n.Operator), n.Left.GoTemplate(), n.Right.GoTemplate())
}

// BooleanLiteralNode represents a boolean literal in the AST.
type BooleanLiteralNode struct {
	Value bool
}

// Evaluate computes the value of the boolean literal.
func (n *BooleanLiteralNode) Evaluate(ctx context.Context) (any, error) {
	return n.Value, nil
}

func (n *BooleanLiteralNode) String() string {
	return fmt.Sprintf("%v", n.Value)
}

// GoTemplate returns the Go template representation of the boolean literal.
func (n *BooleanLiteralNode) GoTemplate() string {
	if n.Value {
		return "true"
	}
	return "false"
}

type FunctionCallNode struct {
	FunctionName string
	Arguments    []ASTNode
}

func (n *FunctionCallNode) Evaluate(ctx context.Context) (any, error) {
	args := make([]any, len(n.Arguments))
	for i, arg := range n.Arguments {
		val, err := arg.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	f, ok := BuiltinFunctions[n.FunctionName]
	if ok {
		return f(ctx, args...)
	}

	return nil, fmt.Errorf("unknown function: %s", n.FunctionName)
}

func (n *FunctionCallNode) String() string {
	var args []string
	for _, arg := range n.Arguments {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("%s(%s)", n.FunctionName, args)
}

func (n *FunctionCallNode) GoTemplate() string {
	var args []string
	for _, arg := range n.Arguments {
		args = append(args, arg.GoTemplate())
	}
	return fmt.Sprintf("%s %s", n.FunctionName, args)
}

type arrayType string

const (
	arrayTypeInt    arrayType = "int"
	arrayTypeFloat  arrayType = "float"
	arrayTypeString arrayType = "string"
	arrayTypeTime   arrayType = "time"
	arrayTypeAny    arrayType = "any"
)

type ArrayNode struct {
	Elements []ASTNode
	Type     arrayType
}

func (n *ArrayNode) Evaluate(ctx context.Context) (any, error) {
	var elements []any
	for _, element := range n.Elements {
		val, err := element.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		elements = append(elements, val)
	}

	return elements, nil
}

func (n *ArrayNode) String() string {
	var elements []string
	for _, element := range n.Elements {
		elements = append(elements, element.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (n *ArrayNode) GoTemplate() string {
	var elements []string
	for _, element := range n.Elements {
		elements = append(elements, element.GoTemplate())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// applyStringComparison applies the comparison operator to the two strings.
func applyStringComparison(left, right string, op TokenType) bool {
	switch op {
	case TokenTypeEqual:
		return left == right
	case TokenTypeNotEqual:
		return left != right
	case TokenTypeLessThan:
		return left < right
	case TokenTypeLessThanOrEqual:
		return left <= right
	case TokenTypeGreaterThan:
		return left > right
	case TokenTypeGreaterThanOrEqual:
		return left >= right
	}
	return false
}

// applyFloatComparison applies the comparison operator to the two floats.
func applyFloatComparison(left, right float64, op TokenType) bool {
	switch op {
	case TokenTypeEqual:
		return left == right
	case TokenTypeNotEqual:
		return left != right
	case TokenTypeLessThan:
		return left < right
	case TokenTypeLessThanOrEqual:
		return left <= right
	case TokenTypeGreaterThan:
		return left > right
	case TokenTypeGreaterThanOrEqual:
		return left >= right
	}
	return false
}

// applyIntComparison applies the comparison operator to the two integers.
func applyIntComparison(left, right int, op TokenType) bool {
	switch op {
	case TokenTypeEqual:
		return left == right
	case TokenTypeNotEqual:
		return left != right
	case TokenTypeLessThan:
		return left < right
	case TokenTypeLessThanOrEqual:
		return left <= right
	case TokenTypeGreaterThan:
		return left > right
	case TokenTypeGreaterThanOrEqual:
		return left >= right
	}
	return false
}

func applyTimeComparison(left, right time.Time, op TokenType) bool {
	switch op {
	case TokenTypeEqual:
		return left.Equal(right)
	case TokenTypeNotEqual:
		return !left.Equal(right)
	case TokenTypeLessThan:
		return left.Before(right)
	case TokenTypeLessThanOrEqual:
		return left.Before(right) || left.Equal(right)
	case TokenTypeGreaterThan:
		return left.After(right)
	case TokenTypeGreaterThanOrEqual:
		return left.After(right) || left.Equal(right)
	}
	return false
}
