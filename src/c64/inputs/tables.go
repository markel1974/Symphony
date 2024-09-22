package inputs

/*
C64 keyboard matrix:

Bit   7   6   5   4   3   2   1   0
0    CUD  F5  F3  F1  F7 CLR RET DEL
1    SHL  E   S   Z   4   A   W   3
2     X   T   F   C   6   D   R   5
3     V   U   H   B   8   G   Y   7
4     N   O   K   M   0   J   I   9
5     ,   @   :   .   -   L   P   +
6     /   ^   =  SHR HOM  ;   *   ?
7    R/S  Q   C= SPC  2  CTL  <-  1
*/

const (
	KeyJFire = iota
	KeyJUp
	KeyJDown
	KeyJLeft
	KeyJRight
	KeyJUpLeft
	KeyJUpRight
	KeyJDownLeft
	KeyJDownRight
	KeyJCenter
)

const (
	VKReturn = -1024 + iota
	VKBack
	VKSpace
	VKEscape
	VKTab
	VKDelete
	VKShift
	VKControl
	VKInsert
	VKHome
	VKEnd
	VKPrior
	VKNext
	VKUp
	VKDown
	VKLeft
	VKRight
	VKF1
	VKF2
	VKF3
	VKF4
	VKF5
	VKF6
	VKF7
	VKF8
	VKSemicolon
	VKEqual
	VKComma
	VKPeriod
	VKMinus
	VKSlash
	VKBracketLeft
	VKBackslash
	VKQuote
	VKAsterisk
	VKGrave
	VKPlus
)
