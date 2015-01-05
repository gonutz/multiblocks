package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestBoardSizeIsSetPerPlayerCount(t *testing.T) {
	logic := NewLogic(alwaysReturn(block()))
	logic.SetBoardSizeForPlayerCount(1, BoardSize{10, 18})
	logic.SetBlockStartPositions(1, []Point{{0, 0}})
	logic.SetBoardSizeForPlayerCount(2, BoardSize{12, 20})
	logic.SetBlockStartPositions(2, []Point{{0, 0}, {0, 0}})

	logic.StartNewGame(1)
	if w, h := logic.Board().Size(); w != 10 || h != 18 {
		t.Error("1 player game had size", w, h)
	}
	logic.StartNewGame(2)
	if w, h := logic.Board().Size(); w != 12 || h != 20 {
		t.Error("2 player game had size", w, h)
	}
}

func TestBlockStartPositionsAreSetPerPlayerCount(t *testing.T) {
	logic := NewLogic(alwaysReturn(block(0, 0)))

	logic.SetBoardSizeForPlayerCount(1, BoardSize{5, 3})
	logic.SetBlockStartPositions(1, []Point{{3, 2}})
	logic.StartNewGame(1)
	checkGame(t, logic, "1 player game",
		"...0.",
		".....",
		".....",
	)

	logic.SetBoardSizeForPlayerCount(2, BoardSize{7, 3})
	logic.SetBlockStartPositions(2, []Point{{1, 2}, {4, 0}})
	logic.StartNewGame(2)
	checkGame(t, logic, "2 player game",
		".0.....",
		".......",
		"....1..",
	)
}

func TestBlocksAreXCenteredOnStartPosition(t *testing.T) {
	b := block(0, 0)
	logic := NewLogic(func() Block { return b.Copy() })
	logic.SetBoardSizeForPlayerCount(1, BoardSize{7, 1})
	logic.SetBlockStartPositions(1, []Point{{3, 0}})

	b = block(0, 0)
	logic.StartNewGame(1)
	checkGame(t, logic, "width = 1", "...0...")

	b = block(0, 0, 1, 0)
	logic.StartNewGame(1)
	checkGame(t, logic, "width = 2", "..00...")

	b = block(0, 0, 1, 0, 2, 0)
	logic.StartNewGame(1)
	checkGame(t, logic, "width = 3", "..000..")
}

func TestNewGameCreatesNewPreviewBlock(t *testing.T) {
	l := NewLogic(increasingYBlocks(1))
	l.SetBlockStartPositions(1, []Point{{10, 10}})
	l.StartNewGame(1)
	checkBlocksEqual(t, l.PreviewBlocks(), block(0, 2))
}

func increasingYBlocks(firstY int) BlockFactory {
	y := firstY - 1
	return func() Block {
		y++
		return Block{Points: []Point{{0, y}}}
	}
}

func TestOnUpdate_LeftInputIsHandled(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{4, 1}, []Point{{1, 0}, {3, 0}})
	logic.StartNewGame(2)
	checkGame(t, logic, "original",
		".0.1",
	)
	logic.Update(InputEvent{0, LeftPressed})
	logic.Update(InputEvent{1, LeftPressed})
	checkGame(t, logic, "both left",
		"0.1.",
	)
	logic.Update(InputEvent{0, LeftPressed})
	logic.Update(InputEvent{1, LeftPressed})
	checkGame(t, logic, "2x both left",
		"01..",
	)
	logic.Update(InputEvent{0, LeftPressed})
	logic.Update(InputEvent{1, LeftPressed})
	checkGame(t, logic, "3x both left",
		"01..",
	)
}

func TestOnUpdate_RightInputIsHandled(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{4, 1}, []Point{{0, 0}, {2, 0}})
	logic.StartNewGame(2)
	checkGame(t, logic, "original",
		"0.1.",
	)
	logic.Update(InputEvent{0, RightPressed}, InputEvent{0, RightReleased})
	logic.Update(InputEvent{1, RightPressed}, InputEvent{0, RightReleased})
	checkGame(t, logic, "both right",
		".0.1",
	)
	logic.Update(InputEvent{0, RightPressed}, InputEvent{0, RightReleased})
	logic.Update(InputEvent{1, RightPressed}, InputEvent{0, RightReleased})
	checkGame(t, logic, "2x both right",
		"..01",
	)
	logic.Update(InputEvent{0, RightPressed}, InputEvent{0, RightReleased})
	logic.Update(InputEvent{1, RightPressed}, InputEvent{0, RightReleased})
	checkGame(t, logic, "3x both right",
		"..01",
	)
}

