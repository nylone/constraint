package model

type Mark int8

const (
	NoMark Mark = iota
	P1Mark
	P2Mark
)

type Field [][]Mark

type Pos struct {
	X, Y int
}

type Model struct {
	Size    uint
	Field   Field
	History []Pos
}

func NewModel(size uint) Model {
	field := make(Field, size)
	for i := range field {
		field[i] = make([]Mark, size)
	}

	return Model{
		Size:    size,
		Field:   field,
		History: make([]Pos, 0),
	}
}

func (model *Model) PutMark(p Pos, m Mark) {
	model.Field[p.X][p.Y] = m
	model.History = append(model.History, p)
}
