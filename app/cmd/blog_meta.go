package cmd

// BlogMeta represents blog meta information
type BlogMeta struct {
	Title       string
	Description string
	Link        string
}

// BlogLink set blog meta link
func (m *BlogMeta) BlogLink(url string) string {
	m.Link = url
	return url
}

// BlogTitle set blog meta title
func (m *BlogMeta) BlogTitle(title string) string {
	m.Title = title
	return title
}

// BlogDescription set blog meta description
func (m *BlogMeta) BlogDescription(desc string) string {
	m.Description = desc
	return desc
}
