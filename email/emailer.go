package email

import(
   models "github.com/jelgar/login/models"
   config "github.com/jelgar/login/config"
)

func send() *models.ApiError {
    auth := smtp.PlainAuth ("", config.MailAuthUser, config.MailAuthPass, config.MailHost)
}
