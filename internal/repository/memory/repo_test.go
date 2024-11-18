package memory

import (
	"testing"
	"time"

	"github.com/vodolaz095/asset_storage/internal/model"
	"github.com/vodolaz095/asset_storage/internal/repository"
)

func TestRepo(t *testing.T) {
	t.Skipf("not implemented yet")
	repo := Repository{
		Users: map[string]model.User{
			"alice": model.User{
				ID:           1,
				Login:        "alice",
				PasswordHash: "5ebe2294ecd0e0f08eab7690d2a6ee69",
				CreatedAt:    time.Now(),
			},
		},
		Sessions: nil,
		Assets:   nil,
		Active:   true,
	}
	repository.TestRepository(t, &repo, &repo, &repo)
}
