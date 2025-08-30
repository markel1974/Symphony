package stub

import "fmt"

type Printer interface {
	Print()
}

type User struct {
	Name string
}

type Article struct {
	Title string
}

func (u User) Print() {
	fmt.Println("User:", u.Name)
}

func (a Article) Print() {
	fmt.Println("Article:", a.Title)
}

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
