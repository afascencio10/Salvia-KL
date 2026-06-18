package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Obtener la ruta donde está el ejecutable para guardar el log ahí mismo
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("No se pudo obtener la ruta del ejecutable: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Configurar logs para que se escriban en backup.log
	logPath := filepath.Join(exeDir, "backup.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.Println("Advertencia: No se pudo crear backup.log, usando salida estándar")
	}

	log.Println("=== Iniciando proceso de respaldo (Local) ===")

	// 1. Cargar variables de entorno relativas al ejecutable
	envPath := filepath.Join(exeDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Advertencia: No se pudo cargar archivo %s, usando variables del sistema\n", envPath)
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := getEnv("DB_NAME", "salvia")
	salviaDir := getEnv("SALVIA_DIR", `C:\salvia`)
	backupDir := getEnv("BACKUP_DIR", `C:\Backups`)
	encKey := os.Getenv("ENCRYPTION_KEY")
	pgDumpPath := getEnv("PG_DUMP_PATH", "pg_dump")

	if len(encKey) != 32 {
		log.Fatal("Error: ENCRYPTION_KEY debe tener exactamente 32 caracteres (AES-256).")
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Fatalf("Error al crear directorio de respaldos: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	tempSqlFile := filepath.Join(backupDir, fmt.Sprintf("dump_%s.sql", timestamp))
	tempZipFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.zip", timestamp))
	finalEncFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.zip.enc", timestamp))

	// 2. Extraer Base de Datos
	log.Println("Iniciando volcado de la base de datos...")
	if err := runPgDump(pgDumpPath, dbHost, dbPort, dbUser, dbPassword, dbName, tempSqlFile); err != nil {
		log.Fatalf("Fallo en pg_dump: %v", err)
	}
	log.Println("Volcado de base de datos exitoso.")

	// 3. Crear ZIP con base de datos y carpeta de binarios
	log.Println("Comprimiendo archivos...")
	if err := createZip(tempZipFile, tempSqlFile, salviaDir); err != nil {
		// Intentar limpiar SQL antes de salir
		os.Remove(tempSqlFile)
		log.Fatalf("Fallo al crear archivo ZIP: %v", err)
	}
	log.Println("Compresión exitosa.")

	// 4. Encriptar el archivo ZIP
	log.Println("Encriptando archivo...")
	if err := encryptFile(tempZipFile, finalEncFile, encKey); err != nil {
		os.Remove(tempSqlFile)
		os.Remove(tempZipFile)
		log.Fatalf("Fallo al encriptar archivo: %v", err)
	}
	log.Println("Encriptación exitosa.")

	// 5. Limpieza de temporales
	log.Println("Limpiando archivos temporales...")
	_ = os.Remove(tempSqlFile)
	_ = os.Remove(tempZipFile)

	// 6. Retención local: Mantener solo el último backup (.enc)
	log.Println("Aplicando retención local...")
	if err := cleanupOldBackups(backupDir); err != nil {
		log.Printf("Advertencia: No se pudieron limpiar backups antiguos: %v", err)
	}

	log.Println("=== Respaldo completado exitosamente ===")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func runPgDump(pgDumpCmd, host, port, user, password, dbname, outFile string) error {
	// Preparar PGPASSWORD
	os.Setenv("PGPASSWORD", password)
	defer os.Unsetenv("PGPASSWORD")

	cmd := exec.Command(pgDumpCmd, "-h", host, "-p", port, "-U", user, "-d", dbname, "-F", "p", "-f", outFile)
	
	// Capturar salida de error en caso de fallo
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error de pg_dump: %s, salida: %s", err.Error(), string(output))
	}
	return nil
}

func createZip(zipFile string, sqlFile string, sourceDir string) error {
	archive, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Agregar el archivo SQL
	if err := addFileToZip(zipWriter, sqlFile, filepath.Base(sqlFile)); err != nil {
		return err
	}

	// Agregar el directorio recursivamente
	err = filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Crear la ruta interna en el zip
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		
		zipPath := filepath.Join(filepath.Base(sourceDir), relPath)

		// Convertir separadores de Windows a formato ZIP
		zipPath = filepath.ToSlash(zipPath)

		return addFileToZip(zipWriter, path, zipPath)
	})

	return err
}

func addFileToZip(zipWriter *zip.Writer, filePath string, zipInternalPath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		// Si un archivo está bloqueado (ej. logs en uso), podríamos omitirlo o fallar
		// Para propósitos de este script, logueamos la advertencia y continuamos
		log.Printf("Advertencia: No se pudo agregar archivo al zip: %s (%v)\n", filePath, err)
		return nil
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = zipInternalPath
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func encryptFile(inputFile string, outputFile string, key string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	// GCM requiere nonce (IV)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Para archivos grandes, leer en memoria todo no es viable.
	// NOTA: Para este MVP leeremos todo el archivo. 
	// En producción para archivos inmensos es mejor usar streams o bloques cifrados (ej. age o AES-CTR).
	plaintext, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Escribir nonce seguido del ciphertext
	if _, err := out.Write(nonce); err != nil {
		return err
	}
	if _, err := out.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func cleanupOldBackups(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var encFiles []fs.FileInfo
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip.enc") {
			info, err := f.Info()
			if err == nil {
				encFiles = append(encFiles, info)
			}
		}
	}

	if len(encFiles) <= 1 {
		return nil // No hay nada que limpiar
	}

	// Ordenar por tiempo de modificación descendente (más nuevo primero)
	sort.Slice(encFiles, func(i, j int) bool {
		return encFiles[i].ModTime().After(encFiles[j].ModTime())
	})

	// Eliminar todos excepto el más nuevo (índice 0)
	for i := 1; i < len(encFiles); i++ {
		filePath := filepath.Join(dir, encFiles[i].Name())
		log.Printf("Eliminando backup antiguo: %s\n", encFiles[i].Name())
		os.Remove(filePath)
	}

	return nil
}
