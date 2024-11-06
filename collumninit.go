package main

// Provides the mock data to fill the kanban board

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	// Init To Do
	todos, err := taskByStatus(todo)
	if err != nil {
		panic(err)
	}
	b.cols[todo].list.Title = "To Do"
	for _, t := range todos {
		b.cols[todo].list.InsertItem(len(b.cols[todo].list.Items())-1, t)
	}
	// Init in progress
	inprogresses, err := taskByStatus(inProgress)
	if err != nil {
		panic(err)
	}
	b.cols[inProgress].list.Title = "In Progress"
	for _, t := range inprogresses {
		b.cols[inProgress].list.InsertItem(len(b.cols[inProgress].list.Items())-1, t)
	}
	// Init done
	dones, err := taskByStatus(done)
	if err != nil {
		panic(err)
	}
	b.cols[done].list.Title = "Done"
	for _, t := range dones {
		b.cols[done].list.InsertItem(len(b.cols[done].list.Items())-1, t)
	}
}

func (b *Board) resetLists() {
	for col := range b.cols {
		for i := len(b.cols[col].list.Items()) - 1; i >= 0; i-- {
			b.cols[col].list.RemoveItem(i)
		}
	}

	todos, err := taskByStatus(todo)
	if err != nil {
		panic(err)
	}
	for _, t := range todos {
		b.cols[todo].list.InsertItem(len(b.cols[todo].list.Items())-1, t)
	}

	inprogresses, err := taskByStatus(inProgress)
	if err != nil {
		panic(err)
	}
	for _, t := range inprogresses {
		b.cols[inProgress].list.InsertItem(len(b.cols[inProgress].list.Items())-1, t)
	}

	dones, err := taskByStatus(done)
	if err != nil {
		panic(err)
	}
	for _, t := range dones {
		b.cols[done].list.InsertItem(len(b.cols[done].list.Items())-1, t)
	}
}
