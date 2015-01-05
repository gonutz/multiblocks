package game

type Logic struct {
	blockFactory          BlockFactory
	physics               *physics
	previewBlocks         []Block
	dropTimer             DropTimer
	sizes                 [5]BoardSize
	startPositions        [5][]Point
	playerCount           int
	hasDroppedThisFrame   []bool
	lineAnimation         LineAnimation
	fullLines             []int
	leftKeys              []*repeatableKey
	rightKeys             []*repeatableKey
	downKeys              []*repeatableKey
	initialLeftRightDelay int
	initialDownDelay      int
	shortLeftRightDelay   int
	shortDownDelay        int
	scorer                Scorer
	soundPlayer           GameSoundPlayer
}

func NewLogic(f BlockFactory) *Logic {
	return &Logic{blockFactory: f}
}

func (l *Logic) SetDropTimer(timer DropTimer) {
	l.dropTimer = timer
}

func (l *Logic) SetLineAnimation(a LineAnimation) {
	l.lineAnimation = a
}

func (l *Logic) SetScorer(s Scorer) {
	l.scorer = s
}

func (l *Logic) SetSoundPlayer(s GameSoundPlayer) {
	l.soundPlayer = s
}

func (l *Logic) SetBoardSizeForPlayerCount(players int, size BoardSize) {
	l.sizes[players] = size
}

func (l *Logic) SetBlockStartPositions(players int, start []Point) {
	l.startPositions[players] = start
}

// Setting key delays to 0 means that the key is repeated on every update.
// Setting it to 1 means 1 update between repeats, and so on.
func (l *Logic) SetInitialLeftRightKeyDelay(delay int) {
	l.initialLeftRightDelay = delay
}

func (l *Logic) SetShortLeftRightKeyDelay(delay int) {
	l.shortLeftRightDelay = delay
}

func (l *Logic) SetInitialDownKeyDelay(delay int) {
	l.initialDownDelay = delay
}

func (l *Logic) SetShortDownKeyDelay(delay int) {
	l.shortDownDelay = delay
}

func (l *Logic) StartNewGame(players int) {
	l.playerCount = players
	l.hasDroppedThisFrame = make([]bool, players)
	l.physics = newPhysics(l.sizes[players], BlockCount(players))
	l.physics.AddCollisionObserver(l)
	if l.soundPlayer != nil {
		l.physics.AddCollisionObserver(l.soundPlayer)
		l.physics.AddBlockMoveObserver(l.soundPlayer)
	}
	l.createBlocks()
	l.createRepeatableKeys()
}

func (l *Logic) createBlocks() {
	l.previewBlocks = make([]Block, l.playerCount)
	for i := range l.previewBlocks {
		l.previewBlocks[i] = l.blockFactory()
	}
	for i := 0; i < l.playerCount; i++ {
		l.resetBlockToPreview(i)
	}
	if l.dropTimer != nil {
		l.dropTimer.Reset()
	}
}

func (l *Logic) createRepeatableKeys() {
	l.leftKeys = l.makeKeys(l.initialLeftRightDelay, l.shortLeftRightDelay)
	l.rightKeys = l.makeKeys(l.initialLeftRightDelay, l.shortLeftRightDelay)
	l.downKeys = l.makeKeys(l.initialDownDelay, l.shortDownDelay)
}

func (l *Logic) makeKeys(initialDelay, fastDelay int) []*repeatableKey {
	keys := make([]*repeatableKey, l.playerCount)
	for i := 0; i < l.playerCount; i++ {
		keys[i] = newRepeatableKey(initialDelay, fastDelay)
	}
	return keys
}

func (l *Logic) startPositionFor(index int) Point {
	return Point{0, 0}
}

func (l *Logic) Blocks() []Block {
	return l.physics.Blocks()
}

func (l *Logic) PreviewBlocks() []Block {
	return l.previewBlocks
}

func (l *Logic) Board() Board {
	return l.physics.Board()
}

func (l *Logic) Update(events ...InputEvent) {
	if l.lineAnimation != nil && l.lineAnimation.IsRunning() {
		l.handleReleaseEvents(events...)
		l.lineAnimation.Update()
		return
	}
	l.giveScoresForFullLines()
	l.resetPreviouslyDroppedBlocks()
	l.removeFullLines()
	l.handleInputEvents(events...)
	l.dropBlocksIfTimeForIt()
	l.checkCompleteLines()
	if l.soundPlayer != nil {
		l.soundPlayer.PlaySounds()
	}
}

func (l *Logic) giveScoresForFullLines() {
	if l.scorer != nil {
		lines := make([][]int, l.playerCount)
		l.fillWithPlayerToLineInfo(lines)
		l.scorer.LinesRemoved(&playerToLines{lines})
	}
}

func (l *Logic) fillWithPlayerToLineInfo(lines [][]int) {
	for _, line := range l.fullLines {
		for player := 0; player < l.playerCount; player++ {
			if l.playerIsDroppedInLine(player, line) {
				lines[player] = append(lines[player], line)
			}
		}
	}
}

func (l *Logic) playerIsDroppedInLine(player, line int) bool {
	return l.hasDroppedThisFrame[player] && l.blockIsInLine(player, line)
}

