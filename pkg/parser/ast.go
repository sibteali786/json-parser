// pkg/parser/ast.go
package parser

// JSONValue represents any JSON value
type JSONValue interface {
	String() string
}

// JSONObject represents a JSON object
type JSONObject struct {
	Pairs map[string]JSONValue
}

func (obj *JSONObject) String() string {
	return "JSONObject"
}

// JSONString represents a JSON string
type JSONString struct {
	Value string
}

func (s *JSONString) String() string {
	return s.Value
}

// JSONNumber represents a JSON number
type JSONNumber struct {
	Value string
}

func (n *JSONNumber) String() string {
	return n.Value
}

// JSONBoolean represents a JSON boolean
type JSONBoolean struct {
	Value bool
}

func (b *JSONBoolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// JSONNull represents JSON null
type JSONNull struct{}

func (n *JSONNull) String() string {
	return "null"
}

// JSONArray represents a JSON array (we'll need this for later steps)
type JSONArray struct {
	Elements []JSONValue
}

func (a *JSONArray) String() string {
	return "JSONArray"
}