func TestOnUpdate_DownInputIsHandled(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{2, 4}, []Point{{0, 3}, {0, 1}})
	logic.StartNewGame(2)
	checkGame(t, logic, "original",
		"0.",
		"..",
		"1.",
		"..",
	)
	logic.Update(InputEvent{0, DownPressed}, InputEvent{0, DownReleased})
	logic.Update(InputEvent{1, DownPressed}, InputEvent{0, DownReleased})
	checkGame(t, logic, "both down",
		"..",
		"0.",
		"..",
		"1.",
	)
	logic.Update(InputEvent{0, DownPressed}, InputEvent{0, DownReleased})
	checkGame(t, logic, "2x both down",
		"..",
		"..",
		"0.",
		"1.",
	)
}

func TestOnUpdate_RotationIsHandled(t *testing.T) {
	b := block(0, 0)
	b.RotationDeltas = [][]Point{[]Point{{1, 1}}}
	logic := NewLogic(alwaysReturn(b))
	logic.SetBoardSizeForPlayerCount(2, BoardSize{3, 3})
	logic.SetBlockStartPositions(2, []Point{{0, 0}, {1, 1}})
	logic.StartNewGame(2)

	logic.Update(InputEvent{0, RotateRight})
	checkGame(t, logic, "right rotation blocked",
		"...",
		".1.",
		"0..",
	)

	logic.Update(InputEvent{1, RotateRight}, InputEvent{0, RotateRight})
	checkGame(t, logic, "rotated both right",
		"..1",
		".0.",
		"...",
	)

	logic.Update(InputEvent{0, RotateLeft})
	checkGame(t, logic, "0 rotated left",
		"..1",
		"...",
		"0..",
	)
}

func TestNewGameResetsDropTimer(t *testing.T) {
	logic, spy := createSpyDropTimerLogic()
	checkInt(t, spy.reset, 0, "should not be reset yet")
	logic.StartNewGame(1)
	checkInt(t, spy.reset, 1, "should be reset")
	logic.StartNewGame(1)
	checkInt(t, spy.reset, 2, "should be reset again")
}

func TestOnUpdateTheDropTimerIsUpdated(t *testing.T) {
	logic, spy := createSpyDropTimerLogic()
	logic.StartNewGame(1)
	checkInt(t, spy.updated, 0, "updated")
	logic.Update()
	checkInt(t, spy.updated, 1, "updated")
	logic.Update()
	checkInt(t, spy.updated, 2, "updated")
}

func createSpyDropTimerLogic() (*Logic, *spyDropTimer) {
	logic := NewLogic(alwaysReturn(block()))
	spy := &spyDropTimer{}
	logic.SetDropTimer(spy)
	logic.SetBoardSizeForPlayerCount(1, BoardSize{1, 1})
	logic.SetBlockStartPositions(1, []Point{{0, 0}})
	return logic, spy
}

func TestDropIsDoneAfterOtherEventsIfTimeForIt(t *testing.T) {
	logic := createSingleBlockGame(3, BoardSize{4, 2}, []Point{{0, 1}, {1, 1}, {1, 0}})
	timer := &spyDropTimer{}
	logic.SetDropTimer(timer)
	logic.StartNewGame(3)
	checkGame(t, logic, "original",
		"01..",
		".2..",
	)
	logic.Update()
	checkGame(t, logic, "no drop yet",
		"01..",
		".2..",
	)
	timer.isTimeForDrop = true
	logic.Update(InputEvent{2, RightPressed})
	checkGame(t, logic, "after move and drop",
		"....",
		"012.",
	)
}

func TestOnceDroppedIntoGroundBlockDoesNotMoveAnymore(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{2, 2}, []Point{{0, 1}})
	logic.StartNewGame(1)
	checkGame(t, logic, "original",
		"0.",
		"..",
	)
	logic.Update(
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{0, RightPressed},
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{0, LeftPressed})
	checkGame(t, logic, "hit the ground at 1,0",
		"..",
		".0",
	)
}

