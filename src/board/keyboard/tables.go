package keyboard

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
	KEY_F9        = 256
	KEY_F10       = 257
	KEY_F11       = 258
	KEY_F12       = 259
	KEY_FIRE      = 260
	KEY_JUP       = 261
	KEY_JDN       = 262
	KEY_JLF       = 263
	KEY_JRT       = 264
	KEY_JUPLF     = 265
	KEY_JUPRT     = 266
	KEY_JDNLF     = 267
	KEY_JDNRT     = 268
	KEY_CENTER    = 269
	KEY_NUMLOCK   = 270
	KEY_KPPLUS    = 271
	KEY_KPMINUS   = 272
	KEY_KPMULT    = 273
	KEY_KPDIV     = 274
	KEY_KPENTER   = 275
	KEY_KPPERIOD  = 276
	KEY_PAUSE     = 277
	KEY_ALTENTER  = 278
	KEY_CTRLENTER = 279
)

const (
	JOYSTICK_SENSITIVITY = 40     // % of live range
	JOYSTICK_MIN         = 0x0000 // min value of range
	JOYSTICK_MAX         = 0xffff // max value of range
	JOYSTICK_RANGE       = JOYSTICK_MAX - JOYSTICK_MIN
)

const (
	MAX_STORAGE_SIZE = 1024
	CMD_COUNT        = 0
)

const (
	VK_return   = -1
	VK_back     = -2
	VK_space    = -3
	VK_escape   = -4
	VK_tab      = -5
	VK_delete   = -6
	VK_shift    = -7
	VK_control  = -8
	VK_menu     = -9
	VK_insert   = -10
	VK_numlock  = -11
	VK_capital  = -12
	VK_home     = -13
	VK_end      = -14
	VK_prior    = -15
	VK_next     = -16
	VK_up       = -17
	VK_down     = -18
	VK_left     = -19
	VK_right    = -20
	VK_numpad0  = -21
	VK_numpad1  = -22
	VK_numpad2  = -23
	VK_numpad3  = -24
	VK_numpad4  = -25
	VK_numpad5  = -26
	VK_numpad6  = -27
	VK_numpad7  = -28
	VK_numpad8  = -29
	VK_numpad9  = -30
	VK_f1       = -31
	VK_f2       = -32
	VK_f3       = -33
	VK_f4       = -34
	VK_f5       = -35
	VK_f6       = -36
	VK_f7       = -37
	VK_f8       = -38
	VK_f9       = -39
	VK_f10      = -40
	VK_f11      = -41
	VK_f12      = -42
	VK_clear    = -43
	VK_pause    = -44
	VK_multiply = -45
	VK_divide   = -46
	VK_subtract = -47
	VK_add      = -48
	VK_decimal  = -49

	VK_semicolon   = 0xba
	VK_equal       = 0xbb
	VK_comma       = 0xbc
	VK_period      = 0xbe
	VK_minus       = 0xbd
	VK_slash       = 0xbf
	VK_bracketleft = 0xdb
	VK_backslash   = 0xdc
	VK_quote       = 0xde
	VK_asterisk    = 0xdd
	VK_grave       = 0xc0
	VK_plus        = 0xc1
)
