package utils

import (
	"github.com/gotk3/gotk3/gtk"
)

type PairStore struct {
	variables     []*KeyValuePair
	container     *gtk.Grid
	onAfterUpdate func()
}

func newPairStore() *PairStore {
	return &PairStore{
		variables: make([]*KeyValuePair, 0),
	}
}

func (s *PairStore) afterUpdate() {
	if s.onAfterUpdate != nil {
		s.onAfterUpdate()
	}
}

func (s *PairStore) Delete(index int) {
	s.variables = append(s.variables[:index], s.variables[index+1:]...)

	ClearChildren(s.container.Container)
	showVariables(s, s.container)
	s.container.ShowAll()

	s.afterUpdate()
}

func (s *PairStore) addBeforeStart(key, value string) {
	s.variables = append(s.variables, &KeyValuePair{
		Key:   key,
		Value: value,
	})
}

func (s *PairStore) Add(key, value string) {
	s.variables = append(s.variables, &KeyValuePair{
		Key:   key,
		Value: value,
	})
	ClearChildren(s.container.Container)
	showVariables(s, s.container)
	s.container.ShowAll()

	s.afterUpdate()
}

func (s *PairStore) Get(key string) string {
	for _, variable := range s.variables {
		if variable.Key == key {
			return variable.Value
		}
	}
	return ""
}

func (s *PairStore) List() []*KeyValuePair {
	return s.variables
}

func (s *PairStore) ToMap() map[string]string {
	variables := make(map[string]string)
	for _, variable := range s.variables {
		variables[variable.Key] = variable.Value
	}
	return variables
}

type KeyValuePair struct {
	Key   string
	Value string
}

type KeyValueEditor struct {
	*gtk.ScrolledWindow
	Store *PairStore
}

func (kve *KeyValueEditor) OnAfterUpdate(handler func()) {
	kve.Store.onAfterUpdate = handler
}

func NewKeyValueEditor(values map[string]string) *KeyValueEditor {
	container, err := gtk.GridNew()
	HandleErr(err)

	container.SetHExpand(true)

	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	container.SetRowSpacing(5)
	container.SetColumnSpacing(5)
	container.SetColumnHomogeneous(false)

	store := newPairStore()

	for key, value := range values {
		store.addBeforeStart(key, value)
	}

	store.container = container
	showVariables(store, container)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	HandleErr(err)

	scrolledWindow.Add(container)

	return &KeyValueEditor{
		scrolledWindow,
		store,
	}
}

func showVariables(store *PairStore, container *gtk.Grid) {
	lastRow := 0
	for i, variable := range store.List() {
		variableCopy := variable
		idxCopy := i

		keyEntry, valueEntry := newVariableEntry(store, variableCopy)

		container.Attach(keyEntry, 0, i, 3, 1)
		container.Attach(valueEntry, 3, i, 3, 1)

		deleteBtn := Must(
			gtk.ButtonNewFromIconName("user-trash-symbolic", gtk.ICON_SIZE_BUTTON),
		)
		deleteBtn.Connect("clicked", func() {
			store.Delete(idxCopy)
		})
		container.Attach(deleteBtn, 6, i, 1, 1)

		lastRow = i
	}

	addBtn, err := gtk.ButtonNewWithLabel("Add")
	HandleErr(err)

	addBtn.Connect("clicked", func() {
		store.Add("", "")
	})

	container.Attach(addBtn, 0, lastRow+1, 2, 1)
}

func newVariableEntry(store *PairStore, variable *KeyValuePair) (*gtk.Entry, *gtk.Entry) {
	keyEntry := Must(gtk.EntryNew())

	keyEntry.SetPlaceholderText("Key")
	keyEntry.SetHExpand(true)

	valueEntry := Must(gtk.EntryNew())

	valueEntry.SetPlaceholderText("Value")
	valueEntry.SetHExpand(true)

	keyEntry.SetText(variable.Key)
	valueEntry.SetText(variable.Value)

	keyEntry.Connect("changed", func() {
		var err error
		variable.Key, err = keyEntry.GetText()
		HandleErr(err)
		store.afterUpdate()
	})

	valueEntry.Connect("changed", func() {
		var err error
		variable.Value, err = valueEntry.GetText()
		HandleErr(err)
		store.afterUpdate()
	})

	return keyEntry, valueEntry
}