func TestDroppedBlocksAreSolidifiedAndReset(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{4, 2}, []Point{{0, 1}, {1, 1}})
	logic.StartNewGame(2)
	checkGame(t, logic, "original",
		"01..",
		"....",
	)
	logic.Update(
		InputEvent{1, RightPressed}, InputEvent{1, RightReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
	)
	checkGame(t, logic, "all dropped",
		"....",
		"0.1.",
	)
	logic.Update()
	checkGame(t, logic, "all solidified and reset",
		"01..",
		"0.1.",
	)
}

func TestResetBlockMovesUpIfThereIsACollision(t *testing.T) {
	logic := NewLogic(alwaysReturn(block(0, 0, 0, 1)))
	logic.SetBoardSizeForPlayerCount(2, BoardSize{2, 6})
	logic.SetBlockStartPositions(2, []Point{{0, 2}, {1, 2}})
	logic.StartNewGame(2)
	checkGame(t, logic, "original",
		"..",
		"..",
		"01",
		"01",
		"..",
		"..",
	)
	logic.Update(
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{1, LeftPressed}, InputEvent{1, LeftReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{0, RightPressed}, InputEvent{0, RightReleased},
	)
	checkGame(t, logic, "1 hit ground and 0 in 1's position",
		"..",
		"..",
		".0",
		".0",
		"1.",
		"1.",
	)

	logic.Update()
	checkGame(t, logic, "1 reset higher because 0 blocked its position",
		".1",
		".1",
		".0",
		".0",
		"1.",
		"1.",
	)
}

func TestNewGameClearsBoard(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{3, 2}, []Point{{0, 1}, {1, 1}})
	logic.SetDropTimer(&spyDropTimer{isTimeForDrop: true})
	logic.StartNewGame(2)
	checkGame(t, logic, "original",
		"01.",
		"...",
	)
	logic.Update()
	logic.Update()
	logic.Update()
	checkGame(t, logic, "all solidified and reset",
		"01.",
		"01.",
	)

	logic.StartNewGame(2)
	checkGame(t, logic, "clear board after restart",
		"01.",
		"...",
	)
}

func TestOnlyEventsForExistingPlayersAreHandled(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{3, 1}, []Point{{1, 0}})
	logic.StartNewGame(1)
	logic.Update(InputEvent{0, RightPressed}, InputEvent{1, LeftPressed})
	checkGame(t, logic, "player 1 does not exist", "..0")
}

func TestFullLineStartsAnimation(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{1, 2}, []Point{{0, 1}})
	spy := &spyLineAnimation{}
	logic.SetLineAnimation(spy)
	logic.StartNewGame(1)
	logic.Update(InputEvent{0, DownPressed})
	if len(spy.lines) != 0 {
		t.Error("animation started too early, lines are", spy.lines)
	}
	logic.Update(InputEvent{0, DownPressed})
	if len(spy.lines) != 1 {
		t.Error("animation not started, lines are", spy.lines)
	}
	checkIntsEqual(t, spy.lines[0], []int{0}, "lines")
}

func TestWhileAnimationRunsBlocksDoNotMove(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{1, 3}, []Point{{0, 2}})
	spy := &spyLineAnimation{running: false}
	logic.SetLineAnimation(spy)
	logic.StartNewGame(1)
	logic.Update(InputEvent{0, DownPressed})
	checkGame(t, logic, "animation not running",
		".",
		"0",
		".",
	)

	spy.running = true
	logic.Update(InputEvent{0, DownPressed})
	checkGame(t, logic, "no drop while animation runs",
		".",
		"0",
		".",
	)
	if spy.updated != 1 {
		t.Error("animation not updated while running")
	}
}

func TestFullLinesAreRemovedAfterAnimation(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{2, 3}, []Point{{0, 2}, {1, 2}})
	spy := &spyLineAnimation{running: false}
	logic.SetLineAnimation(spy)
	logic.StartNewGame(2)
	logic.Update(
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
	)
	checkGame(t, logic, "animation not running",
		"..",
		"..",
		"01",
	)
	logic.Update()
	checkGame(t, logic, "animation over and line removed",
		"01",
		"..",
		"..",
	)
}

