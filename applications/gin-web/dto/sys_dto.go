package dto

type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
	VerifyKey  string `json:"verify_key" binding:"required"`
}

type LoginResponse struct {
	UserId       int32  `json:"userId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type ArcoRouteResponse struct {
	Name       string              `json:"name"`
	Key        string              `json:"key"`
	Icon       string              `json:"icon"`
	BreadCrumb bool                `json:"breadcrumb"`
	Children   []ArcoRouteResponse `json:"children"`
}

type VerifyData struct {
	VerifyKey   string `json:"verify_key"`
	VerifyImage string `json:"verify_image"`
}

type GetUserInfoResponse struct {
	UserId   int64    `json:"userId"`
	UserName string   `json:"userName"`
	UserRole []string `json:"userRole"`
	Role     string   `json:"role"`
}

type RouteListResponse struct {
	ID        int64               `json:"id"`
	Key       int64               `json:"key"`
	Name      string              `json:"name"`
	ParentID  int64               `json:"parent_id"`
	Label     string              `json:"label"`
	Path      string              `json:"path"`
	ApiPath   string              `json:"api_path"`
	Component string              `json:"component"`
	Icon      string              `json:"icon"`
	Sequence  int32               `json:"sequence"`
	Type      int32               `json:"type"`
	Status    int32               `json:"status"`
	Children  []RouteListResponse `json:"children"`
}

type UserListResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleID   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
	Status   int8   `json:"status"`
	IsSuper  int8   `json:"is_super"`
	LastTime string `json:"last_time"`
	LastIp   string `json:"last_ip"`
}

type UserUpdateRequest struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleId   int64  `json:"role_id"`
	Status   int32  `json:"status"`
	IsSuper  int32  `json:"is_super"`
}
type RouteUpdateRequest struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	ParentId  int64  `json:"parent_id"`
	Label     string `json:"label" binding:"required"`
	Path      string `json:"path"`
	ApiPath   string `json:"api_path"`
	Icon      string `json:"icon"`
	Sequence  int32  `json:"sequence" binding:"required"`
	Type      int32  `json:"type" binding:"required"`
	Status    int32  `json:"status"`
	Component string `json:"component"`
}

type RouteDeleteRequest struct {
	Id int64 `json:"id"`
}
