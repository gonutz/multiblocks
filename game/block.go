package game

// Block is a game piece consising of several (usually four) pieces. It contains
// the coordinates of its pieces (Points) and all possible rotations, encoded in
// RotationDeltas. These are the deltas that have to be added to the points to
// get the next rotation.
type Block struct {
	Points         []Point
	RotationDeltas [][]Point
	rotation       int
}

type Point struct{ X, Y int }

// Size calculates the current maximum x and y spread. It can change depending
// on the current rotation of the Block.
func (b *Block) Size() (w, h int) {
	if len(b.Points) == 0 {
		return 0, 0
	}
	first := b.Points[0]
	minX, maxX := first.X, first.X
	minY, maxY := first.Y, first.Y
	for _, p := range b.Points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return maxX - minX + 1, maxY - minY + 1
}

// RotateRight applies the next RotationDeltas values to the Block's Points.
func (b *Block) RotateRight() {
	if len(b.RotationDeltas) > 0 {
		for i := range b.Points {
			b.Points[i].X += b.RotationDeltas[b.rotation][i].X
			b.Points[i].Y += b.RotationDeltas[b.rotation][i].Y
		}
		b.increaseRotation()
	}
}

func (b *Block) increaseRotation() {
	b.rotation = (b.rotation + 1) % len(b.RotationDeltas)
}

// RotateRight applies the previous RotationDeltas values to the Block's Points.
func (b *Block) RotateLeft() {
	if len(b.RotationDeltas) > 0 {
		b.decreaseRotation()
		for i := range b.Points {
			b.Points[i].X -= b.RotationDeltas[b.rotation][i].X
			b.Points[i].Y -= b.RotationDeltas[b.rotation][i].Y
		}
	}
}

func (b *Block) decreaseRotation() {
	b.rotation = (b.rotation + len(b.RotationDeltas) - 1) % len(b.RotationDeltas)
}

// MoveBy shifts all Block's Points by the given amounts.
func (b *Block) MoveBy(dx, dy int) {
	for i := range b.Points {
		b.Points[i].X += dx
		b.Points[i].Y += dy
	}
}

// Copy creates an exact copy of the Block but with newly created arrays so that
// changing the copy does not change the original.
func (b *Block) Copy() Block {
	c := Block{}

	c.Points = make([]Point, len(b.Points))
	copy(c.Points, b.Points)

	c.RotationDeltas = make([][]Point, len(b.RotationDeltas))
	for i := range c.RotationDeltas {
		c.RotationDeltas[i] = make([]Point, len(b.RotationDeltas[i]))
		copy(c.RotationDeltas[i], b.RotationDeltas[i])
	}

	c.rotation = b.rotation

	return c
}
