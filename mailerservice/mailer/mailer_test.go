package mailer_test

import (
	"../mailer"
	"log"
	"testing"
)

func TestPostman(t *testing.T) {
	//
	//
	//
	postman, err := mailer.NewPostman(
		map[string]struct{ FilePath, Subject string }{
			"sign_up": {
				FilePath: "C:/PetFinderRepos/petfinder/authservice/app/assets/mail/auth/sign_up.html",
				Subject:  "Вы зарегистрировались в Pet Finder",
			},
		},
		"smtp.gmail.com",
		"465",
		"587",
		nil,
		mailer.AddBox(
			"auth_service_gmail",
			"walkmanmail19@gmail.com",
			"QUADRopheniamail12345",
			"",
		),
	)
	//
	//
	//
	if err != nil {
		t.Fatal("Unreachable 1 , ", err)
		return
	}
	//
	// "anikyev95@gmail.com"
	//
	err = postman.SendLetter(
		"auth_service_gmail",
		"sign_up",
		"walkmanmail19@gmail.com",
		struct {
			Name, Telephone, Email string
		}{
			Name:      "Vlad",
			Telephone: "+7-937-458-1886",
			Email:     "anikyev95@gmail.com",
		},
		mailer.TCP,
	)
	if err != nil {
		t.Fatal("Unreachable 3 , ", err)
		return
	}
	//
	//
	//
	log.Println("Finish.")
}
