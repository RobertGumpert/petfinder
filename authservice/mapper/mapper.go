package mapper

import (
	"authservice/entity"
	"authservice/pckg/regularexp"
	"strings"
)

type UserListViewModel struct {
	Users []*UserViewModel `json:"users"`
}

func (m *UserListViewModel) Mapper1(list []entity.User) *UserListViewModel {
	m.Users = make([]*UserViewModel, 0)
	for _, model := range list {
		m.Users = append(m.Users, new(UserViewModel).Mapper(&model))
	}
	return m
}

func (m *UserListViewModel) Mapper2(list []*UserViewModel) *UserListViewModel {
	m.Users = make([]*UserViewModel, 0)
	for _, model := range list {
		m.Users = append(m.Users, model)
	}
	return m
}

//---------------------------------------------------

type UserViewModel struct {
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func (m *UserViewModel) Mapper(user *entity.User) *UserViewModel {
	m.UserID = user.UserID
	m.Name = user.Name
	m.Telephone = user.Telephone
	m.Email = user.Email
	return m
}

//---------------------------------------------------

type FindUserViewModel struct {
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func (m *FindUserViewModel) Validator() error {
	if m.UserID == 0 {
		return ErrorNonValidData
	}
	name := strings.TrimSpace(m.Name)
	m.Email = strings.TrimSpace(m.Email)
	m.Telephone = strings.TrimSpace(m.Telephone)
	if name == "" && m.Email != "" && m.Telephone != "" && m.UserID == 0 {
		return ErrorNonValidData
	}
	if m.Email != "" && !regularexp.EmailValid(m.Email) {
		return ErrorNonValidData
	}
	if m.Telephone != "" && !regularexp.TelephoneValid(m.Telephone) {
		return ErrorNonValidData
	}
	return nil
}

func (m *FindUserViewModel) Mapper() *entity.User {
	return &entity.User{
		UserID:    m.UserID,
		Telephone: m.Telephone,
		Email:     m.Email,
		Name:      m.Name,
	}
}

//---------------------------------------------------

type UpdateUserViewModel struct {
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func (m *UpdateUserViewModel) Validator() error {
	if m.UserID == 0 {
		return ErrorNonValidData
	}
	name := strings.TrimSpace(m.Name)
	m.Email = strings.TrimSpace(m.Email)
	m.Telephone = strings.TrimSpace(m.Telephone)
	if name == "" && m.Telephone != "" && m.Email != "" && m.UserID == 0 {
		return ErrorNonValidData
	}
	if m.Email != "" && !regularexp.EmailValid(m.Email) {
		return ErrorNonValidData
	}
	if m.Telephone != "" && !regularexp.TelephoneValid(m.Telephone) {
		return ErrorNonValidData
	}
	return nil
}

func (m *UpdateUserViewModel) Mapper() *entity.User {
	return &entity.User{
		UserID:    m.UserID,
		Telephone: m.Telephone,
		Email:     m.Email,
		Name:      m.Name,
	}
}

//---------------------------------------------------

type UpdateAvatarViewModel struct {
	UserID    uint64 `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
}

func (m *UpdateAvatarViewModel) Validator() error {
	if m.UserID == 0 {
		return ErrorNonValidData
	}
	m.AvatarUrl = strings.TrimSpace(m.AvatarUrl)
	if m.AvatarUrl == "" {
		return ErrorNonValidData
	}
	return nil
}

func (m *UpdateAvatarViewModel) Mapper() (uint64, map[string]interface{}) {
	return m.UserID, map[string]interface{}{
		"avatar_url": m.AvatarUrl,
	}
}

//---------------------------------------------------

type RegisterUserViewModel struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func (m *RegisterUserViewModel) Mapper() *entity.User {
	return &entity.User{
		Telephone: m.Telephone,
		Password:  m.Password,
		Email:     m.Email,
		Name:      m.Name,
	}
}

func (m *RegisterUserViewModel) Validator() error {
	m.Email = strings.TrimSpace(m.Email)
	m.Telephone = strings.TrimSpace(m.Telephone)
	m.Password = strings.TrimSpace(m.Password)
	name := strings.TrimSpace(m.Name)
	if m.Email == "" || m.Telephone == "" || m.Password == "" || name == "" {
		return ErrorNonValidData
	}
	if !regularexp.EmailValid(m.Email) || !regularexp.TelephoneValid(m.Telephone) {
		return ErrorNonValidData
	}
	return nil
}

//---------------------------------------------------

type AuthorizationUserViewModel struct {
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func (m *AuthorizationUserViewModel) Validator() error {
	if m.UserID == 0 {
		return ErrorNonValidData
	}
	email := strings.TrimSpace(m.Email)
	telephone := strings.TrimSpace(m.Telephone)
	password := strings.TrimSpace(m.Password)
	if email == "" || telephone == "" || password == "" {
		return ErrorNonValidData
	}
	if !regularexp.EmailValid(email) || !regularexp.TelephoneValid(telephone) {
		return ErrorNonValidData
	}
	return nil
}

func (m *AuthorizationUserViewModel) Mapper() *entity.User {
	return &entity.User{
		UserID:    m.UserID,
		Telephone: m.Telephone,
		Password:  m.Password,
		Email:     m.Email,
	}
}

//---------------------------------------------------

type ResetPasswordViewModel struct {
	ResetToken string `json:"reset_token"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
}

func (m *ResetPasswordViewModel) Mapper() *entity.User {
	return &entity.User{
		Telephone: m.Telephone,
		Password:  m.Password,
		Email:     m.Email,
	}
}

func (m *ResetPasswordViewModel) Validator(checkPassToken bool) error {
	if checkPassToken {
		if strings.TrimSpace(m.Password) == "" || strings.TrimSpace(m.ResetToken) == "" {
			return ErrorNonValidData
		}
	}
	email := strings.TrimSpace(m.Email)
	telephone := strings.TrimSpace(m.Telephone)
	if email == "" || telephone == "" {
		return ErrorNonValidData
	}
	if !regularexp.EmailValid(email) || !regularexp.TelephoneValid(telephone) {
		return ErrorNonValidData
	}
	return nil
}

//---------------------------------------------------

type IsAuthorizedViewModel struct {
	Access string `json:"access"`
}

func (m *IsAuthorizedViewModel) Validator() error {
	if strings.TrimSpace(m.Access) == "" {
		return ErrorNonValidData
	}
	return nil
}

//---------------------------------------------------

type NewAccessTokenViewModel struct {
	Access    string `json:"access"`
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

func (m *NewAccessTokenViewModel) Validator() error {
	if m.UserID == 0 {
		return ErrorNonValidData
	}
	token := strings.TrimSpace(m.Access)
	email := strings.TrimSpace(m.Email)
	telephone := strings.TrimSpace(m.Telephone)
	if email == "" || telephone == "" || token == "" {
		return ErrorNonValidData
	}
	if !regularexp.EmailValid(email) || !regularexp.TelephoneValid(telephone) {
		return ErrorNonValidData
	}
	return nil
}

func (m *NewAccessTokenViewModel) Mapper() *entity.User {
	return &entity.User{
		UserID:    m.UserID,
		Telephone: m.Telephone,
		Email:     m.Email,
	}
}
