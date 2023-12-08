package core

type PolicyData struct {
	Policy	Policy	`redis:"policy" json:"policy"`
}

type Policy struct {
	Name        string    `redis:"name" json:"name"`
	Description string    `redis:"description" json:"description"`
	Verb    	[]string  `redis:"verbs" json:"verbs"`
}

type RoleData struct {
	Role	Role		`redis:"role" json:"role"`
}

type Resource struct {
	Name        string    `redis:"name" json:"name"`
	Description string    `redis:"description" json:"description"`
}

type Role struct {
	Name        string    `redis:"name" json:"name"`
	Description string    `redis:"description" json:"description"`
	Policies    []string  `redis:"policies" json:"policies"`
	Resources   []string  `redis:"resources" json:"resources"`
}
