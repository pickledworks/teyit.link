package database

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
	"gitlab.com/nod/teyit/link/utils"
	"log"
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
	FailedAt        *time.Time
	DeletedAt       *time.Time
}

type ArchivePublic struct {
	ArchiveID  uuid.UUID         `json:"archive_id"`
	Slug       string            `json:"slug"`
	Meta       ArchivePublicMeta `json:"meta"`
	Screenshot string            `json:"screenshot"`
	RequestUrl string            `json:"request_url"`
	ArchivedAt time.Time         `json:"archived_at"`
}

type ArchivePublicMeta struct {
	Title       string
	Description string
}

type ArchiveSearchParams struct {
	Query      string
	Before     time.Time
	After      time.Time
	RequestUrl string
	CountOnly  bool
}

type CheckPreviousArchivesResponse struct {
	Count          int        `json:"count"`
	LastArchivedAt *time.Time `json:"last_archived_at"`
}

var UrlValidationError = errors.New("invalidurl")

func CreateArchive(requestUrl string) (*Archive, error) {
	if requestUrl == "" {
		return nil, UrlValidationError
	}

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

	go func() {
		now := time.Now()
		result, err := utils.RunArchiveLambda(archive.ArchiveID, archive.RequestUrl)

		if err != nil {
			log.Println("Error", err)
			archive.FailedAt = &now
		} else {
			archive.MetaTitle = result.Title
			archive.MetaDescription = result.Description
			archive.Image = result.Image
			archive.ArchivedAt = &now
		}

		SaveArchive(&archive)
	}()

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

	if err := GetDB().Where("slug = ?", slug).First(&archive).Error; err != nil {
		return nil, ArchiveNotFoundError
	} else {
		return &archive, nil
	}
}

func SaveArchive(archive *Archive) {
	GetDB().Save(&archive)
}

func CountArchivesByRequestUrl(requestUrl string) (CheckPreviousArchivesResponse, error) {
	var archives []Archive
	var last Archive

	db := GetDB().Select("archived_at").Where("request_url = ?", requestUrl)
	db.Find(&archives)

	count := len(archives)
	if count > 0 {
		last = archives[0]
	} else {
		nilTime := time.Time{}
		last = Archive{ArchivedAt: &nilTime}
	}

	return CheckPreviousArchivesResponse{
		count,
		last.ArchivedAt,
	}, nil
}

func GetArchiveAsArchivePublic(archive *Archive) ArchivePublic {
	public := ArchivePublic{
		ArchiveID: archive.ArchiveID,
		Slug: archive.Slug,
		RequestUrl: archive.RequestUrl,
	}

	if archive.ArchivedAt != nil {
		public.Meta = ArchivePublicMeta{archive.MetaTitle, archive.MetaDescription}
		public.ArchivedAt = *archive.ArchivedAt
	}

	return public
}
