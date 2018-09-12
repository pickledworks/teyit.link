package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
	"gitlab.com/nod/teyit/link/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//go:generate goqueryset -in archive.go

// Archive model is the core model
// gen:qs
type Archive struct {
	ID              uint      `gorm:"primary_key"`
	RequestUrl      string    `gorm:"size:2048",valid:"requrl,required"`
	CallbackUrl     string    `gorm:"-",valid:"url",valid:"requrl,required"`
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

func (a *Archive) GetAsPublic() ArchivePublic {
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

var UrlValidationError = errors.New("invalid request_url or callback_url")
var SelfArchiveError = errors.New("teyit.link urls can't be archived")

func CreateArchive(archive *Archive) (*Archive, error) {
	// We only validate the URL's
	isValidArchive, err := govalidator.ValidateStruct(&archive)
	if err != nil || isValidArchive != true {
		return nil, UrlValidationError
	}
	log.Println(isValidArchive, err)

	// URL's are already validated so we can ignore the error here
	url, _ := url.Parse(archive.RequestUrl)
	if strings.Contains(url.Host, "teyit.link") {
		return nil, SelfArchiveError
	}

	archive.ArchiveID = uuid.NewV4()
	archive.Slug = GenerateArchiveSlug()

	db := GetDB()
	db.NewRecord(&archive)
	db.Create(&archive)

	go func() {
		now := time.Now()
		result, err := utils.RunArchiveLambda(archive.ArchiveID, archive.RequestUrl)

		if err != nil {
			log.Println("Error running lambda", err)
			archive.FailedAt = &now
		} else {
			archive.MetaTitle = result.Title
			archive.MetaDescription = result.Description
			archive.MetaImage = result.Image
			archive.ArchivedAt = &now
		}

		db.Save(&archive)

		if archive.CallbackUrl != "" {
			body, err := json.Marshal(archive.GetAsPublic())
			req, err := http.NewRequest("POST", archive.CallbackUrl, bytes.NewBuffer(body))
			req.Header.Set("X-Teyit-Link-Version", "v2.0.0")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("error when triggering callback url", err)
			} else {
				defer resp.Body.Close()
			}
		}
	}()

	return archive, nil
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

var ArchiveNotFoundError = errors.New("archive not found")

func GetArchive(slug string) (*Archive, error) {
	var archive Archive

	if err := GetDB().Where("slug = ?", slug).First(&archive).Error; err != nil {
		return nil, ArchiveNotFoundError
	} else {
		return &archive, nil
	}
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
