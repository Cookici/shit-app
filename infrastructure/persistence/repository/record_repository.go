package repository

import (
	"context"
	"errors"
	"fmt"
	"record-project/domain/entity"
	"record-project/domain/repository"
	"record-project/infrastructure/persistence/model"
	"sort"
	"time"

	"gorm.io/gorm"
)

// recordRepository 记录仓储实现
type recordRepository struct {
	db *gorm.DB
}

// NewRecordRepository 创建记录仓储
func NewRecordRepository(db *gorm.DB) repository.RecordRepository {
	return &recordRepository{db: db}
}

// FindByID 根据ID查找记录
func (r *recordRepository) FindByID(ctx context.Context, id uint64) (*entity.Record, error) {
	var recordModel model.Record
	if err := r.db.WithContext(ctx).First(&recordModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return recordModel.ToEntity(), nil
}

// FindByUserID 根据用户ID查找记录
func (r *recordRepository) FindByUserID(ctx context.Context, userID uint64, page, size int) ([]*entity.Record, int64, error) {
	var recordModels []model.Record
	var total int64

	offset := (page - 1) * size

	if err := r.db.WithContext(ctx).Model(&model.Record{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("record_time DESC").Offset(offset).Limit(size).Find(&recordModels).Error; err != nil {
		return nil, 0, err
	}

	records := make([]*entity.Record, len(recordModels))
	for i, recordModel := range recordModels {
		records[i] = recordModel.ToEntity()
	}

	return records, total, nil
}

// FindByDateRange 根据日期范围查找记录
func (r *recordRepository) FindByDateRange(ctx context.Context, userID uint64, start, end time.Time, page, size int) ([]*entity.Record, error) {
	var recordModels []model.Record

	// 计算分页偏移量
	offset := (page - 1) * size

	// 添加分页参数
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND record_time BETWEEN ? AND ?", userID, start, end).
		Order("record_time DESC").
		Offset(offset).
		Limit(size)

	if err := query.Find(&recordModels).Error; err != nil {
		return nil, err
	}

	records := make([]*entity.Record, len(recordModels))
	for i, recordModel := range recordModels {
		records[i] = recordModel.ToEntity()
	}

	return records, nil
}

func (r *recordRepository) CountRecordsByDateRange(ctx context.Context, userID uint64, start time.Time, end time.Time) (int64, error) {
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Record{}).
		Where("user_id = ? AND record_time BETWEEN ? AND ?", userID, start, end)

	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Save 保存记录
func (r *recordRepository) Save(ctx context.Context, record *entity.Record) error {
	var recordModel model.Record
	recordModel.FromEntity(record)
	if err := r.db.WithContext(ctx).Create(&recordModel).Error; err != nil {
		return err
	}
	record.ID = recordModel.ID
	return nil
}

// Update 更新记录
func (r *recordRepository) Update(ctx context.Context, record *entity.Record) error {
	var recordModel model.Record
	recordModel.FromEntity(record)
	result := r.db.WithContext(ctx).Model(&model.Record{}).Where("id = ?", record.ID).Updates(map[string]interface{}{
		"user_id":      recordModel.UserID,
		"record_time":  recordModel.RecordTime,
		"duration":     recordModel.Duration,
		"poop_type_id": recordModel.PoopTypeID,
		"note":         recordModel.Note,
		"updated_at":   time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete 删除记录
func (r *recordRepository) Delete(ctx context.Context, id uint64) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除关联的标签关系
		if err := tx.Where("record_id = ?", id).Delete(&model.RecordTag{}).Error; err != nil {
			return err
		}

		// 再删除记录本身
		return tx.Delete(&model.Record{}, id).Error
	})
}

// GetRecordWithTags 获取记录及其关联的标签
func (r *recordRepository) GetRecordWithTags(ctx context.Context, id uint64) (*entity.Record, error) {
	var recordModel model.Record
	if err := r.db.WithContext(ctx).First(&recordModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	record := recordModel.ToEntity()

	// 获取关联的标签
	var recordTags []model.RecordTag
	if err := r.db.WithContext(ctx).Where("record_id = ?", id).Find(&recordTags).Error; err != nil {
		return nil, err
	}

	if len(recordTags) > 0 {
		var tagIDs []uint64
		for _, rt := range recordTags {
			tagIDs = append(tagIDs, rt.TagID)
		}

		var tagModels []model.Tag
		if err := r.db.WithContext(ctx).Where("id IN ?", tagIDs).Find(&tagModels).Error; err != nil {
			return nil, err
		}

		tags := make([]*entity.Tag, len(tagModels))
		for i, tm := range tagModels {
			tags[i] = tm.ToEntity()
		}

		record.Tags = tags
	}

	// 获取关联的用户信息
	var userModel model.User
	if err := r.db.WithContext(ctx).First(&userModel, record.UserID).Error; err == nil {
		record.User = userModel.ToEntity()
	}

	// 获取关联的屎的类型信息
	if record.PoopTypeID > 0 {
		var poopTypeModel model.PoopType
		if err := r.db.WithContext(ctx).First(&poopTypeModel, record.PoopTypeID).Error; err == nil {
			record.PoopType = poopTypeModel.ToEntity()
		}
	}

	return record, nil
}

// CreateWithTags 使用事务创建记录并关联标签
func (r *recordRepository) CreateWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64, recordTagRepo repository.RecordTagRepository) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建记录仓储的事务版本
		txRecordRepo := &recordRepository{db: tx}

		// 保存记录
		if err := txRecordRepo.Save(ctx, record); err != nil {
			return err
		}

		// 如果有标签，保存标签关联
		if len(tagIDs) > 0 {
			// 检查recordTagRepo类型
			if _, ok := recordTagRepo.(*recordTagRepository); !ok {
				return fmt.Errorf("无法创建记录标签仓储的事务版本")
			}

			// 使用事务创建新的recordTagRepository实例
			txRecordTagRepo := &recordTagRepository{db: tx}

			// 保存标签关联
			if err := txRecordTagRepo.SaveRecordTags(ctx, record.ID, tagIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateWithTags 使用事务更新记录并关联标签
func (r *recordRepository) UpdateWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64, recordTagRepo repository.RecordTagRepository) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建记录仓储的事务版本
		txRecordRepo := &recordRepository{db: tx}

		// 更新记录
		if err := txRecordRepo.Update(ctx, record); err != nil {
			return err
		}

		// 检查recordTagRepo类型
		if _, ok := recordTagRepo.(*recordTagRepository); !ok {
			return fmt.Errorf("无法创建记录标签仓储的事务版本")
		}

		// 使用事务创建新的recordTagRepository实例
		txRecordTagRepo := &recordTagRepository{db: tx}

		// 保存标签关联
		if err := txRecordTagRepo.SaveRecordTags(ctx, record.ID, tagIDs); err != nil {
			return err
		}

		return nil
	})
}

// GetGlobalRanking 获取全局排行榜（按记录次数排序）
func (r *recordRepository) GetGlobalRanking(ctx context.Context, start, end time.Time, limit int) ([]*entity.RankingItem, error) {
	var rankingItems []*entity.RankingItem

	err := r.db.WithContext(ctx).Model(&model.Record{}).
		Select("user_id, COUNT(*) as record_count, SUM(duration) as total_duration").
		Where("record_time BETWEEN ? AND ?", start, end).
		Group("user_id").
		Order("record_count DESC").
		Limit(limit).
		Scan(&rankingItems).Error
	if err != nil {
		return nil, err
	}

	// 为每个排名项设置排名
	for i, item := range rankingItems {
		item.Rank = uint64(i + 1)
	}

	return rankingItems, nil
}

// GetFriendRanking 获取好友排行榜（按记录次数排序，包含记录为0的用户）
func (r *recordRepository) GetFriendRanking(ctx context.Context, userIDs []uint64, startDate, endDate time.Time, page, pageSize int) ([]*entity.RankingItem, int, error) {
	offset := (page - 1) * pageSize

	// 1. 先获取所有好友在指定时间段内的记录统计
	var recordStats []struct {
		UserID        uint64
		RecordCount   int64
		TotalDuration int64
	}

	err := r.db.WithContext(ctx).Model(&model.Record{}).
		Select("user_id, COUNT(*) as record_count, SUM(duration) as total_duration").
		Where("user_id IN ? AND record_time BETWEEN ? AND ?", userIDs, startDate, endDate).
		Group("user_id").
		Scan(&recordStats).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. 创建用户ID到记录统计的映射
	statsMap := make(map[uint64]struct {
		RecordCount   int64
		TotalDuration int64
	})

	for _, stat := range recordStats {
		statsMap[stat.UserID] = struct {
			RecordCount   int64
			TotalDuration int64
		}{
			RecordCount:   stat.RecordCount,
			TotalDuration: stat.TotalDuration,
		}
	}

	// 3. 为所有用户创建排行项（包括没有记录的用户）
	allRankingItems := make([]*entity.RankingItem, 0, len(userIDs))
	for _, userID := range userIDs {
		item := &entity.RankingItem{
			UserID:        userID,
			RecordCount:   0,
			TotalDuration: 0,
		}

		// 如果用户有记录，更新统计数据
		if stat, exists := statsMap[userID]; exists {
			item.RecordCount = stat.RecordCount
			item.TotalDuration = stat.TotalDuration
		}

		allRankingItems = append(allRankingItems, item)
	}

	// 4. 按记录次数排序
	// 使用自定义排序
	sort.Slice(allRankingItems, func(i, j int) bool {
		return allRankingItems[i].RecordCount > allRankingItems[j].RecordCount
	})

	// 5. 设置排名
	for i := range allRankingItems {
		allRankingItems[i].Rank = uint64(i + 1)
	}

	// 6. 计算总数
	total := len(allRankingItems)

	// 7. 分页
	var pagedItems []*entity.RankingItem
	if offset < len(allRankingItems) {
		end := offset + pageSize
		if end > len(allRankingItems) {
			end = len(allRankingItems)
		}
		pagedItems = allRankingItems[offset:end]
	} else {
		pagedItems = []*entity.RankingItem{}
	}

	return pagedItems, total, nil
}

// GetUsersDailyRecordStats 批量获取指定用户当天的拉屎记录总数和时间
func (r *recordRepository) GetUsersDailyRecordStats(ctx context.Context, userIDs []uint64, date time.Time) (map[uint64]*entity.DailyRecordStats, error) {
	// 设置日期范围为当天的0点到23:59:59
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
	
	// 查询所有符合条件的记录
	var records []model.Record
	if err := r.db.WithContext(ctx).
		Where("user_id IN ? AND record_time BETWEEN ? AND ?", userIDs, startDate, endDate).
		Order("record_time ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	
	// 按用户ID分组统计
	result := make(map[uint64]*entity.DailyRecordStats)
	
	// 初始化每个用户的统计数据
	for _, userID := range userIDs {
		result[userID] = &entity.DailyRecordStats{
			UserID:      userID,
			Date:        startDate,
			Count:       0,
			TotalTime:   0,
			RecordTimes: []time.Time{},
		}
	}
	
	// 统计每个用户的记录
	for _, record := range records {
		stats, exists := result[record.UserID]
		if !exists {
			// 这种情况不应该发生，因为我们已经初始化了所有用户的统计数据
			continue
		}
		
		stats.Count++
		stats.TotalTime += record.Duration
		stats.RecordTimes = append(stats.RecordTimes, record.RecordTime)
	}
	
	return result, nil
}

// 为了兼容接口，保留原来的方法但内部调用新方法
func (r *recordRepository) GetRankingByUserIDs(ctx context.Context, userIDs []uint64, startDate, endDate time.Time, offset, limit int) ([]*entity.RankingItem, int, error) {
	page := offset/limit + 1
	return r.GetFriendRanking(ctx, userIDs, startDate, endDate, page, limit)
}
