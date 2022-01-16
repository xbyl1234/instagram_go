package emali

type Guerrilla struct {
	EmailInfo
}

func (this *Guerrilla) RequireAccount() (string, error) {

}
func (this *Guerrilla) RequireCode(number string) (string, error) {

}

func (this *Guerrilla) ReleaseAccount(number string) error {
	return nil
}

func (this *Guerrilla) BlackAccount(number string) error {
	return nil
}

func (this *Guerrilla) GetBalance() (string, error) {
	return "", nil
}

func (this *Guerrilla) Login() error {

}
