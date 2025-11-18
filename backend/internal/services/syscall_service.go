package services

import (
	// 컨텍스트 추가
	"admin_server/backend/internal/config"
	"admin_server/backend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const SyscallSetKey = "cluster_callable_syscalls"

// SyscallService handles syscall-related operations
type SyscallService struct {
	cfg        *config.Config
	ccslClient *redis.Client // Redis 클라이언트 필드 추가
}

func NewSyscallService(cfg *config.Config, ccslClient *redis.Client) *SyscallService {
	return &SyscallService{
		cfg:        cfg,
		ccslClient: ccslClient,
	}
}

// GetCallableSyscalls retrieves all syscalls that the cluster can call
func (s *SyscallService) GetCallableSyscalls() (*models.CallableSyscallsResponse, error) {
	log.Println("Getting callable syscalls from CCSL Redis") // 로그 수정

	ctx := context.Background()

	// Redis Set에서 모든 멤버(시스템콜 이름)를 가져옵니다.
	syscalls, err := s.ccslClient.SMembers(ctx, SyscallSetKey).Result()
	if err != nil {
		log.Printf("ERROR: Failed to retrieve syscalls from Redis: %v", err)
		return nil, fmt.Errorf("failed to retrieve syscalls from Redis: %w", err)
	}

	// Redis에서 가져온 문자열 목록을 응답 모델로 변환합니다.
	resultSyscalls := make([]models.Syscall, len(syscalls))
	for i, name := range syscalls {
		resultSyscalls[i] = models.Syscall{
			Name: name,
			// 인자와 설명은 'Source of Truth' Redis에서 가져오거나
			// 정적 분석 단계에서 채워져야 하지만, 현재는 TBD로 처리합니다.
			Args:        []models.SyscallArg{},
			Description: "TBD",
		}
	}

	return &models.CallableSyscallsResponse{
		TotalCount: len(resultSyscalls),
		Syscalls:   resultSyscalls,
	}, nil
}
