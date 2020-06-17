package models

type User struct {
	Id           string  `json:",omitempty"`
	User         string `json:",omitempty"`
	UserPassword string `json:",omitempty"`

	FirstName string `json:",omitempty"`

	LastName string `json:",omitempty"`
	Email string `json:",omitempty"`
	Phone string `json:",omitempty"`
	AccountStatus string `json:",omitempty"`
	Domain string `json:",omitempty"`
	Roles []interface{} `json:",omitempty"`
	
}

// type Roles struct {
// 	RoleId:
	
// }



func NewUser (id, user, userPassword, firstName,lastName,email,phone,accountStatus,domain string,roles []interface{}) *User {
	
	return &User{
		Id:  id,
		User: user,
		UserPassword:  userPassword,
		FirstName:  firstName,
		LastName:  lastName,
		Email: email,
		Phone:phone,
		AccountStatus: accountStatus,
		Domain:       domain,
		Roles:roles,

	}
}

func (user *User) ToMap() (map[string]interface{}, error) {
	userAttributeMap := make(map[string]interface{})
	A(userAttributeMap, "id", user.Id)
	A(userAttributeMap, "username", user.User)
	A(userAttributeMap, "password", user.UserPassword)
	A(userAttributeMap, "firstName", user.FirstName)
	A(userAttributeMap, "lastName", user.LastName)
	A(userAttributeMap, "emailAddress", user.Email)
	A(userAttributeMap, "phoneNumber", user.Phone)
	A(userAttributeMap, "accountStatus", user.AccountStatus)
	A(userAttributeMap, "domainId", user.Domain)
	A(userAttributeMap, "roles", user.Roles)


	return userAttributeMap, nil
}