# Expronaut: Magic in Expression Evaluation

Welcome to Expronaut, a fun and dynamic Go package designed for parsing and evaluating expressions with a touch of magic. 
Whether you're crafting conditions for template rendering or just dabbling in the alchemy of logic and math, Expronaut is your companion on this adventurous journey through the realms of syntax and semantics.

## Features

- **Dynamic Expression Parsing:** Dive into expressions with variables, nested properties, and an assortment of operators. From simple arithmetic to complex boolean logic, Expronaut ~~understands it all~~, tries to understand.
- **Flexible Variable Context:** Whether you're working with flat landscapes or exploring the depths of nested objects, Expronaut navigates through your data with ease, bringing context to your expressions.
- **Customizable Evaluation:** Tailor the evaluation context to suit your adventure. Pass in variables and watch as Expronaut conjures up the results you seek.

## Getting Started

Embark on your journey with Expronaut by incorporating it into your Go projects. Here's how to get started:

### Installation

Ensure you have Go installed on your system, then fetch the Expronaut package:

```sh
go get github.com/donseba/expronaut
```

## Usage Examples

### complex example

```go
func TestOrderOfOperations(t *testing.T) {
	input := "-2 + 3 * 4 - 5 // 2 ^ 2 << 1 >> 2 % 3"

	out, err := Evaluate(context.TODO(), input)
	if err != nil {
		t.Error(err)
	}

	expected := 4
	if !equalNumber(out, expected) {
		t.Errorf("expected %v, got %v", expected, out)
	}
}
```
This tree visually represents how the operations in the expression are structured and the order in which they would be evaluated, starting from the bottom operations moving up.
```markdown
                              RIGHT_SHIFT
                              /          \
                             /            \
                          LEFT_SHIFT       MODULO
                           /     \          /   \
                          /       \        /     \
                       MINUS       1      2       3
                      /    \
                     /      \
                  PLUS     DIVIDE_INTEGER
                  /   \        /        \
                -2    MULTIPLY          EXPONENT
                        /  \            /      \
                       3    4          2        2
```


### Converting expressions to Go Template Strings
please not this only works for the most basic expressions, for more complex expressions there is still lost to do to make it work.

```go
func TestNewParserVariableBool(t *testing.T) {
   input := `false == foo || ( bar >= baz )`

   lexer := NewLexer(input)

   p := NewParser(lexer)
   tree := p.Parse()

   t.Logf(`%s`, tree.GoTemplate()) // or ( eq false .foo ) ( ge .bar .baz )
}

```

Expronaut transforms your intricate expressions into results or Go template strings, ready for dynamic rendering.

## Supported Operators

### Arithmetic Operators

- **+ (Addition):** Adds two numbers.
- **- (Subtraction):** Subtracts the second number from the first.
- **\* (Multiplication):** Multiplies two numbers.
- **/ (Division):** Divides the first number by the second. Performs floating-point division.
- **// (Integer Division):** Divides the first number by the second, discarding any remainder to return an integer result.
- **% (Modulo):** Returns the remainder of dividing the first number by the second.
- **^ (Exponentiation):** Raises the first number to the power of the second. ( ** is a valid alternative.)

### Bitwise Operators

- **<< (Left Shift):** Shifts the first operand left by the number of bits specified by the second operand.
- **>> (Right Shift):** Shifts the first operand right by the number of bits specified by the second operand.

### Logical Operators

- **&& (Logical AND):** Returns true if both operands are true.
- **|| (Logical OR):** Returns true if at least one of the operands is true.

### Comparison Operators

- **== (Equal):** Returns true if the operands are equal.
- **!= (Not Equal):** Returns true if the operands are not equal.
- **< (Less Than):** Returns true if the first operand is less than the second.
- **<= (Less Than or Equal To):** Returns true if the first operand is less than or equal to the second.
- **> (Greater Than):** Returns true if the first operand is greater than the second.
- **>= (Greater Than or Equal To):** Returns true if the first operand is greater than or equal to the second.

### Builtin functions 
- **sqrt (Square Root):** Calculates the square root of a number (Considered as a function call, sqrt(x)).
- **date (Date):** Returns the current date (Considered as a function call, date()). **"2006-01-02"** is the format.
- **time (Time):** Returns the current time (Considered as a function call, time()). **"15:04"** is the format.
- **datetime (Date Time):** Returns the current date and time (Considered as a function call, datetime()). **"2006-01-02 15:04"** is the format.
- **reduce (Reduce):** Reduces a list of numbers to a single value (Considered as a function call, `reduce(int[1,2,3,4,5],"add", 0)`). The second argument is the function to apply to the list. The first argument is the list of numbers.
- **map (Map):** Applies a function to each element of a list (Considered as a function call, `map(int[1,2,3,4,5], double)`). The second argument is the function to apply to the list. The first argument is the list of numbers.
- 


