// v, err := models.GetSeriesChaptersById(id)
// if err != nil {
// 	c.Data["json"] = err.Error()
// 	c.Ctx.Output.SetStatus(404)
// 	c.ServeJSON()
// 	return
// } else {
// 	if !v.Uploaded {
// 		c.Data["json"] = "Error not published"
// 		c.Ctx.Output.SetStatus(401)
// 		c.ServeJSON()
// 	}
// 	if match := c.Ctx.Input.Header("If-Modified-Since"); match != "" {
// 		last, err := time.Parse(http.TimeFormat, match)
// 		if err == nil {

// 			if last == v.Updated {
// 				c.Ctx.Output.SetStatus(http.StatusNotModified)
// 				return
// 			}
// 		}
// 	}

// 		if b, err := json.Marshal(web); err != nil {
// 			log.Fatal(err)
// 		} else {
// 			if err := cache.Put(cacheName, b, time.Hour*24); err != nil {
// 				log.Println(err)
// 			}
// 			logs.Debug("Write cache: " + cacheName)
// 		}
// 		models.AddChapterViews(view)
// 		if v, err := models.GetSeriesById(v.Series.Id); err == nil {
// 			es.UpdateSeriesById(v)
// 		}

// 	}()
// 	logs.Critical(web)

// 	c.Ctx.ResponseWriter.Header().Set("Last-Modified", v.Updated.Format(http.TimeFormat))
// 	c.Ctx.ResponseWriter.Header().Set("Expires", time.Now().Add(time.Minute*5).Format(http.TimeFormat))
// 	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	c.Ctx.ResponseWriter.Header().Set("Etag", etags.Generate(cacheName, true))
// 	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "public, must-revalidate, max-age=240")
// 	if err := json.NewEncoder(c.Ctx.ResponseWriter).Encode(web); err != nil {
// 		panic(err)
// 	}

package main

import (
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
	"github.com/mangacat/micro-services/utils/webpub"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var metadata *webpub.WebPubMetadata
	var authors []*webpub.SchemaThing
	var artists []*webpub.SchemaThing
	var translators []*webpub.SchemaThing
	var genres []*webpub.WebPubSubject

	if len(v.Series.Tags) != 0 {

		for _, l := range v.Series.Tags {
			k := &webpub.WebPubSubject{
				Name: l.TagName,
				Code: l.TagNamespace,
			}
			genres = append(genres, k)
		}
	}
	if len(v.Series.People) != 0 {

		for _, v := range v.Series.People {
			k := &webpub.SchemaThing{
				SchemaType: "Person",
				Name:       v.People.Name,
				Alt:        v.People.AlternativeNames,
				Identifier: string(v.Id),
				Url:        "/people/" + string(v.Series.Id),
			}
			authors = append(authors, k)

		}
	}

	if len(v.Groups) != 0 {
		for _, v := range v.Groups {
			k := &webpub.SchemaThing{
				Name:       v.Name,
				Identifier: "/groups/" + string(v.Id),
			}
			translators = append(translators, k)

		}
	}
	issue, err := strconv.ParseFloat(v.ChapterNumberAbsolute, 64)
	if err != nil {
		log.Println(err)
		issue = 0
	}
	var title string
	if v.VolumeNumber != "" {
		title = title + fmt.Sprintf("Vol %s ", v.VolumeNumber)
	}
	if v.ChapterNumberVolume != "" {
		title = title + fmt.Sprintf("Ch %s ", v.ChapterNumberVolume)
	}
	if v.Title != "" && v.ChapterNumberVolume == "" && v.VolumeNumber == "" {

		title = title + fmt.Sprintf("%s", v.Title)
	}
	if v.Title != "" {
		title = title + fmt.Sprintf("- %s", v.Title)
	}
	if title == "" {
		title = v.Series.Name
	}

	if len(v.Series.Description) > 200 {
		v.Series.Description = truncate.TruncateString(v.Series.Description, 200)
	}

	metadata = &webpub.WebPubMetadata{

		Title:                title,
		Subtitle:             title,
		IssueNumber:          issue,
		Identifier:           "urn:uuid:" + v.Hash,
		Description:          v.Series.Description,
		Expires:              time.Now().Add(time.Hour * 3),
		Published:            v.TimeUploaded,
		Free:                 true,
		Provider:             "manga.sh",
		BelongsTo:            &webpub.WebPubOwnership{},
		PageCount:            uint64(len(v.Files)),
		Image:                v.Series.CoverImage,
		Author:               authors,
		Artist:               artists,
		Translator:           translators,
		Language:             langs.Langs[v.Language],
		Genre:                genres,
		Direction:            v.Series.Direction,
		AccessibilitySummary: "Sequence of images containing drawings with text",
		AccessMode:           "visual",
		AccessibilityControl: []string{"fullKeyboardControl", "fullMouseControl", "fullTouchControl"},
	}
	if v.Series.Direction == "ttb" {
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
		Name:       v.Series.Name,
		Identifier: fmt.Sprintf("%s/series/%d/%s", os.Getenv("APP_LINK"), v.Series.Id, slugify.Slugify(v.Series.Name)),
	})

	var pages []*webpub.WebPubLink
	for _, k := range v.Files {
		link := &webpub.WebPubLink{
			Type:   k.Type,
			Width:  uint(k.Width),
			Height: uint(k.Height),
		}
		if k.Batoto {
			hash := strings.Replace(v.Hash, "comics/", "", -1)
			link.Link = fmt.Sprintf("%s%s/%s", os.Getenv("cdn_url"), hash, k.UUID)
		} else {

			link.Link = fmt.Sprintf("%schapters/%s/%s", os.Getenv("cdn_url"), v.Hash, k.UUID)
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
		Link:     fmt.Sprintf("%s/series_chapters/%d", os.Getenv("cdn_url"), v.Id),
		Type:     "application/webpub+json",
	})
	links = append(links, &webpub.WebPubLink{
		Relation: "cover",
		Link:     v.Series.CoverImage,
		Type:     mime.TypeByExtension(filepath.Ext(strings.Replace(v.Series.CoverImage, "covers/", "", -1))),
	})

	web := webpub.GenerateComicIssueWebPub(metadata, links, pages)
	b, err := json.Marshal(web)
	if err != nil {
		log.Println(err)
	}
	w.Write(b)
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
