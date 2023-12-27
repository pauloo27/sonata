package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

type VariablesStore struct {
	variables []*KeyValuePair
	container *gtk.Grid
}

func newVariablesStore() *VariablesStore {
	return &VariablesStore{
		variables: make([]*KeyValuePair, 0),
	}
}

func (s *VariablesStore) Add(key, value string) {
	s.variables = append(s.variables, &KeyValuePair{
		Key:   key,
		Value: value,
	})
	utils.ClearChildren(s.container.Container)
	showVariables(s, s.container)
	s.container.ShowAll()
}

func (s *VariablesStore) Get(key string) string {
	for _, variable := range s.variables {
		if variable.Key == key {
			return variable.Value
		}
	}
	return ""
}

func (s *VariablesStore) List() []*KeyValuePair {
	return s.variables
}

type KeyValuePair struct {
	Key   string
	Value string
}

func newVariablesContainer(store *ContentStore) *gtk.ScrolledWindow {
	container, err := gtk.GridNew()
	utils.HandleErr(err)

	container.SetHExpand(true)

	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	container.SetRowSpacing(5)
	container.SetColumnSpacing(5)
	container.SetColumnHomogeneous(true)

	store.VarStore.container = container
	showVariables(store.VarStore, container)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleErr(err)

	scrolledWindow.Add(container)

	return scrolledWindow
}

func showVariables(varStore *VariablesStore, container *gtk.Grid) {
	lastRow := 0
	for i, variable := range varStore.List() {
		keyEntry, err := gtk.EntryNew()
		utils.HandleErr(err)

		variableCopy := variable

		keyEntry.SetPlaceholderText("Key")
		keyEntry.SetHExpand(true)

		valueEntry, err := gtk.EntryNew()
		utils.HandleErr(err)

		valueEntry.SetPlaceholderText("Value")
		valueEntry.SetHExpand(true)

		keyEntry.SetText(variable.Key)
		valueEntry.SetText(variable.Value)

		keyEntry.Connect("changed", func() {
			var err error
			variableCopy.Key, err = keyEntry.GetText()
			utils.HandleErr(err)
		})

		valueEntry.Connect("changed", func() {
			var err error
			variableCopy.Value, err = valueEntry.GetText()
			utils.HandleErr(err)
		})

		container.Attach(keyEntry, 0, i, 1, 1)
		container.Attach(valueEntry, 1, i, 1, 1)

		lastRow = i
	}

	addBtn, err := gtk.ButtonNewWithLabel("Add")
	utils.HandleErr(err)

	addBtn.Connect("clicked", func() {
		varStore.Add("", "")
	})

	container.Attach(addBtn, 0, lastRow+1, 2, 1)
}
