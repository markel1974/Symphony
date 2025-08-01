package process

type commandSorterByName []*Command

// Len returns the number of elements in the commandSorterByName collection.
func (c commandSorterByName) Len() int { return len(c) }

// Swap exchanges the elements at indices i and j in the receiver.
func (c commandSorterByName) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less compares the names of two Command elements at indexes i and j and returns true if the first is less than the second.
func (c commandSorterByName) Less(i, j int) bool { return c[i].Name() < c[j].Name() }
