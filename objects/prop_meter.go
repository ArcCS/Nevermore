package objects

type Meter struct{
	Max int64
	Current int64
}

func (m *Meter) Add(val int64){
	if m.Current+val < m.Max{
		m.Current += val
	}else{
		m.Current = m.Max
	}
}

func (m *Meter) Subtract(val int64){
	 if m.Current-val > 0 {
	 	m.Current -= val
	 }else{
	 	m.Current = 0
	 }
}