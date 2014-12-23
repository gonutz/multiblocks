package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestNoBlocksYieldEmptyBoard(t *testing.T) {
	p = newPhysics(BoardSize{4, 3}, BlockCount(0))
	checkBlocks(t, "no blocks",
		"....",
		"....",
		"....")
	checkBoard(t, "empty board",
		"....",
		"....",
		"....")
}

func TestBlocksCanBeAboveTheBoardTop(t *testing.T) {
	p = newPhysics(BoardSize{5, 4}, BlockCount(2))
	p.SetBlock(0, T_at(0, 0))
	p.SetBlock(1, I_at(3, 2))
	checkBlocks(t, "two blocks",
		"...1.",
		"...1.",
		"000..",
		".0...")
}

func TestBlocksCanMoveLeft(t *testing.T) {
	p = newPhysics(BoardSize{5, 2}, BlockCount(1))
	p.SetBlock(0, T_at(1, 0))
	checkBlocks(t, "before",
		".000.",
		"..0..")

	p.MoveLeft(0)
	checkBlocks(t, "moved left",
		"000..",
		".0...")
}

func TestBlocksDoNotMoveIntoLeftWall(t *testing.T) {
	p = newPhysics(BoardSize{4, 2}, BlockCount(1))
	p.SetBlock(0, T_at(0, 0))
	p.MoveLeft(0)
	checkBlocks(t, "left wall hit",
		"000.",
		".0..")
}

func TestBlockDoesNotHitBlockLeftOfIt(t *testing.T) {
	p = newPhysics(BoardSize{4, 2}, BlockCount(2))
	p.SetBlock(0, T_at(0, 0))
	p.SetBlock(1, I_at(3, 0))
	p.MoveLeft(1)
	checkBlocks(t, "left block hit",
		"0001",
		".0.1")
}

func TestBlocksCanMoveRight(t *testing.T) {
	p = newPhysics(BoardSize{5, 2}, BlockCount(1))
	p.SetBlock(0, T_at(1, 0))
	checkBlocks(t, "before",
		".000.",
		"..0..")

	p.MoveRight(0)
	checkBlocks(t, "moved right",
		"..000",
		"...0.")
}

func TestBlocksDoNotMoveIntoRightWall(t *testing.T) {
	p = newPhysics(BoardSize{4, 2}, BlockCount(1))
	p.SetBlock(0, T_at(1, 0))
	p.MoveRight(0)
	checkBlocks(t, "right wall hit",
		".000",
		"..0.")
}

func TestBlockDoesNotHitBlockRightOfIt(t *testing.T) {
	p = newPhysics(BoardSize{4, 2}, BlockCount(2))
	p.SetBlock(0, T_at(0, 0))
	p.SetBlock(1, I_at(3, 0))
	p.MoveRight(0)
	checkBlocks(t, "right block hit",
		"0001",
		".0.1")
}

func TestBlocksCanMoveDown(t *testing.T) {
	p = newPhysics(BoardSize{3, 3}, BlockCount(2))
	p.SetBlock(1, I_at(1, 1))
	checkBlocks(t, "before",
		".1.",
		".1.",
		"...")

	p.MoveDown(1)
	checkBlocks(t, "moved down",
		".1.",
		".1.",
		".1.")
}

func TestBlocksDoNotMoveIntoGround(t *testing.T) {
	p = newPhysics(BoardSize{3, 5}, BlockCount(3))
	p.SetBlock(2, I_at(1, 0))
	p.MoveDown(2)
	checkBlocks(t, "move into ground",
		"...",
		".2.",
		".2.",
		".2.",
		".2.")
}

func TestBlockDoesNotMoveDownIntoOtherBlock(t *testing.T) {
	p = newPhysics(BoardSize{3, 7}, BlockCount(3))
	p.SetBlock(0, T_at(0, 0))
	p.SetBlock(1, T_at(0, 2))
	p.MoveDown(1)
	checkBlocks(t, "move down into block",
		"...",
		"111",
		".1.",
		"000",
		".0.")
}

func TestWallHitIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{4, 5}, BlockCount(2))
	spy := &spyCollisionObserver{}
	p.AddCollisionObserver(spy)
	p.SetBlock(0, T_at(0, 0))
	p.SetBlock(1, T_at(1, 2))
	// ....
	// .111
	// ..1.
	// 000.
	// .0..
	p.MoveRight(0)
	checkIntsEqual(t, spy.horizontalHits, []int{}, "right move hit nothing")
	p.MoveRight(0)
	checkIntsEqual(t, spy.horizontalHits, []int{0}, "right move hit wall")
	p.MoveLeft(1)
	checkIntsEqual(t, spy.horizontalHits, []int{0}, "left move hit nothing")
	p.MoveLeft(1)
	checkIntsEqual(t, spy.horizontalHits, []int{0, 1}, "left move hit wall")
	checkIntsEqual(t, spy.blockHits, []int{}, "blocks were not hit")
	checkIntsEqual(t, spy.groundHits, []int{}, "no block hit the ground")
}

func TestBlockHittingOtherBlockIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{6, 5}, BlockCount(4))
	spy := &spyCollisionObserver{}
	p.AddCollisionObserver(spy)
	p.SetBlock(0, I_at(0, 1))
	p.SetBlock(1, I_at(5, 1))
	p.SetBlock(2, T_at(1, 0))
	p.SetBlock(3, T_at(1, 3))
	// 0333.1
	// 0.3..1
	// 0....1
	// 0222.1
	// ..2...
	p.MoveRight(3)
	checkIntsEqual(t, spy.blockHits, []int{}, "right move hit nothing")
	p.MoveRight(3)
	checkIntsEqual(t, spy.blockHits, []int{3}, "right move hit block")
	p.MoveLeft(1)
	checkIntsEqual(t, spy.blockHits, []int{3, 1}, "left move hit block")
	p.MoveLeft(3)
	checkIntsEqual(t, spy.blockHits, []int{3, 1}, "left move hit nothing")
	p.MoveDown(3)
	checkIntsEqual(t, spy.blockHits, []int{3, 1}, "down move hit nothing")
	p.MoveDown(3)
	checkIntsEqual(t, spy.blockHits, []int{3, 1, 3}, "down move hit block")
	checkIntsEqual(t, spy.horizontalHits, []int{}, "wall was not hit")
	checkIntsEqual(t, spy.groundHits, []int{}, "ground was not hit")
}

func TestBlockHittingGroundIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{2, 3}, BlockCount(2))
	spy := &spyCollisionObserver{}
	p.AddCollisionObserver(spy)
	p.SetBlock(0, I_at(0, 0))
	p.SetBlock(1, I_at(1, 1))
	// 01
	// 01
	// 0.
	p.MoveDown(0)
	checkIntsEqual(t, spy.groundHits, []int{0}, "ground was hit")
	p.MoveDown(1)
	checkIntsEqual(t, spy.groundHits, []int{0}, "ground was not hit")
	p.MoveDown(1)
	checkIntsEqual(t, spy.groundHits, []int{0, 1}, "ground was hit again")
	checkIntsEqual(t, spy.horizontalHits, []int{}, "wall was not hit")
	checkIntsEqual(t, spy.blockHits, []int{}, "blocks were not hit")
}

func TestMovingAndRotatingCanBeObserved(t *testing.T) {
	p := newPhysics(BoardSize{30, 20}, BlockCount(2))
	p.SetBlock(0, I_at(5, 5))
	p.SetBlock(1, T_at(15, 5))
	spy := &spyBlockMoveObserver{}
	p.AddBlockMoveObserver(spy)
	p.MoveLeft(0)
	p.MoveRight(1)
	p.MoveDown(0)
	p.RotateLeft(0)
	p.RotateRight(1)
	p.DropBlocks([]int{0, 1})
	expected := "0 horizontal 1 horizontal 0 down 0 rotated 1 rotated 0 down 1 down "
	if spy.log != expected {
		t.Errorf("expected log was\n'%s'\nbut actual log was\n'%s'", expected, spy.log)
	}
}

