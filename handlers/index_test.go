package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDatabase struct {
	//counts
	getCalledCount        int
	getAllKeysCalledCount int
	setCalledCount        int
	deleteCalledCount     int

	//arguments
	getKeyArg    string
	setKeyArg    string
	setValueArg  interface{}
	deleteKeyArg string

	//options
	isEmpty           bool
	shouldError       bool
	getShouldError    bool
	deleteShouldError bool
}

func (m *mockDatabase) Get(key string) (interface{}, error) {
	m.getCalledCount++
	m.getKeyArg = key

	if m.shouldError || m.getShouldError {
		return nil, errors.New("error")

	}

	if key == "not-found" {
		return nil, nil
	}

	return "hello", nil
}

func (m *mockDatabase) GetAllKeys() ([]string, error) {
	m.getAllKeysCalledCount++

	if m.shouldError {
		return nil, errors.New("error")

	}

	if m.isEmpty {
		return []string{}, nil
	}

	return []string{"hello", "world"}, nil
}

func (m *mockDatabase) Set(key string, value interface{}) error {
	m.setCalledCount++
	m.setKeyArg = key
	m.setValueArg = value

	if m.shouldError {
		return errors.New("error")
	}

	return nil
}

func (m *mockDatabase) Delete(key string) error {
	m.deleteCalledCount++
	m.deleteKeyArg = key

	if m.shouldError || m.deleteShouldError {
		return errors.New("error")
	}

	return nil
}

func TestIndexHandler(t *testing.T) {
	t.Run("IndexHandler should return a function", func(t *testing.T) {
		d := &mockDatabase{}
		h := IndexHandler(d)
		if h == nil {
			t.Error("IndexHandler returned nil")
		}
	})
}