func TestLinesAboveBoardCanNotBeFull(t *testing.T) {
	logic := NewLogic(alwaysReturn(block(0, 0, 0, 1)))
	logic.SetBoardSizeForPlayerCount(2, BoardSize{2, 1})
	logic.SetBlockStartPositions(2, []Point{{0, 0}, {1, 0}})
	logic.StartNewGame(2)
	checkGame(t, logic, "blocks are 1 higher than board",
		"01",
	)
	logic.Update(InputEvent{0, DownPressed}, InputEvent{1, DownPressed})
	checkGame(t, logic, "blocks are dropped",
		"01",
	)
}

func TestRightKeyIsRepeatedWhileNotReleased(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{6, 1}, []Point{{0, 0}})
	logic.SetInitialLeftRightKeyDelay(3)
	logic.SetShortLeftRightKeyDelay(2)
	logic.StartNewGame(1)
	checkGame(t, logic, "initial",
		"0.....",
	)
	logic.Update(InputEvent{0, RightPressed})
	checkGame(t, logic, "on right down move right",
		".0....",
	)
	logic.Update()
	logic.Update()
	logic.Update()
	checkGame(t, logic, "not moved right again before initial repeat dealay",
		".0....",
	)
	logic.Update()
	checkGame(t, logic, "moved right after initial repeat dealay",
		"..0...",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "not moved right again before fast repeat delay",
		"..0...",
	)
	logic.Update()
	checkGame(t, logic, "moved right after fast repeat delay",
		"...0..",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "not moved right again yet",
		"...0..",
	)
	logic.Update()
	checkGame(t, logic, "moved right second time after fast repeat delay",
		"....0.",
	)
	logic.Update(InputEvent{0, RightReleased})
	logic.Update()
	logic.Update()
	logic.Update()
	checkGame(t, logic, "not moved after right has been released",
		"....0.",
	)
}

func TestPressingRightAgainWhileAlreadyDownDoesNotMoveAgain(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{6, 1}, []Point{{0, 0}})
	logic.SetInitialLeftRightKeyDelay(3)
	logic.SetShortLeftRightKeyDelay(2)
	logic.StartNewGame(1)
	checkGame(t, logic, "initial",
		"0.....",
	)
	logic.Update(InputEvent{0, RightPressed})
	checkGame(t, logic, "on right down move right",
		".0....",
	)
	logic.Update(InputEvent{0, RightPressed}, InputEvent{0, RightPressed})
	checkGame(t, logic, "not moving right before repeat or key up and down again",
		".0....",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "not moved before initial key repeat delay",
		".0....",
	)
	logic.Update()
	checkGame(t, logic, "moved right after initial key repeat delay",
		"..0...",
	)
}

func TestAllPlayersAreRightRepeated(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{6, 1}, []Point{{0, 0}, {1, 0}})
	logic.SetInitialLeftRightKeyDelay(0)
	logic.SetShortLeftRightKeyDelay(0)
	logic.StartNewGame(2)
	checkGame(t, logic, "initial",
		"01....",
	)
	logic.Update(InputEvent{1, RightPressed})
	checkGame(t, logic, "right after initial",
		"0.1...",
	)
	logic.Update()
	checkGame(t, logic, "right repeated",
		"0..1..",
	)
	logic.Update()
	checkGame(t, logic, "right repeated twice",
		"0...1.",
	)
	logic.Update()
}

func TestLeftKeyIsRepeated(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{6, 1}, []Point{{5, 0}, {4, 0}})
	logic.SetInitialLeftRightKeyDelay(3)
	logic.SetShortLeftRightKeyDelay(2)
	logic.StartNewGame(2)
	checkGame(t, logic, "initial",
		"....10",
	)
	logic.Update(InputEvent{1, LeftPressed})
	checkGame(t, logic, "moved on left pressed",
		"...1.0",
	)
	logic.Update()
	logic.Update()
	logic.Update()
	checkGame(t, logic, "initial delay not over",
		"...1.0",
	)
	logic.Update()
	checkGame(t, logic, "moved after initial delay",
		"..1..0",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "fast delay not over",
		"..1..0",
	)
	logic.Update()
	checkGame(t, logic, "moved after fast delay",
		".1...0",
	)
}

