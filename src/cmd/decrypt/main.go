package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("=== Herramienta de Desencriptación de Respaldos Salvia ===")

	if len(os.Args) < 2 {
		fmt.Println("Uso: decrypt_salvia.exe <archivo.zip.enc>")
		fmt.Println("Ejemplo: decrypt_salvia.exe backup_20260618_105047.zip.enc")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	if !strings.HasSuffix(inputFile, ".enc") {
		log.Fatalf("Error: El archivo debe tener la extensión .enc")
	}

	outputFile := strings.TrimSuffix(inputFile, ".enc")

	key := "S4lvia_9XmK2pQ7vNc8RjT5wLd1HyZ3!"

	if len(key) != 32 {
		log.Fatalf("Error: ENCRYPTION_KEY no encontrada en el .env o no tiene 32 caracteres.")
	}

	fmt.Printf("Desencriptando %s...\n", inputFile)
	if err := decryptFile(inputFile, outputFile, key); err != nil {
		log.Fatalf("Fallo crítico al desencriptar: %v", err)
	}

	fmt.Printf("¡Éxito! Archivo recuperado en: %s\n", outputFile)
}

func decryptFile(inputFile, outputFile, key string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de entrada: %v", err)
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo de salida: %v", err)
	}
	defer out.Close()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Leer el nonce (primeros 12 bytes)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(in, nonce); err != nil {
		return fmt.Errorf("error leyendo nonce: %v", err)
	}

	// Leer el resto (ciphertext)
	ciphertext, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("error leyendo ciphertext: %v", err)
	}

	// Desencriptar
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("error al desencriptar (¿llave incorrecta o archivo corrupto?): %v", err)
	}

	// Escribir el zip final
	if _, err := out.Write(plaintext); err != nil {
		return fmt.Errorf("error escribiendo el archivo desencriptado: %v", err)
	}

	return nil
}
