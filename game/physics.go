package game

import "sort"

const Empty = -1

type physics struct {
	boardWidth, boardHeight int
	blocks                  []Block
	collisionObservers      []BlockCollisionObserver
	moveObservers           []BlockMoveObserver
	board                   board
}

type BoardSize struct {
	Width, Height int
}

type BlockCount int

func newPhysics(size BoardSize, blocks BlockCount) *physics {
	return &physics{
		boardWidth:  size.Width,
		boardHeight: size.Height,
		blocks:      make([]Block, blocks),
		board:       newBoard(size.Width, size.Height),
	}
}

func (p *physics) AddCollisionObserver(o BlockCollisionObserver) {
	p.collisionObservers = append(p.collisionObservers, o)
}

func (p *physics) AddBlockMoveObserver(o BlockMoveObserver) {
	p.moveObservers = append(p.moveObservers, o)
}

func (p *physics) Board() Board {
	return p.board
}

func (p *physics) Blocks() []Block {
	return p.blocks
}

func (p *physics) SetBlock(index int, block Block) {
	p.blocks[index] = block
}

func (p *physics) MoveLeft(block int) bool {
	return p.moveBlockX(block, -1)
}

func (p *physics) MoveRight(block int) bool {
	return p.moveBlockX(block, 1)
}

func (p *physics) moveBlockX(block, dx int) bool {
	p.blocks[block].MoveBy(dx, 0)
	if p.isInWall(block) || p.isInSolidPartOfBoard(block) {
		p.blocks[block].MoveBy(-dx, 0)
		p.notifyOfLeftRightHit(block)
		return false
	} else if p.isInOtherBlock(block) {
		p.blocks[block].MoveBy(-dx, 0)
		p.notifyOfBlockHit(block)
		return false
	} else {
		p.notifyOfHorizontalMove(block)
		return true
	}
}

func (p *physics) isInWall(block int) bool {
	for _, point := range p.blocks[block].Points {
		if point.X < 0 || point.X >= p.boardWidth {
			return true
		}
	}
	return false
}

func (p *physics) isInSolidPartOfBoard(block int) bool {
	for _, point := range p.blocks[block].Points {
		if p.board.isBlocked(point.X, point.Y) {
			return true
		}
	}
	return false
}

func (p *physics) isInOtherBlock(block int) bool {
	for other := range p.blocks {
		if other != block && p.blocksCollide(block, other) {
			return true
		}
	}
	return false
}

func (p *physics) blocksCollide(a, b int) bool {
	for _, p1 := range p.blocks[a].Points {
		for _, p2 := range p.blocks[b].Points {
			if p1.X == p2.X && p1.Y == p2.Y {
				return true
			}
		}
	}
	return false
}

func (p *physics) notifyOfLeftRightHit(block int) {
	for _, o := range p.collisionObservers {
		o.BlockHitLeftOrRight(block)
	}
}

func (p *physics) notifyOfBlockHit(block int) {
	for _, o := range p.collisionObservers {
		o.BlockHitOtherBlock(block)
	}
}

func (p *physics) notifyOfHorizontalMove(block int) {
	for _, o := range p.moveObservers {
		o.BlockMovedHorizontally(block)
	}
}

func (p *physics) MoveDown(block int) {
	p.blocks[block].MoveBy(0, -1)
	if p.isInGround(block) || p.isInSolidPartOfBoard(block) {
		p.blocks[block].MoveBy(0, 1)
		p.notifyOfGroundHit(block)
	} else if p.isInOtherBlock(block) {
		p.blocks[block].MoveBy(0, 1)
		p.notifyOfBlockHit(block)
	} else {
		p.notifyOfDownMove(block)
	}
}

func (p *physics) isInGround(block int) bool {
	for _, p := range p.blocks[block].Points {
		if p.Y < 0 {
			return true
		}
	}
	return false
}

func (p *physics) notifyOfGroundHit(block int) {
	for _, o := range p.collisionObservers {
		o.BlockHitGround(block)
	}
}

func (p *physics) notifyOfDownMove(block int) {
	for _, o := range p.moveObservers {
		o.BlockMovedDown(block)
	}
}

func (phy *physics) CopyBlockToBoard(block int) {
	copyBlockToBoard(phy.blocks[block], phy.board, block)
}

