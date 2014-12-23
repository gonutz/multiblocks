package game

type blockFactory struct{}

func NewBlockFactory() blockFactory { return blockFactory{} }

func (blockFactory) CreateO() Block {
	return Block{Points: []Point{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}}
}

func (blockFactory) CreateI() Block {
	return Block{Points: []Point{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}, RotationDeltas: [][]Point{
		[]Point{
			{1, -1},
			{0, 0},
			{-1, 1},
			{-2, 2},
		}, []Point{
			{-1, 1},
			{0, 0},
			{1, -1},
			{2, -2},
		},
	}}
}

func (blockFactory) CreateL() Block {
	return Block{Points: []Point{
		{2, 1},
		{1, 1},
		{0, 1},
		{0, 0},
	}, RotationDeltas: [][]Point{
		[]Point{
			{-1, -1},
			{0, 0},
			{1, 1},
			{0, 2},
		}, []Point{
			{-1, 1},
			{0, 0},
			{1, -1},
			{2, 0},
		}, []Point{
			{1, 1},
			{0, 0},
			{-1, -1},
			{0, -2},
		}, []Point{
			{1, -1},
			{0, 0},
			{-1, 1},
			{-2, 0},
		},
	}}
}

func (blockFactory) CreateJ() Block {
	return Block{Points: []Point{
		{0, 1},
		{1, 1},
		{2, 1},
		{2, 0},
	}, RotationDeltas: [][]Point{
		[]Point{
			{1, 1},
			{0, 0},
			{-1, -1},
			{-2, 0},
		}, []Point{
			{1, -1},
			{0, 0},
			{-1, 1},
			{0, 2},
		}, []Point{
			{-1, -1},
			{0, 0},
			{1, 1},
			{2, 0},
		}, []Point{
			{-1, 1},
			{0, 0},
			{1, -1},
			{0, -2},
		},
	}}
}

func (blockFactory) CreateT() Block {
	return Block{Points: []Point{
		{1, 1},
		{0, 1},
		{1, 0},
		{2, 1},
	}, RotationDeltas: [][]Point{
		[]Point{
			{0, 0},
			{1, 1},
			{-1, 1},
			{-1, -1},
		}, []Point{
			{0, 0},
			{1, -1},
			{1, 1},
			{-1, 1},
		}, []Point{
			{0, 0},
			{-1, -1},
			{1, -1},
			{1, 1},
		}, []Point{
			{0, 0},
			{-1, 1},
			{-1, -1},
			{1, -1},
		},
	}}
}

func (blockFactory) CreateS() Block {
	return Block{Points: []Point{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	}, RotationDeltas: [][]Point{
		[]Point{
			{1, 0},
			{0, 1},
			{-1, 0},
			{-2, 1},
		}, []Point{
			{-1, 0},
			{0, -1},
			{1, 0},
			{2, -1},
		},
	}}
}

func (blockFactory) CreateZ() Block {
	return Block{Points: []Point{
		{0, 1},
		{1, 1},
		{1, 0},
		{2, 0},
	}, RotationDeltas: [][]Point{
		[]Point{
			{1, 1},
			{0, 0},
			{-1, 1},
			{-2, 0},
		}, []Point{
			{-1, -1},
			{0, 0},
			{1, -1},
			{2, 0},
		},
	}}
}
