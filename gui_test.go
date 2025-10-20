package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdateMsgInput(t *testing.T) {
	// a an empty initial model
	m := newModel([]list.Item{})
	m.chosenPrefix = true
	m.chosenScope = true

	// an Enter key press
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.updateMsgInput(msg)
	m = updatedModel.(*model)

	// test that the user cannot proceed
	if m.chosenMsg {
		t.Error("should not be able to proceed with empty commit msg")
	}
}
