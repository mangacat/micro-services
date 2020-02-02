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
		client := graphql.NewClient("http://qkg7t3nssg.lb.c1.gra.k8s.ovh.net/v1/graphql")

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
		req.Header.Set("x-hasura-admin-secret", "")

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
			Link:        &feeds.Link{Href: "https://rss.manga.cat/v1/rss"},
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
				Link:        &feeds.Link{Href: fmt.Sprintf("https://manga.cat/read/%s", chapter.Hash)},
				Created:     chapter.TimeUploaded,
			}
			feed.Items = append(feed.Items, item)
		}
		// switch c.GetString("type") {
		// case "atom":
		// 	if err := feed.WriteAtom(w); err != nil {
		// 		panic(err)
		// 	}
		// case "json":
		// 	if err := feed.WriteJSON(w); err != nil {
		// 		panic(err)
		// 	}
		// case "rss":
		// default:
		if err := feed.WriteAtom(w); err != nil {
			panic(err)
		}
		// }

	}
	// target := os.Getenv("TARGET")
	// if target == "" {
	// 	target = "World"
	// }

	// // fmt.Fprintf(w, "Hello %s!\n", target)
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
