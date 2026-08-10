package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"
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

	// Dominio del remitente, usado para construir el Message-ID
	var fromDomain string = i.FromEMail
	if at := strings.LastIndex(i.FromEMail, "@"); at != -1 {
		fromDomain = i.FromEMail[at+1:]
	}

	// Crear el mensaje. Las cabeceras To, Date y Message-ID son requeridas por
	// los principales proveedores (Gmail/Outlook); sin ellas el correo se
	// clasifica como spam o se rechaza. Separador de cabeceras \r\n (RFC 5322).
	message := []byte("From: " + i.FromName + " <" + i.FromEMail + ">\r\n" +
		"To: " + strings.Join(i.To, ", ") + "\r\n" +
		"Subject: " + i.Subject + "\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"Message-ID: " + fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), fromDomain) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
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
