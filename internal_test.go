package fuze

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"reflect"
	"testing"
)

func TestGetParameters(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectParams  map[int]string
		expectPathEls []string
	}{
		{
			name:          "default path",
			path:          "/home/flat",
			expectPathEls: []string{"home", "flat"},
		},
		{
			name:          "default path with parameters",
			path:          "/home/{number}/flat/{flatNumber}",
			expectParams:  map[int]string{1: "number", 3: "flatNumber"},
			expectPathEls: []string{"home", "flat"},
		},
		{
			name:          "short path",
			path:          "/{id}",
			expectParams:  map[int]string{0: "id"},
			expectPathEls: []string{},
		},
		{
			name:          "all parameters",
			path:          "/{id}/{age}/{country}",
			expectParams:  map[int]string{0: "id", 1: "age", 2: "country"},
			expectPathEls: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params, pathEls := transformPath(tc.path)
			for k, v := range params {
				expVal, ok := tc.expectParams[k]
				require.True(t, ok)
				require.Equal(t, expVal, v)
			}
			require.True(t, reflect.DeepEqual(tc.expectPathEls, pathEls))
		})
	}
}

func TestParseQueryParameter(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		parameters map[int]string
		expect     map[string]string
	}{
		{
			name:       "default path",
			path:       "/home/77/flat/7",
			parameters: map[int]string{1: "number", 3: "flatNumber"},
			expect:     map[string]string{"number": "77", "flatNumber": "7"},
		},
		{
			name:       "short path",
			path:       "/10",
			parameters: map[int]string{0: "id"},
			expect:     map[string]string{"id": "10"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := parseQueryParameter(&url.URL{Path: tc.path}, tc.parameters)
			for k, v := range res {
				expVal, ok := tc.expect[k]
				require.True(t, ok)
				require.Equal(t, expVal, v)
			}
		})
	}
}

func TestFindLikePath(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		paths      map[string]HandlerStruct
		pathExists bool
	}{
		{
			name: "default path",
			path: "/home/77/flat/7",
			paths: map[string]HandlerStruct{"/home/{number}/flat/{flatNumber}": {
				pathElements: []string{"home", "flat"},
				pathParameters: map[int]string{
					1: "number",
					3: "flatNumber",
				},
			}, "/home/{number}": {
				pathElements:   []string{"home"},
				pathParameters: map[int]string{1: "number"},
			}},
			pathExists: true,
		},
		{
			name: "one redundant element in the path",
			path: "/home/77/flat/7/floor",
			paths: map[string]HandlerStruct{"/home/{number}/flat/{flatNumber}": {
				pathElements: []string{"home", "flat"},
				pathParameters: map[int]string{
					1: "number",
					3: "flatNumber",
				},
			}},
			pathExists: false,
		},
		{
			name: "the same path without parameters",
			path: "/home/flat/floor",
			paths: map[string]HandlerStruct{"/home/{number}/flat/{flatNumber}": {
				pathElements: []string{"home", "flat"},
				pathParameters: map[int]string{
					1: "number",
					3: "flatNumber",
				},
			}},
			pathExists: false,
		},
		{
			name: "short with parameter",
			path: "/5",
			paths: map[string]HandlerStruct{"/{id}": {
				pathElements: []string{},
				pathParameters: map[int]string{
					0: "id",
				},
			}},
			pathExists: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, p, ok := findLikePath(tc.path, tc.paths)
			require.Equal(t, tc.pathExists, ok)
			if tc.pathExists {
				require.NotEmpty(t, p)
			}
		})
	}
}

func BenchmarkFindLikePath(b *testing.B) {
	tests := []struct {
		name       string
		path       string
		paths      map[string]HandlerStruct
		pathExists bool
	}{
		{
			name: "default path",
			path: "/home/77/flat/7",
			paths: map[string]HandlerStruct{"/home/{number}/flat/{flatNumber}": {
				pathElements: []string{"home", "flat"},
				pathParameters: map[int]string{
					1: "number",
					3: "flatNumber",
				},
			}},
			pathExists: true,
		},
		{
			name: "one redundant element in the path",
			path: "/home/77/flat/7/floor",
			paths: map[string]HandlerStruct{"/home/{number}/flat/{flatNumber}": {
				pathElements: []string{"home", "flat"},
				pathParameters: map[int]string{
					1: "number",
					3: "flatNumber",
				},
			}},
			pathExists: false,
		},
		{
			name: "the same path without parameters",
			path: "/home/flat/floor",
			paths: map[string]HandlerStruct{"/home/{number}/flat/{flatNumber}": {
				pathElements: []string{"home", "flat"},
				pathParameters: map[int]string{
					1: "number",
					3: "flatNumber",
				},
			}},
			pathExists: false,
		},
		{
			name: "short with parameter",
			path: "/5",
			paths: map[string]HandlerStruct{"/{id}": {
				pathElements: []string{},
				pathParameters: map[int]string{
					0: "id",
				},
			}},
			pathExists: true,
		},
	}

	b.ResetTimer()
	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				findLikePath(tc.path, tc.paths)
			}
		})
	}
}

func TestRemoveFirstSlash(t *testing.T) {
	path := "/path/123"
	removeFirstSlash(&path)

	require.Equal(t, "path/123", path)
}
