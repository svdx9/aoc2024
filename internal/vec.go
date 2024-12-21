package internal

import "math"

type Vec struct {
	X int
	Y int
}

func (v Vec) AddVec(a Vec) Vec {
	return Vec{
		X: v.X + a.X,
		Y: v.Y + a.Y,
	}
}
func (v Vec) SubVec(a Vec) Vec {
	return Vec{
		X: v.X - a.X,
		Y: v.Y - a.Y,
	}
}

func (v Vec) Reflect() Vec {
	return Vec{
		X: v.X * -1,
		Y: v.Y * -1,
	}
}

func (v Vec) Distance(o Vec) float64 {
	x := o.X - v.X
	y := o.Y - v.Y
	return math.Sqrt(float64((x * x) + (y * y)))
}
