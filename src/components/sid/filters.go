package mos6581

type FilterType int

const (
	FilterNone = FilterType(iota)
	FilterLp
	FilterBp
	FilterLpBp
	FilterHp
	FilterNotch
	FilterHpBp
	FilterAll
)

func calcResonanceLp(X float64) float64 {
	v := 227.755 - 1.7635*X - 0.0176385*X*X + 0.00333484*X*X*X - 9.05683e-6*X*X*X*X
	return v
}

func calcResonanceHp(X float64) float64 {
	v := 366.374 - 14.0052*X + 0.603212*X*X - 0.000880196*X*X*X
	return v
}
