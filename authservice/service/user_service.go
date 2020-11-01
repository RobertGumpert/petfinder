package service

import (
	"authservice/entity"
	"authservice/pckg/jwtoken"
	"authservice/pckg/runtimeinfo"
	"authservice/repository"
	"context"
	"encoding/base64"
	"log"
	"strings"
	"time"
)

type User struct {
	jwtToken                   *jwtoken.JwtTokenMember
	lifetimeResetPasswordToken time.Duration
}

func NewUserService(uniqueSignature []byte, accessTokenLifetime time.Duration, refreshTokenLifetime time.Duration, lifetimeResetPasswordToken time.Duration) *User {
	return &User{
		jwtToken: jwtoken.NewJwtTokenConstructor(
			uniqueSignature,
			accessTokenLifetime,
			refreshTokenLifetime,
		),
		lifetimeResetPasswordToken: lifetimeResetPasswordToken,
	}
}

func (u *User) Register(inputViewModel *RegisterUserViewModel, db repository.UserRepository, ctx context.Context) (*UserViewModel, error) {
	if err := inputViewModel.Validator(); err != nil {
		return nil, err
	}
	model := inputViewModel.Mapper()
	outputViewModel := new(UserViewModel)
	err := db.Create(model, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, ErrorBadDataOperation
	}
	user, err := db.EntityGet(model, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, ErrorBadDataOperation
	}
	return outputViewModel.Mapper(user), nil
}

// case: польователь не авторизован:
//   > access токен, который прислал пользователь, просрочен.
//   > refresh токен в базе пустой.
//   > refresh токен в базе не пустой, но просрочен.
//   > payload'ы токенов не совпадают.
// case: польователь авторизован:
//   > access токен, который прислал пользователь, не просрочен.
//   > refresh токен в базе не пустой.
//   > refresh токен в базе не пустой и не просрочен.
//   > payload'ы токенов совпадают.
func (u *User) Authorized(inputViewModel *AuthorizationUserViewModel, db repository.UserRepository, ctx context.Context) (access, refresh string, outputViewModel *UserViewModel, err error) {
	access, refresh = "", ""
	outputViewModel = new(UserViewModel)
	if err = inputViewModel.Validator(); err != nil {
		return access, refresh, outputViewModel, err
	}
	userEntity, err := db.EntityGet(inputViewModel.Mapper(), ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return access, refresh, outputViewModel, ErrorNonExistUser
	}
	if strings.TrimSpace(userEntity.RefreshToken) != "" {
		if _, err := u.tokenIsExpire(userEntity.RefreshToken); err != nil {
			if err := u.updateRefreshToken(userEntity, nil, db, ctx); err != nil {
				go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
				return access, refresh, outputViewModel, ErrorBadDataOperation
			}
			return access, refresh, outputViewModel, ErrorNonValidRefreshToken
		}
		return access, refresh, outputViewModel, ErrorNonValidRefreshToken
	}
	access, refresh, err = u.createAuthorizationTokens(userEntity)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return access, refresh, outputViewModel, ErrorNonValidData
	}
	if err := u.updateRefreshToken(userEntity, refresh, db, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return access, refresh, outputViewModel, ErrorBadDataOperation
	}
	return access, refresh, outputViewModel.Mapper(userEntity), nil
}

