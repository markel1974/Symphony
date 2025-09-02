package tests

import "fmt"

type Attacker interface {
	Attack() string
}

type Knight struct {
	Power int
}

func (k Knight) Attack() string {
	return fmt.Sprintf("Knight attacks with power %d", k.Power)
}

type Mage struct {
	Spell string
}

func (m Mage) Attack() string {
	return fmt.Sprintf("Mage casts %s", m.Spell)
}

// Questa funzione usa un type switch per determinare l'azione.
func performAttack(a Attacker) string {
	switch v := a.(type) {
	case Knight:
		// 'v' qui è di tipo Knight
		return v.Attack()
	case Mage:
		// 'v' qui è di tipo Mage
		return v.Attack()
	default:
		return "Unknown attacker"
	}
}

func main() {
	fmt.Println("--- Running Test: Type Switch ---")

	k := Knight{Power: 15}
	m := Mage{Spell: "Fireball"}

	// Eseguiamo il type switch con entrambi i tipi
	result1 := performAttack(k)
	result2 := performAttack(m)

	finalValue := result1 + " | " + result2
	expectedValue := "Knight attacks with power 15 | Mage casts Fireball"

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Type switch statement worked correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Mismatch in type switch logic.\nGot: %s\nExpected: %s\n", finalValue, expectedValue)
	}
}
