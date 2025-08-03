package ui

import (
	"gophkeeper/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_startManageMetadata_NilDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: nil,
		},
	}

	result, cmd := ui.startManageMetadata()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
}

func TestUIController_startManageMetadata_WithMetadata(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
		},
	}

	result, cmd := ui.startManageMetadata()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMetadataList, ui.state)
	assert.Equal(t, 0, ui.itemCtrl.itemMetaCtrl.metadataMenu)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.currentMetaValue)
	assert.Len(t, ui.itemCtrl.itemMetaCtrl.metadataKeys, 2)
	assert.Contains(t, ui.itemCtrl.itemMetaCtrl.metadataKeys, "key1")
	assert.Contains(t, ui.itemCtrl.itemMetaCtrl.metadataKeys, "key2")
}

func TestUIController_startManageMetadata_EmptyMetadata(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{},
				},
			},
		},
	}

	result, cmd := ui.startManageMetadata()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMetadataList, ui.state)
	assert.Len(t, ui.itemCtrl.itemMetaCtrl.metadataKeys, 0)
}

func TestUIController_handleMetadataListInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleMetadataListInput_Back(t *testing.T) {
	ui := &UIController{
		state: stateMetadataList,
	}

	// Test esc
	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)

	// Test b
	ui.state = stateMetadataList
	model, cmd = ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
}

func TestUIController_handleMetadataListInput_Navigation(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 1,
				metadataKeys: []string{"key1", "key2"},
			},
		},
	}

	// Test up
	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemCtrl.itemMetaCtrl.metadataMenu)

	// Test k
	ui.itemCtrl.itemMetaCtrl.metadataMenu = 1
	model, cmd = ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemCtrl.itemMetaCtrl.metadataMenu)

	// Test down
	model, cmd = ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.itemCtrl.itemMetaCtrl.metadataMenu)

	// Test j
	ui.itemCtrl.itemMetaCtrl.metadataMenu = 0
	model, cmd = ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.itemCtrl.itemMetaCtrl.metadataMenu)
}

func TestUIController_handleMetadataListInput_EnterEditExisting(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{"key1": "value1"},
				},
			},
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 0,
				metadataKeys: []string{"key1"},
			},
		},
	}

	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateEditMetadataValue, ui.state)
	assert.Equal(t, "key1", ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Equal(t, "value1", ui.input)
}

func TestUIController_handleMetadataListInput_EnterAddNew(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 1, // Index equals length (add option)
				metadataKeys: []string{"key1"},
			},
		},
		input: "some-input",
	}

	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddMetadataKey, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleMetadataListInput_Delete(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 0,
				metadataKeys: []string{"key1"},
			},
		},
	}

	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateConfirmDeleteMetadata, ui.state)
	assert.Equal(t, "key1", ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Equal(t, 0, ui.confirmChoice)
}

func TestUIController_handleMetadataListInput_DeleteInvalidIndex(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 5, // Index out of bounds
				metadataKeys: []string{"key1"},
			},
		},
		state: stateMetadataList,
	}

	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMetadataList, ui.state) // Should remain unchanged
}

func TestUIController_handleMetadataListInput_Add(t *testing.T) {
	ui := &UIController{
		input: "some-input",
	}

	model, cmd := ui.handleMetadataListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddMetadataKey, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_startAddMetadata(t *testing.T) {
	ui := &UIController{
		input: "some-input",
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey:   "old-key",
				currentMetaValue: "old-value",
			},
		},
	}

	result, cmd := ui.startAddMetadata()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddMetadataKey, ui.state)
	assert.Empty(t, ui.input)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.currentMetaValue)
}

func TestUIController_handleAddMetadataKeyInput_Enter_ValidKey(t *testing.T) {
	ui := &UIController{
		input: "  new-key  ",
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{"existing-key": "value"},
				},
			},
		},
	}

	model, cmd := ui.handleAddMetadataKeyInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "new-key", ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Empty(t, ui.input)
	assert.Equal(t, stateAddMetadataValue, ui.state)
}