func TestNewGameHasEmptyBoard(t *testing.T) {
	p = newPhysics(BoardSize{3, 2}, BlockCount(0))
	checkBoard(t, "empty map",
		"...",
		"...")
}

func TestBlocksCanBeCopiedToBoard(t *testing.T) {
	p = newPhysics(BoardSize{4, 5}, BlockCount(2))
	p.SetBlock(0, T_at(0, 1))
	p.SetBlock(1, I_at(3, 0))
	p.CopyBlockToBoard(0)
	checkBoard(t, "copied T block",
		"....",
		"....",
		"000.",
		".0..",
		"....",
	)
	p.CopyBlockToBoard(1)
	checkBoard(t, "copied T block",
		"....",
		"...1",
		"0001",
		".0.1",
		"...1",
	)
}

func TestCopyingBlockOutsideBoardDoesNothing(t *testing.T) {
	p = newPhysics(BoardSize{3, 2}, BlockCount(1))
	blockBoardWith(0, []Point{{1, 1}, {-1, 0}, {0, 10}})
}

func TestBlockDoesNotMoveHorizontallyIntoSolidBoard(t *testing.T) {
	p = newPhysics(BoardSize{5, 2}, BlockCount(2))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, T_at(1, 0))
	blockBoardWith(1, []Point{{0, 1}, {4, 1}})
	checkBoardAndBlock := func(msg string) {
		checkBlocks(t, "blocks: "+msg,
			".000.",
			"..0..",
		)
		checkBoard(t, "board: "+msg,
			"1...1",
			".....",
		)
	}
	checkBoardAndBlock("original state")

	p.MoveLeft(0)
	checkBoardAndBlock("moved left into solid board")
	checkIntsEqual(t, o.horizontalHits, []int{0}, "board hit moving left")

	p.MoveRight(0)
	checkBoardAndBlock("moved right into solid board")
	checkIntsEqual(t, o.horizontalHits, []int{0, 0}, "board hit moving right")
}

func blockBoardWith(block int, at []Point) {
	p.SetBlock(block, Block{Points: at})
	p.CopyBlockToBoard(block)
	p.SetBlock(block, Block{})
}

func TestBlockDoesNotMoveDownIntoSolidBoard(t *testing.T) {
	p = newPhysics(BoardSize{3, 3}, BlockCount(2))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, T_at(0, 1))
	blockBoardWith(1, []Point{{1, 0}})
	checkBoardAndBlock := func(msg string) {
		checkBlocks(t, "blocks: "+msg,
			"000",
			".0.",
			"...",
		)
		checkBoard(t, "board: "+msg,
			"...",
			"...",
			".1.",
		)
	}
	checkBoardAndBlock("original state")

	p.MoveDown(0)
	checkBoardAndBlock("moved down into solid board")
	checkIntsEqual(t, o.groundHits, []int{0}, "ground hit")
}

func TestRemovingTopLineFillsItWithEmptyBlocks(t *testing.T) {
	p = newPhysics(BoardSize{4, 2}, BlockCount(1))
	blockBoardWith(0, []Point{{0, 0}, {0, 1}, {3, 1}})
	checkBoard(t, "before removing lines",
		"0..0",
		"0...",
	)
	p.RemoveLines(1)
	checkBoard(t, "line 1 removed",
		"....",
		"0...",
	)
}

func TestRemovingConsecutiveTopLinesClearsThem(t *testing.T) {
	p = newPhysics(BoardSize{4, 2}, BlockCount(1))
	blockBoardWith(0, []Point{{0, 0}, {0, 1}, {3, 1}})
	checkBoard(t, "before removing lines",
		"0..0",
		"0...",
	)
	p.RemoveLines(1, 0)
	checkBoard(t, "all lines removed",
		"....",
		"....",
	)
}

func TestRemovingLineDropsLinesAbove(t *testing.T) {
	p = newPhysics(BoardSize{4, 3}, BlockCount(1))
	blockBoardWith(0, []Point{{0, 0}, {0, 1}, {3, 1}, {2, 2}})
	checkBoard(t, "before removing lines",
		"..0.",
		"0..0",
		"0...",
	)
	p.RemoveLines(0)
	checkBoard(t, "line 0 removed",
		"....",
		"..0.",
		"0..0",
	)
}

