package github

type GHArchiveFormat int

const (
	AFTarball GHArchiveFormat = iota
	AFZipball
)

var archiveFormatValues = map[GHArchiveFormat]string{
	AFTarball: "tarball",
	AFZipball: "zipball",
}

var archiveFormatFileEx = map[GHArchiveFormat]string{
	AFTarball: ".tar.gz",
	AFZipball: ".zip",
}

type GHRepo struct {
	Id    int
	Owner struct {
		Login             string
		Id                int
		AvatarUrl         string
		GravatarId        string
		Url               string
		HtmlUrl           string
		FollowersUrl      string
		FollowingUrl      string
		GistsUrl          string
		StarredUrl        string
		SubscriptionsUrl  string
		OrganizationsUrl  string
		ReposUrl          string
		EventsUrl         string
		ReceivedEventsUrl string
		Type              string
		SiteAdmin         bool
	}
	Name             string
	FullName         string
	Description      string
	Private          bool
	Fork             bool
	Url              string
	HtmlUrl          string `json:"html_url"`
	ArchiveUrl       string `json:"archive_url"`
	AssigneesUrl     string
	BlobsUrl         string
	BranchesUrl      string
	CloneUrl         string
	CollaboratorsUrl string
	CommentsUrl      string
	CommitsUrl       string
	CompareUrl       string
	ContentsUrl      string
	ContributorsUrl  string
	DeploymentsUrl   string
	DownloadsUrl     string
	EventsUrl        string
	ForksUrl         string
	GitCommitsUrl    string
	GitRefsUrl       string
	GitTagsUrl       string
	GitUrl           string
	HooksUrl         string
	IssueCommentUrl  string
	IssueEventsUrl   string
	IssuesUrl        string
	KeysUrl          string
	LabelsUrl        string
	LanguagesUrl     string
	MergesUrl        string
	MilestonesUrl    string
	MirrorUrl        string
	NotificationsUrl string
	PullsUrl         string
	ReleasesUrl      string
	SshUrl           string
	StargazersUrl    string
	StatusesUrl      string
	SubscribersUrl   string
	SubscriptionUrl  string
	SvnUrl           string
	TagsUrl          string
	TeamsUrl         string
	TreesUrl         string
	Homepage         string
	Language         string
	ForksCount       int
	StargazersCount  int
	WatchersCount    int
	Size             int
	DefaultBranch    string `json:"default_branch"`
	OpenIssuesCount  int
	Topics           []string
	HasIssues        bool
	HasWiki          bool
	HasPages         bool
	HasDownloads     bool
	PushedAt         string
	CreatedAt        string
	UpdatedAt        string
	Permissions      map[string]bool
	AllowRebaseMerge bool
	AllowSquashMerge bool
	AllowMergeCommit bool
	SubscribersCount int
	NetworkCount     int
}
