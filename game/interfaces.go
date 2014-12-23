package game

type Board interface {
	Size() (w, h int)
	At(x, y int) int // the origin (0,0) is the bottom-left field
	Copy() Board
	SetAt(x, y, setTo int)
}

const NoPlayer = -1

type BlockFactory func() Block

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

type PlayersToLines interface {
	LinesForPlayer(p int) []int
}

type Scorer interface {
	LinesRemoved(lines PlayersToLines)
}
