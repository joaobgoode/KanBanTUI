package main

// Provides the mock data to fill the kanban board

func (b *Board) initColums() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	b.cols[todo].list.Title = "To Do"
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[done].list.Title = "Done"
}

func (b *Board) clearsList() {
	for col := range b.cols {
		for i := len(b.cols[col].list.Items()) - 1; i >= 0; i-- {
			b.cols[col].list.RemoveItem(i)
		}
	}
}

func (b *Board) fillLists() {
	// get todos from the db
	todos, err := taskByStatus(todo)
	if err != nil {
		panic(err)
	}
	// sort.Sort(todos)
	// fill the list with the todos
	for _, t := range todos {
		b.cols[todo].list.InsertItem(-1, t)
	}
	// get in progress tasks from the db
	inprogresses, err := taskByStatus(inProgress)
	if err != nil {
		panic(err)
	}
	// sort.Sort(inprogresses)
	// fill the list with the in progress tasks
	for _, t := range inprogresses {
		b.cols[inProgress].list.InsertItem(-1, t)
	}
	// get done tasks from the db
	dones, err := taskByStatus(done)
	if err != nil {
		panic(err)
	}
	// sort.Sort(dones)
	// fill the list with the done tasks
	for _, t := range dones {
		b.cols[done].list.InsertItem(-1, t)
	}
}

func (b *Board) resetLists() {
	b.clearsList()
	b.fillLists()
}

func (b *Board) initLists() {
	b.initColums()
	b.fillLists()
}