// case: польователь не авторизован:
//   > access токен, который прислал пользователь, просрочен.
//   > refresh токен в базе пустой.
//   > refresh токен в базе не пустой, но просрочен.
//   > payload'ы токенов не совпадают.
// case: польователь авторизован:
//   > access токен, который прислал пользователь, не просрочен.
//   > refresh токен в базе не пустой.
//   > refresh токен в базе не пустой и не просрочен.
//   > payload'ы токенов совпадают.
func (u *User) IsAuthorized(inputViewModel *IsAuthorizedViewModel, db repository.UserRepository, ctx context.Context) (*UserViewModel, error) {
	if err := inputViewModel.Validator(); err != nil {
		return nil, err
	}
	accessPayload, err := u.tokenIsExpire(inputViewModel.Access)
	if err != nil {
		return nil, ErrorNonValidAccessToken
	}
	userEntity, err := u.userFromPayload(accessPayload, db, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user not found", err, "]")
		return nil, err
	}
	if strings.TrimSpace(userEntity.RefreshToken) == "" {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user refresh token is empty]")
		return nil, ErrorNonValidRefreshToken
	}
	refreshPayload, err := u.tokenIsExpire(userEntity.RefreshToken)
	if err != nil {
		if err := u.updateRefreshToken(userEntity, nil, db, ctx); err != nil {
			go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return nil, ErrorBadDataOperation
		}
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user refresh token is expire]")
		return nil, ErrorNonValidRefreshToken
	}
	if err := u.compareTokensPayload(accessPayload, refreshPayload); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[token isn't compared]")
		return nil, err
	}
	return new(UserViewModel).Mapper(userEntity), nil
}

// case: выдать новый токен:
//   > access токен, который прислал пользователь, просрочен.
//   > refresh токен в базе не пустой.
//   > refresh токен в базе не пустой и не просрочен.
//   > payload'ы токенов совпадают.
// case: не выдать новый токен:
//   > access токен, который прислал пользователь, не просрочен.
//   > refresh токен в базе пустой.
//   > refresh токен в базе не пустой, но просрочен.
//   > payload'ы токенов не совпадают.
func (u *User) UpdateAccessToken(inputViewModel *NewAccessTokenViewModel, db repository.UserRepository, ctx context.Context) (string, *UserViewModel, error) {
	if err := inputViewModel.Validator(); err != nil {
		return "", nil, err
	}
	accessPayload, err := u.tokenIsExpire(inputViewModel.Access)
	if err == nil {
		return "", nil, ErrorNonValidAccessToken
	}
	userEntity, err := db.EntityGet(inputViewModel.Mapper(), ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user not found ", err, "]")
		return "", nil, ErrorNonExistUser
	}
	if strings.TrimSpace(userEntity.RefreshToken) == "" {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user refresh token is empty]")
		return "", nil, ErrorNonValidRefreshToken
	}
	refreshPayload, err := u.tokenIsExpire(userEntity.RefreshToken)
	if err != nil {
		if err := u.updateRefreshToken(userEntity, nil, db, ctx); err != nil {
			go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return "", nil, ErrorBadDataOperation
		}
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user refresh token is expire]")
		return "", nil, ErrorNonValidRefreshToken
	}
	if err := u.compareTokensPayload(accessPayload, refreshPayload); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[token isn't compared]")
		return "", nil, err
	}
	access, _, err := u.createAuthorizationTokens(userEntity)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return "", nil, ErrorBadDataOperation
	}
	return access, new(UserViewModel).Mapper(userEntity), nil
}

func (u *User) Update(inputViewModel *UpdateUserViewModel, db repository.UserRepository, ctx context.Context) (*UserViewModel, error) {
	if err := inputViewModel.Validator(); err != nil {
		return nil, err
	}
	parseUserEntity := inputViewModel.Mapper()
	err := db.EntityUpdate(parseUserEntity, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, ErrorBadDataOperation
	}
	updateUserEntity, err := db.EntityGet(parseUserEntity, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, ErrorNonExistUser
	}
	return new(UserViewModel).Mapper(updateUserEntity), nil
}

func (u *User) UpdateAvatar(inputViewModel *UpdateAvatarViewModel, db repository.UserRepository, ctx context.Context) error {
	if err := inputViewModel.Validator(); err != nil {
		return err
	}
	id, mp := inputViewModel.Mapper()
	err := db.MapUpdate(id, mp, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return ErrorBadDataOperation
	}
	return nil
}

