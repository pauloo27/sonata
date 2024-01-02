package project

import (
	"github.com/pauloo27/sonata/gui/utils"
)

func newVariablesContainer(store *ProjectStore) *utils.KeyValueEditor {
	kvEditor := utils.NewKeyValueEditor(store.DraftRequest.Variables)
	store.VarStore = kvEditor.Store

	return kvEditor
}
