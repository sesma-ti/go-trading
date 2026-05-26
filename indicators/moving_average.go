package indicators

type MovingAverage struct {
	values []float64
	size   int
	sum    float64
	pos    int
	filled bool
}

func NewMovingAverage(size int) *MovingAverage {
	return &MovingAverage{
		values: make([]float64, size),
		size:   size,
	}
}

func (m *MovingAverage) Add(value float64) float64 {
	if m.filled {
		m.sum -= m.values[m.pos]
	}

	m.values[m.pos] = value
	m.sum += value

	m.pos = (m.pos + 1) % m.size

	if m.pos == 0 {
		m.filled = true
	}

	count := m.pos
	if m.filled {
		count = m.size
	}

	return m.sum / float64(count)
}
