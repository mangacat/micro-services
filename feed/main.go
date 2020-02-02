package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/machinebox/graphql"
)

type SeriesChapters struct {
	SeriesChapters []struct {
		ID                   int    `json:"id"`
		Hash                 string `json:"hash"`
		Language             string `json:"language"`
		SeriesChaptersSeries struct {
			Name       string `json:"name"`
			ID         int    `json:"id"`
			CoverImage string `json:"cover_image"`
		} `json:"series_chapters_series"`
		ChapterNumberVolume   string    `json:"chapter_number_volume"`
		ChapterNumberAbsolute string    `json:"chapter_number_absolute"`
		VolumeNumber          string    `json:"volume_number"`
		TimeUploaded          time.Time `json:"time_uploaded"`
		Title                 string    `json:"title"`
	} `json:"series_chapters"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("helloworld: received a request")
	if r.Method == "GET" {
		// create a client (safe to share across requests)
		client := graphql.NewClient(os.Getenv("HASURA_URL"))

		// make a request
		req := graphql.NewRequest(`
			query {
				series_chapters(limit: 25, order_by: {time_uploaded: desc}, where: {delay: {}, published: {_eq: true}}) {
					id
					hash
					language
					series_chapters_series {
						name
						id
						cover_image
					}
					chapter_number_volume
					chapter_number_absolute
					volume_number
					time_uploaded
					title
				}
			}
		`)

		// set any variables
		// req.Var("key", "value")

		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))

		// define a Context for the request
		ctx := context.Background()

		// run it and capture the response
		var respData SeriesChapters
		if err := client.Run(ctx, req, &respData); err != nil {
			log.Fatal(err)
		}
		fmt.Println(respData)
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Manga Cat Recent Chapters",
			Link:        &feeds.Link{Href: "https://rss."},
			Description: "Manga Cat Recent Chapters",
			Created:     now,
		}

		for _, chapter := range respData.SeriesChapters {
			fmt.Println(chapter)
			var title string
			if chapter.SeriesChaptersSeries.Name != "" {
				title = title + fmt.Sprintf("%s ", chapter.SeriesChaptersSeries.Name)
			}
			if chapter.VolumeNumber != "" {
				title = title + fmt.Sprintf("Vol %s ", chapter.VolumeNumber)
			}
			if chapter.ChapterNumberVolume != "" {
				title = title + fmt.Sprintf("Ch %s ", chapter.ChapterNumberVolume)
			}
			if chapter.Title != "" {
				title = title + fmt.Sprintf("- %s", chapter.Title)
			}
			fmt.Println(title)

			item := &feeds.Item{
				Title:       title,
				Description: title,
				Link:        &feeds.Link{Href: fmt.Sprintf("%s/read/%s", os.Getenv("MANGA_FRONTEND_URL"), chapter.Hash)},
				Created:     chapter.TimeUploaded,
			}
			feed.Items = append(feed.Items, item)
		}

		keys, ok := r.URL.Query()["type"]

		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}

		switch keys[0] {
		case "atom":
			if err := feed.WriteAtom(w); err != nil {
				log.Println(err)
			}
		case "json":
			if err := feed.WriteJSON(w); err != nil {
				log.Println(err)
			}
		case "rss":
		default:
			if err := feed.WriteRss(w); err != nil {
				log.Println(err)
			}
		}

	}
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/chapters", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
