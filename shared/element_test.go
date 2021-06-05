package shared

import (
	"testing"

	"github.com/tj/assert"
)

func TestElement(t *testing.T) {
	t.Run("NewElement", func(t *testing.T) {
		el := NewElement("test", 10)
		t.Run("Creates new element", func(t *testing.T) {
			assert.Equal(t, "test", el.Name)
			assert.Equal(t, 10., el.Val)
		})
	})
}

func TestElements(t *testing.T) {
	t.Run("Given Elements", func(t *testing.T) {
		el := NewElements()
		t.Run("Calling Add", func(t *testing.T) {
			el.Add("test", 10)
			t.Run("Adds the element to the list", func(t *testing.T) {
				assert.Equal(t, 1, el.Len())
			})
			t.Run("After adding more elements", func(t *testing.T) {
				el.Add("test3", 13)
				el.Add("test2", 12)
				el.Add("test1", 11)
				t.Run("Calling Index on present element", func(t *testing.T) {
					index, found := el.Index("test2")
					t.Run("Returns the correct index", func(t *testing.T) {
						assert.Equal(t, 2, index)
					})
					t.Run("Returns positive found", func(t *testing.T) {
						assert.True(t, found)
					})
				})
				t.Run("Calling Index on missing element", func(t *testing.T) {
					_, found := el.Index("test111")
					t.Run("Returns not found", func(t *testing.T) {
						assert.False(t, found)
					})
				})
				t.Run("After Sort", func(t *testing.T) {
					el.Sort()
					t.Run("Elements are sorted", func(t *testing.T) {
						index, _ := el.Index("test3")
						assert.Equal(t, 3, index)
						index2, _ := el.Index("test1")
						assert.Equal(t, 1, index2)
					})
				})
				t.Run("Having second set of elements", func(t *testing.T) {
					el2 := NewElements()
					el2.Add("test3", 113)
					el2.Add("test2", 112)
					el2.Add("test1", 111)
					el2.Add("test4", 444)
					t.Run("SumMerge with coef 2", func(t *testing.T) {
						el.SumMerge(el2, 2)
						t.Run("Returns correct elements", func(t *testing.T) {
							index, found := el.Index("test1")
							assert.True(t, found)
							assert.Equal(t, 1, index)
							assert.Equal(t, 233., el[index].Val)
						})
						t.Run("New elements are added", func(t *testing.T) {
							index, found := el.Index("test4")
							assert.True(t, found)
							assert.Equal(t, 888., el[index].Val)
						})
					})
				})
			})
		})
	})
}
