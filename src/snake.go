package main

type Snake struct {
	Parts []Point
}

// NewSnake - конструктор змеи
func NewSnake() *Snake {
	snake := &Snake{}

	return snake
}

// Len - для определения длины змеи
func (s *Snake) Len() int {
	return len(s.Parts)
}

// Head - для определения положения головы змеи
func (s *Snake) Head() Point {
	if s.Len() == 0 {
		return Point{X: -1, Y: -1}
	}
	return s.Parts[0]
}

// Tail - для определения положения хвоста
func (s *Snake) Tail() Point {
	if s.Len() == 0 {
		return Point{X: -1, Y: -1}
	}
	return s.Parts[len(s.Parts)-1]
}

// Add - для еды и увеличения змеи
func (s *Snake) Add(p Point) {
	s.Parts = append([]Point{p}, s.Parts...) //змея растёт с головы, что выглядит более реалистично

}

// IsSnake - проверка, наткнулась ли змея на себя
func (s *Snake) IsSnake(p Point) bool {
	for i := range s.Parts {
		if s.Parts[i] == p {
			return true
		}
	}
	return false
}

// CutIfSnake - обрезание змеи если она наткнулась на себя
func (s *Snake) CutIfSnake(p Point) bool {
	i := 0
	for ; i < len(s.Parts); i++ {
		if s.Parts[i] == p {
			break
		}
	}

	s.Parts = s.Parts[0:i]
	return i >= len(s.Parts)
}

// Reset - начальное положение змеи
func (s *Snake) Reset() {
	sx, sy, l := 1, 1, 5
	for i := l - 1; i >= 0; i-- {
		s.Parts = append(s.Parts, Point{X: float64(sx + i), Y: float64(sy)})
	}
}

// Move - передвижение змеи
func (s *Snake) Move(d Dir) {
	lastP := s.Parts[0]
	s.Parts[0] = d.Exec(s.Parts[0])
	for i := range s.Parts[1:] {
		s.Parts[i+1], lastP = lastP, s.Parts[i+1]
	}
}
