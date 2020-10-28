package regularexp

import "regexp"

// https://play.golang.org/p/jYuNBChInG7
//
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// https://play.golang.org/p/myjcnn_sqgd
//
var telephoneRegex = regexp.MustCompile("[+]?[]?[0-9]+-?[0-9]{3}-[0-9]{3}-?[0-9]{4}")

func EmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func TelephoneValid(e string) bool {
	return telephoneRegex.MatchString(e)
}
