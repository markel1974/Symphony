package sdk

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/objects"
)

// Definizioni dei flag Z80 per chiarezza
const (
	FlagC  = 1 << 0 // Carry Flag
	FlagN  = 1 << 1 // Add/Subtract Flag
	FlagPV = 1 << 2 // Parity/Overflow Flag
	FlagH  = 1 << 4 // Half Carry Flag
	FlagZ  = 1 << 6 // Zero Flag
	FlagS  = 1 << 7 // Sign Flag
)

func init() {
	//RegisterPackage(NewZ80)
}

// Z80Package incapsula le operazioni della CPU Z80 come funzioni di libreria.
type Z80Package struct {
	container map[string]objects.IObject
}

// NewZ80 crea e registra il pacchetto SDK per l'emulazione Z80.
func NewZ80(factory objects.IGateKeeper) *Z80Package /* IPackage */ {
	z := &Z80Package{}
	container := []objects.IObject{
		// La nostra funzione ALU centrale. Prende l'operazione come primo argomento.
		factory.NewFuncImport(objects.FrameStatic, "alu", 4, z.alu),
	}

	fmt.Println(container)
	z.container = make(map[string]objects.IObject)
	//z.container = BuildContainer(container, nil)
	return z
}

func (z *Z80Package) Name() string { return "z80" }

func (z *Z80Package) Get(name string) (objects.IObject, bool) {
	v, ok := z.container[name]
	return v, ok
}

func (z *Z80Package) alu(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	return 0, nil, fmt.Errorf("not implemented")
}

/*

// alu emula la Arithmetic Logic Unit dello Z80.
// Esegue un'operazione e aggiorna i registri A e F direttamente.
// Argomenti:
// args[0]: IObject (Int) -> L'operatore (es. objects.OperatorAdd).
// args[1]: IObject (Int) -> Il registro A.
// args[2]: IObject (Int) -> Il registro F (dei flag).
// args[3]: IObject (Int) -> L'operando (il valore del secondo registro, es. B).
func (z *Z80Package) alu(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	// 1. Validazione e recupero degli argomenti
	if len(args) != 4 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}

	op, ok := gk.ToInt64(args[0])
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(0, "int", args[0].TypeName())
	}

	regA, ok := args[1].(*objects.Int)
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(1, "int", args[1].TypeName())
	}

	regF, ok := args[2].(*objects.Int)
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(2, "int", args[2].TypeName())
	}

	operand, ok := gk.ToInt64(args[3])
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(3, "int", args[3].TypeName())
	}

	// Converte in tipi a 8-bit per i calcoli
	valA := uint8(regA.AsInt64())
	valOperand := uint8(operand)
	var result uint8
	var newFlags uint8

	// 2. Esecuzione dell'operazione e calcolo dei flag
	switch objects.ArithmeticOperator(op) {
	case objects.OperatorAdd:
		// Calcolo del risultato e del carry
		res16 := uint16(valA) + uint16(valOperand)
		result = uint8(res16)

		// Calcolo dei flag per ADD
		if res16 > 0xFF {
			newFlags |= FlagC
		} // Carry
		if ((valA & 0x0F) + (valOperand & 0x0F)) > 0x0F {
			newFlags |= FlagH
		} // Half Carry
		// Overflow: (segno op1 == segno op2) && (segno risultato != segno op1)
		if ((valA ^ result) & (valOperand ^ result) & 0x80) != 0 {
			newFlags |= FlagPV
		}

	case objects.OperatorSub: // Usato sia per SUB che per CP
		// Calcolo del risultato e del borrow (carry)
		res16 := uint16(valA) - uint16(valOperand)
		result = uint8(res16)
		newFlags |= FlagN // N è sempre impostato per le sottrazioni

		// Calcolo dei flag per SUB
		if res16 > 0xFF {
			newFlags |= FlagC
		} // Carry (borrow)
		if (valA & 0x0F) < (valOperand & 0x0F) {
			newFlags |= FlagH
		} // Half Carry (borrow)
		// Overflow: (segno op1 != segno op2) && (segno risultato != segno op1)
		if ((valA ^ valOperand) & (valA ^ result) & 0x80) != 0 {
			newFlags |= FlagPV
		}

	case objects.OperatorAnd:
		result = valA & valOperand
		newFlags |= FlagH // AND imposta sempre Half Carry
		// Parity
		if bits.OnesCount8(result)%2 == 0 {
			newFlags |= FlagPV
		}

	case objects.OperatorOr:
		result = valA | valOperand
		// Parity
		if bits.OnesCount8(result)%2 == 0 {
			newFlags |= FlagPV
		}

	case objects.OperatorXor:
		result = valA ^ valOperand
		// Parity
		if bits.OnesCount8(result)%2 == 0 {
			newFlags |= FlagPV
		}
	}

	// Flag comuni a tutte le operazioni ALU
	if result == 0 {
		newFlags |= FlagZ
	} // Zero Flag
	if (result & 0x80) != 0 {
		newFlags |= FlagS
	} // Sign Flag

	// 3. Aggiornamento diretto dei registri globali
	// Per CP (Compare), il risultato non viene salvato in A
	if objects.ArithmeticOperator(op) != objects.OperatorSub || regA.Name() != "CP_TEMP" {
		regA.SetValue(int64(result))
	}
	regF.SetValue(int64(newFlags))

	return 0, gk.UndefinedValue(), nil
}

*/

