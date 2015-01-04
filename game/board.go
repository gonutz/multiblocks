package game

// board encodes the game field, a rectangular area where solid block pieces are
// stored.
type board [][]int

func newBoard(w, h int) board {
	b := make([][]int, h)
	for y := range b {
		b[y] = make([]int, w)
		for x := range b[y] {
			b[y][x] = NoPlayer
		}
	}
	return b
}

func (b board) isBlocked(x, y int) bool {
	return y < len(b) && b[y][x] != NoPlayer
}

func (b board) Size() (w, h int) {
	if len(b) == 0 {
		return 0, 0
	}
	return len(b[0]), len(b)
}

func (b board) At(x, y int) int {
	return b[y][x]
}

func (b board) SetAt(x, y, setTo int) {
	b[y][x] = setTo
}

// Copy creates a new board copying the original arrays so that changing the
// copy does not change the orignial.
func (b board) Copy() Board {
	c := make([][]int, len(b))
	for i := range c {
		c[i] = make([]int, len(b[i]))
		copy(c[i], b[i])
	}
	return board(c)
}
