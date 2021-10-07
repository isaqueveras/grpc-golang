package user

type UserReq struct {
	Name string `json:"name" binding:"required"`
	Age  int32  `json:"age" binding:"required"`
}

type UserRes struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Age  *int32  `json:"age,omitempty"`
}
