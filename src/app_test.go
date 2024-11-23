package main

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/runner"
)

type mockDB struct {
	execFunc  func(query string, args ...interface{}) (sql.Result, error)
	queryFunc func(query string, args ...interface{}) *sql.Row
	closeFunc func() error
}

func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.execFunc(query, args...)
}

func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.queryFunc(query, args...)
}

func (m *mockDB) Close() error {
	return m.closeFunc()
}

type mockResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (m mockResult) LastInsertId() (int64, error) {
	return m.lastInsertId, nil
}

func (m mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func TestPushHandler(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		mockDB       *mockDB
		expectedCode int
		expectedBody string
	}{
		{
			name:  "Valid value",
			value: "42",
			mockDB: &mockDB{
				execFunc: func(query string, args ...interface{}) (sql.Result, error) {
					return mockResult{1, 1}, nil
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: "Success\n",
		},
		{
			name:         "Missing value",
			value:        "",
			mockDB:       &mockDB{},
			expectedCode: http.StatusBadRequest,
			expectedBody: "missing value query parameter\n",
		},
		{
			name:         "Invalid value",
			value:        "not_a_number",
			mockDB:       &mockDB{},
			expectedCode: http.StatusBadRequest,
			expectedBody: "value should be integer\n",
		},
		{
			name:  "Database error",
			value: "42",
			mockDB: &mockDB{
				execFunc: func(query string, args ...interface{}) (sql.Result, error) {
					return nil, errors.New("database error")
				},
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "db.Exec INSERT INTO: database error\n",
		},
	}

	for _, tt := range tests {
		runner.Run(t, tt.name, func(t provider.T) {
			db = tt.mockDB
			req, err := http.NewRequest("GET", "/push?value="+tt.value, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(pushHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})

		//t.Run(tt.name, func(t *testing.T) {
		//	db = tt.mockDB
		//	req, err := http.NewRequest("GET", "/push?value="+tt.value, nil)
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//
		//	rr := httptest.NewRecorder()
		//	handler := http.HandlerFunc(pushHandler)
		//	handler.ServeHTTP(rr, req)
		//
		//	if status := rr.Code; status != tt.expectedCode {
		//		t.Errorf("handler returned wrong status code: got %v want %v",
		//			status, tt.expectedCode)
		//	}
		//
		//	if rr.Body.String() != tt.expectedBody {
		//		t.Errorf("handler returned unexpected body: got %v want %v",
		//			rr.Body.String(), tt.expectedBody)
		//	}
		//})
	}
}

func TestAvgHandler(t *testing.T) {
	tests := []struct {
		name         string
		mockDB       *mockDB
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success with data",
			mockDB: &mockDB{

				queryFunc: func(query string, args ...interface{}) *sql.Row {
					return NewRow(15.5)
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: "Average: 15.500000\n",
		},
		{
			name: "No data",
			mockDB: &mockDB{
				queryFunc: func(query string, args ...interface{}) *sql.Row {
					return NewRow(nil)
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: "No data yet\n",
		},
		{
			name: "Database error",
			mockDB: &mockDB{
				queryFunc: func(query string, args ...interface{}) *sql.Row {
					return NewRow(errors.New("database error"))
				},
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "db.QueryRow SELECT AVG: database error\n",
		},
	}

	for _, tt := range tests {
		runner.Run(t, tt.name, func(t provider.T) {
			db = tt.mockDB
			req, err := http.NewRequest("GET", "/avg", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(avgHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})

		//t.Run(tt.name, func(t *testing.T) {
		//	db = tt.mockDB
		//	req, err := http.NewRequest("GET", "/avg", nil)
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//
		//	rr := httptest.NewRecorder()
		//	handler := http.HandlerFunc(avgHandler)
		//	handler.ServeHTTP(rr, req)
		//
		//	if status := rr.Code; status != tt.expectedCode {
		//		t.Errorf("handler returned wrong status code: got %v want %v",
		//			status, tt.expectedCode)
		//	}
		//
		//	if rr.Body.String() != tt.expectedBody {
		//		t.Errorf("handler returned unexpected body: got %v want %v",
		//			rr.Body.String(), tt.expectedBody)
		//	}
		//})
	}
}

func NewRow(value interface{}) *sql.Row {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	switch v := value.(type) {
	case error:
		mock.ExpectQuery(".*").WillReturnError(v)
	default:
		rows := sqlmock.NewRows([]string{"avg"}).AddRow(v)
		mock.ExpectQuery(".*").WillReturnRows(rows)
	}

	return db.QueryRow("dummy query")
}
