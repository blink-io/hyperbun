package schema

type Table[Model any, Columns any] struct {
	Model       *Model
	PrimaryKeys []string
	Schema      string
	Alias       string
	Name        string
	Columns     Columns
}
