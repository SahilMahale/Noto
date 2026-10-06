package models

type NoteRequest struct {
	Title string `json:"title" xml:"title" form:"title"`
	Body  string `json:"body" xml:"body" form:"body"`
}

type NotePatchRequest struct {
	Title string `json:"title,omitempty" xml:"title" form:"title"`
	Body  string `json:"body,omitempty" xml:"body" form:"body"`
}
type NotesRequest []NoteRequest

type NoteResp struct {
	NoteID string `json:"noteID" xml:"noteID" form:"noteID"`
	Title  string `json:"title" xml:"title" form:"title"`
	Body   string `json:"body" xml:"body" form:"body"`
}
type NotePatchResp struct {
	NoteID string `json:"noteID" xml:"noteID" form:"noteID"`
	Title  string `json:"title,omitempty" xml:"title" form:"title"`
	Body   string `json:"body,omitempty" xml:"body" form:"body"`
}

type NotesResp []NoteResp

type UserSignup struct {
	Username string `json:"user" xml:"user" form:"user"`
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"pass" xml:"pass" form:"pass"`
}
type UserSignin struct {
	Username string `json:"user" xml:"user" form:"user"`
	Password string `json:"pass" xml:"pass" form:"pass"`
}
type JwtResp struct {
	Authtoken string `json:"auth_token"`
	Refreshtoken string `json:"refresh_token"`
}