##  The Complementary BuiltInFunctions (bifs): A Swashbuckling Toolkit

To enhance your Go template experience, Expronaut comes with complementary bifs. This map is a collection of arithmetic functions tailored to work within your templates, allowing you to perform calculations directly:

```go
func init() {
    BuiltinFunctions = bif{}
    
    b := BuiltinFunctions
    
    b["abs"] = b.Abs
    b["add"] = b.Add
    b["date"] = b.Date
    b["datetime"] = b.DateTime
    b["div"] = b.Div
    b["divint"] = b.DivInt
    b["double"] = b.Double
    b["exp"] = b.Exp
    b["filter"] = b.Filter
    b["len"] = b.Len
    b["map"] = b.Map
    b["max"] = b.Max
    b["min"] = b.Min
    b["mod"] = b.Mod
    b["mul"] = b.Mul
    b["reduce"] = b.Reduce
    b["sub"] = b.Sub
    b["sqrt"] = b.Sqrt
    b["sum"] = b.Sum
    b["time"] = b.Time
}
```

These functions are designed to handle both integers and floats, ensuring type compatibility and smooth sailing:

- **abs:** Unveils the absolute essence of a value, stripping away the veil of negativity.
- **date:** Translates a temporal sequence into a calendar date, anchoring fleeting moments.
- **datetime:** Merges date and time, capturing the full spectrum of a moment's presence.
- **div:** Divides two numerical values.
- **divint:** Divides with the precision of integers, discarding any fractional whispers.
- **double:** Echoes a value into twice its magnitude, reflecting its potential.
- **exp:** Elevates numbers to the power of another, scaling the heights of exponential growth.
- **filter:** Sifts through collections with a discerning eye, selecting only those that resonate.
- **len:** Measures the length, revealing the extent of data's expanse.
- **map:** Transforms each element with a spell of modification, rebirthing them anew.
- **max:** Ascends to the peak, finding the pinnacle value in a sea of numbers.
- **min:** Delves into the depths, uncovering the lowest ebb amidst numerical waves.
- **mul:** Fuses values in a dance of multiplication, celebrating their combined strength.
- **reduce:** Weaves through an array with a thread of operation, binding it into a single essence.
- **sub:** Draws apart numbers, navigating the distance between their values.
- **sqrt:** Unravels the square, bringing forth the root from the depths of its square cloister.
- **sum:** Gathers scattered numbers into a collective embrace, uniting them into one.
- **time:** Captures the flow of seconds, minutes, and hours, crystallizing them into a timestamp.

## Example: Summoning Arithmetic in Templates

```go
templ, err := template.New("example").Funcs(funcMap).Parse(`
    Result: {{ add .a .b }} | {{ div .c .d }}
    `)

    if err != nil {
        log.Fatalf("Failed to parse template: %v", err)
    }

data := map[string]any{
    "a": 7,
    "b": 5,
    "c": 10,
    "d": 2,
}

var wr bytes.Buffer
if err := templ.Execute(&wr, data); err != nil {
    log.Fatalf("Failed to execute template: %v", err)
}

fmt.Println(wr.String())
// Output: Result: 12 | 5
```

## Embark on an Expronaut Adventure

Step into the realm of Expronaut, a haven where extensive documentation, comprehensive test cases, and illustrative examples shine a light on the boundless capabilities of this enchanting tool.

Born from the essential need to decode and execute expressions within the intricate universe of PHP Blade templates in Go, Expronaut transcends its initial purpose, morphing into a vessel of exploration and discovery. It's not just a tool; it's a gateway to solving complex challenges and savoring the thrill of expression evaluation, all within the rich landscape of Go programming.

With Expronaut, embark on a voyage where logic seamlessly melds with magic, crafting a world brimming with limitless possibilities. Here, expressions aren't just evaluated—they're brought to life, setting the stage for an epic journey of coding sorcery and innovation.

Let your curiosity be your compass as you navigate through the wonders of Expronaut. Here, in the confluence of practicality and imagination, your Go projects will find their wings. The quest begins now.

## what's next
Integration with LLM's is one of the next steps, this will allow for a more dynamic and flexible way of working with expressions.
It would be nice to do something like : 
```go
	input := `ai("gpt", "is the following 42?", ( 21 + 21 ) )`
```

## License

Expronaut is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
