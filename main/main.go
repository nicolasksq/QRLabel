package main

import (
	"log"
	"os/exec"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

func main() {
	// --- Paso 1: Generar el código QR ---
	// Definimos la URL (o ruta) del PDF.
	// Por ejemplo, si el PDF se aloja en línea, se usaría la URL pública.
	pdfURL := "https://example.com/output.pdf" // Reemplaza con la URL real de tu PDF

	// Generamos la imagen QR y la guardamos como "qr.png"
	if err := qrcode.WriteFile(pdfURL, qrcode.Medium, 256, "qr.png"); err != nil {
		log.Fatalf("Error generando el QR: %v", err)
	}

	// --- Paso 2: Crear el PDF e insertar el QR ---
	// Se crea un nuevo PDF en formato A4
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Insertamos la imagen del QR en el PDF.
	// Ajusta las coordenadas (x, y) y el ancho según tus necesidades.
	pdf.Image("qr.png", 10, 10, 50, 0, false, "", 0, "")

	// Para hacer el área del QR clickeable, definimos un enlace.
	link := pdf.AddLink()
	// Se superpone un área rectangular (en este caso, sobre la imagen del QR)
	// y se asocia el enlace.
	pdf.Link(10, 10, 50, 50, link)
	pdf.SetLink(link, 0, 0)

	// Añadimos un texto descriptivo
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 60, "Escanea el QR para acceder al PDF")

	// Guardamos el PDF generado
	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
		log.Fatalf("Error guardando el PDF: %v", err)
	}

	// --- Paso 3: Adjuntar archivos al PDF ---
	// En este ejemplo se utilizan archivos adicionales que queremos adjuntar,
	// por ejemplo "archivo1.txt" y "imagen2.jpg".
	adjuntos := []string{"archivo1.txt", "imagen2.jpg"} // Reemplaza con tus archivos

	// Se recorre la lista de adjuntos y se usa pdfcpu para agregarlos al PDF.
	for _, archivo := range adjuntos {
		// El comando invocado es similar a: pdfcpu attach add output.pdf archivo
		cmd := exec.Command("pdfcpu", "attach", "add", "output.pdf", archivo)
		if err := cmd.Run(); err != nil {
			log.Printf("Error adjuntando el archivo %s: %v", archivo, err)
		} else {
			log.Printf("Archivo %s adjuntado exitosamente", archivo)
		}
	}

	log.Println("Proceso completado: se ha generado el PDF con el QR y los archivos adjuntos.")
}
