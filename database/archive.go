package database

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
	"gitlab.com/nod/teyit/link/utils"
	"log"
	"net/url"
	"strings"
	"time"
)

//go:generate goqueryset -in archive.go

// Archive model is the core model
// gen:qs
type Archive struct {
	ID              uint      `gorm:"primary_key"`
	RequestUrl      string    `gorm:"size:2048",valid:"url"`
	ArchiveID       uuid.UUID `gorm:"unique_index"`
	Slug            string    `gorm:"unique_index"`
	MetaTitle       string
	MetaDescription string `gorm:"size:2048"`
	MetaImage       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ArchivedAt      *time.Time
	FailedAt        *time.Time
	DeletedAt       *time.Time `sql:"index"`
}


func (a Archive) GetAsPublic() ArchivePublic {
	public := ArchivePublic{
		Slug:       a.Slug,
		RequestUrl: a.RequestUrl,
	}

	if a.ArchivedAt != nil {
		public.Meta = ArchivePublicMeta{a.MetaTitle, a.MetaDescription, a.MetaImage}
		public.ArchivedAt = utils.NullableTime{Time: *a.ArchivedAt}
	}

	return public
}


type ArchivePublic struct {
	Slug       string             `json:"slug"`
	RequestUrl string             `json:"request_url"`
	Meta       ArchivePublicMeta  `json:"meta"`
	ArchivedAt utils.NullableTime `json:"archived_at"`
}

type ArchivePublicMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ArchiveSearchParams struct {
	Query      string
	Before     time.Time
	After      time.Time
	RequestUrl string
	CountOnly  bool
}

type CheckPreviousArchivesResponse struct {
	Count       int     `json:"count"`
	LastArchive Archive `json:"last_archived_at"`
}

var UrlValidationError = errors.New("invalidurl")

func CreateArchive(requestUrl string) (*Archive, error) {
	if requestUrl == "" {
		return nil, UrlValidationError
	}

	url, err := url.Parse(requestUrl)

	if err != nil {
		return nil, UrlValidationError
	}

	if strings.Contains(url.Host, "teyit.link") {
		return nil, UrlValidationError
	}

	archiveId := uuid.NewV4()

	archive := Archive{
		ArchiveID:  archiveId,
		RequestUrl: requestUrl,
		Slug:       GenerateArchiveSlug(),
	}

	isValidRequestUrl, err := govalidator.ValidateStruct(archive)
	if err != nil {
		return nil, UrlValidationError
	}

	if isValidRequestUrl != true {
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
			archive.MetaImage = result.Image
			archive.ArchivedAt = &now
		}

		SaveArchive(&archive)
	}()

	return &archive, nil
}

func GenerateArchiveSlug() string {
	slug := utils.RandString(7)

	var count int
	GetDB().Model(Archive{}).Where("slug = ?", slug).Limit(1).Count(&count)

	if count > 0 {
		return GenerateArchiveSlug()
	}

	return slug
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

	db.Where("archived_at is not null").Order("created_at desc").Find(&archives)

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
	var archives []*Archive

	db := GetDB().Where("request_url = ?", requestUrl).Order("created_at desc")
	db.Find(&archives)

	count := len(archives)

	resp := CheckPreviousArchivesResponse{Count: count}

	if count > 0 {
		resp.LastArchive = *archives[0]
	}

	return resp, nil
}

