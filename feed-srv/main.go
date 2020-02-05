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

type SeriesChapter struct {
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
}
type SeriesChapters struct {
	SeriesChapters []SeriesChapter `json:"series_chapters"`
}

type PrivateSeriesChapters struct {
	UsersFollowingSeries []struct {
		SeriesChapter SeriesChapter `json:"users_following_series_chapters"`
	} `json:"users_following_series"`
}

const (
	RSS_PUBLIC_GRAPHQL string = `
query SeriesChapters {
  series_chapters(limit: 25, order_by: {time_uploaded: desc}, where: {published: {_eq: true}}) {
    id
    hash
    language
    chapter_number_volume
    chapter_number_absolute
    volume_number
    time_uploaded
    title
    series_chapters_series {
      id
      cover_image
      name
    }
  }
}
	`
	RSS_PRIVATE_GRAPHQL string = `
query($token: uuid!) {
  users_following_series(order_by: {users_following_series_chapters: {time_uploaded: desc}}, limit: 10, where: {users_following_series_chapters: {published: {_eq: true}}, users_following_series_user: {rss_token: {_eq: $token}}}) {
    users_following_series_chapters {
      id
      hash
      language
      chapter_number_volume
      chapter_number_absolute
      volume_number
      time_uploaded
      title
      series_chapters_series {
        id
        name
        cover_image
      }
    }
  }
}


	`
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("feed-srv: received a request")
	if r.Method == "GET" {
		var token string = ""
		tokens, ok := r.URL.Query()["token"]
		if ok || len(tokens) != 0 {
			token = tokens[0]
			fmt.Println(token)
		}

		// create a client (safe to share across requests)
		client := graphql.NewClient(os.Getenv("HASURA_URL"))

		// var graphql.Req
		var req *graphql.Request
		var respData SeriesChapters

		// make a request
		fmt.Println(RSS_PUBLIC_GRAPHQL)
		if token != "" {
			fmt.Println("uwu")
			req = graphql.NewRequest(RSS_PRIVATE_GRAPHQL)
			req.Var("token", token)
			req.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))

			// define a Context for the request
			ctx := context.Background()
			var resp PrivateSeriesChapters
			// run it and capture the response
			if err := client.Run(ctx, req, &resp); err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp)
			for _, k := range resp.UsersFollowingSeries {
				respData.SeriesChapters = append(respData.SeriesChapters, k.SeriesChapter)
			}
		} else {
			req = graphql.NewRequest(RSS_PUBLIC_GRAPHQL)
			req.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))

			// define a Context for the request
			ctx := context.Background()

			// run it and capture the response
			if err := client.Run(ctx, req, &respData); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println(respData)
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Manga Cat Recent Chapters",
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/rss.xml", os.Getenv("MANGA_FRONTEND_URL"))},
			Description: "Manga Cat Recent Chapters",
			Created:     now,
		}

		for _, chapter := range respData.SeriesChapters {
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
			keys = append(keys, "rss")
			// return
		}
		fmt.Println(keys)

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
			if err := feed.WriteRss(w); err != nil {
				log.Println(err)
			}
		}

	}
}

func main() {
	log.Print("feed-srv: starting server...")

	http.HandleFunc("/feed.xml", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