func TestGetHandler(t *testing.T) {
	tt := []struct {
		name                    string
		request                 *http.Request
		expectedArgument        string
		expectedGetCalledCount  int
		expectedGetAllKeysCount int
		expectedResponseCode    int
		expectedResponseBody    string
		isDbEmpty               bool
		shouldError             bool
	}{
		{
			name:                    "Should Call Get with Correct Key",
			request:                 httptest.NewRequest(http.MethodGet, "/test", nil),
			expectedArgument:        "test",
			expectedGetCalledCount:  1,
			expectedGetAllKeysCount: 0,
			expectedResponseCode:    http.StatusOK,
			expectedResponseBody:    "\"hello\"\n",
			isDbEmpty:               false,
		},
		{
			name:                    "Should Return 404 if Key Not Found",
			request:                 httptest.NewRequest(http.MethodGet, "/not-found", nil),
			expectedArgument:        "not-found",
			expectedGetCalledCount:  1,
			expectedGetAllKeysCount: 0,
			expectedResponseCode:    http.StatusNotFound,
			expectedResponseBody:    "",
			isDbEmpty:               false,
		},
		{
			name:                    "Get With No Path Parameter Should Return All Keys",
			request:                 httptest.NewRequest(http.MethodGet, "/", nil),
			expectedGetCalledCount:  0,
			expectedGetAllKeysCount: 1,
			expectedResponseCode:    http.StatusOK,
			expectedResponseBody:    "[\"hello\",\"world\"]\n",
			isDbEmpty:               false,
		},
		{
			name:                    "Should return empty JSON array if database is empty",
			request:                 httptest.NewRequest(http.MethodGet, "/", nil),
			expectedGetCalledCount:  0,
			expectedGetAllKeysCount: 1,
			expectedResponseCode:    http.StatusOK,
			expectedResponseBody:    "[]\n",
			isDbEmpty:               true,
		},
		{
			name:                    "Should Return 500 if Database Returns Error Getting All Keys",
			request:                 httptest.NewRequest(http.MethodGet, "/", nil),
			expectedGetCalledCount:  0,
			expectedGetAllKeysCount: 1,
			expectedResponseCode:    http.StatusInternalServerError,
			expectedResponseBody:    "error - getting all keys\n",
			shouldError:             true,
		},
		{
			name:                    "Should Return 500 if Database Returns Error Getting Single Key",
			request:                 httptest.NewRequest(http.MethodGet, "/test", nil),
			expectedArgument:        "test",
			expectedGetCalledCount:  1,
			expectedGetAllKeysCount: 0,
			expectedResponseCode:    http.StatusInternalServerError,
			expectedResponseBody:    "error - getting key\n",
			shouldError:             true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			d := &mockDatabase{}
			d.isEmpty = tc.isDbEmpty
			d.shouldError = tc.shouldError
			w := httptest.NewRecorder()
			getHandler(d, w, tc.request)

			if d.getCalledCount != tc.expectedGetCalledCount {
				t.Errorf("Get called count: got %d, want %d", d.getCalledCount, tc.expectedGetCalledCount)
			}

			if d.getAllKeysCalledCount != tc.expectedGetAllKeysCount {
				t.Errorf("GetAllKeys called count: got %d, want %d", d.getAllKeysCalledCount, tc.expectedGetAllKeysCount)
			}

			if d.getKeyArg != tc.expectedArgument {
				t.Errorf("Get called with wrong argument: got %s, want %s", d.getKeyArg, tc.expectedArgument)
			}

			if w.Code != tc.expectedResponseCode {
				t.Errorf("Response code: got %d, want %d", w.Code, tc.expectedResponseCode)
			}

			if w.Body.String() != tc.expectedResponseBody {
				t.Errorf("Response body: got %s, want %s", w.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}

func TestPutHandler(t *testing.T) {
	tt := []struct {
		name                   string
		request                *http.Request
		expectedSetCalledCount int
		expectedSetKey         string
		expectedSetValue       interface{}
		expectedResponseCode   int
		expectedResponseBody   string
		noKeySet               bool
		shouldError            bool
	}{
		{
			name:                   "Should Call Set with Correct Key and Value",
			request:                httptest.NewRequest(http.MethodPut, "/test", bytes.NewBufferString("hello")),
			expectedSetCalledCount: 1,
			expectedSetKey:         "test",
			expectedSetValue:       "hello",
			expectedResponseCode:   http.StatusOK,
			expectedResponseBody:   "",
		},
		{
			name:                 "Should Return 400 if No Key Provided",
			request:              httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("hello")),
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: "error - no key provided\n",
			noKeySet:             true,
		},
		{
			name:                   "Should Return 500 if Database Returns Error NON-JSON Value",
			request:                httptest.NewRequest(http.MethodPut, "/test", bytes.NewBufferString("hello")),
			expectedSetCalledCount: 1,
			expectedSetKey:         "test",
			expectedSetValue:       "hello",
			expectedResponseCode:   http.StatusInternalServerError,
			expectedResponseBody:   "error - putting kv pair\n",
			shouldError:            true,
		},
		{
			name:                   "Should Return 500 if Database Returns Error JSON Value",
			request:                httptest.NewRequest(http.MethodPut, "/test", bytes.NewBufferString("{\"key\": \"value\"}")),
			expectedSetCalledCount: 1,
			expectedSetKey:         "test",
			expectedSetValue:       map[string]interface{}{"key": "value"},
			expectedResponseCode:   http.StatusInternalServerError,
			expectedResponseBody:   "error - putting json kv pair\n",
			shouldError:            true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			d := &mockDatabase{}
			d.shouldError = tc.shouldError
			w := httptest.NewRecorder()
			putHandler(d, w, tc.request)

			if d.setCalledCount != tc.expectedSetCalledCount {
				t.Errorf("Set called count: got %d, want %d", d.setCalledCount, tc.expectedSetCalledCount)
			}

			if w.Code != tc.expectedResponseCode {
				t.Errorf("Response code: got %d, want %d", w.Code, tc.expectedResponseCode)
			}

			if w.Body.String() != tc.expectedResponseBody {
				t.Errorf("Response body: got %s, want %s", w.Body.String(), tc.expectedResponseBody)
			}

			if tc.noKeySet {
				return
			}

			if d.setKeyArg != tc.expectedSetKey {
				t.Errorf("Set called with wrong key: got %s, want %s", d.setKeyArg, tc.expectedSetKey)
			}

			switch tc.expectedSetValue.(type) {
			case string:
				if d.setValueArg.(string) != tc.expectedSetValue {
					t.Errorf("Set called with wrong value: got %s, want %s", d.setValueArg, tc.expectedSetValue)
				}
			case map[string]interface{}:
				for k, v := range d.setValueArg.(map[string]interface{}) {
					if tc.expectedSetValue.(map[string]interface{})[k] != v {
						t.Errorf("Set called with wrong value: got %v, want %v", d.setValueArg, tc.expectedSetValue)
					}
				}
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	tt := []struct {
		name                      string
		request                   *http.Request
		expectedDeleteCalledCount int
		expectedGetCalledCount    int
		expectedGetArgument       string
		expectedDeleteKey         string
		expectedResponseCode      int
		dbGetShouldFail           bool
		dbDeleteShouldFail        bool
	}{
		{
			name:                      "Should Call Delete with Correct Key",
			request:                   httptest.NewRequest(http.MethodDelete, "/test", nil),
			expectedGetArgument:       "test",
			expectedGetCalledCount:    1,
			expectedDeleteCalledCount: 1,
			expectedDeleteKey:         "test",
			expectedResponseCode:      http.StatusOK,
		},
		{
			name:                 "Should Return 400 if No Key Provided",
			request:              httptest.NewRequest(http.MethodDelete, "/", nil),
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                      "Should Return 404 if Key Not Found",
			request:                   httptest.NewRequest(http.MethodDelete, "/not-found", nil),
			expectedDeleteCalledCount: 0,
			expectedGetArgument:       "not-found",
			expectedGetCalledCount:    1,
			expectedResponseCode:      http.StatusNotFound,
		},
		{
			name:                      "Should Return 500 if Database Returns Error On Get",
			request:                   httptest.NewRequest(http.MethodDelete, "/test", nil),
			expectedGetArgument:       "test",
			expectedGetCalledCount:    1,
			expectedDeleteCalledCount: 0,
			expectedResponseCode:      http.StatusInternalServerError,
			dbGetShouldFail:           true,
		},
		{
			name:                      "Should Return 500 if Database Returns Error On Delete",
			request:                   httptest.NewRequest(http.MethodDelete, "/test", nil),
			expectedGetArgument:       "test",
			expectedGetCalledCount:    1,
			expectedDeleteCalledCount: 1,
			expectedDeleteKey:         "test",
			expectedResponseCode:      http.StatusInternalServerError,
			dbDeleteShouldFail:        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			d := &mockDatabase{}
			d.deleteShouldError = tc.dbDeleteShouldFail
			d.getShouldError = tc.dbGetShouldFail
			w := httptest.NewRecorder()
			deleteHandler(d, w, tc.request)

			if d.deleteCalledCount != tc.expectedDeleteCalledCount {
				t.Errorf("Delete called count: got %d, want %d", d.deleteCalledCount, tc.expectedDeleteCalledCount)
			}

			if d.deleteKeyArg != tc.expectedDeleteKey {
				t.Errorf("Delete called with wrong key: got %s, want %s", d.deleteKeyArg, tc.expectedDeleteKey)
			}

			if w.Code != tc.expectedResponseCode {
				t.Errorf("Response code: got %d, want %d", w.Code, tc.expectedResponseCode)
			}

			if d.getKeyArg != tc.expectedGetArgument {
				t.Errorf("Get called with wrong argument: got %s, want %s", d.getKeyArg, tc.expectedGetArgument)
			}

			if d.getCalledCount != tc.expectedGetCalledCount {
				t.Errorf("Get called count: got %d, want %d", d.getCalledCount, tc.expectedGetCalledCount)
			}
		})
	}
}
