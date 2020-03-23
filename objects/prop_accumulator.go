package objects

type Accumulator struct{
	Value int64
}

func (a *Accumulator) Add(val int64){
	a.Value += val
}

func (a *Accumulator) CanSubtract(val int64) (cansub bool){
	cansub = false
	if a.Value - val >= 0{
		cansub = true
	}
	return
}

func (a *Accumulator) Subtract(val int64){
	a.Value -= val
}

func (a *Accumulator) SubIfCan(val int64) (ok bool){
	ok = false
	if a.CanSubtract(val){
		a.Subtract(val)
		ok = true
	}
	return
}
