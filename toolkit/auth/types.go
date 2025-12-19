package auth

type UserInfo struct {
	UserId   string `json:"userId"`
	TenantId string `json:"tenantId"`
	ClientId string `json:"clientId"`
	Timeout  int32  `json:"timeout"`
	DeptName string `json:"deptName"`
	UsMd5    string `json:"usMd5"`

	DataScope   string   `json:"dataScope"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	DeptIds     []string `json:"deptIds"`
}
