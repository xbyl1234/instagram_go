package main

import (
	"makemoney/tools"
)

func a(e error) {
	e = &tools.MakeMoneyError{}
}
func b() error {
	var f = 0
	defer func() {
		print(f)
	}()
	f = 12
	return nil
}

func main() {
	e := b()
	a(e)
	print(e)
	//url, _ := url.Parse("http://adsa.com/sad?a=1&s=1")
	//print(url)
	//insta := goinsta.New("USERNAME", "PASSWORD")
	//insta.Login()
	//insta.Account.Liked()
}
