package common

type HostContext struct {
	IP       string
	Username string
	Password string
}

type Pagination struct {
	Page     int
	PageSize int
}

type QueryFilter struct {
	// The attributes of DB table returned.
	Fields []string
	// The keyword for fuzzy search.
	Keyword map[string]string
	// Pagination
	Pagination *Pagination
	// The condition to filter the records by query.
	Conditions interface{}
}

type SystemInfo struct {
	ComputerName   string `json:"host_name"`
	Caption        string `json:"os_type"`
	OSArchitecture string `json:"os_arch"`
	OSVersion      string `json:"os_version"`
	BuildNumber    string `json:"build_number"`
}

type DirectoryDetail struct {
	HostIP         string `json:"host_ip"`
	Name           string `json:"name"`
	CreationTime   string `json:"creation_time"`
	LastAccessTime string `json:"last_access_time"`
	LastWriteTime  string `json:"last_write_time"`
	Exist          bool   `json:"exist"`
	FullPath       string `json:"full_path"`
	ParentFullPath string `json:"parent_full_path"`
}

type FailedRESTResponse struct {
	Error string `json:"error"`
}

var SecurityKey = "0123456789ABCDEF0123456789ABCDEF"