func TestUIController_handleAddMetadataKeyInput_Enter_EmptyKey(t *testing.T) {
	ui := &UIController{
		input: "   ",
		state: stateAddMetadataKey,
	}

	model, cmd := ui.handleAddMetadataKeyInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddMetadataKey, ui.state) // Should remain unchanged
}

func TestUIController_handleAddMetadataKeyInput_Enter_ExistingKey(t *testing.T) {
	ui := &UIController{
		input: "existing-key",
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{"existing-key": "value"},
				},
			},
		},
	}

	model, cmd := ui.handleAddMetadataKeyInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	// Should show error message but not change state
}

func TestUIController_handleAddMetadataValueInput_Enter_ValidValue(t *testing.T) {
	ui := &UIController{
		input: "  new-value  ",
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	model, cmd := ui.handleAddMetadataValueInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return saveMetadataCmd
	assert.Equal(t, "new-value", ui.itemCtrl.itemMetaCtrl.currentMetaValue)
}

func TestUIController_handleAddMetadataValueInput_Enter_EmptyValue(t *testing.T) {
	ui := &UIController{
		input: "   ",
		state: stateAddMetadataValue,
	}

	model, cmd := ui.handleAddMetadataValueInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddMetadataValue, ui.state) // Should remain unchanged
}

func TestUIController_startEditMetadata(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{"test-key": "test-value"},
				},
			},
		},
		input: "old-input",
	}

	result, cmd := ui.startEditMetadata("test-key")

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, "test-key", ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Equal(t, "test-value", ui.itemCtrl.itemMetaCtrl.currentMetaValue)
	assert.Equal(t, "test-value", ui.input)
	assert.Equal(t, stateEditMetadataValue, ui.state)
}

func TestUIController_handleEditMetadataValueInput_Enter_ValidValue(t *testing.T) {
	ui := &UIController{
		input: "  updated-value  ",
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	model, cmd := ui.handleEditMetadataValueInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return saveMetadataCmd
}

func TestUIController_startDeleteMetadata(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
	}

	result, cmd := ui.startDeleteMetadata("test-key")

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, "test-key", ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Equal(t, stateConfirmDeleteMetadata, ui.state)
	assert.Equal(t, 0, ui.confirmChoice) // Should be reset to 0
}

func TestUIController_handleConfirmDeleteMetadataInput_Yes(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	model, cmd := ui.handleConfirmDeleteMetadataInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return saveMetadataCmd
	assert.Equal(t, 1, ui.confirmChoice)
}

func TestUIController_handleConfirmDeleteMetadataInput_No(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
		state:         stateConfirmDeleteMetadata,
	}

	model, cmd := ui.handleConfirmDeleteMetadataInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)
	assert.Equal(t, stateMetadataList, ui.state)
}

func TestUIController_handleConfirmDeleteMetadata_Confirm(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	result, cmd := ui.handleConfirmDeleteMetadata()

	assert.Equal(t, ui, result)
	assert.NotNil(t, cmd) // Should return saveMetadataCmd
}

func TestUIController_handleConfirmDeleteMetadata_Cancel(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		state:         stateConfirmDeleteMetadata,
	}

	result, cmd := ui.handleConfirmDeleteMetadata()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMetadataList, ui.state)
}

func TestUIController_saveMetadataCmd(t *testing.T) {
	// This test is complex for unit testing as it requires mock Item service
	// Skipping as this is an integration test
	t.Skip("Requires mock Item service - integration test")
}

func TestUIController_handleMetadataSuccessInput_Enter(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Meta: models.Meta{
					Map: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			itemMetaCtrl: itemMetaCtrl{
				metadataSuccessMsg: "test message",
				metadataMenu:       5,
			},
		},
	}

	model, cmd := ui.handleMetadataSuccessInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMetadataList, ui.state)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.metadataSuccessMsg)
	assert.Equal(t, 0, ui.itemCtrl.itemMetaCtrl.metadataMenu)
	assert.Len(t, ui.itemCtrl.itemMetaCtrl.metadataKeys, 2)
}

