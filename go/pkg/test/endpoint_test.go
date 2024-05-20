package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDynamicID(t *testing.T) {
	s := `{"id": 1232, "name": "jack"}`
	newS := RemoveDynamicField(s, "id")
	assert.Equal(t, `{ "name": "jack"}`, newS)

	cases := []struct {
		name   string
		str    string
		newStr string
	}{
		{"case 1", `{"id": 1232, "name": "jack"}`, `{ "name": "jack"}`},
		{"case 2", `{"name": "jack", "id": 1232}`, `{"name": "jack", }`},
	}
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			newS := RemoveDynamicField(ca.str, "id")
			assert.Equal(t, ca.newStr, newS)
		})
	}
}
