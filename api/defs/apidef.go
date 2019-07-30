package defs

//request model
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd string `json:"pwd"`

}

//data model
type VideoInfo struct {
	Id string
	AuthorId int
	Name string
	DisplayCtime string
}