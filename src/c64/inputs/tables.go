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
	VKMenu
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

/*

const (
// KEY_KPPERIOD = 276
// KEY_ALTENTER = 278
)

	KEY_F9        = 256
	KEY_F10       = 257
	KEY_F11       = 258
	KEY_F12       = 259

	KEY_NUMLOCK   = 270
	KEY_KPPLUS    = 271
	KEY_KPMINUS   = 272
	KEY_KPMULT    = 273
	KEY_KPDIV     = 274
	KEY_KPENTER   = 275
	KEY_PAUSE     = 277
	KEY_CTRLENTER = 279

	//VK_numlock = -11
	//VK_capital = -12
	//VK_numpad0  = -21
	//VK_numpad1  = -22
	//VK_numpad2  = -23
	//VK_numpad3  = -24
	//VK_numpad4  = -25
	//VK_numpad5  = -26
	//VK_numpad6  = -27
	//VK_numpad7  = -28
	//VK_numpad8  = -29
	//VK_numpad9  = -30
	//VK_f9       = -39
	//VK_f10      = -40
	//VK_f11      = -41
	//VK_f12      = -42
	//VK_clear    = -43
	//VK_pause    = -44
	//VK_multiply = -45
	//VK_divide   = -46
	//VK_subtract = -47
	//VK_add      = -48
	//VK_decimal  = -49

*/
