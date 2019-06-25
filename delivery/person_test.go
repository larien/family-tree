package delivery

import (
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/controller"
	"github.com/larien/family-tree/entity"
	"github.com/stretchr/testify/assert"
)

func TestPersonEndpoints(t *testing.T) {
	r, err := repository.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c := controller.New(r)

	router := New(c)

	t.Run("should GET all People", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/person", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should have created resource", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := fmt.Sprintf(`[
			{
				"name": "Anakin"
			},
			{
				"name": "Luke",
				"parents": ["Anakin"]
			}
		]`)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/person", strings.NewReader(payload))
		router.ServeHTTP(w, req)

		var people *entity.Person
		json.NewDecoder(w.Body).Decode(&people)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	r.DB.Session.Close()
	r.DB.Driver.Close()
}