func TestLinesAreRemovedInOrderFromTopToBottom(t *testing.T) {
	p = newPhysics(BoardSize{4, 4}, BlockCount(1))
	blockBoardWith(0, []Point{{0, 0}, {0, 1}, {3, 1}, {2, 2}, {0, 3}, {1, 3}})
	checkBoard(t, "before removing lines",
		"00..",
		"..0.",
		"0..0",
		"0...",
	)
	p.RemoveLines(0, 2)
	checkBoard(t, "lines 0 and 2 removed",
		"....",
		"....",
		"00..",
		"0..0",
	)
}

func TestRemovingLinesMayCauseBlockCollision(t *testing.T) {
	p = newPhysics(BoardSize{4, 3}, BlockCount(2))
	blockBoardWith(0, []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 2}, {1, 2}})
	p.SetBlock(1, Block{Points: []Point{{1, 1}, {2, 1}}})
	checkBoard(t, "board before",
		"00..",
		"....",
		"0000",
	)
	checkBlocks(t, "blocks before",
		"....",
		".11.",
		"....",
	)

	p.RemoveLines(0)
	checkBoard(t, "board after",
		"....",
		"00..",
		"....",
	)
	checkBlocks(t, "blocks after",
		"....",
		"....",
		".11.",
	)
}

func TestRemovingLineMayCauseCollisionChain(t *testing.T) {
	setupUp5x6BoardWith2Ts()
	checkBoard(t, "board before",
		"..0..",
		".....",
		".....",
		".....",
		".....",
		".....",
	)
	checkBlocks(t, "blocks before",
		".....",
		"..111",
		"0001.",
		".0...",
		".....",
		".....",
	)

	p.RemoveLines(1)
	checkBoard(t, "board after",
		".....",
		"..0..",
		".....",
		".....",
		".....",
		".....",
	)
	checkBlocks(t, "blocks after",
		".....",
		".....",
		"..111",
		"0001.",
		".0...",
		".....",
	)
}

func setupUp5x6BoardWith2Ts() {
	p = newPhysics(BoardSize{5, 6}, BlockCount(2))
	blockBoardWith(0, []Point{{2, 5}})
	p.SetBlock(0, T_at(0, 2))
	p.SetBlock(1, T_at(2, 3))
}

func TestCollisionsFromLineRemovalAreObserved(t *testing.T) {
	setupUp5x6BoardWith2Ts()
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.RemoveLines(0)
	checkIntsEqual(t, o.draggedDown, []int{1, 0}, "2 blocks dragged by line removal")
}

func TestEveryLineCollisionIsObservedSeparately(t *testing.T) {
	setupUp5x6BoardWith2Ts()
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.RemoveLines(0, 1)
	checkIntsEqual(t, o.draggedDown, []int{1, 0, 1, 0}, "blocks dragged down twice")
}

func TestMultipleBlocksCanDropAtOnce(t *testing.T) {
	p = newPhysics(BoardSize{4, 3}, BlockCount(2))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, T_at(0, 1))
	p.SetBlock(1, I_at(3, 1))
	checkBlocks(t, "blocks before",
		"0001",
		".0.1",
		"....",
	)

	p.DropBlocks([]int{0, 1})
	checkBlocks(t, "dropped blocks",
		"...1",
		"0001",
		".0.1",
	)
	checkNothingWasHit(t, o)
}

func TestSingleBlocksMayCollideWhileDropping(t *testing.T) {
	p = newPhysics(BoardSize{4, 3}, BlockCount(2))
	blockBoardWith(1, []Point{{1, 0}})
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, T_at(0, 1))
	p.SetBlock(1, I_at(3, 1))
	checkBlocks(t, "blocks before",
		"0001",
		".0.1",
		"....",
	)
	checkBoard(t, "board before",
		"....",
		"....",
		".1..",
	)

	p.DropBlocks([]int{0, 1})
	checkBlocks(t, "only block 1 could drop",
		"0001",
		".0.1",
		"...1",
	)
	checkIntsEqual(t, o.groundHits, []int{0}, "block 0 was blocked")
	checkIntsEqual(t, o.blockHits, []int{}, "no block collision")
	checkIntsEqual(t, o.horizontalHits, []int{}, "no horizontal collision")
}

