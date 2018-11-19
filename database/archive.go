package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"gitlab.com/nod/teyit/link/utils"
	"log"
	"net/http"
	"time"
)

//go:generate goqueryset -in archive.go

// Archive model is the core model
// gen:qs
type Archive struct {
	ID              uint      `gorm:"primary_key"`
	RequestUrl      string    `gorm:"size:2048"`
	CallbackUrl     string    `gorm:"-"`
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
type ArchivePublic struct {
	Slug       string             `json:"slug"`
	RequestUrl string             `json:"request_url"`
	Meta       ArchivePublicMeta  `json:"meta"`
	ArchivedAt utils.NullableTime `json:"archived_at"`
}

type ArchivePublicMeta struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
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

func (a *Archive) Validate() error {
	var err error

	err = utils.ValidateUrl(a.RequestUrl)
	if err != nil {
		return err
	}

	callbackUrl := a.CallbackUrl
	if callbackUrl != "" {
		err = utils.ValidateUrl(callbackUrl)
	}
	if err != nil {
		return err
	}

	return nil
}

func (a *Archive) Save(skipValidation bool) error {
	if skipValidation != true {
		err := a.Validate()
		if err != nil {
			return err
		}
	}

	a.ArchiveID = uuid.NewV4()
	a.Slug = generateArchiveSlug()

	db := GetDB()
	db.NewRecord(&a)
	db.Create(&a)

	go a.RunLambda()
	return nil
}

func (a *Archive) RunLambda() {
	result, err := utils.RunArchiveLambda(a.ArchiveID, a.RequestUrl)
	now := time.Now()

	if err != nil {
		log.Println("Error running lambda", err)
		a.FailedAt = &now
	} else {
		a.MetaTitle = result.Title
		a.MetaDescription = result.Description
		a.MetaImage = result.Image
		a.ArchivedAt = &now
	}

	GetDB().Save(&a)

	if a.CallbackUrl != "" {
		body, err := json.Marshal(a.GetAsPublic())
		req, err := http.NewRequest("POST", a.CallbackUrl, bytes.NewBuffer(body))
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
}

type ArchiveSearchParams struct {
	Query      string
	Before     time.Time
	After      time.Time
	Limit      int
	Offset     int
	RequestUrl string
	CountOnly  bool
}

type ArchiveSearchResults struct {
	Results         []Archive
	Total           int
	PreviousPageUrl string
	NextPageUrl     string
}

type CheckPreviousArchivesResponse struct {
	Count       int     `json:"count"`
	LastArchive Archive `json:"last_archived_at"`
}

func generateArchiveSlug() string {
	slug := utils.RandString(7)

	var count int
	GetDB().Model(Archive{}).Where("slug = ?", slug).Limit(1).Count(&count)

	if count > 0 {
		return generateArchiveSlug()
	}

	return slug
}

func FindArchives(params ArchiveSearchParams) ([]Archive, int, error) {
	var archives []Archive
	var total int
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

	db.Where("archived_at is not null").Order("created_at desc")
	db.Model(&Archive{}).Count(&total)
	db.Limit(params.Limit).Offset(params.Offset).Find(&archives)

	return archives, total, nil
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