func (u *User) Get(inputViewModel *FindUserViewModel, db repository.UserRepository, ctx context.Context) (*UserListViewModel, error) {
	if err := inputViewModel.Validator(); err != nil {
		return nil, err
	}
	parseUserEntity := inputViewModel.Mapper()
	list, err := db.EntityList(parseUserEntity, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, ErrorNonExistUser
	}
	return new(UserListViewModel).Mapper1(list), nil
}

func (u *User) GetResetPasswordToken(inputViewModel *ResetPasswordViewModel, db repository.UserRepository, ctx context.Context) (token string, err error) {
	token = ""
	if err := inputViewModel.Validator(false); err != nil {
		return token, err
	}
	userEntity, err := db.EntityGet(inputViewModel.Mapper(), ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return token, ErrorNonExistUser
	}
	token, err = u.jwtToken.Encode(&jwtoken.Payload{
		FieldFirst:  base64.StdEncoding.EncodeToString([]byte(inputViewModel.Telephone)),
		FieldSecond: userEntity.Password,
	}, u.lifetimeResetPasswordToken)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return token, err
	}
	return token, nil
}

func (u *User) ResetPassword(inputViewModel *ResetPasswordViewModel, db repository.UserRepository, ctx context.Context) (access, refresh string, err error) {
	if err := inputViewModel.Validator(true); err != nil {
		return "", "", err
	}
	tokenPayload, err := u.jwtToken.Decode(inputViewModel.ResetToken)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return "", "", ErrorNonValidData
	}
	telephone, err := base64.StdEncoding.DecodeString(tokenPayload.FieldFirst)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return "", "", ErrorNonValidData
	}
	userEntity, err := db.EntityGet(&entity.User{
		Telephone: string(telephone),
	}, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return "", "", ErrorNonExistUser
	}
	if userEntity.Password != tokenPayload.FieldSecond {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user password not equal password from token]")
		return "", "", ErrorRetryingPasswordChange
	}
	if userEntity.Password == inputViewModel.Password {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[user password equal password from inputViewModel]")
		return "", "", ErrorRetryingPasswordChange
	}
	access, refresh, err = u.createAuthorizationTokens(userEntity)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return "", "", ErrorNonValidData
	}
	err = db.MapUpdate(userEntity.UserID, map[string]interface{}{
		"password":      inputViewModel.Password,
		"refresh_token": refresh,
	}, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return "", "", ErrorBadDataOperation
	}
	return access, refresh, nil
}

func (u *User) createAuthorizationTokens(user *entity.User) (access, refresh string, err error) {
	access, refresh, err = u.jwtToken.CreateTokens(&jwtoken.Payload{
		FieldFirst:  user.Telephone,
		FieldSecond: user.Name,
	})
	if err != nil {
		return access, refresh, ErrorNonValidData
	}
	return access, refresh, nil
}

func (u *User) updateRefreshToken(user *entity.User, refresh interface{}, db repository.UserRepository, ctx context.Context) error {
	return db.MapUpdate(user.UserID, map[string]interface{}{
		"refresh_token": refresh,
	}, ctx)
}

func (u *User) userFromPayload(payload *jwtoken.Payload, db repository.UserRepository, ctx context.Context) (*entity.User, error) {
	user, err := db.EntityGet(&entity.User{
		Telephone: payload.FieldFirst,
		Name:      payload.FieldSecond,
	}, ctx)
	if err != nil {
		return nil, ErrorNonExistUser
	}
	return user, nil
}

func (u *User) tokenIsExpire(token string) (*jwtoken.Payload, error) {
	payload, err := u.jwtToken.Decode(token)
	if err != nil {
		return payload, ErrorAuthorizationTokenExpire
	}
	return payload, nil
}

func (u *User) compareTokensPayload(accessPayload, refreshPayload *jwtoken.Payload) error {
	if (strings.Compare(refreshPayload.FieldFirst, accessPayload.FieldFirst) != 0) ||
		(strings.Compare(refreshPayload.FieldSecond, accessPayload.FieldSecond) != 0) {
		return ErrorAuthorizationTokenExpire
	}
	return nil
}
