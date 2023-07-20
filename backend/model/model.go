package model

type Mark int8

const (
	NoMark Mark = iota
	P1Mark
	P2Mark
)

type Winner int8

const (
	NoWinner Winner = iota
	Player1
	Player2
)

type UpdateMessage struct {
	Pos    Pos
	Mark   Mark
	Winner Winner
}

type Field [][]Mark

type Pos struct {
	X, Y int
}

type Model struct {
	Size      uint
	Field     Field
	History   []Pos
	listeners []chan<- (UpdateMessage)
}

func NewModel(size uint) Model {
	field := make(Field, size)
	for i := range field {
		field[i] = make([]Mark, size)
	}

	return Model{
		Size:  size,
		Field: field,
	}
}

func (model *Model) Close() {
	for _, c := range model.listeners {
		close(c)
	}
}

func (model *Model) AddListener(c chan<- (UpdateMessage)) (Mark, Field) {
	model.listeners = append(model.listeners, c)
	if len(model.listeners) > 2 {
		return NoMark, model.Field
	}
	return (Mark)(len(model.listeners)), model.Field
}

func (model *Model) PutMark(p Pos, m Mark) {
	// edit the model
	model.Field[p.X][p.Y] = m
	model.History = append(model.History, p)

	// update the subscribed entities
	message := UpdateMessage{
		Pos:    p,
		Mark:   m,
		Winner: model.CheckWinner(),
	}
	for _, c := range model.listeners {
		c <- message
	}
}

func (model *Model) CheckWinner() Winner {
	if len(model.History) == 0 {
		return NoWinner
	}
	lastPos := model.History[len(model.History)-1]
	posToCheck := [8]Pos{
		{X: lastPos.X - 1, Y: lastPos.Y - 1},
		{X: lastPos.X + 0, Y: lastPos.Y - 1},
		{X: lastPos.X + 1, Y: lastPos.Y - 1},
		{X: lastPos.X - 1, Y: lastPos.Y + 0},
		{X: lastPos.X + 1, Y: lastPos.Y + 0},
		{X: lastPos.X - 1, Y: lastPos.Y + 1},
		{X: lastPos.X + 0, Y: lastPos.Y + 1},
		{X: lastPos.X + 1, Y: lastPos.Y + 1},
	}
	for _, p := range posToCheck {
		// check if p is out of bounds, if so skip it
		if p.X < 0 || p.X > int(model.Size) || p.Y < 0 || p.Y > int(model.Size) {
			continue
		}
		// p is guaranteed in bounds, check to see if it is occupied
		if model.Field[p.X][p.Y] == NoMark {
			return NoWinner
		}
	}
	// all positions are occupied, there is a winner
	if model.Field[lastPos.X][lastPos.Y] == P1Mark {
		return Player2
	} else {
		return Player1
	}
}
