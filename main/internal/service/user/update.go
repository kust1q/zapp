package user

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func (s *service) Update(ctx context.Context, req *entity.UpdateBio) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.db.UpdateUserBio(ctx, req.UserID, req.Bio)
}
