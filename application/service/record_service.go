package service

import (
	"context"
	"record-project/domain/entity"
	"record-project/domain/repository"
	"time"
)

// RecordService 记录服务接口
type RecordService interface {
	GetRecordByID(ctx context.Context, id uint64) (*entity.Record, error)
	GetRecordsByUserID(ctx context.Context, userID uint64, page, size int) ([]*entity.Record, int64, error)
	GetRecordsByDateRange(ctx context.Context, userID uint64, start, end time.Time, page, size int) ([]*entity.Record, error)
	CreateRecord(ctx context.Context, record *entity.Record) error
	UpdateRecord(ctx context.Context, record *entity.Record) error
	DeleteRecord(ctx context.Context, id uint64) error
	SaveRecordTags(ctx context.Context, recordID uint64, tagIDs []uint64) error
	GetRecordTags(ctx context.Context, recordID uint64) ([]*entity.Tag, error)
	CreateRecordWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64) error
	UpdateRecordWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64) error
	CountRecordsByDateRange(ctx context.Context, userID uint64, start, end time.Time) (int64, error)
	GetGlobalRanking(ctx context.Context, start, end time.Time, limit int) ([]*entity.RankingItem, error)
	GetFriendRanking(ctx context.Context, userIDs []uint64, startDate, endDate time.Time, page, pageSize int) ([]*entity.RankingItem, int, error)
	GetUsersDailyRecordStats(ctx context.Context, userIDs []uint64, date time.Time) (map[uint64]*entity.DailyRecordStats, error)
}

// recordService 记录服务实现
type recordService struct {
	recordRepo    repository.RecordRepository
	recordTagRepo repository.RecordTagRepository
}

// NewRecordService 创建记录服务
func NewRecordService(
	recordRepo repository.RecordRepository,
	recordTagRepo repository.RecordTagRepository,
) RecordService {
	return &recordService{
		recordRepo:    recordRepo,
		recordTagRepo: recordTagRepo,
	}
}

// GetRecordByID 根据ID获取记录
func (s *recordService) GetRecordByID(ctx context.Context, id uint64) (*entity.Record, error) {
	return s.recordRepo.FindByID(ctx, id)
}

// GetRecordsByUserID 根据用户ID获取记录列表
func (s *recordService) GetRecordsByUserID(ctx context.Context, userID uint64, page, size int) ([]*entity.Record, int64, error) {
	return s.recordRepo.FindByUserID(ctx, userID, page, size)
}

// GetRecordsByDateRange 根据日期范围获取记录
func (s *recordService) GetRecordsByDateRange(ctx context.Context, userID uint64, start, end time.Time, page, size int) ([]*entity.Record, error) {
	return s.recordRepo.FindByDateRange(ctx, userID, start, end, page, size)
}

// CreateRecord 创建记录
func (s *recordService) CreateRecord(ctx context.Context, record *entity.Record) error {
	return s.recordRepo.Save(ctx, record)
}

// UpdateRecord 更新记录
func (s *recordService) UpdateRecord(ctx context.Context, record *entity.Record) error {
	return s.recordRepo.Update(ctx, record)
}

// DeleteRecord 删除记录
func (s *recordService) DeleteRecord(ctx context.Context, id uint64) error {
	return s.recordRepo.Delete(ctx, id)
}

// SaveRecordTags 保存记录标签关联
func (s *recordService) SaveRecordTags(ctx context.Context, recordID uint64, tagIDs []uint64) error {
	return s.recordTagRepo.SaveRecordTags(ctx, recordID, tagIDs)
}

// GetRecordTags 获取记录关联的标签
func (s *recordService) GetRecordTags(ctx context.Context, recordID uint64) ([]*entity.Tag, error) {
	return s.recordTagRepo.FindTagsByRecordID(ctx, recordID)
}

func (s *recordService) CreateRecordWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64) error {
	return s.recordRepo.CreateWithTags(ctx, record, tagIDs, s.recordTagRepo)
}

func (s *recordService) UpdateRecordWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64) error {
	return s.recordRepo.UpdateWithTags(ctx, record, tagIDs, s.recordTagRepo)
}

func (s *recordService) CountRecordsByDateRange(ctx context.Context, userID uint64, start, end time.Time) (int64, error) {
	return s.recordRepo.CountRecordsByDateRange(ctx, userID, start, end)
}

// GetGlobalRanking 获取全局排行榜
func (s *recordService) GetGlobalRanking(ctx context.Context, start, end time.Time, limit int) ([]*entity.RankingItem, error) {
	return s.recordRepo.GetGlobalRanking(ctx, start, end, limit)
}

// GetFriendRanking 获取好友排行榜数据
func (s *recordService) GetFriendRanking(ctx context.Context, userIDs []uint64, startDate, endDate time.Time, page, pageSize int) ([]*entity.RankingItem, int, error) {
	return s.recordRepo.GetFriendRanking(ctx, userIDs, startDate, endDate, page, pageSize)
}

// GetUsersDailyRecordStats 批量获取指定用户当天的拉屎记录统计
func (s *recordService) GetUsersDailyRecordStats(ctx context.Context, userIDs []uint64, date time.Time) (map[uint64]*entity.DailyRecordStats, error) {
	return s.recordRepo.GetUsersDailyRecordStats(ctx, userIDs, date)
}
