package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

func TestCafeWhenErr(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	req := []struct {
		Address string
		Status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, item := range req {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", item.Address, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, item.Status, response.Code)
		assert.Equal(t, item.message, strings.TrimSpace(response.Body.String()))
		fmt.Println(response.Body.String())
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	req := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(10000, len(cafeList["moscow"]))},
	}

	for _, item := range req {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&count="+strconv.Itoa(item.count), nil)

		handler.ServeHTTP(response, req)

		responseStr := response.Body.String()
		cafes := strings.Split(responseStr, ",")

		if len(cafes) == 1 && cafes[0] == "" {
			cafes = []string{}
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, item.want, len(cafes))
		fmt.Println(response.Body.String())
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, item := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+item.search, nil)

		handler.ServeHTTP(response, req)

		responseStr := response.Body.String()
		cafes := strings.Split(responseStr, ",")

		if len(cafes) == 1 && cafes[0] == "" {
			cafes = []string{}
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, item.wantCount, len(cafes))

		lowerWord := strings.ToLower(item.search)
		for _, cafe := range cafes {
			cafe = strings.ToLower(cafe)
			assert.True(t, strings.Contains(cafe, lowerWord), "Кафе '%s' должно содержать '%s'", cafe, item.search)
		}
	}
}
