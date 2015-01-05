package game

type repeatableKey struct {
	timer        int
	initialDelay int
	fastDelay    int
	down         bool
}

func newRepeatableKey(initial, fast int) *repeatableKey {
	return &repeatableKey{initialDelay: initial, fastDelay: fast}
}

func (k *repeatableKey) Press() (triggering bool) {
	if k.down {
		return false
	}
	k.down = true
	k.timer = k.initialDelay
	return true
}
func (k *repeatableKey) Release() {
	k.down = false
}

func (k *repeatableKey) Update() (triggering bool) {
	if k.down {
		k.timer--
		if k.timer < 0 {
			k.timer = k.fastDelay
			return true
		}
	}
	return false
}

func (k *repeatableKey) IsDown() bool {
	return k.down
}

func (k *repeatableKey) Blocked() {
	k.timer = 0
}
