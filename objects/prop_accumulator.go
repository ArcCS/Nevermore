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
	if a.Value < 0 {
		a.Value = 0
	}
}

func (a *Accumulator) SubMax(val int, finalVal int) {
	a.Value -= val
	if a.Value < finalVal {
		a.Value = finalVal
	}
}

func (a *Accumulator) SubIfCan(val int) (ok bool) {
	ok = false
	if a.CanSubtract(val) {
		a.Subtract(val)
		ok = true
	}
	return
}
