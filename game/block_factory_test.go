package game

import (
	"fmt"
	"testing"
)

func TestODoesNotRotate(t *testing.T) {
	O := NewBlockFactory().CreateO()
	expected := []Point{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	checkBlockEquals(t, O, "O", expected)
	O.RotateLeft()
	checkBlockEquals(t, O, "left", expected)
	O.RotateRight()
	checkBlockEquals(t, O, "right", expected)
}

func TestIHasTwoRotations(t *testing.T) {
	I := NewBlockFactory().CreateI()
	flat := []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
	up := []Point{{1, -1}, {1, 0}, {1, 1}, {1, 2}}

	checkBlockEquals(t, I, "I", flat)
	I.RotateRight()
	checkBlockEquals(t, I, "1x right", up)
	I.RotateRight()
	checkBlockEquals(t, I, "2x right", flat)

	I.RotateLeft()
	checkBlockEquals(t, I, "1x left", up)
	I.RotateLeft()
	checkBlockEquals(t, I, "2x left", flat)
}

func TestLHasFourRotations(t *testing.T) {
	L := NewBlockFactory().CreateL()
	down := []Point{{2, 1}, {1, 1}, {0, 1}, {0, 0}}
	left := []Point{{1, 0}, {1, 1}, {1, 2}, {0, 2}}
	up := []Point{{0, 1}, {1, 1}, {2, 1}, {2, 2}}
	right := []Point{{1, 2}, {1, 1}, {1, 0}, {2, 0}}

	checkBlockEquals(t, L, "L", down)
	L.RotateRight()
	checkBlockEquals(t, L, "1x right", left)
	L.RotateRight()
	checkBlockEquals(t, L, "2x right", up)
	L.RotateRight()
	checkBlockEquals(t, L, "3x right", right)
	L.RotateRight()
	checkBlockEquals(t, L, "4x right", down)
	L.RotateLeft()
	checkBlockEquals(t, L, "1x left", right)
	L.RotateLeft()
	checkBlockEquals(t, L, "2x left", up)
	L.RotateLeft()
	checkBlockEquals(t, L, "3x left", left)
	L.RotateLeft()
	checkBlockEquals(t, L, "4x left", down)
}

func TestJHasFourRotations(t *testing.T) {
	J := NewBlockFactory().CreateJ()
	down := []Point{{0, 1}, {1, 1}, {2, 1}, {2, 0}}
	left := []Point{{1, 2}, {1, 1}, {1, 0}, {0, 0}}
	up := []Point{{2, 1}, {1, 1}, {0, 1}, {0, 2}}
	right := []Point{{1, 0}, {1, 1}, {1, 2}, {2, 2}}

	checkBlockEquals(t, J, "J", down)
	J.RotateRight()
	checkBlockEquals(t, J, "1x right", left)
	J.RotateRight()
	checkBlockEquals(t, J, "2x right", up)
	J.RotateRight()
	checkBlockEquals(t, J, "3x right", right)
	J.RotateRight()
	checkBlockEquals(t, J, "4x right", down)
	J.RotateLeft()
	checkBlockEquals(t, J, "1x left", right)
	J.RotateLeft()
	checkBlockEquals(t, J, "2x left", up)
	J.RotateLeft()
	checkBlockEquals(t, J, "3x left", left)
	J.RotateLeft()
	checkBlockEquals(t, J, "4x left", down)
}

func TestTHasFourRotations(t *testing.T) {
	T := NewBlockFactory().CreateT()

	down := []Point{{1, 1}, {0, 1}, {1, 0}, {2, 1}}
	left := []Point{{1, 1}, {1, 2}, {0, 1}, {1, 0}}
	up := []Point{{1, 1}, {2, 1}, {1, 2}, {0, 1}}
	right := []Point{{1, 1}, {1, 0}, {2, 1}, {1, 2}}

	checkBlockEquals(t, T, "T", down)
	T.RotateRight()
	checkBlockEquals(t, T, "1x right", left)
	T.RotateRight()
	checkBlockEquals(t, T, "2x right", up)
	T.RotateRight()
	checkBlockEquals(t, T, "3x right", right)
	T.RotateRight()
	checkBlockEquals(t, T, "4x right", down)
	T.RotateLeft()
	checkBlockEquals(t, T, "1x left", right)
	T.RotateLeft()
	checkBlockEquals(t, T, "2x left", up)
	T.RotateLeft()
	checkBlockEquals(t, T, "3x left", left)
	T.RotateLeft()
	checkBlockEquals(t, T, "4x left", down)
}

func TestSHasTwoRotations(t *testing.T) {
	S := NewBlockFactory().CreateS()
	s := []Point{{0, 0}, {1, 0}, {1, 1}, {2, 1}}
	up := []Point{{1, 0}, {1, 1}, {0, 1}, {0, 2}}

	checkBlockEquals(t, S, "S", s)
	S.RotateRight()
	checkBlockEquals(t, S, "1x right", up)
	S.RotateRight()
	checkBlockEquals(t, S, "2x right", s)
	S.RotateLeft()
	checkBlockEquals(t, S, "1x left", up)
	S.RotateLeft()
	checkBlockEquals(t, S, "2x left", s)
}

func TestZHasTwoRotations(t *testing.T) {
	Z := NewBlockFactory().CreateZ()
	z := []Point{{0, 1}, {1, 1}, {1, 0}, {2, 0}}
	up := []Point{{1, 2}, {1, 1}, {0, 1}, {0, 0}}

	checkBlockEquals(t, Z, "Z", z)
	Z.RotateRight()
	checkBlockEquals(t, Z, "1x right", up)
	Z.RotateRight()
	checkBlockEquals(t, Z, "2x right", z)
	Z.RotateLeft()
	checkBlockEquals(t, Z, "1x left", up)
	Z.RotateLeft()
	checkBlockEquals(t, Z, "2x left", z)
}

func checkBlockEquals(t *testing.T, b Block, msg string, expected []Point) {
	actual := fmt.Sprint(b.Points)
	exp := fmt.Sprint(expected)
	if actual != exp {
		t.Error(msg, "\n", exp, "expected but was", "\n", actual)
	}
}
