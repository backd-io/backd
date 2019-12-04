package structs

// Index is a simple struct to set the requirements for indexing
type Index struct {
	Fields map[string]interface{}
	Unique bool
}