func TestDownKeyIsRepeated(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{1, 5}, []Point{{0, 4}, {0, 3}})
	logic.SetInitialDownKeyDelay(3)
	logic.SetShortDownKeyDelay(2)
	logic.StartNewGame(2)
	checkGame(t, logic, "initial",
		"0",
		"1",
		".",
		".",
		".",
	)
	logic.Update(InputEvent{1, DownPressed})
	checkGame(t, logic, "moved on down pressed",
		"0",
		".",
		"1",
		".",
		".",
	)
	logic.Update()
	logic.Update()
	logic.Update()
	checkGame(t, logic, "moved on down pressed",
		"0",
		".",
		"1",
		".",
		".",
	)
	logic.Update()
	checkGame(t, logic, "moved on down pressed",
		"0",
		".",
		".",
		"1",
		".",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "moved on down pressed",
		"0",
		".",
		".",
		"1",
		".",
	)
	logic.Update()
	checkGame(t, logic, "moved on down pressed",
		"0",
		".",
		".",
		".",
		"1",
	)
}

func TestKeyReleaseEventsAreHandledEvenWhenAnimationIsRunning(t *testing.T) {
	logic := createSingleBlockGame(2, BoardSize{8, 3}, []Point{{0, 2}, {7, 0}})
	logic.SetInitialLeftRightKeyDelay(0)
	logic.SetInitialDownKeyDelay(0)
	logic.SetShortLeftRightKeyDelay(0)
	logic.SetShortDownKeyDelay(0)
	animation := &spyLineAnimation{running: false}
	logic.SetLineAnimation(animation)
	logic.StartNewGame(2)
	logic.Update(
		InputEvent{0, DownPressed},
		InputEvent{0, RightPressed},
		InputEvent{1, LeftPressed})
	checkGame(t, logic, "initial",
		"........",
		".0......",
		"......1.",
	)
	animation.running = true
	logic.Update(
		InputEvent{0, DownReleased},
		InputEvent{0, RightReleased},
		InputEvent{1, LeftReleased},
		InputEvent{123, LeftReleased})
	animation.running = false
	logic.Update()
	checkGame(t, logic, "nothing changed",
		"........",
		".0......",
		"......1.",
	)
}

func TestRemovedLinesAreGivenToScorer(t *testing.T) {
	b := block(0, 0, 0, 1)
	logic := NewLogic(func() Block { return b.Copy() })
	logic.SetBoardSizeForPlayerCount(2, BoardSize{2, 3})
	logic.SetBlockStartPositions(2, []Point{{0, 1}, {1, 1}})
	spy := &spyScorer{}
	logic.SetScorer(spy)
	logic.StartNewGame(2)
	logic.Board().SetAt(1, 0, 0)
	logic.Update(
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{0, DownPressed}, InputEvent{0, DownReleased},
		InputEvent{1, DownPressed}, InputEvent{1, DownReleased},
	)
	checkGame(t, logic, "dropped block",
		".1",
		"01",
		"00",
	)
	checkIntsEqual(t, spy.lines.LinesForPlayer(0), []int{},
		"no lines removed yet")
	checkIntsEqual(t, spy.lines.LinesForPlayer(1), []int{},
		"no lines removed yet")

	logic.Update()
	checkIntsEqual(t, spy.lines.LinesForPlayer(0), []int{0, 1},
		"lines 0 and 1 expected full for player 0")
	checkIntsEqual(t, spy.lines.LinesForPlayer(1), []int{1},
		"line 1 expected full for player 1")

	logic.Update()
	checkIntsEqual(t, spy.lines.LinesForPlayer(0), []int{},
		"this frame no lines are full for player 0")
	checkIntsEqual(t, spy.lines.LinesForPlayer(1), []int{},
		"this frame no lines are full for player 1")
}

func TestDownKeyRepeatRateMayDifferFromHorizontalOne(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{5, 4}, []Point{{0, 3}})
	logic.SetInitialDownKeyDelay(3)
	logic.SetInitialLeftRightKeyDelay(2)
	logic.SetShortDownKeyDelay(1)
	logic.SetShortLeftRightKeyDelay(0)
	logic.StartNewGame(1)
	logic.Update(InputEvent{0, DownPressed}, InputEvent{0, RightPressed})
	checkGame(t, logic, "initial",
		".....",
		".0...",
		".....",
		".....",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "no initial delay over",
		".....",
		".0...",
		".....",
		".....",
	)
	logic.Update()
	checkGame(t, logic, "initial left right delay over",
		".....",
		"..0..",
		".....",
		".....",
	)
	logic.Update()
	checkGame(t, logic, "initial down, left and right delay over",
		".....",
		".....",
		"...0.",
		".....",
	)
	logic.Update()
	checkGame(t, logic, "fast left right delay not over",
		".....",
		".....",
		"....0",
		".....",
	)
	logic.Update()
	checkGame(t, logic, "fast left right delay over",
		".....",
		".....",
		".....",
		"....0",
	)
}

