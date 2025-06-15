package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}

}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&count=%d", v.count), nil)

		handler.ServeHTTP(response, req)
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("can't read%v\n", err)
		}

		slice := strings.Split(string(body), ",")
		var sliceCafe []string
		for _, str := range slice {
			if str != "" {
				sliceCafe = append(sliceCafe, str)
			}
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.want, len(sliceCafe))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", v.search), nil)

		handler.ServeHTTP(response, req)
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("can't read%v\n", err)
		}

		slice := strings.Split(string(body), ",")
		var found int
		for _, str := range slice {
			if strings.Contains(strings.ToLower(str), strings.ToLower(v.search)) {
				found++
			}
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.wantCount, found)
	}
}
