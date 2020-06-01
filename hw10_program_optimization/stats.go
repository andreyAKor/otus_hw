package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"

	//"io/ioutil"
	"errors"
	"regexp"
	"strings"
)

//go:generate easyjson -all
type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	var line []byte
	br := bufio.NewReader(r)
	for i := 0; ; i++ {
		line, _, err = br.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			return
		}

		var user User
		if err = user.UnmarshalJSON(line); err != nil {
			return
		}
		result[i] = user
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re := regexp.MustCompile("\\." + domain)

	for _, user := range u {
		if re.Match([]byte(user.Email)) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