func (l *Logic) blockIsInLine(block, line int) bool {
	for _, p := range l.Blocks()[block].Points {
		if p.Y == line {
			return true
		}
	}
	return false
}

type playerToLines struct {
	lines [][]int
}

func (l *playerToLines) LinesForPlayer(p int) []int {
	return l.lines[p]
}

func (l *Logic) resetPreviouslyDroppedBlocks() {
	for b := 0; b < l.playerCount; b++ {
		if l.hasDroppedThisFrame[b] {
			l.physics.CopyBlockToBoard(b)
			l.resetBlockToPreview(b)
			for l.physics.isInOtherBlock(b) {
				l.physics.Blocks()[b].MoveBy(0, 1)
			}
			l.downKeys[b].Release()
			l.hasDroppedThisFrame[b] = false
		}
	}
}

func (l *Logic) resetBlockToPreview(block int) {
	b := l.previewBlocks[block]
	start := l.startPositions[l.playerCount][block]
	w, _ := b.Size()
	b.MoveBy(start.X-w/2, start.Y)
	l.physics.SetBlock(block, b)
	l.previewBlocks[block] = l.blockFactory()
}

func (l *Logic) removeFullLines() {
	l.physics.RemoveLines(l.fullLines...)
}

func (l *Logic) handleReleaseEvents(events ...InputEvent) {
	for _, e := range events {
		if e.Player < l.playerCount {
			switch e.Command {
			case DownReleased:
				l.downKeys[e.Player].Release()
			case LeftReleased:
				l.leftKeys[e.Player].Release()
			case RightReleased:
				l.rightKeys[e.Player].Release()
			}
		}
	}
}

func (l *Logic) handleInputEvents(events ...InputEvent) {
	l.handleKeyRepeatEvents()

	for _, e := range events {
		if e.Player < l.playerCount {
			switch e.Command {

			case DownPressed:
				if !l.hasDroppedThisFrame[e.Player] && l.downKeys[e.Player].Press() {
					l.physics.MoveDown(e.Player)
				}
			case DownReleased:
				l.downKeys[e.Player].Release()

			case LeftPressed:
				if !l.hasDroppedThisFrame[e.Player] && l.leftKeys[e.Player].Press() {
					if !l.physics.MoveLeft(e.Player) {
						l.leftKeys[e.Player].Blocked()
					}
				}
			case LeftReleased:
				l.leftKeys[e.Player].Release()

			case RightPressed:
				if !l.hasDroppedThisFrame[e.Player] && l.rightKeys[e.Player].Press() {
					if !l.physics.MoveRight(e.Player) {
						l.rightKeys[e.Player].Blocked()
					}
				}
			case RightReleased:
				l.rightKeys[e.Player].Release()

			case RotateRight:
				l.physics.RotateRight(e.Player)
			case RotateLeft:
				l.physics.RotateLeft(e.Player)
			}
		}
	}
}

func (l *Logic) handleKeyRepeatEvents() {
	for i := 0; i < l.playerCount; i++ {
		if l.rightKeys[i].Update() {
			l.physics.MoveRight(i)
		}
		if l.leftKeys[i].Update() {
			l.physics.MoveLeft(i)
		}
		if l.downKeys[i].Update() {
			l.physics.MoveDown(i)
		}
	}
}

func (l *Logic) dropBlocksIfTimeForIt() {
	if l.dropTimer != nil {
		l.dropTimer.Update()
		if l.dropTimer.IsTimeToDrop() {
			l.dropBlocksThatDoNotMoveDown()
		}
	}
}

func (l *Logic) dropBlocksThatDoNotMoveDown() {
	l.physics.DropBlocks(l.nonDroppingBlocks())
}

func (l *Logic) nonDroppingBlocks() []int {
	all := make([]int, 0, l.playerCount)
	for i := 0; i < l.playerCount; i++ {
		if !l.downKeys[i].IsDown() {
			all = append(all, i)
		}
	}
	return all
}

func (l *Logic) checkCompleteLines() {
	b := l.Board().Copy()
	l.copyDroppedBlocksToBoard(b)
	w, h := b.Size()
	l.fullLines = make([]int, 0, h)
	for y := 0; y < h; y++ {
		if lineFull(b, y, w) {
			l.fullLines = append(l.fullLines, y)
		}
	}

	if len(l.fullLines) > 0 && l.lineAnimation != nil {
		l.lineAnimation.Start(l.fullLines)
	}
}

func (l *Logic) copyDroppedBlocksToBoard(board Board) {
	for player, b := range l.Blocks() {
		if l.hasDroppedThisFrame[player] {
			copyBlockToBoard(b, board, player)
		}
	}
}

func lineFull(b Board, y, w int) bool {
	for x := 0; x < w; x++ {
		if b.At(x, y) == NoPlayer {
			return false
		}
	}
	return true
}

func (l *Logic) BlockHitGround(block int) {
	l.hasDroppedThisFrame[block] = true
}

func (l *Logic) BlockCouldNotRotate(block int)           {}
func (l *Logic) BlockDraggedDownByLineRemoval(block int) {}
func (l *Logic) BlockHitLeftOrRight(block int)           {}
func (l *Logic) BlockHitOtherBlock(block int)            {}
