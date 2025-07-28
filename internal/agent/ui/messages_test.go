package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessages_init(t *testing.T) {
	m := &Messages{}

	m.init()

	assert.NotNil(t, m.M)
	assert.Empty(t, m.M)
	assert.IsType(t, map[string]string{}, m.M)
}

func TestMessages_Set_NilMap(t *testing.T) {
	m := &Messages{M: nil}

	m.Set("test-key", "test-value")

	assert.NotNil(t, m.M)
	assert.Equal(t, "test-value", m.M["test-key"])
}

func TestMessages_Set_ExistingMap(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"existing-key": "existing-value",
		},
	}

	m.Set("new-key", "new-value")

	assert.Equal(t, "existing-value", m.M["existing-key"])
	assert.Equal(t, "new-value", m.M["new-key"])
}

func TestMessages_Set_OverwriteExisting(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"test-key": "old-value",
		},
	}

	m.Set("test-key", "new-value")

	assert.Equal(t, "new-value", m.M["test-key"])
}

func TestMessages_Set_EmptyKey(t *testing.T) {
	m := &Messages{}

	m.Set("", "test-value")

	assert.Equal(t, "test-value", m.M[""])
}

func TestMessages_Set_EmptyValue(t *testing.T) {
	m := &Messages{}

	m.Set("test-key", "")

	assert.Equal(t, "", m.M["test-key"])
}

func TestMessages_Get_NilMap(t *testing.T) {
	m := &Messages{M: nil}

	result := m.Get("test-key")

	assert.Empty(t, result)
}

func TestMessages_Get_ExistingKey(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"test-key": "test-value",
		},
	}

	result := m.Get("test-key")

	assert.Equal(t, "test-value", result)
}

func TestMessages_Get_NonExistingKey(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"existing-key": "existing-value",
		},
	}

	result := m.Get("non-existing-key")

	assert.Empty(t, result)
}

func TestMessages_Get_EmptyKey(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"": "empty-key-value",
		},
	}

	result := m.Get("")

	assert.Equal(t, "empty-key-value", result)
}

func TestMessages_Clear_NilMap(t *testing.T) {
	m := &Messages{M: nil}

	// Should not panic
	m.Clear("test-key")

	assert.Nil(t, m.M)
}

func TestMessages_Clear_ExistingKey(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	m.Clear("key1")

	assert.NotContains(t, m.M, "key1")
	assert.Contains(t, m.M, "key2")
	assert.Equal(t, "value2", m.M["key2"])
}

func TestMessages_Clear_NonExistingKey(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"existing-key": "existing-value",
		},
	}

	m.Clear("non-existing-key")

	// Should not affect existing keys
	assert.Contains(t, m.M, "existing-key")
	assert.Equal(t, "existing-value", m.M["existing-key"])
}

func TestMessages_Clear_EmptyKey(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"":         "empty-key-value",
			"test-key": "test-value",
		},
	}

	m.Clear("")

	assert.NotContains(t, m.M, "")
	assert.Contains(t, m.M, "test-key")
}

func TestMessages_ClearAll_NilMap(t *testing.T) {
	m := &Messages{M: nil}

	// Should not panic
	m.ClearAll()

	assert.Nil(t, m.M)
}

func TestMessages_ClearAll_ExistingMap(t *testing.T) {
	m := &Messages{
		M: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}

	m.ClearAll()

	assert.NotNil(t, m.M)
	assert.Empty(t, m.M)
	assert.Len(t, m.M, 0)
}

func TestMessages_ClearAll_EmptyMap(t *testing.T) {
	m := &Messages{
		M: map[string]string{},
	}

	m.ClearAll()

	assert.NotNil(t, m.M)
	assert.Empty(t, m.M)
}

// Integration tests - testing multiple operations together
func TestMessages_Integration_SetGetClear(t *testing.T) {
	m := &Messages{}

	// Initially empty
	assert.Empty(t, m.Get("key1"))

	// Set values
	m.Set("key1", "value1")
	m.Set("key2", "value2")

	// Get values
	assert.Equal(t, "value1", m.Get("key1"))
	assert.Equal(t, "value2", m.Get("key2"))

	// Clear one key
	m.Clear("key1")
	assert.Empty(t, m.Get("key1"))
	assert.Equal(t, "value2", m.Get("key2"))

	// Clear all
	m.ClearAll()
	assert.Empty(t, m.Get("key2"))
}

func TestMessages_Integration_MultipleSetOperations(t *testing.T) {
	m := &Messages{}

	// Set multiple values
	m.Set("error", "Login failed")
	m.Set("info", "Processing...")
	m.Set("success", "Operation completed")

	// Verify all values
	assert.Equal(t, "Login failed", m.Get("error"))
	assert.Equal(t, "Processing...", m.Get("info"))
	assert.Equal(t, "Operation completed", m.Get("success"))

	// Update existing value
	m.Set("error", "Updated error message")
	assert.Equal(t, "Updated error message", m.Get("error"))

	// Other values should remain unchanged
	assert.Equal(t, "Processing...", m.Get("info"))
	assert.Equal(t, "Operation completed", m.Get("success"))
}

func TestMessages_Integration_ClearAndSet(t *testing.T) {
	m := &Messages{}

	// Set initial values
	m.Set("key1", "value1")
	m.Set("key2", "value2")

	// Clear all
	m.ClearAll()

	// Set new values after clearing
	m.Set("key3", "value3")

	// Old keys should be empty, new key should have value
	assert.Empty(t, m.Get("key1"))
	assert.Empty(t, m.Get("key2"))
	assert.Equal(t, "value3", m.Get("key3"))
}

func TestMessages_ZeroValue(t *testing.T) {
	var m Messages // Zero value

	// Should work with zero value
	assert.Empty(t, m.Get("test"))

	m.Set("test", "value")
	assert.Equal(t, "value", m.Get("test"))

	m.Clear("test")
	assert.Empty(t, m.Get("test"))

	m.ClearAll()
	assert.Empty(t, m.Get("anything"))
}

func TestMessages_Pointer_NilReceiver(t *testing.T) {
	var m *Messages = nil

	// These should panic with nil receiver
	assert.Panics(t, func() {
		m.Set("key", "value")
	})

	assert.Panics(t, func() {
		m.Get("key")
	})

	assert.Panics(t, func() {
		m.Clear("key")
	})

	assert.Panics(t, func() {
		m.ClearAll()
	})
}