func TestWhileDownKeyPressed_BlocksAreNotDroppedOnTimer(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{1, 3}, []Point{{0, 2}})
	logic.SetInitialDownKeyDelay(10)
	logic.SetDropTimer(&spyDropTimer{isTimeForDrop: true})
	logic.StartNewGame(1)
	checkGame(t, logic, "initial",
		"0",
		".",
		".",
	)
	logic.Update(InputEvent{0, DownPressed})
	checkGame(t, logic, "dropped once because of down key press",
		".",
		"0",
		".",
	)
	logic.Update()
	checkGame(t, logic, "not dropped while down key pressed",
		".",
		"0",
		".",
	)
}

func TestAfterBlockResetTheDownKeyIsReleasedAutomatically(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{2, 3}, []Point{{0, 2}})
	logic.SetInitialDownKeyDelay(0)
	logic.SetShortDownKeyDelay(0)
	logic.StartNewGame(1)
	checkGame(t, logic, "initial",
		"0.",
		"..",
		"..",
	)
	logic.Update(InputEvent{0, DownPressed})
	logic.Update()
	logic.Update()
	checkGame(t, logic, "dropped to ground",
		"..",
		"..",
		"0.",
	)
	logic.Update()
	logic.Update()
	checkGame(t, logic, "stopped dropping after block reset",
		"0.",
		"..",
		"0.",
	)
}

func TestSoundPlayerIsAddedAsBlockObserver(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{2, 2}, []Point{{0, 1}})
	spy := &spySoundPlayer{}
	logic.SetSoundPlayer(spy)
	logic.StartNewGame(1)

	logic.Update(InputEvent{0, DownPressed}, InputEvent{0, LeftPressed})

	if spy.log != "0 down " {
		t.Error("down move not observed, log was:", spy.log)
	}
	checkIntsEqual(t, spy.horizontalHits, []int{0}, "left collision not observed")
	if spy.played != 1 {
		t.Error("not played the sounds", spy.played)
	}
}

func TestHoldingRightMovesAsSoonAsNotBlockedAnymore(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{3, 2}, []Point{{0, 1}})
	logic.SetInitialLeftRightKeyDelay(100)
	logic.SetShortLeftRightKeyDelay(0)
	logic.StartNewGame(1)
	logic.Board().SetAt(1, 1, 0)
	logic.Board().SetAt(2, 1, 0)
	logic.Update(InputEvent{0, RightPressed})
	checkGame(t, logic, "blocked right",
		"000",
		"...",
	)
	logic.Update(InputEvent{0, DownPressed}, InputEvent{0, DownReleased})
	logic.Update()
	checkGame(t, logic, "moved right after blocked",
		".00",
		".0.",
	)
	logic.Update()
	checkGame(t, logic, "continue with fast delay",
		".00",
		"..0",
	)
}

func TestHoldingLeftMovesAsSoonAsNotBlockedAnymore(t *testing.T) {
	logic := createSingleBlockGame(1, BoardSize{3, 2}, []Point{{2, 1}})
	logic.SetInitialLeftRightKeyDelay(100)
	logic.SetShortLeftRightKeyDelay(0)
	logic.StartNewGame(1)
	logic.Board().SetAt(0, 1, 0)
	logic.Board().SetAt(1, 1, 0)
	logic.Update(InputEvent{0, LeftPressed})
	checkGame(t, logic, "blocked left",
		"000",
		"...",
	)
	logic.Update(InputEvent{0, DownPressed}, InputEvent{0, DownReleased})
	logic.Update()
	checkGame(t, logic, "moved left after blocked",
		"00.",
		".0.",
	)
	logic.Update()
	checkGame(t, logic, "continue with fast delay",
		"00.",
		"0..",
	)
}

func TestAdjacentBlocksAreFollowedIfHoldingRightDown(t *testing.T) {
	// TODO find a better name for this test (more descriptive)
	// if a block is moving right repeatedly and another block is in its way
	// than whenever that block moves the one behind it should move as well.
	// This way you stick to the other block as long as it moves "slower" than
	// you or with the same speed (when it also repeatedly goes the direction).
	// Note: this may require a function similar to that for dropping blocks
	// at the same time since e.g. two blocks moving right repeatedly might
	// cause problems depending on which one is moved first, ideally they really
	// move at the same time but this has the same problems as with the
	// dropping.
}