func TestCollisionChainMayOccurWhileDropping(t *testing.T) {
	p = newPhysics(BoardSize{5, 7}, BlockCount(4))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, T_at(0, 0))
	p.SetBlock(1, T_at(1, 2))
	p.SetBlock(2, I_at(1, 3))
	p.SetBlock(3, I_at(4, 2))
	checkBlocks(t, "blocks before",
		".2...",
		".2..3",
		".2..3",
		".1113",
		"..1.3",
		"000..",
		".0...",
	)

	p.DropBlocks([]int{0, 1, 2})
	checkBlocks(t, "blocks before",
		".2...",
		".2..3",
		".2..3",
		".1113",
		"..1.3",
		"000..",
		".0...",
	)
	checkIntsEqual(t, o.groundHits, []int{0, 1, 2}, "ground hits")
	checkIntsEqual(t, o.blockHits, []int{}, "no block collision")
	checkIntsEqual(t, o.horizontalHits, []int{}, "no horizontal collision")
}

func TestBlockCanBeRotated(t *testing.T) {
	p = newPhysics(BoardSize{3, 3}, BlockCount(1))
	p.SetBlock(0, rotatingPointAt(0, 0))
	checkBlocks(t, "original",
		"...",
		"...",
		"0..",
	)
	p.RotateRight(0)
	checkBlocks(t, "rotated",
		"...",
		".0.",
		"...",
	)
	p.RotateLeft(0)
	checkBlocks(t, "rotated back",
		"...",
		"...",
		"0..",
	)
}

func TestWallBlockingRotationIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{2, 2}, BlockCount(1))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, rotatingPoint(0, 1, -1, 0))
	p.RotateRight(0)
	checkBlocks(t, "hit left wall",
		"0.",
		"..",
	)
	checkIntsEqual(t, o.rotationHits, []int{0}, "hit wall rotating right")

	p.SetBlock(0, rotatingPoint(0, 1, 1, 0))
	p.RotateLeft(0)
	checkBlocks(t, "hit left wall",
		"0.",
		"..",
	)
	checkIntsEqual(t, o.rotationHits, []int{0, 0}, "hit wall rotating left")
}

func TestGroundBlockingRotationIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{2, 2}, BlockCount(1))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, rotatingPoint(1, 0, 0, -1))
	p.RotateRight(0)
	checkBlocks(t, "hit ground",
		"..",
		".0",
	)
	checkIntsEqual(t, o.rotationHits, []int{0}, "hit ground rotating right")
}

func TestBoardBlockingRotationIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{2, 2}, BlockCount(1))
	blockBoardWith(0, []Point{{0, 1}})
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, rotatingPoint(0, 0, 0, 1))
	p.RotateRight(0)
	checkBoard(t, "board",
		"0.",
		"..",
	)
	checkBlocks(t, "hit board",
		"..",
		"0.",
	)
	checkIntsEqual(t, o.rotationHits, []int{0}, "hit board rotating right")
}

func TestOtherBlockBlockingRotationIsObserved(t *testing.T) {
	p = newPhysics(BoardSize{2, 2}, BlockCount(2))
	o := &spyCollisionObserver{}
	p.AddCollisionObserver(o)
	p.SetBlock(0, rotatingPoint(0, 0, 0, 0))
	p.SetBlock(1, rotatingPoint(1, 1, -1, -1))
	p.RotateRight(1)
	checkBlocks(t, "hit board",
		".1",
		"0.",
	)
	checkIntsEqual(t, o.rotationHits, []int{1}, "hit other block rotating right")
}

// auxiliary test variables and functions start here

var p *physics

