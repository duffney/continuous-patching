package main

type Package struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PackageType  string `json:"package_type"`
	Owner        Owner  `json:"owner"`
	VersionCount int    `json:"version_count"`
	Visibility   string `json:"visibility"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	HTMLURL      string `json:"html_url"`
}

type Owner struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Tag struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	URL            string   `json:"url"`
	PackageHTMLURL string   `json:"package_html_url"`
	License        string   `json:"license"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	Descrption     string   `json:"description"`
	HTMLURL        string   `json:"html_url"`
	Metadata       Metadata `json:"metadata"`
}

type Metadata struct {
	PackageType string `json:"package_type"`
	Container   struct {
		Tags []string `json:"tags"`
	} `json:"container"`
}

// "metadata": {
//       "package_type": "container",
//       "container": {
//         "tags": [
//           "latest"
//         ]
//       }
//     }
//   },
