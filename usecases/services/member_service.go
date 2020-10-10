package services

import (
	"context"

	"book-management-system/entities/models"
	"book-management-system/repositories"
	"book-management-system/repositories/mysql"
)

// MemberService handle business logic related to book
type MemberService interface {
	GetMembers(context.Context) (models.Members, error)
	CreateMember(context.Context, *models.Member) error
	UpdateMember(context.Context, *models.Member) error
}

type memberService struct {
	MySQLMemberRepository mysql.MemberRepository
}

// NewMemberService returns MemberService
func NewMemberService(repo *repositories.Repository) MemberService {
	return &memberService{
		MySQLMemberRepository: repo.MySQLMemberRepository,
	}
}

func (svc *memberService) GetMembers(ctx context.Context) (models.Members, error) {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	return svc.MySQLMemberRepository.GetAll(ctx)
}

func (svc *memberService) CreateMember(ctx context.Context, member *models.Member) error {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	err := svc.MySQLMemberRepository.CreateMember(ctx, member)
	if err != nil {
		return err
	}

	return nil
}

func (svc *memberService) UpdateMember(ctx context.Context, member *models.Member) error {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	err := svc.MySQLMemberRepository.UpdateMember(ctx, member)
	if err != nil {
		return err
	}

	return nil
}
