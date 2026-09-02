package models

type IdpUser struct {
	Name    string   `json:"name"`
	Enabled *bool    `json:"enabled,omitempty"`
	Groups  []string `json:"groups,omitempty"`
}

type IdpUserResponse struct {
	Code int32   `json:"code"`
	User IdpUser `json:"user"`
}

type IdpAddUserRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

type IdpUpdateUserRequest struct {
	Password string `json:"password,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

type IdpGroup struct {
	Name    string   `json:"name"`
	Comment string   `json:"comment,omitempty"`
	Users   []string `json:"users,omitempty"`
}

type IdpGroupResponse struct {
	Code  int32    `json:"code"`
	Group IdpGroup `json:"group"`
}

type IdpAddGroupRequest struct {
	Group   string `json:"group"`
	Comment string `json:"comment,omitempty"`
}

type IdpGroupMembershipChangeRequest struct {
	UsersToAdd    []string `json:"usersToAdd,omitempty"`
	UsersToRemove []string `json:"usersToRemove,omitempty"`
}
