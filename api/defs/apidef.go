package defs

//request model
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd string `json:"pwd"`

}

//response model
type SignedUp struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}


//data model
type VideoInfo struct {
	Id string
	AuthorId int
	Name string
	DisplayCtime string
}

type Comment struct {
	Id string
	VideoId string
	Author string
	Content string
}

//session struct
type SimpleSession struct {
	Username string
	TTL int64
}