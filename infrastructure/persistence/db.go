package persistence

import (
	"fmt"
	"log"
	"record-project/infrastructure/config"
	"record-project/infrastructure/persistence/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Database 数据库连接
type Database struct {
	DB *gorm.DB
}

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.Config) (*Database, error) {
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.DB.GetDSN()), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.DB.MaxLifetime)

	return &Database{DB: db}, nil
}

// InitData 初始化数据
func (d *Database) InitData() error {
	// 初始化屎的类型数据
	return d.initPoopTypes()
}

// initPoopTypes 初始化屎的类型数据
func (d *Database) initPoopTypes() error {
	// 检查是否已有数据
	var count int64
	d.DB.Model(&model.PoopType{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 布里斯托尔大便分类法的7种类型
	poopTypes := []model.PoopType{
		{
			Name:             "第一型",
			Description:      "分离的硬块，像坚果一样（难以排出）",
			Color:            "深棕色",
			HealthIndication: "严重便秘",
		},
		{
			Name:             "第二型",
			Description:      "香肠状但结块的",
			Color:            "棕色",
			HealthIndication: "轻微便秘",
		},
		{
			Name:             "第三型",
			Description:      "像香肠但有裂缝的表面",
			Color:            "棕色",
			HealthIndication: "正常",
		},
		{
			Name:             "第四型",
			Description:      "像香肠或蛇一样，光滑而柔软",
			Color:            "棕色",
			HealthIndication: "正常理想状态",
		},
		{
			Name:             "第五型",
			Description:      "有清晰边缘的软块",
			Color:            "棕色",
			HealthIndication: "轻微腹泻",
		},
		{
			Name:             "第六型",
			Description:      "边缘模糊的松软块，糊状大便",
			Color:            "浅棕色",
			HealthIndication: "腹泻",
		},
		{
			Name:             "第七型",
			Description:      "水样，没有固体，完全液体",
			Color:            "黄色或浅棕色",
			HealthIndication: "严重腹泻",
		},
	}

	for _, poopType := range poopTypes {
		if err := d.DB.Create(&poopType).Error; err != nil {
			return fmt.Errorf("初始化屎的类型数据失败: %w", err)
		}
	}

	log.Println("初始化屎的类型数据完成")
	return nil
}
