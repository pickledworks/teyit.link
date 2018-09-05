package database

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
	"gitlab.com/nod/teyitlink-web/utils"
	"time"
)

type Archive struct {
	ArchiveID       uuid.UUID `gorm:"primary_key"`
	Slug            string    `gorm:"unique_index"`
	MetaTitle       string
	MetaDescription string `gorm:"size:2048"`
	Image           string
	RequestUrl      string `gorm:"size:2048",valid:"url"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ArchivedAt      *time.Time
	DeletedAt       *time.Time
}

type ArchiveSearchParams struct {
	Query      string
	Before     time.Time
	After      time.Time
	RequestUrl string
	CountOnly  bool
}

var UrlValidationError = errors.New("invalidurl")

func CreateArchive(requestUrl string) (*Archive, error) {
	archiveId, _ := uuid.NewV4()

	archive := Archive{
		ArchiveID:  archiveId,
		RequestUrl: requestUrl,
		Slug:       utils.RandString(7),
	}

	_, err := govalidator.ValidateStruct(archive)
	if err != nil {
		return nil, UrlValidationError
	}

	db := GetDB()
	db.NewRecord(&archive)
	db.Create(&archive)

	return &archive, nil
}

func FindArchives(params ArchiveSearchParams) ([]Archive, error) {
	var archives []Archive
	db := GetDB()

	if params.Query != "" {
		query := fmt.Sprintf("%%%s%%", params.Query)
		db = db.Where("request_url LIKE ?", query)
		db = db.Or("meta_title LIKE ?", query).Or("meta_description LIKE ?", query)
	}

	if params.RequestUrl != "" {
		db = db.Where("request_url = ?", params.RequestUrl)
	}

	emptyTime := time.Time{}
	if params.Before != emptyTime {
		db = db.Where("archived_at < ?", params.Before)
	}

	if params.After != emptyTime {
		db = db.Where("archived_at > ?", params.After)
	}

	db.Where("archived_at is not null").Find(&archives)

	return archives, nil
}

var ArchiveNotFoundError = errors.New("archivenotfound")

func GetArchive(slug string) (*Archive, error) {
	var archive Archive
	db := GetDB()

	if err := db.Where("slug = ?", slug).First(&archive).Error; err != nil {
		return nil, ArchiveNotFoundError
	} else {
		return &archive, nil
	}
}
