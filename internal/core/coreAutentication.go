package core

type User struct {
    UserId 				string `json:"userId,omitempty"`
    Password  			string `json:"password,omitempty"`
	Token				string `json:"token,omitempty"`
}