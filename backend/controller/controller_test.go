package controller

import (
	"testing"

	"constraint/model"
)

func genController() Controller {
	model := model.NewModel(8)
	return NewController(&model)
}

func TestNewModel(t *testing.T) {
	contr := genController()
	if contr.model == nil {
		t.Fatal("NewController failed to embed model")
	}
}

func TestIsSafePos(t *testing.T) {
	if isSafePos(model.Pos{X: -1, Y: -1}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: -1, Y: +1}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: +1, Y: +1}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: +1, Y: -1}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: -1, Y: 0}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: +1, Y: 0}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: 0, Y: -1}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: 0, Y: +1}, 1) != false {
		t.Fatal("isSafePos failed to protect from unsafe position")
	}
	if isSafePos(model.Pos{X: 0, Y: 0}, 1) != true {
		t.Fatal("isSafePos failed to allow safe position")
	}
}

func TestIsMarkAdjactent(t *testing.T) {
	lastPos := model.Pos{X: 2, Y: 2}
	posToCheckTrue := []model.Pos{
		{X: lastPos.X - 1, Y: lastPos.Y - 1},
		{X: lastPos.X + 0, Y: lastPos.Y - 1},
		{X: lastPos.X + 1, Y: lastPos.Y - 1},

		{X: lastPos.X - 1, Y: lastPos.Y + 0},
		{X: lastPos.X + 1, Y: lastPos.Y + 0},

		{X: lastPos.X - 1, Y: lastPos.Y + 1},
		{X: lastPos.X + 0, Y: lastPos.Y + 1},
		{X: lastPos.X + 1, Y: lastPos.Y + 1},
	}
	for _, p := range posToCheckTrue {
		if isMarkAdjacent(p, lastPos) != true {
			t.Fatal("isMarkAdjacent failed to recognise valid position")
		}
	}
	posToCheckFalse := []model.Pos{
		{X: lastPos.X - 2, Y: lastPos.Y - 2},
		{X: lastPos.X - 2, Y: lastPos.Y - 1},
		{X: lastPos.X - 2, Y: lastPos.Y + 0},
		{X: lastPos.X - 2, Y: lastPos.Y + 1},
		{X: lastPos.X - 2, Y: lastPos.Y + 2},

		{X: lastPos.X + 2, Y: lastPos.Y - 2},
		{X: lastPos.X + 2, Y: lastPos.Y - 1},
		{X: lastPos.X + 2, Y: lastPos.Y + 0},
		{X: lastPos.X + 2, Y: lastPos.Y + 1},
		{X: lastPos.X + 2, Y: lastPos.Y + 2},

		{X: lastPos.X - 1, Y: lastPos.Y - 2},
		{X: lastPos.X + 0, Y: lastPos.Y - 2},
		{X: lastPos.X + 1, Y: lastPos.Y - 2},

		{X: lastPos.X - 1, Y: lastPos.Y + 2},
		{X: lastPos.X + 0, Y: lastPos.Y + 2},
		{X: lastPos.X + 1, Y: lastPos.Y + 2},
	}
	for _, p := range posToCheckFalse {
		if isMarkAdjacent(p, lastPos) != false {
			t.Fatal("isMarkAdjacent failed to recognise invalid position")
		}
	}
}

func TestAddMark(t *testing.T) {
	contr := genController()

	// test first round errors
	if contr.AddMark(model.Pos{X: 0, Y: 0}, model.P2Mark) == nil {
		t.Fatal("AddMark failed to recognise wrong turn")
	}
	if contr.AddMark(model.Pos{X: -1, Y: 0}, model.P1Mark) == nil {
		t.Fatal("AddMark failed to recognise unsafe position")
	}
	if contr.AddMark(model.Pos{X: 0, Y: 0}, model.P1Mark) != nil {
		t.Fatal("AddMark failed to allow legal first move")
	}

	// test all other edge cases
	if contr.AddMark(model.Pos{X: 0, Y: 0}, model.NoMark) == nil {
		t.Fatal("AddMark failed to recognise spectator")
	}
	if contr.AddMark(model.Pos{X: 0, Y: 0}, model.P2Mark) == nil {
		t.Fatal("AddMark failed to recognise invalid position")
	}
	if contr.AddMark(model.Pos{X: -1, Y: 0}, model.P2Mark) == nil {
		t.Fatal("AddMark failed to recognise unsafe position")
	}
	if contr.AddMark(model.Pos{X: 0, Y: 0}, model.P1Mark) == nil {
		t.Fatal("AddMark failed to recognise wrong turn")
	}

	// continue game
	if contr.AddMark(model.Pos{X: 1, Y: 0}, model.P2Mark) != nil {
		t.Fatal("AddMark failed to allow legal move")
	}
	if contr.AddMark(model.Pos{X: 1, Y: 1}, model.P1Mark) != nil {
		t.Fatal("AddMark failed to allow legal move")
	}
	if contr.AddMark(model.Pos{X: 1, Y: 2}, model.P2Mark) != nil {
		t.Fatal("AddMark failed to allow legal move")
	}
	if contr.AddMark(model.Pos{X: 0, Y: 2}, model.P1Mark) != nil {
		t.Fatal("AddMark failed to allow legal move")
	}
	if contr.AddMark(model.Pos{X: 0, Y: 1}, model.P2Mark) != nil {
		t.Fatal("AddMark failed to allow legal move")
	}

	// check game over condition
	if contr.AddMark(model.Pos{X: 3, Y: 3}, model.P1Mark) == nil {
		t.Fatal("AddMark failed to recognise game over")
	}
}
