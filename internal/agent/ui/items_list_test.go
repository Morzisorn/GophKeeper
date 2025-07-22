package ui

import (
	"fmt"
	"gophkeeper/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_handleViewItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 5,
		},
		state: stateMenuLoggedIn,
	}

	result, cmd := ui.handleViewAllItems()

	assert.Equal(t, ui, result)
	assert.NotNil(t, cmd) // Should return loadItemsCmd
	assert.Equal(t, stateItemsList, ui.state)
	assert.Equal(t, 0, ui.itemCtrl.currentItem) // Should be reset
	// Should set loading message
	assert.Equal(t, "Loading items...", ui.messages.Get("info"))
}

func TestUIController_loadItemsCmd(t *testing.T) {
	// Этот тест сложен для unit-тестирования, так как требует mock Item service
	// Пропускаем, так как это интеграционный тест
	t.Skip("Requires mock Item service - integration test")
}

func TestUIController_handleItemsListInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleItemsListInput_Escape(t *testing.T) {
	ui := &UIController{
		state:        stateItemsList,
		loggedInMenu: 5,
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleItemsListInput_Navigation_Up(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 2,
		},
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyUp})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.itemCtrl.currentItem)

	// Test k key
	model, cmd = ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemCtrl.currentItem)

	// Test that it doesn't go below 0
	model, cmd = ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyUp})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemCtrl.currentItem)
}

func TestUIController_handleItemsListInput_Navigation_Down(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{Name: "item1"},
				{Name: "item2"},
				{Name: "item3"},
			},
		},
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.itemCtrl.currentItem)

	// Test j key
	model, cmd = ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 2, ui.itemCtrl.currentItem)

	// Test that it doesn't go beyond last item
	model, cmd = ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 2, ui.itemCtrl.currentItem) // Should remain at last index
}

func TestUIController_handleItemsListInput_Navigation_Down_EmptyItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items:       []models.EncryptedItem{},
		},
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemCtrl.currentItem) // Should remain unchanged
}

func TestUIController_handleItemsListInput_Enter_ValidItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					ID:   [16]byte{1, 2, 3},
					Name: "test-item",
				},
			},
		},
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return result from handleViewItemDetails
	assert.Equal(t, stateItemDetails, ui.state)
	assert.NotNil(t, ui.itemCtrl.selectedItem)
}

func TestUIController_handleItemsListInput_Enter_EmptyItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items:       []models.EncryptedItem{},
		},
		state: stateItemsList,
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemsList, ui.state) // Should remain unchanged
}

func TestUIController_handleItemsListInput_Enter_InvalidIndex(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 5, // Index out of bounds
			items: []models.EncryptedItem{
				{Name: "item1"},
			},
		},
		state: stateItemsList,
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemsList, ui.state) // Should remain unchanged
}

func TestUIController_handleItemsListInput_Refresh(t *testing.T) {
	ui := &UIController{
		state: stateItemsList,
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return loadItemsCmd
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handleItemsListInput_OtherKey(t *testing.T) {
	ui := &UIController{
		state: stateItemsList,
		itemCtrl: itemCtrl{
			currentItem: 1,
		},
	}

	model, cmd := ui.handleItemsListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemsList, ui.state)   // Should remain unchanged
	assert.Equal(t, 1, ui.itemCtrl.currentItem) // Should remain unchanged
}

func TestUIController_itemsListView_EmptyItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			items: []models.EncryptedItem{},
		},
	}

	view := ui.itemsListView()

	assert.Contains(t, view, "Your Items")
	assert.Contains(t, view, "No items found")
	assert.Contains(t, view, "r to refresh")
	assert.Contains(t, view, "Esc to go back")
	assert.Contains(t, view, "q to quit")
	assert.NotContains(t, view, "↑/↓ to navigate")
}

func TestUIController_itemsListView_WithItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 1,
			items: []models.EncryptedItem{
				{
					Name: "Credentials Item",
					Type: models.ItemTypeCREDENTIALS,
				},
				{
					Name: "Text Item",
					Type: models.ItemTypeTEXT,
				},
				{
					Name: "Card Item",
					Type: models.ItemTypeCARD,
				},
			},
		},
	}

	view := ui.itemsListView()

	assert.Contains(t, view, "Your Items")
	assert.Contains(t, view, "Credentials Item (CREDENTIALS)")
	assert.Contains(t, view, "Text Item (TEXT)")
	assert.Contains(t, view, "Card Item (CARD)")
	assert.Contains(t, view, "↑/↓ to navigate")
	assert.Contains(t, view, "Enter to view details")
	assert.Contains(t, view, "r to refresh")
	assert.Contains(t, view, "Esc to go back")
}

func TestUIController_itemsListView_ItemSelection(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0, // First item selected
			items: []models.EncryptedItem{
				{
					Name: "Selected Item",
					Type: models.ItemTypeTEXT,
				},
				{
					Name: "Unselected Item",
					Type: models.ItemTypeCARD,
				},
			},
		},
	}

	view := ui.itemsListView()

	assert.Contains(t, view, "Selected Item")
	assert.Contains(t, view, "Unselected Item")
	// The view should show the selected item differently (though we can't test styling directly)
	assert.Contains(t, view, "Selected Item (TEXT)")
	assert.Contains(t, view, "Unselected Item (CARD)")
}

func TestUIController_itemsListView_SingleItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name: "Only Item",
					Type: models.ItemTypeBINARY,
				},
			},
		},
	}

	view := ui.itemsListView()

	assert.Contains(t, view, "Your Items")
	assert.Contains(t, view, "Only Item (BINARY)")
	assert.Contains(t, view, "↑/↓ to navigate")
	assert.Contains(t, view, "Enter to view details")
}

func TestUIController_itemsListView_ManyItems(t *testing.T) {
	items := make([]models.EncryptedItem, 10)
	for i := 0; i < 10; i++ {
		items[i] = models.EncryptedItem{
			Name: fmt.Sprintf("Item %d", i+1),
			Type: models.ItemTypeTEXT,
		}
	}

	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 5, // Middle item selected
			items:       items,
		},
	}

	view := ui.itemsListView()

	assert.Contains(t, view, "Your Items")
	assert.Contains(t, view, "Item 1 (TEXT)")
	assert.Contains(t, view, "Item 6 (TEXT)") // Selected item
	assert.Contains(t, view, "Item 10 (TEXT)")
	assert.Contains(t, view, "↑/↓ to navigate")
}

func TestUIController_itemsListView_DifferentItemTypes(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{Name: "Creds", Type: models.ItemTypeCREDENTIALS},
				{Name: "Note", Type: models.ItemTypeTEXT},
				{Name: "Payment", Type: models.ItemTypeCARD},
				{Name: "File", Type: models.ItemTypeBINARY},
			},
		},
	}

	view := ui.itemsListView()

	assert.Contains(t, view, "Creds (CREDENTIALS)")
	assert.Contains(t, view, "Note (TEXT)")
	assert.Contains(t, view, "Payment (CARD)")
	assert.Contains(t, view, "File (BINARY)")
}
