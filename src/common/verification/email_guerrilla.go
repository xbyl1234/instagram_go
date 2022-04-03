package verification

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"makemoney/common"
	"makemoney/common/log"
	"time"
)

const sqlRequireEmail = "SELECT mail_id,subject from new where new.to = ? and is_new is null ORDER BY date desc LIMIT 1"
const sqlUpdateEmail = "update new set is_new ='0' where mail_id=?"

type Guerrilla struct {
	MysqlDB  *sqlx.DB
	Domain   string
	MysqlUrl string
	EmailInfo
}

func (this *Guerrilla) RequireAccount() (string, error) {
	return common.GenString(common.CharSet_abc+common.CharSet_123, 10) + "@" + this.Domain, nil
}

type EmailResult struct {
	MailId  int    `db:"mail_id"`
	Subject string `db:"subject"`
}

func (this *Guerrilla) RequireCode(email string) (string, error) {
	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		var result []EmailResult
		err := this.MysqlDB.Select(&result, sqlRequireEmail, email)
		if err != nil {
			log.Warn("select email db error: %v", err)
		} else {
			if len(result) == 1 {
				_, err = this.MysqlDB.Exec(sqlUpdateEmail, result[0].MailId)
				if err != nil {
					log.Warn("update email db error: %v", err)
				}
				return common.GetCode(result[0].Subject), nil
			}
		}

		log.Warn("wait for email %s code...", email)
		time.Sleep(this.RetryDelayDuration)
	}

	return "", &common.MakeMoneyError{ErrStr: "require code timeout", ErrType: common.RecvPhoneCodeError}
}

func (this *Guerrilla) Login() error {
	return nil
}