/*
// Esempio nel transpiler Z80
// ...
// Prepara i parametri per la chiamata a z80.push16
h.scopes.SymbolEmitGet(symbolDellIndirizzoDiRitorno) // Arg 0: valueToPush
h.scopes.SymbolEmitGet(simboloDelRegistroSP)        // Arg 1: sp
h.scopes.SymbolEmitGet(simboloDellaMemoria)         // Arg 2: memoryArray

// Chiama la funzione della libreria
h.imports.SymbolEmit("z80", "push16")
// ...

*/

/*



func (z *Z80Package) push16(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	// 1. Validazione degli argomenti
	if len(args) != 3 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}

	valueToPush, ok := gk.ToInt64(args[0])
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(0, "int", args[0].TypeName())
	}

	sp, ok := args[1].(*objects.Int)
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(1, "int", args[1].TypeName())
	}

	memory, ok := args[2].(*objects.Array)
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(2, "array", args[2].TypeName())
	}

	// 2. Logica di emulazione del PUSH
	currentSP := sp.AsInt64()

	// Push del byte HI
	currentSP--
	hiByte := (valueToPush >> 8) & 0xFF
	memory.SetValue(int(currentSP), gk.NewInt(frame, hiByte))

	// Push del byte LO
	currentSP--
	loByte := valueToPush & 0xFF
	memory.SetValue(int(currentSP), gk.NewInt(frame, loByte))

	// 3. Aggiornamento diretto del registro SP globale
	sp.SetValue(currentSP)

	// 4. La funzione non ha un valore di ritorno significativo,
	//    poiché ha modificato lo stato globale.
	return 0, gk.UndefinedValue(), nil
}

// pop16 emula l'istruzione POP di uno Z80 per un valore a 16 bit.
// Modifica direttamente SP e restituisce il valore letto.
// Argomenti attesi:
// args[0]: IObject (Int) -> Il registro SP (che verrà modificato).
// args[1]: IObject (Array) -> L'array della MEMORY da cui leggere.
func (z *Z80Package) pop16(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	// 1. Validazione degli argomenti
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}

	sp, ok := args[0].(*objects.Int)
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(0, "int", args[0].TypeName())
	}

	memory, ok := args[1].(*objects.Array)
	if !ok {
		return 0, nil, objects.NewInvalidArgumentError(1, "array", args[1].TypeName())
	}

	// 2. Logica di emulazione del POP
	currentSP := sp.AsInt64()

	// Pop del byte LO
	loByteObj, err := memory.index(int(currentSP))
	if err != nil {
		return 0, nil, err
	}
	currentSP++

	// Pop del byte HI
	hiByteObj, err := memory.index(int(currentSP))
	if err != nil {
		return 0, nil, err
	}
	currentSP++

	// 3. Aggiornamento diretto del registro SP globale
	sp.SetValue(currentSP)

	// 4. Ricombina i byte e restituisci il valore a 16 bit.
	reconstructedValue := (hiByteObj.AsInt64() << 8) | loByteObj.AsInt64()

	return 1, gk.NewInt(frame, reconstructedValue), nil
}

*/

// Firma della funzione SDK alu
// args[0]: operatore
// args[1]: *objects.Array -> lo stato della CPU
// args[2]: *objects.Array -> lo stato della MEMORY (se necessario)
// args[3]: *objects.Int -> l'operando
//func (z *Z80Package) alu(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
//   cpu, _ := args[1].(*objects.Array)
//
//    // Accesso ultra-veloce tramite indice!
//    valA_obj, _ := cpu.index(A_IDX)
//    valF_obj, _ := cpu.index(F_IDX)
//
//    // ... logica ...

//    // Scrittura ultra-veloce tramite indice!
//    cpu.SetValue(A_IDX, /* nuovo valore A */)
//cpu.SetValue(F_IDX, /* nuovo valore F */ )
//return 0, gk.UndefinedValue(), nil
//}
