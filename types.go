package main

// UserDeprecated holds properties of old user
type UserDeprecated struct {
	ID            int     `json:"id"`
	Username      string  `json:"username,omitempty"`
	Email         string  `json:"email,omitempty"`
	Lastname      string  `json:"lastname,omitempty"`
	Firstname     string  `json:"firstname,omitempty"`
	Mobile        string  `json:"mobile,omitempty"`
	NeedSMS       int     `json:"need_sms,omitempty"`
	Phone         string  `json:"phone,omitempty"`
	Street        string  `json:"street,omitempty"`
	ZIP           string  `json:"zip,omitempty"`
	City          string  `json:"city,omitempty"`
	DateOfBirth   string  `json:"date_of_birth,omitempty"`
	DateOfEntry   string  `json:"date_of_entry,omitempty"`
	DateOfExit    string  `json:"date_of_exit,omitempty"`
	State         int     `json:"state,omitempty"`
	Credit        float32 `json:"credit,omitempty"`
	CreditDate    string  `json:"credit_date,omitempty"`
	CreditComment string  `json:"credit_comment,omitempty"`
	IBAN          string  `json:"iban,omitempty"`
	BIC           string  `json:"bic,omitempty"`
	SEPA          string  `json:"sepa,omitempty"`
	RememberToken string  `json:"remember_token,omitempty"`
	Additionals   string  `json:"additionals,omitempty"`
	Comment       string  `json:"comment,omitempty"`
	GroupComment  string  `json:"group_comment,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
	LastLogin     string  `json:"last_login,omitempty"`
	Password      string  `json:"password,omitempty"`
}

// User hold props of new user
type User struct {
	ID              int    `json:"id"`
	State           string `json:"state"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Street          string `json:"street"`
	StreetNumber    string `json:"streetNumber"`
	Zip             string `json:"zip"`
	Country         string `json:"country"`
	Birthday        string `json:"birthday"`
	EntranceDate    string `json:"entranceDate"`
	LeavingDate     string `json:"leavingDate,omitempty"`
	AdditionalInfos string `json:"additionalInfos,omitempty"`
	LastActivityAt  string `json:"lastActivityAt,omitempty"`
	CreatedAt       string `json:"createdAt"`
	CreatedBy       int    `json:"createdBy"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
	UpdatedBy       int    `json:"updatedBy,omitempty"`
	UpdateComment   string `json:"updateComment,omitempty"`
	Password        string `json:"password"`
}
