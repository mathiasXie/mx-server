package dto

type AddUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	//Age      int    `json:"age" binding:"gt=18"`
}

//type LoginResponse struct {
//	Id    int64  `json:"id"`
//	Token string `json:"token"`
//}

type UserResponse struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
}

type GetAllUserRequest struct {
	Page int64 `json:"page"  binding:"required"`
}

type GetAllUserResponse struct {
	UserList []*UserResponse `json:"user_list"`
	Total    int64           `json:"total"`
}
