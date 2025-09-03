package sources

import "fmt"

const static = 10

// --- Definizione dei Tipi (Structs) ---
// Il compilatore supporta la definizione di tipi struct.
//

type SimpleTask struct {
	Title     string
	Completed bool
}

type Project struct {
	ProjectName string
	Tasks       []SimpleTask
	IsArchived  bool
}

// --- Metodi associati ai Tipi ---
// Il compilatore gestisce correttamente i metodi con ricevitori (receiver).
// In questo caso, usiamo un puntatore per modificare lo stato interno.
//

func (t *SimpleTask) MarkAsComplete() {
	t.Completed = true
}

func (t *SimpleTask) Display() {
	status := "[ ]"
	if t.Completed {
		status = "[x]"
	}
	fmt.Println(status, "Task:", t.Title)
}

func (p *Project) Display() {
	status := "Active"
	if p.IsArchived {
		status = "Archived"
	}
	fmt.Println("--- Project:", p.ProjectName, "(", status, ") ---")
	// Il compilatore supporta i cicli 'range' su slice/array.
	//
	for _, taskIt := range p.Tasks {
		taskIt.Display()
	}
	fmt.Println("--------------------")
}

// --- Polimorfismo (Implicito) ---
// Questa funzione accetta un array di interfacce (implicite) che hanno
// il metodo Display(). Il sistema a oggetti della VM lo permette.

func ProcessItems(items []interface{}) {
	fmt.Println("\n=> Processing all items...")
	for _, item := range items {
		fmt.Println("Item:", item)
		// La VM risolverà dinamicamente la chiamata al metodo corretto.
		//item.Display()
		//TODO
	}
}

// --- Closure e Funzioni di Ordine Superiore ---
// Questa funzione restituisce un'altra funzione (una closure).
// La closure "cattura" la variabile 'status' dal suo ambiente esterno.
// Il compilatore e la VM gestiscono correttamente le "free variables".
//

func TaskFilterGenerator7(status bool) func(SimpleTask) bool {
	return func(taskClosure SimpleTask) bool {
		return taskClosure.Completed == status
	}
}

// --- Funzione con Multi-Valore di Ritorno ---
// Il compilatore supporta funzioni che restituiscono valori multipli.
func GetSystemStatus() (string, int) {
	return "Operational", 200
}

func main() {
	// --- Creazione e Assegnazione ---
	// Il compilatore gestisce l'inizializzazione di struct (CompositeLit)
	// e l'assegnazione con ':='.
	//
	//task1 := &SimpleTask{Title: "Buy milk", Completed: false}
	task1 := SimpleTask{Title: "Buy milk", Completed: false}
	task2 := SimpleTask{Title: "Write compiler example", Completed: true}

	// --- Modifica tramite Puntatore ---
	// La chiamata al metodo modifica l'oggetto originale 'task1'.
	task1.MarkAsComplete()

	project := &Project{
		ProjectName: "VM Development",
		Tasks: []SimpleTask{
			{Title: "Implement bytecode parsing", Completed: true},
			{Title: "Design object system", Completed: true},
			{Title: "Write complex example", Completed: false},
		},
		IsArchived: false,
	}

	// Il compilatore gestisce l'assegnazione a più variabili.
	//
	status, code := GetSystemStatus()
	fmt.Println("System Status:", status, "- Code:", code)

	// --- Esecuzione Polimorfica ---
	// Creiamo un array contenente tipi diversi.
	allItems := []interface{}{task1, &task2, project}
	ProcessItems(allItems)

	// --- Utilizzo della Closure ---
	fmt.Println("\n=> Filtering for incomplete tasks...")

	// Creiamo un filtro specifico per i task non completati.
	//filterIncomplete := TaskFilterGenerator7(false)

	//allTasks := []SimpleTask{ *task1, task2, project.Tasks[0], project.Tasks[1], project.Tasks[2] }
	//allTasks := []SimpleTask{ task1, task2, project.Tasks[0], project.Tasks[1], project.Tasks[2] }
	allTasks := []SimpleTask{task1, task2, task1}
	for _, taskT := range allTasks {
		// Applichiamo il filtro.
		//if filterIncomplete(taskT) {
		taskT.Display()
		//}
	}
}
