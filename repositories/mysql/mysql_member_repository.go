package mysql

import (
	"context"

	"gorm.io/gorm"

	"book-management-system/entities/models"
)

// MemberRepository handle sql query to members table
type MemberRepository interface {
	GetAll(context.Context) (models.Members, error)
	CreateMember(context.Context, *models.Member) error
	UpdateMember(context.Context, *models.Member) error
}

type memberRepository struct {
	db *gorm.DB
}

// NewMemberRepository returns new MemberRepository
func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{
		db: db,
	}
}

func (repo *memberRepository) GetAll(ctx context.Context) (models.Members, error) {
	var members models.Members

	query := repo.db.WithContext(ctx).
		Find(&members)
	return members, query.Error
}

func (repo *memberRepository) CreateMember(ctx context.Context, member *models.Member) error {
	query := repo.db.WithContext(ctx).
		Create(member)
	return query.Error
}

func (repo *memberRepository) UpdateMember(ctx context.Context, member *models.Member) error {
	query := repo.db.WithContext(ctx).
		Updates(member)
	return query.Error
}
