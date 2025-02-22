package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

// convertFilesToPDF genera un PDF a partir de los archivos subidos y le inserta un QR que linkea al PDF.
func convertFilesToPDF(uploadDir, pdfOutput string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		return fmt.Errorf("error al leer el directorio uploads: %w", err)
	}

	allowedExt := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
	}

	yPos := 10.0
	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if !allowedExt[ext] {
				log.Printf("Archivo %s no es una imagen soportada, se omitirá", file.Name())
				continue
			}
			imgPath := filepath.Join(uploadDir, file.Name())
			// Intenta insertar la imagen y loguea si hay error interno de gofpdf.
			pdf.Image(imgPath, 10, yPos, 50, 0, false, "", 0, "")
			log.Printf("Imagen %s insertada en el PDF", file.Name())
			yPos += 60
		}
	}

	// Define la URL pública para acceder al PDF
	pdfURL := "http://localhost:8080/" + filepath.Base(pdfOutput)
	qrPath := "./static/qr.png"

	if err := qrcode.WriteFile(pdfURL, qrcode.Medium, 256, qrPath); err != nil {
		return fmt.Errorf("error generando el QR: %w", err)
	}
	log.Printf("QR generado en %s", qrPath)

	pdf.Image(qrPath, 150, 10, 50, 0, false, "", 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.Text(150, 70, "Escanea para ver el PDF")

	if err := pdf.OutputFileAndClose(pdfOutput); err != nil {
		return fmt.Errorf("error al guardar el PDF: %w", err)
	}
	log.Printf("PDF guardado en %s", pdfOutput)
	return nil
}

// uploadHandler procesa la subida de archivos, genera el PDF y crea el QR.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parseamos el formulario con un límite de 10 MB.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	// Obtenemos los archivos (el input se llama "files").
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No se seleccionaron archivos", http.StatusBadRequest)
		return
	}

	// Creamos el directorio donde se guardarán los archivos subidos.
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Error al crear el directorio de uploads", http.StatusInternalServerError)
		return
	}

	// Guardamos cada archivo subido.
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println("Error al abrir archivo:", err)
			continue
		}
		defer file.Close()

		dstPath := filepath.Join(uploadDir, fileHeader.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			log.Println("Error al crear archivo:", err)
			continue
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Println("Error al guardar archivo:", err)
			continue
		}
		log.Printf("Archivo %s subido exitosamente", fileHeader.Filename)
	}

	// Convertimos los archivos subidos a PDF y generamos el QR para linkearlo.
	pdfOutput := "./static/output.pdf"
	if err := convertFilesToPDF(uploadDir, pdfOutput); err != nil {
		http.Error(w, "Error al convertir a PDF", http.StatusInternalServerError)
		return
	}

	// Establece el Content-Type a HTML para que el navegador renderice el contenido
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	response := `Archivos convertidos a PDF exitosamente.<br>
		Accede a tu PDF <a href="/output.pdf">aquí</a> o escanea el siguiente QR:<br>
		<img src="/qr.png" alt="QR Code">`
	fmt.Fprint(w, response)
}

func main() {
	// Servimos archivos estáticos desde la carpeta "static" (incluye index.html, output.pdf, qr.png, etc.)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Ruta para el manejo de la subida de archivos.
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
