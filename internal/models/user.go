package models

import "github.com/okta/okta-sdk-golang/v5/okta"

type CreateUserRequest struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Password    string `json:"password,omitempty"`
	MobilePhone string `json:"mobilePhone,omitempty"`
}

type UserResponse struct {
	ID        string      `json:"id"`
	Status    string      `json:"status"`
	Created   string      `json:"created"`
	Activated string      `json:"activated,omitempty"`
	Profile   UserProfile `json:"profile"`
}

type UserProfile struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Email       string `json:"email"`
	Login       string `json:"login"`
	MobilePhone string `json:"mobilePhone,omitempty"`
}

func ConvertSDKUserToUserResponse(sdkUser *okta.User) *UserResponse {
	if sdkUser == nil {
		return nil
	}

	resp := &UserResponse{
		ID:        sdkUser.GetId(),
		Status:    sdkUser.GetStatus(),
		Created:   sdkUser.Created.String(),
		Activated: sdkUser.GetActivated().String(),
	}

	if sdkUser.Profile != nil {
		resp.Profile.FirstName = sdkUser.Profile.GetFirstName()
		resp.Profile.LastName = sdkUser.Profile.GetLastName()
		resp.Profile.Email = sdkUser.Profile.GetEmail()
		resp.Profile.Login = sdkUser.Profile.GetLogin()
		resp.Profile.MobilePhone = sdkUser.Profile.GetMobilePhone()
	}

	return resp
}