func copyBlockToBoard(block Block, board Board, player int) {
	w, h := board.Size()
	for _, p := range block.Points {
		if p.X >= 0 && p.Y >= 0 && p.X < w && p.Y < h {
			board.SetAt(p.X, p.Y, player)
		}
	}
}

func (p *physics) RemoveLines(lines ...int) {
	sortDescending(lines)
	for _, line := range lines {
		p.removeLine(line)
	}
}

func sortDescending(ints []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(ints)))
}

func (p *physics) removeLine(line int) {
	for y := line; y < p.boardHeight-1; y++ {
		copy(p.board[y], p.board[y+1])
	}
	p.clearTopRow()
	p.resolveLineRemovalCollisions()
}

func (p *physics) clearTopRow() {
	top := p.boardHeight - 1
	for x := 0; x < p.boardWidth; x++ {
		p.board[top][x] = NoPlayer
	}
}

func (p *physics) resolveLineRemovalCollisions() {
	collided, ok := p.findGroundAndBoardHits(-1)
	moreCollisions := true
	for moreCollisions {
		moreCollisions, collided, ok = p.findMoreCollisions(collided, ok, -1)
	}
	for _, block := range collided {
		p.notifyOfDragDown(block)
	}
}

func (p *physics) notifyOfDragDown(block int) {
	for _, o := range p.collisionObservers {
		o.BlockDraggedDownByLineRemoval(block)
	}
}

func (p *physics) DropBlocks(blocks []int) {
	p.moveBlocksDown(blocks)
	collided, ok := p.findGroundAndBoardHits(1)
	moreCollisions := true
	for moreCollisions {
		moreCollisions, collided, ok = p.findMoreCollisions(collided, ok, 1)
	}
	p.notifyOfDropCollisions(collided)
	for _, block := range ok {
		p.notifyOfDownMove(block)
	}
}

func (p *physics) moveBlocksDown(blocks []int) {
	for _, block := range blocks {
		p.blocks[block].MoveBy(0, -1)
	}
}

func (p *physics) findGroundAndBoardHits(moveYBackBy int) (collided, ok []int) {
	for block := range p.blocks {
		if p.isInGround(block) || p.isInSolidPartOfBoard(block) {
			collided = append(collided, block)
			p.blocks[block].MoveBy(0, moveYBackBy)
		} else {
			ok = append(ok, block)
		}
	}
	return
}

func (p *physics) findMoreCollisions(collided, ok []int, moveYBackBy int) (
	moreCollisions bool, nowCollided, stillOk []int) {
	nowCollided = collided
	for _, block := range ok {
		if p.collidesWithAnyOf(block, collided) {
			nowCollided = append(nowCollided, block)
			p.blocks[block].MoveBy(0, moveYBackBy)
			moreCollisions = true
		} else {
			stillOk = append(stillOk, block)
		}
	}
	return
}

func (p *physics) collidesWithAnyOf(block int, others []int) bool {
	for _, other := range others {
		if p.blocksCollide(block, other) {
			return true
		}
	}
	return false
}

func (p *physics) notifyOfDropCollisions(blocks []int) {
	for _, block := range blocks {
		p.notifyOfGroundHit(block)
	}
}

func (p *physics) RotateRight(block int) {
	p.blocks[block].RotateRight()
	p.handleRotationCollision(block, p.blocks[block].RotateLeft)
}

func (p *physics) RotateLeft(block int) {
	p.blocks[block].RotateLeft()
	p.handleRotationCollision(block, p.blocks[block].RotateRight)
}

func (p *physics) handleRotationCollision(block int, reset func()) {
	if p.isInWall(block) || p.isInGround(block) ||
		p.isInSolidPartOfBoard(block) || p.isInOtherBlock(block) {
		reset()
		p.notifyOfRotationHit(block)
	} else {
		p.notifyOfRotation(block)
	}
}

func (p *physics) notifyOfRotationHit(block int) {
	for _, o := range p.collisionObservers {
		o.BlockCouldNotRotate(block)
	}
}

func (p *physics) notifyOfRotation(block int) {
	for _, o := range p.moveObservers {
		o.BlockRotated(block)
	}
}
