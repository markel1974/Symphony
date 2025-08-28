package stub

import "fmt"

// 1. Definizione dell'interfaccia
type Printer interface {
	Print()
}

// 2. Definizione di due struct
type User struct {
	Name string
}

type Article struct {
	Title string
}

// 3. Implementazione dell'interfaccia per User
func (u User) Print() {
	fmt.Println("User:", u.Name)
}

// 4. Implementazione dell'interfaccia per Article
func (a Article) Print() {
	fmt.Println("Article:", a.Title)
}

// Funzione che accetta l'interfaccia
func DoPrint(p Printer) {
	// 6. Chiamata polimorfica
	p.Print()
}

func main() {
	u := User{Name: "Mario"}
	a := Article{Title: "Interfaces in Go"}
	var p1 Printer
	p1 = u
	var p2 Printer = a
	DoPrint(p1) // Dovrebbe stampare "User: Mario"
	DoPrint(p2) // Dovrebbe stampare "Article: Interfaces in Go"
}
