package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"daytask/internal/http-server/handlers/task/save"
	"daytask/internal/http-server/handlers/task/save/mocks"
	"daytask/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		title      string
		owner     string
		date       string
		respError string
		mockError error
	}{
		{
			title:  "Success",
			owner: "test_owner",
			date:   "2024-02-01",
		},
		{
			title:  "Empty owner",
			owner: 	"",
			date:   "2024-02-02",
		},
		{
			title:      "Empty date",
			date:       "",
			owner:     "some_owner",
			respError: "field Date is a required field",
		},
		{
			title:      "Invalid date",
			date:       "some invalid date",
			owner:     "some_owner",
			respError: "field Date is not valid",
		},
		{
			title:      "SaveTask Error",
			owner:     "test_owner",
			date:       "2024-02-03",
			respError: "failed to save task",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			taskSaverMock := mocks.NewTASKSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				taskSaverMock.On("SaveTask",mock.AnythingOfType("string"), mock.AnythingOfType("string"), tc.date).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), taskSaverMock)

			input := fmt.Sprintf(`{"title": "%s", "owner": "%s", "date": "%s"}`, tc.title, tc.owner, tc.date)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
