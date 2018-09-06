package main

import (
	"gitlab.com/nod/teyit/link/database"
	"log"
)

func main() {
	db := database.InitDB("mysql", "")
	database.Migrate(db)
	defer db.Close()

	resp, err := CountArchivesByRequestUrl("http://nod.digital")
	if err != nil {
		log.Println("hey")
	} else {
		log.Println(resp)
	}

}

func CountArchivesByRequestUrl(requestUrl string) (database.CheckPreviousArchivesResponse, error) {
	var archives []database.Archive

	db := database.GetDB()
	db.Where("request_url = ?", requestUrl).Find(&archives)

	last := archives[0]

	log.Println("last", last)
	log.Println("archives", archives)
	return database.CheckPreviousArchivesResponse{
		len(archives),
		last.ArchivedAt,
	}, nil
}
