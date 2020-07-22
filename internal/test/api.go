package test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rdnply/backend-trainee-assignment/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type APITestCase struct {
	Name         string
	Method       string
	URL          string
	Body         string
	Handler      http.HandlerFunc
	WantStatus   int
	WantResponse string
}

func Endpoint(t *testing.T, tc APITestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		req, err := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))
		if err != nil {
			t.Fatalf("can't create test request %v", err)
		}
		res := httptest.NewRecorder()

		tc.Handler.ServeHTTP(res, req)
		assert.Equal(t, tc.WantStatus, res.Code, wrongCode(res.Code, tc.WantStatus))

		if tc.WantResponse != "" {
			pattern := strings.Trim(tc.WantResponse, "*")
			if pattern != tc.WantResponse {
				assert.Contains(t, res.Body.String(), pattern, wrongBody(res.Body.String(), pattern))
			} else {
				assert.JSONEq(t, tc.WantResponse, res.Body.String(), "response mismatch")
			}
		}

	})
}

func wrongCode(want int, actual int) string {
	return fmt.Sprintf("returned wrong status code: got %v, want %v", want, actual)
}

func wrongBody(want string, actual string) string {
	return fmt.Sprintf("returned unexpected body: got %v\n, want %v", want, actual)
}

func Logger() logger.Logger {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: true,
	}

	logger, err := logger.New(config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatal("could not instantiate logger: ", err)
	}

	return logger
}
