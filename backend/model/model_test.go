package model

import "testing"

func TestNewModel(t *testing.T) {
	model := NewModel(8)
	for _, r := range model.Field {
		for _, m := range r {
			if m != NoMark {
				t.Fatal("NewModel failed to create clean field")
			}
		}
	}
}
func TestPutMark(t *testing.T) {
	model := NewModel(8)
	model.PutMark(Pos{X: 0, Y: 0}, P1Mark)
	if model.Field[0][0] != P1Mark {
		t.Fatal("PutMark failed to update field")
	}
}

func TestCheckWinner(t *testing.T) {
	model := NewModel(8)
	if model.CheckWinner() != NoWinner {
		t.Fatal("CheckWinner failed to identify no winner")
	}
	model.PutMark(Pos{X: 0, Y: 0}, P1Mark)
	if model.CheckWinner() != NoWinner {
		t.Fatal("CheckWinner failed to identify no winner")
	}
	for x, r := range model.Field {
		for y := range r {
			model.PutMark(Pos{x, y}, P1Mark)
		}
	}
	if model.CheckWinner() != Player2 {
		t.Fatal("CheckWinner failed to identify P2 winner")
	}
	for x, r := range model.Field {
		for y := range r {
			model.PutMark(Pos{x, y}, P2Mark)
		}
	}
	if model.CheckWinner() != Player1 {
		t.Fatal("CheckWinner failed to identify P1 winner")
	}
}
