package objects

type Accumulator struct {
	Value int
}

func (a *Accumulator) Add(val int) {
	a.Value += val
}

func (a *Accumulator) CanSubtract(val int) (cansub bool) {
	cansub = false
	if a.Value-val >= 0 {
		cansub = true
	}
	return
}

func (a *Accumulator) Subtract(val int) {
	a.Value -= val
}

func (a *Accumulator) SubIfCan(val int) (ok bool) {
	ok = false
	if a.CanSubtract(val) {
		a.Subtract(val)
		ok = true
	}
	return
}
