package pla

type TriggerData struct {
	id   int
	addr uint16
	w    WriteFn
}

func NewTriggerData(id int, addr uint16, fn WriteFn) *TriggerData {
	return &TriggerData{
		id:   id,
		addr: addr,
		w:    fn,
	}
}

func (td *TriggerData) Exec(data uint8) {
	td.w(td.addr, data)
}

func (td *TriggerData) GetId() int {
	return td.id
}

type Trigger struct {
	container []*TriggerData
	addr      uint16
	idx       int
}

func NewTrigger(addr uint16) *Trigger {
	return &Trigger{
		idx:       0,
		addr:      addr,
		container: nil,
	}
}

func (wt *Trigger) Add(fn WriteFn) int {
	id := wt.idx
	wt.container = append(wt.container, NewTriggerData(id, wt.addr, fn))
	wt.idx++
	return id
}

func (wt *Trigger) Remove(id int) {
	for idx, f := range wt.container {
		if id == f.GetId() {
			wt.container = append(wt.container[:idx], wt.container[idx+1:]...)
			break
		}
	}
}

func (wt *Trigger) Exec(data uint8) {
	for _, f := range wt.container {
		f.Exec(data)
	}
}

type WriteTriggers struct {
	triggers []*Trigger
}

func NewWriteTriggers(r int) *WriteTriggers {
	wt := &WriteTriggers{triggers: nil}
	wt.triggers = make([]*Trigger, r)
	return wt
}

func (wt *WriteTriggers) Add(addr uint16, fn WriteFn) int {
	t := wt.triggers[addr]
	if t == nil {
		t = NewTrigger(addr)
		wt.triggers[addr] = t
	}
	return t.Add(fn)
}

func (wt *WriteTriggers) Remove(addr uint16, id int) {
	if t := wt.triggers[addr]; t != nil {
		t.Remove(id)
	}
}

func (wt *WriteTriggers) Exec(addr uint16, data uint8) {
	if t := wt.triggers[addr]; t != nil {
		t.Exec(data)
	}
}
