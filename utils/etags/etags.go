package etags

// Generate an Etag for given sring. Allows specifying whether to generate weak
// Etag or not as second parameter
func Generate(str string, weak bool) string {
	if weak {
		return "W/" + `"` + str + `"`
	}

	return str
}