func TestUIController_handleMetadataErrorInput_Enter(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataErrorMsg: "test error",
			},
		},
	}

	model, cmd := ui.handleMetadataErrorInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMetadataList, ui.state)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.metadataErrorMsg)
}

// View tests
func TestUIController_metadataListView_NilDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: nil,
		},
	}

	view := ui.metadataListView()

	assert.Equal(t, "No item selected", view)
}

func TestUIController_metadataListView_WithMetadata(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-item",
				Meta: models.Meta{
					Map: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 0,
				metadataKeys: []string{"key1", "key2"},
			},
		},
	}

	view := ui.metadataListView()

	assert.Contains(t, view, "Manage Metadata: test-item")
	assert.Contains(t, view, "Current metadata:")
	assert.Contains(t, view, "key1: value1")
	assert.Contains(t, view, "key2: value2")
	assert.Contains(t, view, "+ Add new metadata")
	assert.Contains(t, view, "Enter to edit/add")
}

func TestUIController_metadataListView_EmptyMetadata(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-item",
			},
			itemMetaCtrl: itemMetaCtrl{
				metadataMenu: 0,
				metadataKeys: []string{},
			},
		},
	}

	view := ui.metadataListView()

	assert.Contains(t, view, "Manage Metadata: test-item")
	assert.Contains(t, view, "No metadata found")
	assert.Contains(t, view, "+ Add new metadata")
}

func TestUIController_addMetadataKeyView(t *testing.T) {
	ui := &UIController{
		input: "test-key",
	}

	view := ui.addMetadataKeyView()

	assert.Contains(t, view, "Add Metadata - Enter Key")
	assert.Contains(t, view, "test-key")
	assert.Contains(t, view, "█")
	assert.Contains(t, view, "Esc to go back")
}

func TestUIController_addMetadataValueView(t *testing.T) {
	ui := &UIController{
		input: "test-value",
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	view := ui.addMetadataValueView()

	assert.Contains(t, view, "Add Metadata - Enter Value")
	assert.Contains(t, view, "Key: test-key")
	assert.Contains(t, view, "test-value")
	assert.Contains(t, view, "█")
}

func TestUIController_editMetadataValueView(t *testing.T) {
	ui := &UIController{
		input: "updated-value",
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	view := ui.editMetadataValueView()

	assert.Contains(t, view, "Edit Metadata - Edit Value")
	assert.Contains(t, view, "Key: test-key")
	assert.Contains(t, view, "updated-value")
	assert.Contains(t, view, "█")
}

func TestUIController_confirmDeleteMetadataView(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				currentMetaKey: "test-key",
			},
		},
	}

	view := ui.confirmDeleteMetadataView()

	assert.Contains(t, view, "Confirm Delete Metadata")
	assert.Contains(t, view, "delete metadata key 'test-key'")
	assert.Contains(t, view, "[ No ]")
	assert.Contains(t, view, "[ Yes ]")
	assert.Contains(t, view, "y/n for quick choice")
}

func TestUIController_metadataSuccessView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataSuccessMsg: "Metadata updated successfully!",
			},
		},
	}

	view := ui.metadataSuccessView()

	assert.Contains(t, view, "Metadata Updated Successfully")
	assert.Contains(t, view, "Metadata updated successfully!")
	assert.Contains(t, view, "Enter to continue")
}

func TestUIController_metadataErrorView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemMetaCtrl: itemMetaCtrl{
				metadataErrorMsg: "Failed to save metadata",
			},
		},
	}

	view := ui.metadataErrorView()

	assert.Contains(t, view, "Metadata Error")
	assert.Contains(t, view, "Failed to save metadata")
	assert.Contains(t, view, "Enter to try again")
}
