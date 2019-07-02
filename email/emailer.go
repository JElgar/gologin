package email

import(
   models "github.com/jelgar/login/models"
   config "github.com/jelgar/login/config"

   smtp "net/smtp"

    "bytes"
    "fmt"
    "html/template"
)

var auth smtp.Auth

func Send(name string, email string, url string, template string) *models.ApiError {
    auth = smtp.PlainAuth ("", config.MailAuthUser, config.MailAuthPass, config.MailHost)
   
    emailData := struct {
        Name    string
        URL     string
    }{
        Name: name,
        URL: url,
    }

    r := NewRequest([]string{email}, "Hello", "Hello world")
    if err := r.ParseTemplate(template, emailData); err == nil {
        ok, err := r.SendEmail()
        if err != nil { panic(err.Err) }
        fmt.Print("mail sent ")
        fmt.Println(ok)
        return nil
    } else { return &models.ApiError{err, err.Message, 500} }
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, *models.ApiError) {
    mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
    subject := "Subject: " + r.subject + "\n"
    msg := []byte(subject + mime + "\n" + r.body)
//    addr := config.MailHost + ":" + config.MailHostPort

    if err := smtp.SendMail("smtp.gmail.com:587", auth, config.MailAuthUser, r.to, msg); err != nil {
		return false, &models.ApiError{err, "Email Failed to send", 500}
	}

	//if err := smtp.SendMail(addr, auth, config.MailAuthUser, r.to, msg); err != nil {
//		return false, &models.ApiError{err, "Email Failed to send", 500}
//	}
	return true, nil

}


// Need to work out how this works as i just copied it :D
func (r *Request) ParseTemplate(templateFileName string, data interface{}) *models.ApiError {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return &models.ApiError{err, "Error getting email template", 500}
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return &models.ApiError{err, "Template Issue", 500}
	}
	r.body = buf.String()
	return nil
}
