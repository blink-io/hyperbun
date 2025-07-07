package schema

type Schema[Model any, Columns any] struct {
	PrimaryKeys PrimaryKeys
	Label       Label
	Alias       Alias
	Table       Table
	Model       *Model
	Columns     Columns
}

type PrimaryKeys []string

type Alias string

func (t Alias) String() string {
	return string(t)
}

type Label string

func (t Label) String() string {
	return string(t)
}

type Table string

func (t Table) String() string {
	return string(t)
}

type Column string

func (t Column) String() string {
	return string(t)
}
