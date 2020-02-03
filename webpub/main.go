package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Machiel/slugify"
	"github.com/go-http-utils/etag"
	"github.com/gorilla/mux"
	"github.com/machinebox/graphql"

	"github.com/mangacat/micro-services/utils/langs"
	"github.com/mangacat/micro-services/utils/truncate"
	"github.com/mangacat/micro-services/utils/webpub"
)

type SeriesChapters struct {
	SeriesChapters []struct {
		ID                   int    `json:"id"`
		Hash                 string `json:"hash"`
		Language             string `json:"language"`
		SeriesChaptersSeries struct {
			Name        string `json:"name"`
			ID          int    `json:"id"`
			CoverImage  string `json:"cover_image"`
			Description string `json:"description"`
			TagsSeries  []struct {
				TagsSeries struct {
					TagName      string `json:"tag_name"`
					TagNamespace string `json:"tag_namespace"`
				} `json:"tags_series"`
			} `json:"tags_series"`
			PeopleSeries []struct {
				PeopleSeries struct {
					ID               int    `json:"id"`
					AlternativeNames string `json:"alternative_names"`
					Name             string `json:"name"`
				} `json:"people_series"`
			} `json:"people_series"`
			Direction string `json:"direction"`
		} `json:"series_chapters_series"`
		ChapterNumberVolume   string    `json:"chapter_number_volume"`
		ChapterNumberAbsolute string    `json:"chapter_number_absolute"`
		VolumeNumber          string    `json:"volume_number"`
		TimeUploaded          time.Time `json:"time_uploaded"`
		Title                 string    `json:"title"`
		SeriesChaptersFiles   []struct {
			Hash      string    `json:"hash"`
			Extension string    `json:"extension"`
			Batoto    bool      `json:"batoto"`
			Type      string    `json:"type"`
			UUID      string    `json:"uuid"`
			Width     int       `json:"width"`
			Height    int       `json:"height"`
			ID        int       `json:"id"`
			Created   time.Time `json:"created"`
			Order     int       `json:"order"`
		} `json:"series_chapters_files"`
		GroupsSeriesChapters []struct {
			GroupsScanlationSeriesChaptersGroups struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"groups_scanlation_series_chapters_groups"`
		} `json:"groups_series_chapters"`
		Published bool `json:"published"`
	} `json:"series_chapters"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var metadata *webpub.WebPubMetadata
	var authors []*webpub.SchemaThing
	var artists []*webpub.SchemaThing
	var translators []*webpub.SchemaThing
	var genres []*webpub.WebPubSubject
	client := graphql.NewClient(os.Getenv("HASURA_URL"))

	// make a request
	req := graphql.NewRequest(`
			query($key: Int!,) {
				series_chapters(limit: 1, where: {id: {_eq: $key}, published: {_eq: true}}) {
					id
					hash
					language
					series_chapters_series {
						name
						id
						cover_image
						description
						tags_series {
							tags_series {
							tag_name
							tag_namespace
							}
						}
						people_series {
							people_series {
							id
							alternative_names
							name
							}
						}
						direction
					}
					chapter_number_volume
					chapter_number_absolute
					volume_number
					time_uploaded
					title
					series_chapters_files(order_by: {order: asc}) {
					hash
					extension
					batoto
					type
					uuid
					width
					height
					id
					created
					order
					}
					groups_series_chapters {
						groups_scanlation_series_chapters_groups {
							name
							id
						}
					}
					published
				}
			}
		`)

	vars := mux.Vars(r)
	hash := vars["hash"]
	hashInt, err := strconv.Atoi(hash)
	// set any variables
	req.Var("key", hashInt)

	req.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData SeriesChapters
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	if len(respData.SeriesChapters) == 0 {
		log.Fatal("not found")
		return
	}
	chapter := respData.SeriesChapters[0]

	if len(chapter.SeriesChaptersSeries.TagsSeries) != 0 {

		for _, l := range chapter.SeriesChaptersSeries.TagsSeries {
			// fmt.Println(l)
			k := &webpub.WebPubSubject{
				Name: l.TagsSeries.TagName,
				Code: l.TagsSeries.TagNamespace,
			}
			genres = append(genres, k)
		}
	}
	if len(chapter.SeriesChaptersSeries.PeopleSeries) != 0 {
		for _, l := range chapter.SeriesChaptersSeries.PeopleSeries {
			k := &webpub.SchemaThing{
				SchemaType: "Person",
				Name:       l.PeopleSeries.Name,
				Alt:        l.PeopleSeries.AlternativeNames,
				Identifier: string(l.PeopleSeries.ID),
				Url:        "/people/" + string(chapter.SeriesChaptersSeries.ID),
			}
			authors = append(authors, k)

		}
	}

	if len(chapter.GroupsSeriesChapters) != 0 {
		for _, l := range chapter.GroupsSeriesChapters {
			k := &webpub.SchemaThing{
				Name:       l.GroupsScanlationSeriesChaptersGroups.Name,
				Identifier: "/groups/" + string(l.GroupsScanlationSeriesChaptersGroups.ID),
			}
			translators = append(translators, k)

		}
	}
	issue, err := strconv.ParseFloat(chapter.ChapterNumberAbsolute, 64)
	if err != nil {
		log.Println(err)
		issue = 0
	}
	var title string
	if chapter.VolumeNumber != "" {
		title = title + fmt.Sprintf("Vol %s ", chapter.VolumeNumber)
	}
	if chapter.ChapterNumberVolume != "" {
		title = title + fmt.Sprintf("Ch %s ", chapter.ChapterNumberVolume)
	}
	if chapter.Title != "" && chapter.ChapterNumberVolume == "" && chapter.VolumeNumber == "" {

		title = title + fmt.Sprintf("%s", chapter.Title)
	}
	if chapter.Title != "" {
		title = title + fmt.Sprintf("- %s", chapter.Title)
	}
	if title == "" {
		title = chapter.SeriesChaptersSeries.Name
	}

	description := chapter.SeriesChaptersSeries.Description
	if len(chapter.SeriesChaptersSeries.Description) > 200 {
		description = truncate.TruncateString(description, 200)
	}

	metadata = &webpub.WebPubMetadata{

		Title:                title,
		Subtitle:             title,
		IssueNumber:          issue,
		Identifier:           "urn:uuid:" + chapter.Hash,
		Description:          description,
		Expires:              time.Now().Add(time.Hour * 3),
		Published:            chapter.TimeUploaded,
		Free:                 true,
		Provider:             "manga.sh",
		BelongsTo:            &webpub.WebPubOwnership{},
		PageCount:            uint64(len(chapter.SeriesChaptersFiles)),
		Image:                chapter.SeriesChaptersSeries.CoverImage,
		Author:               authors,
		Artist:               artists,
		Translator:           translators,
		Language:             langs.Langs[chapter.Language],
		Genre:                genres,
		Direction:            chapter.SeriesChaptersSeries.Direction,
		AccessibilitySummary: "Sequence of images containing drawings with text",
		AccessMode:           "visual",
		AccessibilityControl: []string{"fullKeyboardControl", "fullMouseControl", "fullTouchControl"},
	}
	if chapter.SeriesChaptersSeries.Direction == "ttb" {
		metadata.Rendition = &webpub.WebPubRendition{
			Overflow: "scrolled-continuous",
			Fit:      "width",
		}

	} else {
		metadata.Rendition = &webpub.WebPubRendition{
			Layout:      "pre-paginated",
			Orientation: "portrait",
			Spread:      "landscape",
		}
	}
	metadata.BelongsTo = &webpub.WebPubOwnership{}
	metadata.BelongsTo.Series = append(metadata.BelongsTo.Series, &webpub.WebPubOwner{
		Name:       chapter.SeriesChaptersSeries.Name,
		Identifier: fmt.Sprintf("%s/series/%d/%s", os.Getenv("APP_LINK"), chapter.SeriesChaptersSeries.ID, slugify.Slugify(chapter.SeriesChaptersSeries.Name)),
	})
	fmt.Println(metadata)

	var pages []*webpub.WebPubLink
	for _, k := range chapter.SeriesChaptersFiles {
		link := &webpub.WebPubLink{
			Type:   k.Type,
			Width:  uint(k.Width),
			Height: uint(k.Height),
		}
		if k.Batoto {
			hash := strings.Replace(chapter.Hash, "comics/", "", -1)
			link.Link = fmt.Sprintf("%s%s/%s", os.Getenv("cdn_url"), hash, k.UUID)
		} else {

			link.Link = fmt.Sprintf("%schapters/%s/%s", os.Getenv("cdn_url"), chapter.Hash, k.UUID)
		}
		pages = append(pages, link)

	}
	// view := &models.ChapterViews{
	// 	Chapter:   v,
	// 	TimeStamp: time.Now(),
	// }
	var links []*webpub.WebPubLink
	links = append(links, &webpub.WebPubLink{
		Relation: "self",
		Link:     fmt.Sprintf("%s/series_chapters/%d", os.Getenv("cdn_url"), chapter.ID),
		Type:     "application/webpub+json",
	})
	links = append(links, &webpub.WebPubLink{
		Relation: "cover",
		Link:     chapter.SeriesChaptersSeries.CoverImage,
		Type:     mime.TypeByExtension(filepath.Ext(strings.Replace(chapter.SeriesChaptersSeries.CoverImage, "covers/", "", -1))),
	})

	web := webpub.GenerateComicIssueWebPub(metadata, links, pages)
	b, err := json.Marshal(web)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Last-Modified", time.Now().AddDate(0, 0, -1).Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().Add(time.Minute*5).Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "public, must-revalidate, max-age=240")
	w.Write(b)
}

func main() {
	log.Print("helloworld: starting server...")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/series_chapters/{hash}", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), etag.Handler(r, false)))
}
