package webpub

import (
	"time"
)

// https://github.com/readium/webpub-manifest

const (
	WEBPUB_CONTEXT          = "https://readium.org/webpub-manifest/context.jsonld"
	COMIC_ISSUE_TYPE        = "ComicIssue"
	WEBPUB_DIRECTION_RTL    = "rtl"
	WEBPUB_DIRECTION_LTR    = "ltr"
	COMIC_DEFAULT_DIRECTION = WEBPUB_DIRECTION_RTL
)

type SchemaThing struct {
	SchemaType string `json:"@type"`
	Name       string `json:"name"`
	Alt        string `json:"alternateName,omitempty"`
	Identifier string `json:"identifier"`
	Url        string `json:"url,omitempty,string"`
}

type WebPubSubject struct {
	Name   string `json:"name"`
	SortAs string `json:"sortAs,omitempty"`
	Scheme string `json:"scheme,omitempty"`
	Code   string `json:"code,omitempty"`
}

type WebPubOwner struct {
	Name       string  `json:"name"`
	Identifier string  `json:"identifier"`
	SortAs     string  `json:"sortAs,omitempty"`
	Position   float64 `json:"position"`
}

type WebPubOwnership struct {
	Collection []*WebPubOwner `json:"collection,omitempty"`
	Series     []*WebPubOwner `json:"series,omitempty"`
}

type WebPubRendition struct {
	Overflow    string `json:"overflow,omitempty"`
	Layout      string `json:"layout"`
	Orientation string `json:"orientation"`
	Spread      string `json:"spread"`
	Fit         string `json:"fit"`
}

type WebPubMetadata struct {
	SchemaType  string  `json:"@type"` // Probably ComicIssue
	Identifier  string  `json:"identifier"`
	Title       string  `json:"title"`
	Subtitle    string  `json:"subtitle,omitempty"`
	Description string  `json:"description,omitempty"`
	IssueNumber float64 `json:"issueNumber,omitempty"`

	// Contributors
	Author      []*SchemaThing `json:"author,omitempty"`
	Artist      []*SchemaThing `json:"artist,omitempty"`
	Editor      []*SchemaThing `json:"editor,omitempty"`
	Translator  []*SchemaThing `json:"translator,omitempty"`
	Illustrator []*SchemaThing `json:"illustrator,omitempty"`
	Letterer    []*SchemaThing `json:"letterer,omitempty"`
	Penciler    []*SchemaThing `json:"penciler,omitempty"`
	Colorist    []*SchemaThing `json:"colorist,omitempty"`
	Inker       []*SchemaThing `json:"inker,omitempty"`
	Narrator    []*SchemaThing `json:"narrator,omitempty"`

	Publisher            []*SchemaThing   `json:"publisher,omitempty"`
	Imprint              []*SchemaThing   `json:"imprint,omitempty"`
	AccessMode           string           `json:"accessMode,omitempty"`
	AccessibilityControl []string         `json:"accessibilityControl,omitempty"`
	AccessibilitySummary string           `json:"accessibilitySummary,omitempty"`
	Free                 bool             `json:"isAccessibleForFree"`
	Provider             string           `json:"provider,omitempty"`
	Published            time.Time        `json:"published,omitempty"`
	Modified             time.Time        `json:"modified,omitempty"`
	Expires              time.Time        `json:"expires,omitempty"`
	Language             string           `json:"language,omitempty"`
	Genre                []*WebPubSubject `json:"subject,omitempty"`
	BelongsTo            *WebPubOwnership `json:"belongsTo,omitempty"`
	Direction            string           `json:"readingProgression,omitempty"`
	PageCount            uint64           `json:"numberOfPages,omitempty"`
	Image                string           `json:"image,omitempty"`
	Thumbnail            string           `json:"thumbnailUrl,omitempty"`
	Rendition            *WebPubRendition `json:"rendition,omitempty"`
}

type WebPubEncryptionObject struct {
	Algorithm      string `json:"algorithm"`
	Compression    string `json:"compression,omitempty"`
	OriginalLength uint64 `json:"original-length,omitempty"`
	Profile        string `json:"profile,omitempty"`
	Scheme         string `json:"scheme,omitempty"`
}

type WebPubLinkProperties struct {
	Orientation  string                  `json:"orientation,omitempty"`
	Page         string                  `json:"page,omitempty"`
	Overflow     string                  `json:"overflow,omitempty"`
	Spread       string                  `json:"spread,omitempty"`
	Layout       string                  `json:"layout,omitempty"`
	MediaOverlay string                  `json:"media-overlay,omitempty"`
	Contains     string                  `json:"contains,omitempty"`
	Encrypted    *WebPubEncryptionObject `json:"encrypted,omitempty"`
}

type WebPubLink struct {
	Link  string `json:"href"`
	Type  string `json:"type"` // Mimetype
	Title string `json:"title,omitempty"`

	// https://github.com/readium/webpub-manifest/blob/master/relationships.md
	// Can be: alternate, contents, cover, manifest, search, self
	Relation string `json:"rel,omitempty"`

	Properties *WebPubLinkProperties `json:"properties,omitempty"`
	Height     uint                  `json:"height,omitempty"`
	Width      uint                  `json:"width,omitempty"`
	Duration   float64               `json:"duration,omitempty"`
	Templated  bool                  `json:"templated,omitempty"`
}

type WebPub struct {
	Context   string          `json:"@context"` // See above const
	Metadata  *WebPubMetadata `json:"metadata"`
	Links     []*WebPubLink   `json:"links"`
	Spine     []*WebPubLink   `json:"readingOrder"`
	Resources []*WebPubLink   `json:"resources,omitempty"`
} // TODO https://github.com/readium/webpub-manifest/blob/master/extensions/epub.md#collection-roles

func GenerateComicIssueWebPub(metadata *WebPubMetadata, links []*WebPubLink, pages []*WebPubLink) *WebPub {
	pub := &WebPub{}
	pub.Context = WEBPUB_CONTEXT

	metadata.SchemaType = COMIC_ISSUE_TYPE
	metadata.Provider = "manga.cat"
	metadata.Rendition = &WebPubRendition{
		Layout:      "pre-paginated",
		Orientation: "auto",
		Spread:      "landscape",
	}

	pub.Metadata = metadata
	pub.Spine = pages
	pub.Links = links
	return pub
}
