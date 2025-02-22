package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// uploadFileToDrive sube el PDF a Google Drive y lo hace público, devolviendo la URL pública.
func uploadFileToDrive(pdfPath string) (string, error) {
	ctx := context.Background()

	// Inicializa el servicio de Drive con las credenciales de la cuenta de servicio.
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("service_account.json"), option.WithScopes(drive.DriveFileScope))
	if err != nil {
		return "", fmt.Errorf("no se pudo crear el servicio de Drive: %v", err)
	}

	// Abre el archivo PDF.
	f, err := os.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el PDF: %v", err)
	}
	defer f.Close()

	// Define los metadatos del archivo a subir.
	fileMetadata := &drive.File{
		Name:     "output.pdf",
		MimeType: "application/pdf",
	}

	// Sube el archivo a Drive.
	file, err := srv.Files.Create(fileMetadata).Media(f).Do()
	if err != nil {
		return "", fmt.Errorf("error subiendo el archivo a Drive: %v", err)
	}

	// Cambia los permisos para que cualquiera con el enlace pueda leer el archivo.
	perm := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	if _, err := srv.Permissions.Create(file.Id, perm).Do(); err != nil {
		return "", fmt.Errorf("error al actualizar los permisos: %v", err)
	}

	// Construye la URL pública.
	publicURL := "https://drive.google.com/uc?export=view&id=" + file.Id
	return publicURL, nil
}

// generaQR crea un código QR con la URL proporcionada y lo guarda en la ruta qrPath.
func generaQR(url, qrPath string) error {
	return qrcode.WriteFile(url, qrcode.Medium, 256, qrPath)
}

func main() {
	// Ruta del PDF que quieres subir (asegúrate de haberlo generado previamente).
	pdfPath := "./static/output.pdf"

	// Subir el PDF a Google Drive y obtener la URL pública.
	publicURL, err := uploadFileToDrigit initve(pdfPath)
	if err != nil {
		log.Fatalf("Error al subir a Drive: %v", err)
	}
	log.Printf("Archivo subido, URL pública: %s\n", publicURL)

	// Generar el QR a partir de la URL pública.
	qrPath := "./static/qr.png"
	if err := generaQR(publicURL, qrPath); err != nil {
		log.Fatalf("Error al generar el QR: %v", err)
	}
	log.Printf("Código QR generado en: %s\n", qrPath)

	// Aquí podrías continuar con la lógica de tu aplicación,
	// por ejemplo, mostrando una respuesta al usuario o redirigiendo.
	fmt.Printf("PDF disponible en: %s\n", publicURL)
}