func checkBlocks(t *testing.T, msg string, blockMap ...string) {
	for invY, row := range blockMap {
		y := len(blockMap) - 1 - invY
		for x := range row {
			switch {
			case row[x] == '.':
				checkEmpty(t, msg, x, y)
			case row[x] >= '0' && row[x] <= '9':
				checkBlock(t, msg, x, y, int(row[x])-'0')
			default:
				panic(msg + ": unknown rune")
			}
		}
	}
}

func checkEmpty(t *testing.T, msg string, x, y int) {
	for blockIndex, b := range p.Blocks() {
		for _, p := range b.Points {
			if p.X == x && p.Y == y {
				t.Error(msg+":", x, y, "blocked by", blockIndex)
			}
		}
	}
}

func checkBlock(t *testing.T, msg string, x, y int, block int) {
	for _, p := range p.Blocks()[block].Points {
		if p.X == x && p.Y == y {
			return
		}
	}
	t.Error(msg+": no block", block, "at", x, y)
}

func T_at(x, y int) Block {
	return Block{Points: []Point{
		{x + 1, y + 0},
		{x + 0, y + 1},
		{x + 1, y + 1},
		{x + 2, y + 1}}}
}

func I_at(x, y int) Block {
	return Block{Points: []Point{
		{x, y + 0},
		{x, y + 1},
		{x, y + 2},
		{x, y + 3}}}
}

type spyCollisionObserver struct {
	horizontalHits []int
	blockHits      []int
	groundHits     []int
	draggedDown    []int
	rotationHits   []int
}

func (spy *spyCollisionObserver) BlockHitLeftOrRight(block int) {
	spy.horizontalHits = append(spy.horizontalHits, block)
}

func (spy *spyCollisionObserver) BlockHitOtherBlock(block int) {
	spy.blockHits = append(spy.blockHits, block)
}

func (spy *spyCollisionObserver) BlockHitGround(block int) {
	spy.groundHits = append(spy.groundHits, block)
}

func (spy *spyCollisionObserver) BlockDraggedDownByLineRemoval(block int) {
	spy.draggedDown = append(spy.draggedDown, block)
}

func (spy *spyCollisionObserver) BlockCouldNotRotate(block int) {
	spy.rotationHits = append(spy.rotationHits, block)
}

func checkIntsEqual(t *testing.T, actual, expected []int, msg string) {
	if fmt.Sprint(expected) != fmt.Sprint(actual) {
		t.Error(msg, expected, "expected but was", actual)
	}
}

func checkBoard(t *testing.T, msg string, fields ...string) {
	expected := "\n" + strings.Join(fields, "\n") + "\n"
	board := p.Board()
	actual := "\n"
	w, h := board.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			field := board.At(x, h-1-y)
			if field == NoPlayer {
				actual += "."
			} else {
				actual += fmt.Sprint(field)
			}
		}
		actual += "\n"
	}
	if actual != expected {
		t.Error(msg, expected, "expected but was", actual)
	}
}

func checkNothingWasHit(t *testing.T, o *spyCollisionObserver) {
	checkIntsEqual(t, o.blockHits, nil, "no block hits")
	checkIntsEqual(t, o.groundHits, nil, "no ground hits")
	checkIntsEqual(t, o.horizontalHits, nil, "no horizontal hits")
}

func rotatingPointAt(x, y int) Block {
	return Block{
		Points:         []Point{{x, y}},
		RotationDeltas: [][]Point{[]Point{{1, 1}}}}
}

func rotatingPoint(x, y, rotX, rotY int) Block {
	return Block{
		Points:         []Point{{x, y}},
		RotationDeltas: [][]Point{[]Point{{rotX, rotY}}}}
}

type spyBlockMoveObserver struct {
	log string
}

func (spy *spyBlockMoveObserver) BlockMovedHorizontally(block int) {
	spy.log += fmt.Sprintf("%v horizontal ", block)
}

func (spy *spyBlockMoveObserver) BlockMovedDown(block int) {
	spy.log += fmt.Sprintf("%v down ", block)
}

func (spy *spyBlockMoveObserver) BlockRotated(block int) {
	spy.log += fmt.Sprintf("%v rotated ", block)
}
