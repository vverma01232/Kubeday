package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"github.com/skip2/go-qrcode"
)

func generateQRCode(url string, size int) ([]byte, error) {
	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	return qrCode.PNG(size)
}

// Base64Encode encodes a byte slice as a base64 string
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://www.initializ.ai/early-access"
	size := 256 // Set the desired size of the QR code

	// Generate the QR code
	qrCodeImage, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Define the data to be passed to the template
	data := struct {
		Title   string
		Message string
		QRCode  []byte
	}{
		Title:   "initializ",
		Message: "Be First in Line for Early Access to Initializ",
		QRCode:  qrCodeImage,
	}

	// Parse the HTML template with the custom Base64Encode function
	tmpl, err := template.New("index").Funcs(template.FuncMap{"Base64Encode": Base64Encode}).Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>{{.Title}}</title>
	<style>
		body {
			background-color: black;
			color: white;
			font-family: 'Poppins', Arial, sans-serif;
			text-align: center;
			display: flex;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			height: 100vh; /* Set to 100% of the viewport height */
			margin: 0;
		}
		.header {
			display: flex;
			flex-direction: row;
			align-items: center;
			margin-bottom: 25px;
		}
		.logo {
			width: auto;
			height: 50px; /* Match the font size of the heading */
			margin-right: 10px; /* Added margin to separate logo and heading */
		}
		.heading {
			font-size: 50px;
			color: #ce4da4; /* Updated color code */
			font-weight: bold;
		}
		.message {
			font-size: 24px;
		}
		.qr-code {
			margin-top: 20px;
		}
		.link {
			margin-top: 10px;
			font-size: 20px;
		}
		.link a {
			color: white; /* Set the color to white */
			text-decoration: none; 
		}
	</style>
	<link rel="stylesheet" type="text/css" href="/fonts/poppins.css"> <!-- Link to Poppins font CSS file -->
</head>
<body>
	<div class="header">
		<img src="https://assets-global.website-files.com/64fe21b82836dd338c66292b/64fe28b48c5ca2fdb2e85d04_Logo%20only.svg" alt="Logo" class="logo" /> <!-- Replace with the actual path to your image -->
		<div class="heading">{{.Title}}</div>
	</div>
	<p class="message">{{.Message}}</p>
	<div class="qr-code">
		<img src="data:image/png;base64,{{Base64Encode .QRCode}}" alt="QR Code" />
	</div>
	<div class="link">
    <a href="https://www.initializ.ai" target="_blank">initializ.ai</a>
</div>
</body>
</html>
`)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data and write the response
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", HomePageHandler)
	fmt.Println("server up and running on port 8080")
	http.ListenAndServe(":8080", nil)
	fmt.Println("server up and running on port 8080")
}
