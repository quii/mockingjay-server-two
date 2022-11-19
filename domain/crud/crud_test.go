package crud_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/crud"
)

func TestCRUD(t *testing.T) {
	repo := crud.NewCRUD[string, int]()

	t.Run("basic create and read", func(t *testing.T) {
		assert.NoError(t, repo.Create("hello", 1))

		t.Run("and read all", func(t *testing.T) {
			all, err := repo.GetAll()
			assert.NoError(t, err)
			assert.Equal(t, len(all), 1)
			assert.Equal(t, all[0], 1)
		})

		t.Run("and get by id", func(t *testing.T) {
			got, exists, err := repo.GetByID("hello")
			assert.NoError(t, err)
			assert.True(t, exists)
			assert.Equal(t, 1, got)
		})
	})

	t.Run("multiple items", func(t *testing.T) {
		assert.NoError(t, repo.Create("goodbye", 2))

		t.Run("and can read all", func(t *testing.T) {
			all, err := repo.GetAll()
			assert.NoError(t, err)
			assert.Equal(t, len(all), 2)
		})

		t.Run("and can get by id", func(t *testing.T) {
			got, exists, err := repo.GetByID("goodbye")
			assert.NoError(t, err)
			assert.True(t, exists)
			assert.Equal(t, 2, got)
		})

		t.Run("can delete by id", func(t *testing.T) {
			assert.NoError(t, repo.Delete("goodbye"))
			all, err := repo.GetAll()
			assert.NoError(t, err)
			assert.Equal(t, len(all), 1)
			assert.Equal(t, all[0], 1)
		})
	})

	t.Run("get by returns not found when something doesn't exist", func(t *testing.T) {
		_, found, err := repo.GetByID("what")
		assert.NoError(t, err)
		assert.False(t, found)
	})
}
