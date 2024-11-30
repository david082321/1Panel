package repo

import (
	"context"

	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/global"
	"gorm.io/gorm"
)

type BackupRepo struct{}

type IBackupRepo interface {
	ListRecord(opts ...DBOption) ([]model.BackupRecord, error)
	PageRecord(page, size int, opts ...DBOption) (int64, []model.BackupRecord, error)
	CreateRecord(record *model.BackupRecord) error
	DeleteRecord(ctx context.Context, opts ...DBOption) error
	UpdateRecord(record *model.BackupRecord) error
	WithByDetailName(detailName string) DBOption
	WithByFileName(fileName string) DBOption
	WithByCronID(cronjobID uint) DBOption
	WithFileNameStartWith(filePrefix string) DBOption
}

func NewIBackupRepo() IBackupRepo {
	return &BackupRepo{}
}

func (u *BackupRepo) ListRecord(opts ...DBOption) ([]model.BackupRecord, error) {
	var users []model.BackupRecord
	db := global.DB.Model(&model.BackupRecord{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&users).Error
	return users, err
}

func (u *BackupRepo) PageRecord(page, size int, opts ...DBOption) (int64, []model.BackupRecord, error) {
	var users []model.BackupRecord
	db := global.DB.Model(&model.BackupRecord{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (u *BackupRepo) WithByFileName(fileName string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(fileName) == 0 {
			return g
		}
		return g.Where("file_name = ?", fileName)
	}
}

func (u *BackupRepo) WithByDetailName(detailName string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(detailName) == 0 {
			return g
		}
		return g.Where("detail_name = ?", detailName)
	}
}

func (u *BackupRepo) WithFileNameStartWith(filePrefix string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("file_name LIKE ?", filePrefix+"%")
	}
}

func (u *BackupRepo) CreateRecord(record *model.BackupRecord) error {
	return global.DB.Create(record).Error
}

func (u *BackupRepo) UpdateRecord(record *model.BackupRecord) error {
	return global.DB.Save(record).Error
}

func (u *BackupRepo) DeleteRecord(ctx context.Context, opts ...DBOption) error {
	return getTx(ctx, opts...).Delete(&model.BackupRecord{}).Error
}

func (u *BackupRepo) WithByCronID(cronjobID uint) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("cronjob_id = ?", cronjobID)
	}
}

func (u *BackupRepo) GetRecord(opts ...DBOption) (*model.BackupRecord, error) {
	var record *model.BackupRecord
	db := global.DB.Model(&model.BackupRecord{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(record).Error
	return record, err
}
