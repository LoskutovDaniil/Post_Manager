package adaptor_test

import (
	"testing"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/adaptor"
)

func TestInMemoryRepository(t *testing.T) {
	t.Parallel()
	runSharedRepositoryTests(t, adaptor.NewInMemoryRepository())
}
