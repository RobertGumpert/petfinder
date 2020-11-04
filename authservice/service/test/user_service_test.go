package test

import (
	"authservice/entity"
	"authservice/mapper"
	"authservice/pckg/jwtoken"
	"authservice/pckg/storage"
	"authservice/repository"
	"authservice/service"
	"log"
	"testing"
	"time"
)

var (
	postgresOrm = storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		"postgres",
		"toster123",
		"pet_finder_user",
		"5432",
		"disable",
	)
	userRepository repository.UserRepository = repository.NewUserGormRepository(
		postgresOrm.DB,
	)
	jwtToken = jwtoken.NewJwtTokenConstructor(
		[]byte("key"),
		5*time.Second,
		10*time.Second,
	)
	userService = service.NewUserService(
		[]byte("key"),
		5*time.Second,
		7*time.Second,
		10*time.Second,
	)
)

func register() (*mapper.UserViewModel, error) {
	return userService.Register(&mapper.RegisterUserViewModel{
		Telephone: "8-953-983-0807",
		Password:  "toster123",
		Email:     "walkmanmail19@gmail.com",
		Name:      "Vlad",
	}, userRepository, nil)
}

func auth(password string, user *mapper.UserViewModel) (access, refresh string, outputViewModel *mapper.UserViewModel, err error) {
	return userService.Authorized(&mapper.AuthorizationUserViewModel{
		UserID:    user.UserID,
		Telephone: user.Telephone,
		Password:  password,
		Email:     user.Email,
	}, userRepository, nil)
}

func newAccess(access string, user *mapper.UserViewModel) (string, *mapper.UserViewModel, error) {
	return userService.UpdateAccessToken(&mapper.NewAccessTokenViewModel{
		Access:    access,
		UserID:    user.UserID,
		Telephone: user.Telephone,
		Email:     user.Email,
	}, userRepository, nil)
}

func TestRegisterFlow(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	postgresOrm.Exec("delete from users where user_id > 0;")
	//
	user, err := register()
	if err != nil {
		log.Fatal(err)
	}
	//
	list, err := userService.Get(&mapper.FindUserViewModel{
		UserID: user.UserID,
	}, userRepository, nil)
	if err != nil {
		log.Fatal(err)
	}
	if *list.Users[0] != *user {
		log.Fatal("Non equals struct's.")
	}
	log.Println(*list.Users[0])
	//
	postgresOrm.Exec("delete from users where user_id > 0;")
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthorizedFlow(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	postgresOrm.Exec("delete from users where user_id > 0;")
	var getEntityByID = func(viewModel *mapper.UserViewModel) *entity.User{
		user, err := userRepository.EntityGet(&entity.User{UserID: viewModel.UserID}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return user
	}
	log.Println("Регистрация...")
	userRegisterViewModel, err := register()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("-> ОК")
	//
	//
	//
	log.Println("Проверка авторизации с ложным токеном...")
	fakeAccessToken, _ := jwtToken.NewAccessToken(&jwtoken.Payload{
		FieldFirst:  userRegisterViewModel.Telephone,
		FieldSecond: userRegisterViewModel.Name,
	})
	_, err = userService.IsAuthorized(&mapper.IsAuthorizedViewModel{Access: fakeAccessToken}, userRepository, nil)
	if err == nil {
		t.Fatal(" -> Is auth. by fake access token. ")
	}
	log.Println(" * Проверка того что refresh токен пустой...")
	userEntity := getEntityByID(userRegisterViewModel)
	if userEntity.RefreshToken != "" {
		t.Fatal(" -> Refresh non empty. ")
	}
	log.Println("-> ОК")

	//
	//
	//
	log.Println("Авторизации...")
	log.Println(" * Авторизации с ложным паролем...")
	_, _, _, err = auth("fake_password", userRegisterViewModel)
	if err == nil {
		t.Fatal("Is auth. by fake password. ")
	}
	log.Print(" Отклонено ")
	log.Println(" * Проверка того что refresh токен пустой...")
	userEntity = getEntityByID(userRegisterViewModel)
	if userEntity.RefreshToken != "" {
		t.Fatal("Refresh non empty. ")
	}
	log.Print(" Пустой ")
	log.Println(" * Авторизации с настоящим паролем...")
	access, _, userViewModel, err := auth("toster123", userRegisterViewModel)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(" Выполнено ")
	log.Println(" * Проверка того что refresh токен не пустой...")
	userEntity = getEntityByID(userViewModel)
	if userEntity.RefreshToken == "" {
		t.Fatal("Refresh is empty after auth. ")
	}
	log.Print(" Не пустой ")
	log.Println(" * Попытка повторной авторизации...")
	_, _, userViewModel, err = auth("toster123", userRegisterViewModel)
	if err == nil {
		t.Fatal(err)
	}
	log.Print(" Отклонено ")
	log.Println(" * Проверка того что refresh токен не пустой...")
	userEntity = getEntityByID(userViewModel)
	if userEntity.RefreshToken == "" {
		t.Fatal("Refresh is empty after auth. ")
	}
	log.Print(" Не пустой ")
	log.Println("-> OK")
	//
	//
	//
	log.Println("Проверка авторизации...")
	userViewModel, err = userService.IsAuthorized(&mapper.IsAuthorizedViewModel{Access: access}, userRepository, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(" * Проверка того что refresh токен не пустой...")
	userEntity = getEntityByID(userViewModel)
	if userEntity.RefreshToken == "" {
		t.Fatal("Refresh is empty after auth. ")
	}
	log.Print(" Не пустой ")
	log.Println("-> OK")
	log.Println("Time out...")
	time.Sleep(6*time.Second)
	log.Println("Повторная проверка авторизации...")
	userViewModel, err = userService.IsAuthorized(&mapper.IsAuthorizedViewModel{Access: access}, userRepository, nil)
	if err == nil {
		t.Fatal("Access token isn't expire. ")
	}
	log.Println(" * Проверка того что refresh токен не пустой...")
	userEntity = getEntityByID(userRegisterViewModel)
	if userEntity.RefreshToken == "" {
		t.Fatal("Refresh is empty after auth. ")
	}
	log.Print(" Не пустой ")
	log.Println(" * Выдать новый токен...")
	access, userViewModel, err = newAccess(access, userRegisterViewModel)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(" Выдан ")
	log.Println(" * Проверка того что refresh токен не пустой...")
	userEntity = getEntityByID(userViewModel)
	if userEntity.RefreshToken == "" {
		t.Fatal("Refresh is empty after getting new access token. ")
	}
	log.Print(" Не пустой ")
	log.Println(" * Повторно выдать новый токен...")
	_, _, err = newAccess(access, userRegisterViewModel)
	if err == nil {
		t.Fatal("Повторная выдача access токена. ")
	}
	log.Print(" Отклонено ")
	log.Println("Time out...")
	time.Sleep(3*time.Second)
	userViewModel, err = userService.IsAuthorized(&mapper.IsAuthorizedViewModel{Access: access}, userRepository, nil)
	if err == nil {
		t.Fatal("Refresh token isn't expire. ")
	}
	userEntity = getEntityByID(userRegisterViewModel)
	if userEntity.RefreshToken != "" {
		t.Fatal("Refresh isn't empty. ")
	}
	//
	//
	//
	postgresOrm.Exec("delete from users where user_id > 0;")
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}
