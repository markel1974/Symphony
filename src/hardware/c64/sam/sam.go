package sam

import (
	"bufio"
	"fmt"
	"os"
)

type SAM struct {
	_the_token  Token
	_the_number uint16
	_the_string string
	reader      *bufio.Reader
	line        Line
}

func New() *SAM {

	/*
		C64SAM::C64SAM(C64Model* c64, C1541Model* c1541) :
			TheCPU(c64->_board->GetCPU()),
				TheCPU1541(c1541->_board->GetCPU()),
				TheVIC(c64->_board->GetVIC()),
				TheSID(c64->_board->GetSID()),
				TheCIA1(c64->_board->GetCIA1()),
				TheCIA2(c64->_board->GetCIA2()) {
				#ifdef _WIN32
				if (!__consoleEnabled) {
					RedirectIOToConsole();
					__consoleEnabled = true;
				}
				#endif
			}

	*/

	return &SAM{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (s *SAM) readLine() (string, error) {
	txt, err := s.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return txt, err

	//s._in_ptr = _input

	//char* output = fgets(_in_ptr, INPUT_LENGTH, _fin);
	//if (!output)
	//perror("Error reading stdinput")
}

func (s *SAM) Exec() {

	done := true
	var c rune

	// Get CPU registers and current memory configuration
	//TheCPU->getStateMachine(&R64);
	//TheCPU->_ExtConfig = (~R64.ddr | R64.pr) & 7;
	//TheCPU1541->getStateMachine(&R1541);

	//_fin = stdin;
	//_fout = stdout;
	//_ferr = stdout;
	//_access_1541 = false;
	//_address = R64.pc;

	fmt.Printf("\n *** SAM - Simple Assembler and Monitor ***\n ***         Press 'h' for help         ***\n\n")
	//display_registers();

	for done {
		//if (_access_1541) {
		//fprintf(_ferr, "1541> ");
		//} else {
		fmt.Printf("C64> ")

		//fflush(_ferr);
		l, err := s.readLine()
		if err != nil {
			fmt.Printf("error reading line %s\n", err.Error())
			continue
		}
		s.line.Set(l)

		for {
			if c = s.line.getChar(); c != ' ' {
				break
			}
		}

		switch c {
		case 'a': // Assemble
			s.getToken()
			//s.assemble();
			break

		case 'b': // Binary dump
			s.getToken()
			//binary_dump();
			break

		case 'c': // Compare
			s.getToken()
			//compare();
			break

		case 'd': // Disassemble
			s.getToken()
			//disassemble();
			break

		case 'e': // Interrupt vectors
			//int_vectors();
			break

		case 'f': // Fill
			s.getToken()
			//fill();
			break

		case 'h': // Help
			//help();
			break

		case 'i': // ASCII dump
			s.getToken()
			//ascii_dump();
			break

		case 'k': // Memory configuration
			s.getToken()
			//mem_config();
			break

		case 'l': // Load data
			s.getToken()
			//load_data();
			break

		case 'm': // Memory dump
			s.getToken()
			//memory_dump();
			break

		case 'n': // Screen code dump
			s.getToken()
			//screen_dump()
			break

		case 'o': // Redirect output
			s.getToken()
			//redir_output()
			break

		case 'p': // Sprite dump
			s.getToken()
			//sprite_dump()
			break

		case 'r': // Registers
			//get_reg_token()
			//registers()
			break

		case 's': // Save data
			s.getToken()
			//save_data()
			break

		case 't': // Transfer
			s.getToken()
			//transfer()
			break

		case 'v': // View machine state
			//view_state()
			break

		case 'x': // Exit
			done = true
			break

		case ':': // Change memory
			s.getToken()
			//modify()
			break

		case '1': // Switch to 1541 mode
			//_access_1541 = true
			break

		case '6': // Switch to C64 mode
			//_access_1541 = false
			break

		case '?': // Compute expression
			s.getToken()
			//print_expr()
			break

		case '\n': // Blank line
			break

		default: // Unknown command
			//error("Unknown command")
			break
		}
	}

	//if (_fout != _ferr)
	//fclose(_fout);

	// Set CPU registers
	//TheCPU->setStateMachine(&R64);
	//TheCPU1541->setStateMachine(&R1541);
}

func (s *SAM) getToken() (Token, error) {
	var c rune
	// Skip spaces
	for {
		if c = s.line.getChar(); c != ' ' {
			break
		}
	}

	switch c {
	case '\n':
		s._the_token = T_END
		return s._the_token, nil
	case '(':
		s._the_token = T_LPAREN
		return s._the_token, nil
	case ')':
		s._the_token = T_RPAREN
		return s._the_token, nil
	case '+':
		s._the_token = T_ADD
		return s._the_token, nil
	case '-':
		s._the_token = T_SUB
		return s._the_token, nil
	case '*':
		s._the_token = T_MUL
		return s._the_token, nil
	case '/':
		s._the_token = T_DIV
		return s._the_token, nil
	case ',':
		s._the_token = T_COMMA
		return s._the_token, nil
	case '#':
		s._the_token = T_IMMED
		return s._the_token, nil
	case 'x':
		s._the_token = T_X
		return s._the_token, nil
	case 'y':
		s._the_token = T_Y
		return s._the_token, nil
	case 'p':
		if s.line.getChar() == 'c' {
			s._the_token = T_PC
			return s._the_token, nil
		} else {
			s._the_token = T_NULL
			return s._the_token, fmt.Errorf("unrecognized token")
		}
	case 's':
		if s.line.getChar() == 'p' {
			s._the_token = T_SP
			return s._the_token, nil
		} else {
			s._the_token = T_NULL
			return s._the_token, fmt.Errorf("unrecognized token")
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f':
		s.line.putBack(c)
		s._the_number = s.getNumber()
		s._the_token = T_NUMBER
		return s._the_token, nil
	case '"':
		tkn, err := s.get_string(s._the_string)
		if err != nil {
			s._the_token = T_NULL
			return tkn, err
		}
		s._the_token = tkn
		return s._the_token, nil
	default:
		s._the_token = T_NULL
		return s._the_token, fmt.Errorf("unrecognized token")
	}

}
func (s *SAM) getNumber() uint16 {
	/*
		char c;
		uint16 i = 0;

		while (((c = get_char()) >= '0') && (c <= '9') || (c >= 'a') && (c <= 'f'))
		if (c < 'a')
		i = (i << 4) + (c - '0');
		else
		i = (i << 4) + (c - 'a' + 10);

		put_back(c);
		return i;

	*/

	return 0
}

func (s *SAM) get_string(str string) (Token, error) {
	/*

		char c;

		while ((c = get_char()) != '\n') {
		if (c == '"') {
		*str = 0;
		return T_STRING;
		}
		*str++ = c;

	*/
	//error("Unterminated string");
	return T_NULL, fmt.Errorf("unterminated string")
}

/*
func (s *SAM)  get_reg_token()  Token{
char c;

// Skip spaces
while ((c = get_char()) == ' ');

switch (c) {
case '\n':
return _the_token = T_END;
case 'a':
return _the_token = T_A;
case 'd':
if (get_char() == 'r')
return _the_token = T_DR;
else {
error("Unrecognized token");
return _the_token = T_NULL;
}
case 'p':
if ((c = get_char()) == 'c')
return _the_token = T_PC;
else if (c == 'r')
return _the_token = T_PR;
else {
error("Unrecognized token");
return _the_token = T_NULL;
}
case 's':
if (get_char() == 'p')
return _the_token = T_SP;
else {
error("Unrecognized token");
return _the_token = T_NULL;
}
case 'x':
return _the_token = T_X;
case 'y':
return _the_token = T_Y;
default:
error("Unrecognized token");
return _the_token = T_NULL;
}
}





error("Unterminated string");
return T_NULL;
}



//expression = term {(ADD | SUB) term}
//true: OK, false: Error

bool C64SAM::expression(uint16* number) {
uint16 accu, trm;

if (!term(&accu))
return false;

for (;;)
switch (_the_token) {
case T_ADD:
get_token();
if (!term(&trm))
return false;
accu += trm;
break;

case T_SUB:
get_token();
if (!term(&trm))
return false;
accu -= trm;
break;

default:
*number = accu;
return true;
}
}


//term = factor {(MUL | DIV) factor}
 //true: OK, false: Error


bool C64SAM::term(uint16* number) {
uint16 accu, fact;

if (!factor(&accu))
return false;

for (;;)
switch (_the_token) {
case T_MUL:
get_token();
if (!factor(&fact))
return false;
accu *= fact;
break;

case T_DIV:
get_token();
if (!factor(&fact))
return false;
if (fact == 0) {
error("Division by 0");
return false;
}
accu /= fact;
break;

default:
*number = accu;
return true;
}
}



//factor = NUMBER | PC | SP | LPAREN expression RPAREN
//true: OK, false: Error


bool C64SAM::factor(uint16* number) {
switch (_the_token) {
case T_NUMBER:
*number = _the_number;
get_token();
return true;

case T_PC:
get_token();
*number = _access_1541 ? R1541.pc : R64.pc;
return true;

case T_SP:
get_token();
*number = _access_1541 ? R1541.sp : R64.sp;
return true;

case T_LPAREN:
get_token();
if (expression(number))
if (_the_token == T_RPAREN) {
get_token();
return true;
} else {
error("Missing ')'");
return false;
} else {
error("Error in expression");
return false;
}

case T_END:
error("Required argument missing");
return false;

default:
error("'pc', 'sp', '(' or number expected");
return false;
}
}



//address_args = [expression] END
//Read start to "address"
//true: OK, false: Error
bool C64SAM::address_args(void) {
if (_the_token == T_END)
return true;
else {
if (!expression(&_address))
return false;
return _the_token == T_END;
}
}


//range_args = [expression] [[COMMA] expression] END
//Read start address to "address", end address to "end_address"
//true: OK, false: Error
bool C64SAM::range_args(int16 def_range) {
_end_address = _address + def_range;

if (_the_token == T_END)
return true;
else {
if (!expression(&_address))
return false;
_end_address = _address + def_range;
if (_the_token == T_END)
return true;
else {
if (_the_token == T_COMMA) get_token();
if (!expression(&_end_address))
return false;
return _the_token == T_END;
}
}
}


/*
// instr_args = END
// | IMMED NUMBER END
//| NUMBER [COMMA (X | Y)] END
//| LPAREN NUMBER (RPAREN [COMMA Y] | COMMA X RPAREN) END

//Read arguments of a 6510 instruction, determine address and addressing mode
//
//true: OK, false: Error

bool C64SAM::instr_args(uint16* number, char* mode) {
switch (_the_token) {

case T_END:
*mode = A_IMPL;
return true;

case T_IMMED:
get_token();
if (_the_token == T_NUMBER) {
*number = _the_number;
*mode = A_IMM;
get_token();
return _the_token == T_END;
} else {
error("Number expected");
return false;
}

case T_NUMBER:
*number = _the_number;
get_token();
switch (_the_token) {

case T_END:
if (*number < 0x100)
*mode = A_ZERO;
else
*mode = A_ABS;
return true;

case T_COMMA:
get_token();
switch (_the_token) {

case T_X:
get_token();
if (*number < 0x100)
*mode = A_ZEROX;
else
*mode = A_ABSX;
return _the_token == T_END;

case T_Y:
get_token();
if (*number < 0x100)
*mode = A_ZEROY;
else
*mode = A_ABSY;
return _the_token == T_END;

default:
error("Illegal index register");
return false;
}

default:
return false;
}

case T_LPAREN:
get_token();
if (_the_token == T_NUMBER) {
*number = _the_number;
get_token();
switch (_the_token) {

case T_RPAREN:
get_token();
switch (_the_token) {

case T_END:
*mode = A_IND;
return true;

case T_COMMA:
get_token();
if (_the_token == T_Y) {
*mode = A_INDY;
get_token();
return _the_token == T_END;
} else {
error("Only 'y' index register allowed");
return false;
}

default:
error("Illegal characters after ')'");
return false;
}

case T_COMMA:
get_token();
if (_the_token == T_X) {
get_token();
if (_the_token == T_RPAREN) {
*mode = A_INDX;
get_token();
return _the_token == T_END;
} else {
error("')' expected");
return false;
}
} else {
error("Only 'x' index register allowed");
return false;
}

default:
error("')' or ',' expected");
return false;
}
} else {
error("Number expected");
return false;
}

default:
error("'(', '#' or number expected");
return false;
}
}

//Display help
//h


void C64SAM::help(void) {
fprintf(_fout, "a [start]           Assemble\n"
"b [start] [end]     Binary dump\n"
"c start end dest    Compare memory\n"
"d [start] [end]     Disassemble\n"
"e                   Show interrupt vectors\n"
"f start end byte    Fill memory\n"
"i [start] [end]     ASCII/PETSCII dump\n"
"k [config]          Show/set C64 memory configuration\n"
"l start \"file\"      Load data\n"
"m [start] [end]     Memory dump\n"
"n [start] [end]     Screen code dump\n"
"o [\"file\"]          Redirect output\n"
"p [start] [end]     Sprite dump\n"
"r [reg value]       Show/set CPU registers\n"
"s start end \"file\"  Save data\n"
"t start end dest    Transfer memory\n"
"vc1                 View CIA 1 state\n"
"vc2                 View CIA 2 state\n"
"vf                  View 1541 state\n"
"vs                  View SID state\n"
"vv                  View VIC state\n"
"x                   Return to Pet64\n"
": addr {byte}       Modify memory\n"
"1541                Switch to 1541\n"
"64                  Switch to C64\n"
"? expression        Calculate expression\n");
}



//Display/change 6510 registers
 //r [reg value]


void C64SAM::registers(void) {
enum Token the_reg;
uint16 value;

if (_the_token != T_END)
switch (the_reg = _the_token) {
case T_A:
case T_X:
case T_Y:
case T_PC:
case T_SP:
case T_DR:
case T_PR:
get_token();
if (!expression(&value))
return;

switch (the_reg) {
case T_A:
if (_access_1541)
R1541.a = (uint8)value;
else
R64.a = (uint8)value;
break;
case T_X:
if (_access_1541)
R1541.x = (uint8)value;
else
R64.x = (uint8)value;
break;
case T_Y:
if (_access_1541)
R1541.y = (uint8)value;
else
R64.y = (uint8)value;
break;
case T_PC:
if (_access_1541)
R1541.pc = value;
else
R64.pc = value;
break;
case T_SP:
if (_access_1541)
R1541.sp = (value & 0xff) | 0x0100;
else
R64.sp = (value & 0xff) | 0x0100;
break;
case T_DR:
if (!_access_1541)
R64.ddr = (uint8)value;
break;
case T_PR:
if (!_access_1541)
R64.pr = (uint8)value;
break;
default:
break;
}
break;

default:
return;
}

display_registers();
}

void C64SAM::display_registers(void) {
if (_access_1541) {
fprintf(_fout, " PC  A  X  Y   SP  NVDIZC  Instruction\n");
fprintf(_fout, "%04lx %02lx %02lx %02lx %04lx %c%c%c%c%c%c ",
R1541.pc, R1541.a, R1541.x, R1541.y, R1541.sp,
R1541.p & 0x80 ? '1' : '0', R1541.p & 0x40 ? '1' : '0', R1541.p & 0x08 ? '1' : '0',
R1541.p & 0x04 ? '1' : '0', R1541.p & 0x02 ? '1' : '0', R1541.p & 0x01 ? '1' : '0');
disass_line(R1541.pc, SAMReadByte(R1541.pc), SAMReadByte(R1541.pc + 1), SAMReadByte(R1541.pc + 2));
} else {
fprintf(_fout, " PC  A  X  Y   SP  DR PR NVDIZC  Instruction\n");
fprintf(_fout, "%04lx %02lx %02lx %02lx %04lx %02lx %02lx %c%c%c%c%c%c ",
R64.pc, R64.a, R64.x, R64.y, R64.sp, R64.ddr, R64.pr,
R64.p & 0x80 ? '1' : '0', R64.p & 0x40 ? '1' : '0', R64.p & 0x08 ? '1' : '0',
R64.p & 0x04 ? '1' : '0', R64.p & 0x02 ? '1' : '0', R64.p & 0x01 ? '1' : '0');
disass_line(R64.pc, SAMReadByte(R64.pc), SAMReadByte(R64.pc + 1), SAMReadByte(R64.pc + 2));
}
}


//Memory dump
//m [start] [end]


#define MEMDUMP_BPL 16  // Bytes per line

void C64SAM::memory_dump(void) {
bool done = false;
short i;
uint8 mem[MEMDUMP_BPL + 2];
uint8 byte;

mem[MEMDUMP_BPL] = 0;

if (!range_args(16 * MEMDUMP_BPL - 1))  // 16 lines unless end address specified
return;

do {
fprintf(_fout, "%04lx:", _address);
for (i = 0; i < MEMDUMP_BPL; i++, _address++) {
if (_address == _end_address) done = true;

fprintf(_fout, " %02lx", byte = SAMReadByte(_address));
if ((byte >= ' ') && (byte <= '~'))
mem[i] = conv_from_64(byte);
else
mem[i] = '.';
}
fprintf(_fout, "  '%s'\n", mem);
} while (!done && !aborted());
}



//ASCII dump
 //i [start] [end]


#define ASCIIDUMP_BPL 64  // Bytes per line

void C64SAM::ascii_dump(void) {
bool done = false;
short i;
uint8 mem[ASCIIDUMP_BPL + 2];
uint8 byte;

mem[ASCIIDUMP_BPL] = 0;

if (!range_args(16 * ASCIIDUMP_BPL - 1))  // 16 lines unless end address specified
return;

do {
fprintf(_fout, "%04lx:", _address);
for (i = 0; i < ASCIIDUMP_BPL; i++, _address++) {
if (_address == _end_address) done = true;

byte = SAMReadByte(_address);
if ((byte >= ' ') && (byte <= '~'))
mem[i] = conv_from_64(byte);
else
mem[i] = '.';
}
fprintf(_fout, " '%s'\n", mem);
} while (!done && !aborted());
}



//Convert PETSCII->ASCII

char C64SAM::conv_from_64(char c) {
if ((c >= 'A') && (c <= 'Z') || (c >= 'a') && (c <= 'z'))
return c ^ 0x20;
else
return c;
}



//Screen code dump
//n [start] [end]


#define SCRDUMP_BPL 64  // Bytes per line

void C64SAM::screen_dump(void) {
bool done = false;
short i;
uint8 mem[SCRDUMP_BPL + 2];
uint8 byte;

mem[SCRDUMP_BPL] = 0;

if (!range_args(16 * SCRDUMP_BPL - 1))  // 16 Zeilen unless end address specified
return;

do {
fprintf(_fout, "%04lx:", _address);
for (i = 0; i < SCRDUMP_BPL; i++, _address++) {
if (_address == _end_address) done = true;

byte = SAMReadByte(_address);
if (byte < 90)
mem[i] = conv_from_scode(byte);
else
mem[i] = '.';
}
fprintf(_fout, " '%s'\n", mem);
} while (!done && !aborted());
}



//Convert _screen code->ASCII


char C64SAM::conv_from_scode(char c) {
c &= 0x7f;

if (c <= 31)
return c + 64;
else
if (c >= 64)
return c + 32;
else
return c;
}



//Binary dump
//b [start] [end]

void C64SAM::binary_dump(void) {
bool done = false;
char bin[10];

bin[8] = 0;

if (!range_args(7))  // 8 lines unless end address specified
return;

do {
if (_address == _end_address) done = true;

byte_to_bin(SAMReadByte(_address), bin);
fprintf(_fout, "%04lx: %s\n", _address++, bin);
} while (!done && !aborted());
}



//Sprite data dump
//p [start] [end]

void C64SAM::sprite_dump(void) {
bool done = false;
short i;
char bin[10];

bin[8] = 0;

if (!range_args(21 * 3 - 1))  // 21 lines unless end address specified
return;

do {
fprintf(_fout, "%04lx: ", _address);
for (i = 0; i < 3; i++, _address++) {
if (_address == _end_address) done = true;

byte_to_bin(SAMReadByte(_address), bin);
fprintf(_fout, "%s", bin);
}
fputc('\n', _fout);
} while (!done && !aborted());
}


//Convert byte to binary representation

void C64SAM::byte_to_bin(uint8 byte, char* str) {
short i;

for (i = 0; i < 8; i++, byte <<= 1)
if (byte & 0x80)
str[i] = '#';
else
str[i] = '.';
}


//Disassemble
//d [start] [end]

void C64SAM::disassemble(void) {
bool done = false;
short i;
uint8 op[3];
uint16 adr;

if (!range_args(31))  // 32 bytes unless end address specified
return;

do {
fprintf(_fout, "%04lx:", adr = _address);
for (i = 0; i < 3; i++, adr++) {
if (adr == _end_address) done = true;
op[i] = SAMReadByte(adr);
}
_address += (uint16)disass_line(_address, op[0], op[1], op[2]);
} while (!done && !aborted());
}


//Disassemble one instruction, return length
int C64SAM::disass_line(uint16 adr, uint8 op, uint8 lo, uint8 hi) {
char mode = adr_mode[op], mnem = mnemonic[op];

// Display instruction bytes in hex
switch (_adr_length[mode]) {
case 1:
fprintf(_fout, " %02lx       ", op);
break;

case 2:
fprintf(_fout, " %02lx %02lx    ", op, lo);
break;

case 3:
fprintf(_fout, " %02lx %02lx %02lx ", op, lo, hi);
break;
}

// Tag undocumented opcodes with an asterisk
if (mnem > M_ILLEGAL)
fputc('*', _fout);
else
fputc(' ', _fout);

// Print mnemonic
fprintf(_fout, "%c%c%c ", _mnem_1[mnem], _mnem_2[mnem], _mnem_3[mnem]);

// Pring argument
switch (mode) {
case A_IMPL:
break;

case A_ACCU:
fprintf(_fout, "a");
break;

case A_IMM:
fprintf(_fout, "#%02lx", lo);
break;

case A_REL:
fprintf(_fout, "%04lx", ((adr + 2) + (int8)lo) & 0xffff);
break;

case A_ZERO:
fprintf(_fout, "%02lx", lo);
break;

case A_ZEROX:
fprintf(_fout, "%02lx,x", lo);
break;

case A_ZEROY:
fprintf(_fout, "%02lx,y", lo);
break;

case A_ABS:
fprintf(_fout, "%04lx", (hi << 8) | lo);
break;

case A_ABSX:
fprintf(_fout, "%04lx,x", (hi << 8) | lo);
break;

case A_ABSY:
fprintf(_fout, "%04lx,y", (hi << 8) | lo);
break;

case A_IND:
fprintf(_fout, "(%04lx)", (hi << 8) | lo);
break;

case A_INDX:
fprintf(_fout, "(%02lx,x)", lo);
break;

case A_INDY:
fprintf(_fout, "(%02lx),y", lo);
break;
}

fputc('\n', _fout);
return _adr_length[mode];
}


//Assemble
//a [start]


void C64SAM::assemble(void) {
bool done = false;
char c1, c2, c3;
char mnem, mode;
uint8 opcode;
uint16 arg;
int16 rel;

// Read parameters
if (!address_args())
return;

do {
fprintf(_fout, "%04lx> ", _address);
fflush(_ferr);
read_line();

c1 = get_char();
c2 = get_char();
c3 = get_char();

if (c1 != '\n') {

if ((mnem = find_mnemonic(c1, c2, c3)) != M_ILLEGAL) {

get_token();
if (instr_args(&arg, &mode)) {

// Convert A_IMPL -> A_ACCU if necessary
if ((mode == A_IMPL) && find_opcode(mnem, A_ACCU, &opcode))
mode = A_ACCU;

// Handle relative addressing seperately
if (((mode == A_ABS) || (mode == A_ZERO)) && find_opcode(mnem, A_REL, &opcode)) {
mode = A_REL;
rel = arg - (_address + 2) & 0xffff;
if ((rel < -128) || (rel > 127)) {
error("Branch too long");
continue;
} else
arg = rel & 0xff;
}

if (find_opcode(mnem, mode, &opcode)) {

// Print disassembled line
fprintf(_fout, "\v%04lx:", _address);
disass_line(_address, opcode, arg & 0xff, arg >> 8);

switch (_adr_length[mode]) {
case 1:
SAMWriteByte(_address++, opcode);
break;

case 2:
SAMWriteByte(_address++, opcode);
SAMWriteByte(_address++, (uint8)arg);
break;

case 3:
SAMWriteByte(_address++, opcode);
SAMWriteByte(_address++, arg & 0xff);
SAMWriteByte(_address++, arg >> 8);
break;

default:
error("Internal error");
break;
}

} else
error("Addressing mode not supported by instruction");

} else
error("Unrecognized addressing mode");

} else
error("Unknown instruction");

} else			// Input is terminated with a blank line
done = true;
} while (!done);
}



//Find mnemonic code to three letters
//M_ILLEGAL: No matching mnemonic found

char C64SAM::find_mnemonic(char op1, char op2, char op3) {
int i;

for (i = 0; i < M_MAXIMUM; i++) {
if ((_mnem_1[i] == op1) && (_mnem_2[i] == op2) && (_mnem_3[i] == op3)) {
return (char)i;
}
}

return M_ILLEGAL;
}

//Determine opcode of an instruction given mnemonic and addressing mode
//true: OK, false: Mnemonic can't have specified addressing mode
bool C64SAM::find_opcode(char mnem, char mode, uint8* opcode) {
int i;

for (i = 0; i < 256; i++) {
if ((mnemonic[i] == mnem) && (adr_mode[i] == mode)) {
*opcode = (uint8)i;
return true;
}
}

return false;
}



//Show/set memory configuration
//k [config]
void C64SAM::mem_config(void) {
uint16 con;

if (_the_token != T_END)
if (!expression(&con))
return;
else
TheCPU->_ExtConfig = con;
else
con = (uint16)TheCPU->_ExtConfig;

fprintf(_fout, "Configuration: %ld\n", con & 7);
fprintf(_fout, "A000-BFFF: %s\n", (con & 3) == 3 ? "Basic" : "RAM");
fprintf(_fout, "D000-DFFF: %s\n", (con & 3) ? ((con & 4) ? "I/O" : "Char") : "RAM");
fprintf(_fout, "E000-FFFF: %s\n", (con & 2) ? "Kernal" : "RAM");
}



// Fill
 //f start end byte


void C64SAM::fill(void) {
bool done = false;
uint16 adr, end_adr, value;

if (!expression(&adr))
return;
if (!expression(&end_adr))
return;
if (!expression(&value))
return;

do {
if (adr == end_adr) done = true;

SAMWriteByte(adr++, (uint8)value);
} while (!done);
}


//Compare
//c start end dest
void C64SAM::compare(void) {
bool done = false;
uint16 adr, end_adr, dest;
int num = 0;

if (!expression(&adr))
return;
if (!expression(&end_adr))
return;
if (!expression(&dest))
return;

do {
if (adr == end_adr) done = true;

if (SAMReadByte(adr) != SAMReadByte(dest)) {
fprintf(_fout, "%04lx ", adr);
num++;
if (!(num & 7))
fputc('\n', _fout);
}
adr++; dest++;
} while (!done && !aborted());

if (num & 7)
fputc('\n', _fout);
fprintf(_fout, "%ld byte(s) different\n", num);
}



//Transfer memory
//t start end dest
void C64SAM::transfer(void) {
bool done = false;
uint16 adr, end_adr, dest;

if (!expression(&adr))
return;
if (!expression(&end_adr))
return;
if (!expression(&dest))
return;

if (dest < adr)
do {
if (adr == end_adr) done = true;
SAMWriteByte(dest++, SAMReadByte(adr++));
} while (!done);
else {
dest += end_adr - adr;
do {
if (adr == end_adr) done = true;
SAMWriteByte(dest--, SAMReadByte(end_adr--));
} while (!done);
}
}


//Change memory
//: addr {byte}


void C64SAM::modify(void) {
uint16 adr, val;

if (!expression(&adr))
return;

while (_the_token != T_END)
if (expression(&val))
SAMWriteByte(adr++, (uint8)val);
else
return;
}


//Compute and display expression
//? expression

void C64SAM::print_expr(void) {
uint16 val;

if (!expression(&val))
return;

fprintf(_fout, "Hex: %04lx\nDec: %lu\n", val, val);
}



//Redirect output
//o [file]
void C64SAM::redir_output(void) {
// Close old file
if (_fout != _ferr) {
fclose(_fout);
_fout = _ferr;
return;
}

// No argument given?
if (_the_token == T_END)
return;

// Otherwise open file
if (_the_token == T_STRING) {
_fout = fopen(_the_string, "w");
if (!_fout) {
error("Unable to open file");
}
} else
error("'\"' around file name expected");
}


//Display interrupt vectors
void C64SAM::int_vectors(void) {
fprintf(_fout, "        IRQ  BRK  NMI\n");
fprintf(_fout, "%d  : %04lx %04lx %04lx\n",
_access_1541 ? 6502 : 6510,
SAMReadByte(0xffff) << 8 | SAMReadByte(0xfffe),
SAMReadByte(0xffff) << 8 | SAMReadByte(0xfffe),
SAMReadByte(0xfffb) << 8 | SAMReadByte(0xfffa));

if (!_access_1541 && TheCPU->_ExtConfig & 2)
fprintf(_fout, "Kernal: %04lx %04lx %04lx\n",
SAMReadByte(0x0315) << 8 | SAMReadByte(0x0314),
SAMReadByte(0x0317) << 8 | SAMReadByte(0x0316),
SAMReadByte(0x0319) << 8 | SAMReadByte(0x0318));
}



//Display state of custom chips
void C64SAM::view_state(void) {
switch (get_char()) {
case 'c':		// CIA
view_cia_state();
break;

case 's':		// SID
view_sid_state();
break;

case 'v':		// VIC
view_vic_state();
break;

case 'f':		// Floppy
view_1541_state();
break;

default:
error("Unknown command");
break;
}
}

void C64SAM::view_cia_state(void) {
MOS6526Registers cs;

switch (get_char()) {
case '1':
TheCIA1->getStateMachine(&cs);
break;
case '2':
TheCIA2->getStateMachine(&cs);
break;
default:
error("Unknown command");
return;
}

fprintf(_fout, "Timer A  : %s\n", cs._cra & 1 ? "On" : "Off");
fprintf(_fout, " Counter : %04lx  Latch: %04lx\n", cs._ta, cs._latcha);
fprintf(_fout, " Run mode: %s\n", cs._cra & 8 ? "One-shot" : "Continuous");
fprintf(_fout, " Input   : %s\n", cs._cra & 0x20 ? "CNT" : "Phi2");
fprintf(_fout, " Output  : ");
if (cs._cra & 2)
if (cs._cra & 4)
fprintf(_fout, "PB6 Toggle\n\n");
else
fprintf(_fout, "PB6 Pulse\n\n");
else
fprintf(_fout, "None\n\n");

fprintf(_fout, "Timer B  : %s\n", cs._crb & 1 ? "On" : "Off");
fprintf(_fout, " Counter : %04lx  Latch: %04lx\n", cs._tb, cs._latchb);
fprintf(_fout, " Run mode: %s\n", cs._crb & 8 ? "One-shot" : "Continuous");
fprintf(_fout, " Input   : ");
if (cs._crb & 0x40)
if (cs._crb & 0x20)
fprintf(_fout, "Timer A underflow (CNT high)\n");
else
fprintf(_fout, "Timer A underflow\n");
else
if (cs._crb & 0x20)
fprintf(_fout, "CNT\n");
else
fprintf(_fout, "Phi2\n");
fprintf(_fout, " Output  : ");
if (cs._crb & 2)
if (cs._crb & 4)
fprintf(_fout, "PB7 Toggle\n\n");
else
fprintf(_fout, "PB7 Pulse\n\n");
else
fprintf(_fout, "None\n\n");

fprintf(_fout, "TOD         : %lx%lx:%lx%lx:%lx%lx.%lx %s\n",
(cs._tod_hr >> 4) & 1, cs._tod_hr & 0x0f,
(cs._tod_min >> 4) & 7, cs._tod_min & 0x0f,
(cs._tod_sec >> 4) & 7, cs._tod_sec & 0x0f,
cs._tod_10ths & 0x0f, cs._tod_hr & 0x80 ? "PM" : "AM");
fprintf(_fout, "Alarm       : %lx%lx:%lx%lx:%lx%lx.%lx %s\n",
(cs._alm_hr >> 4) & 1, cs._alm_hr & 0x0f,
(cs._alm_min >> 4) & 7, cs._alm_min & 0x0f,
(cs._alm_sec >> 4) & 7, cs._alm_sec & 0x0f,
cs._alm_10ths & 0x0f, cs._alm_hr & 0x80 ? "PM" : "AM");
fprintf(_fout, "TOD input   : %s\n", cs._cra & 0x80 ? "50Hz" : "60Hz");
fprintf(_fout, "Write to    : %s registers\n\n", cs._crb & 0x80 ? "Alarm" : "TOD");

fprintf(_fout, "Serial data : %02lx\n", cs._sdr);
fprintf(_fout, "Serial mode : %s\n\n", cs._cra & 0x40 ? "Output" : "Input");

fprintf(_fout, "Pending int.: ");
dump_cia_ints(cs._icr);
fprintf(_fout, "Enabled int.: ");
dump_cia_ints(cs._int_mask);
}

void C64SAM::dump_cia_ints(uint8 i) {
if (i & 0x1f) {
if (i & 1) fprintf(_fout, "TA ");
if (i & 2) fprintf(_fout, "TB ");
if (i & 4) fprintf(_fout, "Alarm ");
if (i & 8) fprintf(_fout, "Serial ");
if (i & 0x10) fprintf(_fout, "Flag");
} else
fprintf(_fout, "None");
fputc('\n', _fout);
}

void C64SAM::view_sid_state(void) {
MOS6581Register ss;

TheSID->getStateMachine(&ss);

fprintf(_fout, "Voice 1\n");
fprintf(_fout, " Frequency  : %04lx\n", (ss._freq_hi_1 << 8) | ss._freq_lo_1);
fprintf(_fout, " Pulse Width: %04lx\n", ((ss._pw_hi_1 & 0x0f) << 8) | ss._pw_lo_1);
fprintf(_fout, " Env. (ADSR): %lx %lx %lx %lx\n", ss._AD_1 >> 4, ss._AD_1 & 0x0f, ss._SR_1 >> 4, ss._SR_1 & 0x0f);
fprintf(_fout, " Waveform   : ");
dump_sid_waveform(ss._ctrl_1);
fprintf(_fout, " Gate       : %s  Ring mod.: %s\n", ss._ctrl_1 & 0x01 ? "On " : "Off", ss._ctrl_1 & 0x04 ? "On" : "Off");
fprintf(_fout, " Test bit   : %s  Synchron.: %s\n", ss._ctrl_1 & 0x08 ? "On " : "Off", ss._ctrl_1 & 0x02 ? "On" : "Off");
fprintf(_fout, " Filter     : %s\n", ss._res_filt & 0x01 ? "On" : "Off");

fprintf(_fout, "\nVoice 2\n");
fprintf(_fout, " Frequency  : %04lx\n", (ss._freq_hi_2 << 8) | ss._freq_lo_2);
fprintf(_fout, " Pulse Width: %04lx\n", ((ss._pw_hi_2 & 0x0f) << 8) | ss._pw_lo_2);
fprintf(_fout, " Env. (ADSR): %lx %lx %lx %lx\n", ss._AD_2 >> 4, ss._AD_2 & 0x0f, ss._SR_2 >> 4, ss._SR_2 & 0x0f);
fprintf(_fout, " Waveform   : ");
dump_sid_waveform(ss._ctrl_2);
fprintf(_fout, " Gate       : %s  Ring mod.: %s\n", ss._ctrl_2 & 0x01 ? "On " : "Off", ss._ctrl_2 & 0x04 ? "On" : "Off");
fprintf(_fout, " Test bit   : %s  Synchron.: %s\n", ss._ctrl_2 & 0x08 ? "On " : "Off", ss._ctrl_2 & 0x02 ? "On" : "Off");
fprintf(_fout, " Filter     : %s\n", ss._res_filt & 0x02 ? "On" : "Off");

fprintf(_fout, "\nVoice 3\n");
fprintf(_fout, " Frequency  : %04lx\n", (ss._freq_hi_3 << 8) | ss._freq_lo_3);
fprintf(_fout, " Pulse Width: %04lx\n", ((ss._pw_hi_3 & 0x0f) << 8) | ss._pw_lo_3);
fprintf(_fout, " Env. (ADSR): %lx %lx %lx %lx\n", ss._AD_3 >> 4, ss._AD_3 & 0x0f, ss._SR_3 >> 4, ss._SR_3 & 0x0f);
fprintf(_fout, " Waveform   : ");
dump_sid_waveform(ss._ctrl_3);
fprintf(_fout, " Gate       : %s  Ring mod.: %s\n", ss._ctrl_3 & 0x01 ? "On " : "Off", ss._ctrl_3 & 0x04 ? "On" : "Off");
fprintf(_fout, " Test bit   : %s  Synchron.: %s\n", ss._ctrl_3 & 0x08 ? "On " : "Off", ss._ctrl_3 & 0x02 ? "On" : "Off");
fprintf(_fout, " Filter     : %s  Mute     : %s\n", ss._res_filt & 0x04 ? "On" : "Off", ss._mode_vol & 0x80 ? "Yes" : "No");

fprintf(_fout, "\nFilters/Volume\n");
fprintf(_fout, " Frequency: %04lx\n", (ss._fc_hi << 3) | (ss._fc_lo & 0x07));
fprintf(_fout, " Resonance: %lx\n", ss._res_filt >> 4);
fprintf(_fout, " Mode     : ");
if (ss._mode_vol & 0x70) {
if (ss._mode_vol & 0x10) fprintf(_fout, "Low-pass ");
if (ss._mode_vol & 0x20) fprintf(_fout, "Band-pass ");
if (ss._mode_vol & 0x40) fprintf(_fout, "High-pass");
} else
fprintf(_fout, "None");
fprintf(_fout, "\n Volume   : %lx\n", ss._mode_vol & 0x0f);
}

void C64SAM::dump_sid_waveform(uint8 wave) {
if (wave & 0xf0) {
if (wave & 0x10) fprintf(_fout, "Triangle ");
if (wave & 0x20) fprintf(_fout, "Sawtooth ");
if (wave & 0x40) fprintf(_fout, "Rectangle ");
if (wave & 0x80) fprintf(_fout, "Noise");
} else
fprintf(_fout, "None");
fputc('\n', _fout);
}

void C64SAM::view_vic_state(void) {
MOS6569State vs;
short i;

TheVIC->getStateMachine(&vs);

fprintf(_fout, "Raster line       : %04lx\n", vs.raster | ((vs.ctrl1 & 0x80) << 1));
fprintf(_fout, "IRQ raster line   : %04lx\n\n", vs.irq_raster);

fprintf(_fout, "X scroll          : %ld\n", vs.ctrl2 & 7);
fprintf(_fout, "Y scroll          : %ld\n", vs.ctrl1 & 7);
fprintf(_fout, "Horizontal border : %ld columns\n", vs.ctrl2 & 8 ? 40 : 38);
fprintf(_fout, "Vertical border   : %ld rows\n\n", vs.ctrl1 & 8 ? 25 : 24);

fprintf(_fout, "Display mode      : ");
switch (((vs.ctrl1 >> 4) & 6) | ((vs.ctrl2 >> 4) & 1)) {
case 0:
fprintf(_fout, "Standard text\n");
break;
case 1:
fprintf(_fout, "Multicolor text\n");
break;
case 2:
fprintf(_fout, "Standard bitmap\n");
break;
case 3:
fprintf(_fout, "Multicolor bitmap\n");
break;
case 4:
fprintf(_fout, "ECM text\n");
break;
case 5:
fprintf(_fout, "Invalid text (ECM+MCM)\n");
break;
case 6:
fprintf(_fout, "Invalid bitmap (ECM+BMM)\n");
break;
case 7:
fprintf(_fout, "Invalid bitmap (ECM+BMM+ECM)\n");
break;
}
fprintf(_fout, "Sequencer state   : %s\n", vs.display_state ? "Display" : "Idle");
fprintf(_fout, "Bad line state    : %s\n", vs.bad_line ? "Yes" : "No");
fprintf(_fout, "Bad lines enabled : %s\n", vs.bad_line_enable ? "Yes" : "No");
fprintf(_fout, "Video counter     : %04lx\n", vs.vc);
fprintf(_fout, "Video counter base: %04lx\n", vs.vc_base);
fprintf(_fout, "Row counter       : %ld\n\n", vs.rc);

fprintf(_fout, "VIC bank          : %04lx-%04lx\n", vs.bank_base, vs.bank_base + 0x3fff);
fprintf(_fout, "Video matrix base : %04lx\n", vs.matrix_base);
fprintf(_fout, "Character base    : %04lx\n", vs.char_base);
fprintf(_fout, "Bitmap base       : %04lx\n\n", vs.bitmap_base);

fprintf(_fout, "         Spr.0  Spr.1  Spr.2  Spr.3  Spr.4  Spr.5  Spr.6  Spr.7\n");
fprintf(_fout, "Enabled: "); dump_spr_flags(vs.me);
fprintf(_fout, "Data   : %04lx   %04lx   %04lx   %04lx   %04lx   %04lx   %04lx   %04lx\n",
vs.sprite_base[0], vs.sprite_base[1], vs.sprite_base[2], vs.sprite_base[3],
vs.sprite_base[4], vs.sprite_base[5], vs.sprite_base[6], vs.sprite_base[7]);
fprintf(_fout, "MC     : %02lx     %02lx     %02lx     %02lx     %02lx     %02lx     %02lx     %02lx\n",
vs.mc[0], vs.mc[1], vs.mc[2], vs.mc[3], vs.mc[4], vs.mc[5], vs.mc[6], vs.mc[7]);

fprintf(_fout, "Mode   : ");
for (i = 0; i < 8; i++)
if (vs.mmc & (1 << i))
fprintf(_fout, "Multi  ");
else
fprintf(_fout, "Std.   ");

fprintf(_fout, "\nX-Exp. : "); dump_spr_flags(vs.mxe);
fprintf(_fout, "Y-Exp. : "); dump_spr_flags(vs.mye);

fprintf(_fout, "Prio.  : ");
for (i = 0; i < 8; i++)
if (vs.mdp & (1 << i))
fprintf(_fout, "Back   ");
else
fprintf(_fout, "Fore   ");

fprintf(_fout, "\nSS Coll: "); dump_spr_flags(vs.mm);
fprintf(_fout, "SD Coll: "); dump_spr_flags(vs.md);

fprintf(_fout, "\nPending interrupts: ");
dump_vic_ints(vs.irq_flag);
fprintf(_fout, "Enabled interrupts: ");
dump_vic_ints(vs.irq_mask);
}

void C64SAM::dump_spr_flags(uint8 f) {
short i;

for (i = 0; i < 8; i++)
if (f & (1 << i))
fprintf(_fout, "Yes    ");
else
fprintf(_fout, "No     ");

fputc('\n', _fout);
}

void C64SAM::dump_vic_ints(uint8 i) {
if (i & 0x1f) {
if (i & 1) fprintf(_fout, "Raster ");
if (i & 2) fprintf(_fout, "Spr-Data ");
if (i & 4) fprintf(_fout, "Spr-Spr ");
if (i & 8) fprintf(_fout, "Lightpen");
} else
fprintf(_fout, "None");
fputc('\n', _fout);
}

void C64SAM::view_1541_state(void) {
fprintf(_fout, "VIA 1:\n");
fprintf(_fout, " Timer 1 Counter: %04x  Latch: %04x\n", R1541.via1._t1c, R1541.via1._t1l);
fprintf(_fout, " Timer 2 Counter: %04x  Latch: %04x\n", R1541.via1._t2c, R1541.via1._t2l);
fprintf(_fout, " ACR: %02x\n", R1541.via1._acr);
fprintf(_fout, " PCR: %02x\n", R1541.via1._pcr);
fprintf(_fout, " Pending interrupts: ");
dump_via_ints(R1541.via1._ifr);
fprintf(_fout, " Enabled interrupts: ");
dump_via_ints(R1541.via1._ier);

fprintf(_fout, "\nVIA 2:\n");
fprintf(_fout, " Timer 1 Counter: %04x  Latch: %04x\n", R1541.via2._t1c, R1541.via2._t1l);
fprintf(_fout, " Timer 2 Counter: %04x  Latch: %04x\n", R1541.via2._t2c, R1541.via2._t2l);
fprintf(_fout, " ACR: %02x\n", R1541.via2._acr);
fprintf(_fout, " PCR: %02x\n", R1541.via2._pcr);
fprintf(_fout, " Pending interrupts: ");
dump_via_ints(R1541.via2._ifr);
fprintf(_fout, " Enabled interrupts: ");
dump_via_ints(R1541.via2._ier);
}

void C64SAM::dump_via_ints(uint8 i) {
if (i & 0x7f) {
if (i & 0x40) fprintf(_fout, "T1 ");
if (i & 0x20) fprintf(_fout, "T2 ");
if (i & 2) fprintf(_fout, "CA1 ");
if (i & 1) fprintf(_fout, "CA2 ");
if (i & 0x10) fprintf(_fout, "CB1 ");
if (i & 8) fprintf(_fout, "CB2 ");
if (i & 4) fprintf(_fout, "Serial ");
} else
fprintf(_fout, "None");
fputc('\n', _fout);
}



//Load data
//l start "file"


void C64SAM::load_data(void) {
uint16 adr;
FILE* file;
int fc;

if (!expression(&adr))
return;
if (_the_token == T_END) {
error("Missing file name");
return;
}
if (_the_token != T_STRING) {
error("'\"' around file name expected");
return;
}
file = fopen(_the_string, "rb");
if (!file)
error("Unable to open file");
else {
while ((fc = fgetc(file)) != EOF)
SAMWriteByte(adr++, (uint8)fc);
fclose(file);
}
}


//Save data
//s start end "file"
void C64SAM::save_data(void) {
bool done = false;
uint16 adr, end_adr;
FILE* file;

if (!expression(&adr))
return;
if (!expression(&end_adr))
return;
if (_the_token == T_END) {
error("Missing file name");
return;
}
if (_the_token != T_STRING) {
error("'\"' around file name expected");
return;
}

file = fopen(_the_string, "wb");
if (!file)
error("Unable to create file");
else {
do {
if (adr == end_adr) done = true;

fputc(SAMReadByte(adr++), file);
} while (!done);
fclose(file);
}
}

// Access to 6510/6502 address space
uint8 C64SAM::SAMReadByte(uint16 adr) {
if (_access_1541)
return TheCPU1541->ExtReadByte(adr);
else
return TheCPU->ExtReadByte(adr);
}

void C64SAM::SAMWriteByte(uint16 adr, uint8 byte) {
if (_access_1541)
TheCPU1541->ExtWriteByte(adr, byte);
else
TheCPU->ExtWriteByte(adr, byte);
}

*/
