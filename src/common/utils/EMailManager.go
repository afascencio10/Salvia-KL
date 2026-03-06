package utils

import (
	"crypto/tls"
	"net/smtp"
)

type EMail struct {
	Locale       map[string]map[string]string
	FromEMail    string
	FromName     string
	To           []string
	ServerHost   string
	ServerPort   string
	Username     string
	UserPassword string
	Subject      string
	Body         string
}

func (i *EMail) SendEMail() error {
	// Configuración del remitente

	// Configuración del servidor SMTP

	// Configuración del destinatario y contenido del mensaje

	// Crear el mensaje
	message := []byte("From: " + i.FromName + "<" + i.FromEMail + ">" + "\n" +
		//"To: " + to[0] + "\n" +
		"Subject: " + i.Subject + "\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		i.Body)

	// Configuración para la conexión TLS
	serverAddress := i.ServerHost + ":" + i.ServerPort
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         i.ServerHost,
	}

	// Conectar al servidor SMTP usando TLS
	conn, err := tls.Dial("tcp", serverAddress, tlsconfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Crear cliente SMTP
	client, err := smtp.NewClient(conn, i.ServerHost)
	if err != nil {
		return err
	}

	// Autenticarse
	auth := smtp.PlainAuth("", i.Username, i.UserPassword, i.ServerHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Configurar el remitente y destinatario
	if err = client.Mail(i.FromEMail); err != nil {
		return err
	}
	for _, addr := range i.To {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	// Escribir el mensaje
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	// Cerrar la conexión
	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}
