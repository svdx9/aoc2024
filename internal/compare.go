package internal

func Float64AlmostEqual(f1 float64, f2 float64) bool {
	return f2-f1 < 1e-6
}
