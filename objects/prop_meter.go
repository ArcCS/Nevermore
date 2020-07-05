package objects

type Meter struct{
	Max int
	Current int
}

func (m *Meter) Add(val int){
	if m.Current+val <= m.Max{
		m.Current += val
	}else{
		m.Current = m.Max
	}
}

func (m *Meter) Subtract(val int){
	 if m.Current-val >= 0 {
	 	m.Current -= val
	 }else{
	 	m.Current = 0
	 }
}