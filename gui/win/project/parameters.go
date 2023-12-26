package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

type ParameterStore struct {
	parameters []*KeyValuePair
	container  *gtk.Grid
}

func newParametersStore() *ParameterStore {
	return &ParameterStore{
		parameters: make([]*KeyValuePair, 0),
	}
}

func (s *ParameterStore) Add(key, value string) {
	s.parameters = append(s.parameters, &KeyValuePair{
		Key:   key,
		Value: value,
	})
	s.container.GetChildren().Foreach(func(item interface{}) {
		item.(*gtk.Widget).Destroy()
	})
	showParameters(s, s.container)
	s.container.ShowAll()
}

func (s *ParameterStore) Get(key string) string {
	for _, parameter := range s.parameters {
		if parameter.Key == key {
			return parameter.Value
		}
	}
	return ""
}

func (s *ParameterStore) List() []*KeyValuePair {
	return s.parameters
}

type KeyValuePair struct {
	Key   string
	Value string
}

func newParametersContainer(store *ParameterStore) *gtk.ScrolledWindow {
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

	store.container = container
	showParameters(store, container)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleErr(err)

	scrolledWindow.Add(container)

	return scrolledWindow
}

func showParameters(store *ParameterStore, container *gtk.Grid) {
	lastRow := 0
	for i, parameter := range store.List() {
		keyEntry, err := gtk.EntryNew()
		utils.HandleErr(err)

		parameterCopy := parameter

		keyEntry.SetPlaceholderText("Key")
		keyEntry.SetHExpand(true)

		valueEntry, err := gtk.EntryNew()
		utils.HandleErr(err)

		valueEntry.SetPlaceholderText("Value")
		valueEntry.SetHExpand(true)

		keyEntry.SetText(parameter.Key)
		valueEntry.SetText(parameter.Value)

		keyEntry.Connect("changed", func() {
			var err error
			parameterCopy.Key, err = keyEntry.GetText()
			utils.HandleErr(err)
		})

		valueEntry.Connect("changed", func() {
			var err error
			parameterCopy.Value, err = valueEntry.GetText()
			utils.HandleErr(err)
		})

		container.Attach(keyEntry, 0, i, 1, 1)
		container.Attach(valueEntry, 1, i, 1, 1)

		lastRow = i
	}

	addBtn, err := gtk.ButtonNewWithLabel("Add")
	utils.HandleErr(err)

	addBtn.Connect("clicked", func() {
		store.Add("", "")
	})

	container.Attach(addBtn, 0, lastRow+1, 2, 1)
}
