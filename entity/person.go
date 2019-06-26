package entity

// Person represents a person in a family tree. It contains a name,
// their parents and their children.
type Person struct {
	Name string `json:"name"`
	Parents []string `json:"parents, omitempty"`
	Children []string `json:"children, omitempty"`
}

// FamilyTree represents all the people related to a Person in ascendancy.
type FamilyTree struct {
	Name string `json:"name"`
	Relationships []Person `json: "relationships"`
}