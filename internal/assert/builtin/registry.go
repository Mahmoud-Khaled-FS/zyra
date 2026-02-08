package builtin

import (
	"fmt"
)

// EvalFunc defines the signature for all assertion functions.
// 'actual' is the resolved value from the response.
// 'args' are the values passed in the DSL assertion.
type EvalFunc func(actual any, args []any) error

// FunctionRegistry stores all registered assertion functions.
var FunctionRegistry = make(map[string]EvalFunc)

// Register adds a new function to the registry.
// Returns an error if the name already exists.
func Register(name string, fn EvalFunc) error {
	if _, exists := FunctionRegistry[name]; exists {
		return fmt.Errorf("function '%s' already registered", name)
	}
	FunctionRegistry[name] = fn
	return nil
}

// MustRegister panics if registration fails. Useful for built-ins.
func MustRegister(name string, fn EvalFunc) {
	if err := Register(name, fn); err != nil {
		panic(err)
	}
}

// Get retrieves a function by name.
func Get(name string) (EvalFunc, bool) {
	fn, ok := FunctionRegistry[name]
	return fn, ok
}

// InitBuiltin registers all built-in functions.
func InitBuiltin() {
	MustRegister("eq", fnEq)
	MustRegister("ne", fnNe)
	MustRegister("gt", fnGt)
	MustRegister("gte", fnGte)
	MustRegister("lt", fnLt)
	MustRegister("lte", fnLte)
	MustRegister("is", fnIs)
	MustRegister("has", fnHas)
	MustRegister("len", fnLen)
}
