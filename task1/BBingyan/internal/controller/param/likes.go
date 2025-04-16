package param

type UserLikeRequest struct {
	LikedUser string `json:"liked-user"`
}

type PostLikeRequest struct {
	LikedPost int `json:"liked-post"`
}

type UserLikeResponse struct {
	Likes int `json:"likes"`
}