// test helpers start here /////////////////////////////////////////////////////

func createSingleBlockGame(players int, size BoardSize, starts []Point) *Logic {
	logic := NewLogic(alwaysReturn(block(0, 0)))
	logic.SetBoardSizeForPlayerCount(players, size)
	logic.SetBlockStartPositions(players, starts)
	return logic
}

func checkGame(t *testing.T, l *Logic, msg string, boardAndBlocks ...string) {
	actual := gameToString(l)
	expected := "\n" + strings.Join(boardAndBlocks, "\n") + "\n"
	if actual != expected {
		t.Error(msg, expected, "expected but was", actual)
	}
}

func gameToString(l *Logic) string {
	w, h := l.Board().Size()
	fields := make2DByteField(w, h)
	addBoardFields(fields, l.Board(), w, h)
	addBlockFields(fields, l.Blocks(), w, h)
	return concatGame(fields)
}

func make2DByteField(w, h int) [][]byte {
	fields := make([][]byte, h)
	for i := range fields {
		fields[i] = make([]byte, w)
	}
	return fields
}

func addBoardFields(fields [][]byte, board Board, w, h int) {
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if f := board.At(x, y); f != NoPlayer {
				fields[h-1-y][x] = byte(f) + '0'
			} else {
				fields[h-1-y][x] = '.'
			}
		}
	}
}

func addBlockFields(fields [][]byte, blocks []Block, w, h int) {
	for i, b := range blocks {
		for _, p := range b.Points {
			if p.X >= 0 && p.Y >= 0 && p.X < w && p.Y < h {
				fields[h-1-p.Y][p.X] = byte(i) + '0'
			}
		}
	}
}

func concatGame(fields [][]byte) string {
	game := "\n"
	for _, row := range fields {
		game += string(row) + "\n"
	}
	return game
}

func alwaysReturn(b Block) BlockFactory {
	return func() Block {
		return b.Copy()
	}
}

func block(xAndYPairs ...int) Block {
	points := make([]Point, len(xAndYPairs)/2)
	for i := range points {
		points[i] = Point{xAndYPairs[i*2], xAndYPairs[i*2+1]}
	}
	return Block{Points: points}
}

func checkBlocksEqual(t *testing.T, actual []Block, expected ...Block) {
	if len(expected) != len(actual) {
		t.Error(len(expected), "blocks expected but there are", len(actual))
	}
	for i := range expected {
		act := fmt.Sprint(actual[i].Points)
		exp := fmt.Sprint(expected[i].Points)
		if act != exp {
			t.Error("expected", exp, "but block", i, "was", act)
		}
	}
}

func checkIsEmpty(t *testing.T, board Board) {
	w, h := board.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if board.At(x, y) != NoPlayer {
				t.Fatal("board not empty but was", board)
			}
		}
	}
}

func newSpyDropTimer() *spyDropTimer { return &spyDropTimer{} }

type spyDropTimer struct {
	reset         int
	updated       int
	isTimeForDrop bool
}

func (s *spyDropTimer) Reset()             { s.reset++ }
func (s *spyDropTimer) Update()            { s.updated++ }
func (s *spyDropTimer) IsTimeToDrop() bool { return s.isTimeForDrop }

func checkInt(t *testing.T, actual, expected int, msg string) {
	if actual != expected {
		t.Error(msg, ": expected:", expected, "actual:", actual)
	}
}

type spyLineAnimation struct {
	lines   [][]int
	updated int
	running bool
}

func (s *spyLineAnimation) Start(lines []int) {
	s.lines = append(s.lines, lines)
}
func (s *spyLineAnimation) Update()         { s.updated++ }
func (s *spyLineAnimation) IsRunning() bool { return s.running }

type spyScorer struct {
	lines PlayersToLines
}

func (s *spyScorer) LinesRemoved(lines PlayersToLines) { s.lines = lines }

type spySoundPlayer struct {
	spyCollisionObserver
	spyBlockMoveObserver
	played int
}

func (spy *spySoundPlayer) PlaySounds() {
	spy.played++
}
