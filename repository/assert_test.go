package repository

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func assertElementsMatch[T any](t *testing.T, listA, listB []T, eqFunc func(a, b T) bool) {
	t.Helper()	

	extraA, extraB := diffLists(listA, listB, eqFunc)
	if len(extraA) == 0 && len(extraB) == 0 {
		return
	}

	assert.ElementsMatch(t, listA, listB)
}


// this is a copy paste from assert package with generics variant
func diffLists[T any](listA, listB []T, eqFunc func(a, b T) bool) (extraA, extraB []T) {
	aLen := len(listA)
	bLen := len(listB)

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := listA[i]
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if eqFunc(element, listB[j]) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, listB[j])
	}

	return
}