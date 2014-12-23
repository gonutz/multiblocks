package game

import (
	"fmt"
	"testing"
)

func TestMovingBlockChangesPointsXAndYs(t *testing.T) {
	b := Block{Points: []Point{{1, 1}, {2, 3}}}
	b.MoveBy(-4, 10)
	checkBlockEquals(t, b, "moved", []Point{{-3, 11}, {-2, 13}})
}

func TestBlockWithoutPointsHasZeroSize(t *testing.T) {
	checkBlockSize(t, Block{}, 0, 0)
}

func TestOnePointBlockHasSizeOne(t *testing.T) {
	checkBlockSize(t, Block{Points: []Point{{0, 0}}}, 1, 1)
}

func TestSizeIsMinimumBoundingRectAroundPoints(t *testing.T) {
	checkBlockSize(t, Block{Points: []Point{{0, 0}, {1, 0}}}, 2, 1)
	checkBlockSize(t, Block{Points: []Point{{0, 0}, {0, 1}}}, 1, 2)
	checkBlockSize(t, Block{Points: []Point{{-2, -9}, {2, -8}}}, 5, 2)
}

func TestRotatingEmptyBlockDoesNothing(t *testing.T) {
	b := Block{Points: nil}
	b.RotationDeltas = [][]Point{[]Point{{0, 0}}}
	b.RotateLeft()
	b.RotateRight()
	if len(b.Points) != 0 {
		t.Error("not empty after rotation")
	}
}

func TestRotatingRightAppliesNextRotationDeltaToPoints(t *testing.T) {
	b := Block{Points: []Point{{0, 0}}}
	b.RotationDeltas = [][]Point{[]Point{{1, 3}}, []Point{{-5, 2}}}
	b.RotateRight()
	if b.Points[0].X != 0+1 || b.Points[0].Y != 0+3 {
		t.Error("(1 3) expected but was", b.Points[0])
	}
	b.RotateRight()
	if b.Points[0].X != 0+1-5 || b.Points[0].Y != 0+3+2 {
		t.Error("(-4 5) expected but was", b.Points[0])
	}
	b.RotateRight()
	if b.Points[0].X != 0+1-5+1 || b.Points[0].Y != 0+3+2+3 {
		t.Error("(-3 8) expected but was", b.Points[0])
	}
}

func TestRotatingWithNoRotationDeltasDoesNothing(t *testing.T) {
	b := Block{Points: []Point{{0, 0}}, RotationDeltas: nil}
	b.RotateRight()
	if b.Points[0].X != 0 || b.Points[0].Y != 0 {
		t.Error("(0 0) expected but was", b.Points[0])
	}
	b.RotateLeft()
	if b.Points[0].X != 0 || b.Points[0].Y != 0 {
		t.Error("(0 0) expected but was", b.Points[0])
	}
}

func TestEachPointIsRotatedRightByItsDelta(t *testing.T) {
	b := Block{Points: []Point{{0, 0}, {5, 3}}}
	b.RotationDeltas = [][]Point{[]Point{{1, 1}, {-1, -1}}}
	b.RotateRight()
	if b.Points[0].X != 1 || b.Points[0].Y != 1 {
		t.Error("(1 1) expected but was", b.Points[0])
	}
	if b.Points[1].X != 4 || b.Points[1].Y != 2 {
		t.Error("(4 2) expected but was", b.Points[1])
	}
}

func TestRotatingLeftUsesNegativeDeltas(t *testing.T) {
	b := Block{Points: []Point{{0, 0}}}
	b.RotationDeltas = [][]Point{[]Point{{1, 3}}, []Point{{-5, 2}}}
	b.RotateLeft()
	if b.Points[0].X != 0+5 || b.Points[0].Y != 0-2 {
		t.Error("(5 -2) expected but was", b.Points[0])
	}
	b.RotateLeft()
	if b.Points[0].X != 0+5-1 || b.Points[0].Y != 0-2-3 {
		t.Error("(4 -5) expected but was", b.Points[0])
	}
	b.RotateLeft()
	if b.Points[0].X != 0+5-1+5 || b.Points[0].Y != 0-2-3-2 {
		t.Error("(9 -7) expected but was", b.Points[0])
	}
}

func TestEachPointIsRotatedLeftByItsDelta(t *testing.T) {
	b := Block{Points: []Point{{0, 0}, {5, 3}}}
	b.RotationDeltas = [][]Point{[]Point{{1, 1}, {-2, -2}}}
	b.RotateLeft()
	if b.Points[0].X != -1 || b.Points[0].Y != -1 {
		t.Error("(1- -1) expected but was", b.Points[0])
	}
	if b.Points[1].X != 7 || b.Points[1].Y != 5 {
		t.Error("(7 5) expected but was", b.Points[1])
	}
}

func TestRotatingBackAndForthResultsInOriginalBlock(t *testing.T) {
	b := Block{Points: []Point{{0, 0}}}
	b.RotationDeltas = [][]Point{[]Point{{1, 3}}, []Point{{-5, 2}}}
	b.RotateLeft()
	b.RotateRight()
	b.RotateLeft()
	b.RotateLeft()
	b.RotateRight()
	b.RotateRight()
	b.RotateLeft()
	b.RotateRight()
	b.RotateRight()
	b.RotateLeft()
	if b.Points[0].X != 0 || b.Points[0].Y != 0 {
		t.Error("(0 0) expected but was", b.Points[0])
	}
}

func TestBlocksCanBeCopied(t *testing.T) {
	original := createTestBlock()
	expected := fmt.Sprint(original)
	copy := original.Copy()
	actual := fmt.Sprint(copy)
	if expected != actual {
		t.Error("\n", expected, "expected but was\n", actual)
	}
}

func createTestBlock() Block {
	return Block{
		Points: []Point{{1, 2}, {3, 4}},
		RotationDeltas: [][]Point{
			[]Point{{5, 6}, {7, 8}},
			[]Point{{9, 10}, {11, 12}},
		},
		rotation: 2,
	}
}

func TestChangingBlockCopyDoesNotAffectOriginal(t *testing.T) {
	original := createTestBlock()
	expected := fmt.Sprint(original)
	copy := original.Copy()
	copy.Points[0] = Point{-1, -1}
	copy.rotation = -1
	copy.RotationDeltas[0][0] = Point{-1, -1}
	copy.RotationDeltas[1][1] = Point{-1, -1}
	actual := fmt.Sprint(original)
	if expected != actual {
		t.Error("\n", expected, "expected but original changed to\n", actual)
	}
}

func checkBlockSize(t *testing.T, b Block, expectedW, expectedH int) {
	if w, h := b.Size(); w != expectedW || h != expectedH {
		t.Error("size should be", expectedW, expectedH, "but was", w, h)
	}
}
