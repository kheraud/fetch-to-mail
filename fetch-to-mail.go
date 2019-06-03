package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gomail "gopkg.in/gomail.v2"
)

func extractEnv() (smtpHostname, smtpUser, smtpPwd, headerHost, headerForwardedProto, emailFrom, emailTo, emailSubject string) {
	smtpHostname = os.Getenv("SMTP_HOSTNAME")
	smtpUser = os.Getenv("SMTP_USER")
	smtpPwd = os.Getenv("SMTP_PWD")
	headerHost = os.Getenv("HEADER_HOST")
	headerForwardedProto = os.Getenv("HEADER_FWD_PROTO")
	emailFrom = os.Getenv("EMAIL_FROM")
	emailTo = os.Getenv("EMAIL_TO")
	emailSubject = os.Getenv("EMAIL_SUBJECT")

	toCheck := map[string]string{
		"SMTP_HOSTNAME": smtpHostname,
		"SMTP_USER":     smtpUser,
		"SMTP_PWD":      smtpPwd,
		"EMAIL_FROM":    emailFrom,
		"EMAIL_TO":      emailTo,
		"EMAIL_SUBJECT": emailSubject,
	}

	for env, value := range toCheck {
		if len(value) <= 0 {
			log.Fatalln("Valeur incorrecte pour la variable d'environnement ", env, value)
		}
	}

	return
}

func extractParams() string {
	if len(os.Args) <= 1 {
		log.Fatalln("Usage: fetch-to-mail http://targetfetchurl.tld/path/endpoint")
	}

	return os.Args[1]
}

func httpGetBody(url, headerHost, headerForwardedProto string) string {

	client := http.Client{ }

	req, err := http.NewRequest("GET", url, nil)

	if len(headerHost) > 0 {
		req.Host = headerHost
	}

	if len(headerForwardedProto) > 0 {
		req.Header.Add("X-Forwarded-Proto", headerForwardedProto)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln("Impossible d'accéder à l'url", url, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Impossible de lire le contenu de la réponse de l'url", url, err)
	}
	return string(body)
}

func main() {

	smtpHostname, smtpUser, smtpPwd, headerHost, headerForwardedProto, emailFrom, emailTo, emailSubject := extractEnv()

	targetFetch := extractParams()

	fetchedResult := httpGetBody(targetFetch, headerHost, headerForwardedProto)

	if len(fetchedResult) <= 10 {
		log.Fatalln("Le résultat du fetch url devrait être plus long", fetchedResult)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html", fetchedResult)

	d := gomail.NewDialer(smtpHostname, 587, smtpUser, smtpPwd)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		log.Fatalln("Erreurs lors de l'envoi de l'email", err)
	}

	log.Printf("Contenu de la page %s récupéré et envoyé à %s\n", targetFetch, emailTo)
}
