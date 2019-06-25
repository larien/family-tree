package entity

// Person represents a person in a family tree. It contains a name,
// their parents and their children.
type Person struct {
	Name string `json:"name"`
	Parents []string `json:"parents"`
	Children []string `json:"children"`
}