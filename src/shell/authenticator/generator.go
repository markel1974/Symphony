package authenticator

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// Package password provides a library for generating high-entropy random
// password strings via the crypto/rand package.
//
//	res, err := Generate(64, 10, 10, false, false)
//	if err != nil  {
//	  log.Fatal(err)
//	}
//	log.Printf(res)
//
// Most functions are safe for concurrent use.

// _ is a compile-time assertion ensuring Generator implements the PasswordGenerator interface.
var _ PasswordGenerator = (*Generator)(nil)

// PasswordGenerator defines methods for generating secure passwords with specific constraints.
type PasswordGenerator interface {
	Generate(int) (string, error)
}

// LowerLetters is the string containing all lowercase English letters.
// UpperLetters is the string containing all uppercase English letters.
// Digits is the string containing numeric digits 0 through 9.
// Symbols is the string containing various special characters.
const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// Symbols is the list of symbols.
	Symbols = "~!@#$%^&*()_+-={}|[]:<>?,./"
)

// Generator is a struct used for generating random strings with customizable character sets and rules.
type Generator struct {
	lowerLetters string
	upperLetters string
	digits       string
	symbols      string
	reader       io.Reader
}

// GeneratorInput represents the input configuration for creating a Generator.
// LowerLetters defines the set of lowercase letters available for password generation.
// UpperLetters defines the set of uppercase letters available for password generation.
// Digits defines the set of numeric characters available for password generation.
// Symbols defines the set of special characters available for password generation.
// Reader specifies a source of random data, with rand.Reader used by default if nil.
type GeneratorInput struct {
	LowerLetters string
	UpperLetters string
	Digits       string
	Symbols      string
	Reader       io.Reader // rand.Reader by default
}

// NewGenerator creates and initializes a new Generator with the specified GeneratorInput or default values if nil is provided.
// Returns the initialized Generator or an error if any occurs during the setup.
func NewGenerator(i *GeneratorInput) (*Generator, error) {
	if i == nil {
		i = new(GeneratorInput)
	}

	g := &Generator{
		lowerLetters: i.LowerLetters,
		upperLetters: i.UpperLetters,
		digits:       i.Digits,
		symbols:      i.Symbols,
		reader:       i.Reader,
	}

	if g.lowerLetters == "" {
		g.lowerLetters = LowerLetters
	}

	if g.upperLetters == "" {
		g.upperLetters = UpperLetters
	}

	if g.digits == "" {
		g.digits = Digits
	}

	if g.symbols == "" {
		g.symbols = Symbols
	}

	if g.reader == nil {
		g.reader = rand.Reader
	}

	return g, nil
}

// Generate generates a random string based on the specified parameters for length, digits, symbols, case, and repetition.
func (g *Generator) Generate(length int) (string, error) {
	container := g.lowerLetters
	container += g.upperLetters
	container += g.digits
	container += g.symbols
	const minChar = 8
	chars := length
	if chars < minChar {
		chars = minChar
	}
	var result string
	for i := 0; i < chars; i++ {
		ch, err := randomElement(g.reader, container)
		if err != nil {
			return "", err
		}

		result, err = randomInsert(g.reader, result, ch)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

// MustGenerate generates a random string based on the provided parameters and panics if an error occurs during generation.
func (g *Generator) MustGenerate(length int) string {
	res, err := g.Generate(length)
	if err != nil {
		panic(err)
	}
	return res
}

// Generate creates a random string of specified length containing a mix of letters, digits, and symbols based on given constraints.
func Generate(length int) (string, error) {
	gen, err := NewGenerator(nil)
	if err != nil {
		return "", err
	}

	return gen.Generate(length)
}

// randomInsert inserts the string `val` at a random position within the string `s` using the provided random source `reader`.
// If `s` is empty, it returns `val` as the result.
// Returns the resulting string and an error if the random number generation fails.
func randomInsert(reader io.Reader, s, val string) (string, error) {
	if s == "" {
		return val, nil
	}

	n, err := rand.Int(reader, big.NewInt(int64(len(s)+1)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random integer: %w", err)
	}
	i := n.Int64()
	return s[0:i] + val + s[i:], nil
}

// randomElement selects a random character from the given string using a secure random number generator.
// The io.Reader provides the source of randomness, and the function returns the selected character or an error.
func randomElement(reader io.Reader, s string) (string, error) {
	n, err := rand.Int(reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random integer: %w", err)
	}
	return string(s[n.Int64()]), nil
}
