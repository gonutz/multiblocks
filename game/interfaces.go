package game

type Board interface {
	Size() (w, h int)
	At(x, y int) int // the origin (0,0) is the bottom-left field
	Copy() Board
	SetAt(x, y, setTo int)
}

// NoPlayer is used in a Board to signal that a spot is empty.
const NoPlayer = -1

type BlockFactory func() Block

// DropTimer tells the game logic when it is time to drop all blocks at the same
// time. This usually happens regularly and with higher frequency the higher the
// difficulty.
type DropTimer interface {
	Reset()
	Update()
	IsTimeToDrop() bool
}

type BlockCollisionObserver interface {
	BlockHitLeftOrRight(block int)
	BlockHitOtherBlock(block int)
	BlockHitGround(block int)
	BlockDraggedDownByLineRemoval(block int)
	BlockCouldNotRotate(block int)
}

type BlockMoveObserver interface {
	BlockMovedHorizontally(block int)
	BlockMovedDown(block int)
	BlockRotated(block int)
}

type GameSoundPlayer interface {
	BlockCollisionObserver
	BlockMoveObserver
	PlaySounds()
}

type LineAnimation interface {
	Start(lines []int)
	Update()
	IsRunning() bool
}

type Scorer interface {
	LinesRemoved(linesForPlayer [][]int)
}
