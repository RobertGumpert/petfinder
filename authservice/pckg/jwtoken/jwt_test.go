package jwtoken_test

import (
	"authservice/pckg/jwtoken"
	"log"
	"testing"
	"time"
)

var(
	dataToBeWrittenStorage = &jwtoken.Payload{
		FieldFirst:  "999-999-999-999",
		FieldSecond: "Vlad",
	}
)

func TestToken(t *testing.T) {
	jwt := jwtoken.NewJwtTokenConstructor([]byte("token_test"), 6 * time.Second, 10 * time.Second)
	//
	//
	// Создаем пару токенов с разным временем жизни:
	// - access - временный
	// - refresh - длительный
	//
	//
	access, refresh, err := jwt.CreateTokens(dataToBeWrittenStorage)
	if err != nil {
		log.Println("Finish error : ", err)
		return
	}
	//
	//
	// Имитируем время простоя access токена,
	// пока никак не проверяется авторизация пользователя.
	//
	//
	log.Println("Create access : ", access)
	log.Println("Create refresh : ", refresh)
	log.Println("--> TIMEOUT 1 : Wait 2 second...")
	time.Sleep(3 * time.Second)
	//
	//
	// --> 1
	//
	// Проверяем, не протух ли access токен,
	// как следствие проверяем авторизован ли пользователь.
	//
	//
	log.Println("Checkout access token...")
	firstDecodeAccess, err := jwt.Decode(access)
	if err != nil {
		log.Println("Finish error : ", err)
		return
	}
	if *firstDecodeAccess != *dataToBeWrittenStorage {
		log.Println("Payload isn't valid : ", err)
		return
	}
	log.Println("Access token is valid : ", firstDecodeAccess)
	//
	//
	// --> 2
	//
	// Имитируем время простоя access токена,
	// пока никак не проверяется авторизация пользователя.
	//
	//
	log.Println("--> TIMEOUT 2 : Wait 4 second...")
	time.Sleep(4 * time.Second)
	//
	//
	// --> 3
	//
	// Проверяем, не протух ли access токен,
	// как следствие проверяем авторизован ли пользователь.
	//
	//
	log.Println("Checkout access token...")
	secondDecodeAccess, err := jwt.Decode(access)
	if err != nil {
		log.Println("Access token isn't valid : ", err)
	} else {
		log.Println("Access token is valid : ", secondDecodeAccess)
		return
	}
	//
	//
	// --> 4
	//
	// Access токен протух, для получения нового,
	// отправляем refresh.
	// Если refresh токен не протух и при расшифровки,
	// данные из payload действительно валидны, то есть такой
	// пользователь существует, то создаём новый access токен.
	//
	//
	//
	log.Println("Create new access token...")
	firstDecodeRefresh, err := jwt.Decode(refresh)
	if err != nil {
		log.Println("Finish error : ", err)
		return
	}
	if *firstDecodeRefresh != *dataToBeWrittenStorage {
		log.Println("Payload isn't valid : ", err)
		return
	}
	log.Println("Refresh token is valid : ", firstDecodeRefresh)
	newAccess, err := jwt.NewAccessToken(dataToBeWrittenStorage)
	if err != nil {
		log.Println("Finish error : ", err)
		return
	}
	log.Println("New access token : ", newAccess)
	//
	//
	// --> 5
	//
	// Имитируем время простоя access токена,
	// пока никак не проверяется авторизация пользователя.
	//
	//
	log.Println("--> TIMEOUT 3 : Wait 7 second...")
	time.Sleep(7 * time.Second)
	log.Println("Checkout access token...")
	thirdDecodeAccess, err := jwt.Decode(newAccess)
	if err != nil {
		log.Println("Access token isn't valid : ", err)
	} else {
		log.Println("Access token is valid : ", thirdDecodeAccess)
		return
	}
	//
	//
	// --> 6
	//
	// Access токен протух, для получения нового,
	// отправляем refresh.
	// На этот раз refresh токен протух, следовательно
	// пользователю необходимо снова выполнить авторизацию
	// в системе.
	//
	//
	//
	log.Println("Create new access token...")
	secondDecodeRefresh, err := jwt.Decode(refresh)
	if err != nil {
		log.Println("Refresh token isn't valid : ", err)
	}else {
		log.Println("Refresh token is valid : ", secondDecodeRefresh)
		return
	}
	log.Println("Finish ok : need new Sign In!")
}