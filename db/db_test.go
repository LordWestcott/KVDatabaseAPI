package db

import "testing"

func TestNewDatabase(t *testing.T) {
	db := NewDatabase()

	t.Run("NewDatabase should return a non-nil database", func(t *testing.T) {
		if db == nil {
			t.Error("NewDatabase returned nil")
		}
		return
	})

	t.Run("NewDatabase should return a database with an empty non-nil Data map", func(t *testing.T) {
		if db.Data == nil {
			t.Error("NewDatabase returned a database with a nil Data map")
		}

		if len(db.Data) != 0 {
			t.Error("NewDatabase returned a database with a non-empty Data map")
		}
		return
	})
}

func TestInitCheck(t *testing.T) {
	db := NewDatabase()
	db.Data = make(map[string]interface{})
	t.Run("initCheck should return nil if Data is not nil", func(t *testing.T) {
		err := initCheck(db)
		if err != nil {
			t.Errorf("initCheck returned an error: %s", err)
		}
		return
	})

	t.Run("initCheck should return an error if Data is nil", func(t *testing.T) {
		db.Data = nil
		err := initCheck(db)
		if err == nil {
			t.Error("initCheck did not return an error")
		}
		return
	})
}

func TestGetAllKeys(t *testing.T) {
	tt := []struct {
		name      string
		data      map[string]interface{}
		expected  []string
		shouldErr bool
	}{
		{
			name:      "empty db",
			data:      make(map[string]interface{}),
			expected:  []string{},
			shouldErr: false,
		},
		{
			name: "one key",
			data: map[string]interface{}{
				"key": "value",
			},
			expected:  []string{"key"},
			shouldErr: false,
		},
		{
			name: "multiple keys",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			expected:  []string{"key1", "key2", "key3"},
			shouldErr: false,
		},
		{
			name:      "uninitialized db",
			data:      nil,
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewDatabase()
			db.Data = tc.data

			keys, err := db.GetAllKeys()
			if err != nil {
				if !tc.shouldErr {
					t.Errorf("GetAllKeys returned an error: %s", err)
				}
				return
			}

			if len(keys) != len(tc.expected) {
				t.Errorf("GetAllKeys returned %d keys, expected %d", len(keys), len(tc.expected))
			}

			for _, k := range keys {
				var found bool
				for _, e := range tc.expected {
					if k == e {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetAllKeys did not contain expected key: %s", k)
				}
			}
		})
	}
}

func BenchmarkDatabase_GetAllKeys(b *testing.B) {

	db := NewDatabase()
	db.Data = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_, _ = db.GetAllKeys()
	}
}

func TestGet(t *testing.T) {
	tt := []struct {
		name      string
		data      map[string]interface{}
		key       string
		expected  interface{}
		shouldErr bool
	}{
		{
			name: "key exists",
			data: map[string]interface{}{
				"key": "value",
			},
			key:       "key",
			expected:  "value",
			shouldErr: false,
		},
		{
			name: "key does not exist",
			data: map[string]interface{}{
				"key": "value",
			},
			key:       "key2",
			expected:  nil,
			shouldErr: false,
		},
		{
			name: "correct value is selected when multiple keys exist",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			key:       "key2",
			expected:  "value2",
			shouldErr: false,
		},
		{
			name: "numeric value returns correctly",
			data: map[string]interface{}{
				"key": 123,
			},
			key:       "key",
			expected:  123,
			shouldErr: false,
		},
		{
			name: "boolean value returns correctly",
			data: map[string]interface{}{
				"key": true,
			},
			key:       "key",
			expected:  true,
			shouldErr: false,
		},
		{
			name:      "uninitialized db",
			data:      nil,
			key:       "key",
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewDatabase()
			db.Data = tc.data

			value, err := db.Get(tc.key)
			if err != nil {
				if !tc.shouldErr {
					t.Errorf("Get returned an error: %s", err)
				}
				return
			}

			if value != tc.expected {
				t.Errorf("Get returned %v, expected %v", value, tc.expected)
			}
		})
	}
}

func BenchmarkDatabase_Get(b *testing.B) {

	db := NewDatabase()
	db.Data = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_, _ = db.Get("key2")
	}
}

func TestSet(t *testing.T) {
	tt := []struct {
		name      string
		data      map[string]interface{}
		key       string
		value     interface{}
		expected  map[string]interface{}
		shouldErr bool
	}{
		{
			name: "set key functions correctly",
			data: map[string]interface{}{
				"key": "value",
			},
			key:   "key2",
			value: "value2",
			expected: map[string]interface{}{
				"key":  "value",
				"key2": "value2",
			},
			shouldErr: false,
		},
		{
			name: "update key functions correctly",
			data: map[string]interface{}{
				"key": "value",
			},
			key:   "key",
			value: "value2",
			expected: map[string]interface{}{
				"key": "value2",
			},
			shouldErr: false,
		},
		{
			name: "update key functions correctly when multiple keys exist",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			key:   "key2",
			value: "value4",
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value4",
				"key3": "value3",
			},
			shouldErr: false,
		},
		{
			name: "update functions correctly when changing between types",
			data: map[string]interface{}{
				"key": "value",
			},
			key:   "key",
			value: 123,
			expected: map[string]interface{}{
				"key": 123,
			},
			shouldErr: false,
		},
		{
			name:      "uninitialized db",
			data:      nil,
			key:       "key",
			value:     "value",
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewDatabase()
			db.Data = tc.data

			err := db.Set(tc.key, tc.value)
			if err != nil {
				if !tc.shouldErr {
					t.Errorf("Set returned an error: %s", err)
				}
				return
			}

			for k, v := range tc.expected {
				if db.Data[k] != v {
					t.Errorf("Set did not set key correctly: %s", k)
				}
			}
		})
	}
}

func BenchmarkDatabase_Set(b *testing.B) {

	db := NewDatabase()
	db.Data = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_ = db.Set("key2", "value4")
	}
}

func TestDelete(t *testing.T) {
	tt := []struct {
		name      string
		data      map[string]interface{}
		key       string
		expected  map[string]interface{}
		shouldErr bool
	}{
		{
			name: "delete key functions correctly",
			data: map[string]interface{}{
				"key": "value",
			},
			key:       "key",
			expected:  map[string]interface{}{},
			shouldErr: false,
		},
		{
			name: "delete key functions correctly when multiple keys exist",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			key: "key2",
			expected: map[string]interface{}{
				"key1": "value1",
				"key3": "value3",
			},
			shouldErr: false,
		},
		{
			name: "delete key does not throw error when key does not exist",
			data: map[string]interface{}{
				"key": "value",
			},
			key: "key2",
			expected: map[string]interface{}{
				"key": "value",
			},
			shouldErr: false,
		},
		{
			name:      "uninitialized db",
			data:      nil,
			key:       "key",
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewDatabase()
			db.Data = tc.data

			err := db.Delete(tc.key)
			if err != nil {
				if !tc.shouldErr {
					t.Errorf("Delete returned an error: %s", err)
				}
				return
			}

			for k, v := range tc.expected {
				if db.Data[k] != v {
					t.Errorf("Delete did not delete key correctly: %s", k)
				}
			}
		})

	}
}

func BenchmarkDatabase_Delete(b *testing.B) {
	db := NewDatabase()
	db.Data = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_ = db.Delete("key2")
	}
}
