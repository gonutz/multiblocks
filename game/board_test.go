package game

import "testing"

func TestNewBoardHasGivenSize(t *testing.T) {
	b := newBoard(5, 9)
	w, h := b.Size()
	if w != 5 || h != 9 {
		t.Error("wrong size", w, h)
	}
}

func TestNewBoardHasAll_NoPlayer_Fields(t *testing.T) {
	b := newBoard(3, 5)
	for x := 0; x < 3; x++ {
		for y := 0; y < 5; y++ {
			if b.At(x, y) != NoPlayer {
				t.Fatal("board not empty but", b)
			}
		}
	}
}

func Test_player_can_be_set_at_position(t *testing.T) {
	b := newBoard(5, 10)
	b.SetAt(0, 0, 123)
	b.SetAt(4, 9, 999)
	if player := b.At(0, 0); player != 123 {
		t.Error("position not set but was", player)
	}
	if player := b.At(4, 9); player != 999 {
		t.Error("position not set but was", player)
	}
}

func TestPlayerCanBeSetAtPosition(t *testing.T) {
	b := newBoard(5, 10)
	b.SetAt(0, 0, 123)
	b.SetAt(4, 9, 999)
	if player := b.At(0, 0); player != 123 {
		t.Error("position not set but was", player)
	}
	if player := b.At(4, 9); player != 999 {
		t.Error("position not set but was", player)
	}
}

func TestBoardCanBeCopied(t *testing.T) {
	b := newBoard(4, 2)
	b.SetAt(0, 0, 123)
	copy := b.Copy()
	if player := copy.At(0, 0); player != 123 {
		t.Error("copy differes from original:", copy)
	}
}

func TestChangingBoardCopyDoesNotChangeOriginal(t *testing.T) {
	original := newBoard(4, 5)
	original.SetAt(2, 3, 555)
	copy := original.Copy()
	copy.SetAt(2, 3, -10)
	if player := original.At(2, 3); player != 555 {
		t.Error("original changed to", player)
	}
}
