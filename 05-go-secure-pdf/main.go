package main

import (
	"go-secure-pdf/helpers"
	pdfutils "go-secure-pdf/pdf-utils"
	"log"
	"os"
)

func main() {
	// PDF Protection Service in Go

	// pdf file paths
	securedPdfPath, unlockedPdfPath := "./secured.pdf", "./unlocked.pdf"
	helloWorldPdfBase64 := "JVBERi0xLjEKJcKlwrHDqwoKMSAwIG9iagogIDw8IC9UeXBlIC9DYXRhbG9nCiAgICAgL1BhZ2VzIDIgMCBSCiAgPj4KZW5kb2JqCgoyIDAgb2JqCiAgPDwgL1R5cGUgL1BhZ2VzCiAgICAgL0tpZHMgWzMgMCBSXQogICAgIC9Db3VudCAxCiAgICAgL01lZGlhQm94IFswIDAgMzAwIDE0NF0KICA+PgplbmRvYmoKCjMgMCBvYmoKICA8PCAgL1R5cGUgL1BhZ2UKICAgICAgL1BhcmVudCAyIDAgUgogICAgICAvUmVzb3VyY2VzCiAgICAgICA8PCAvRm9udAogICAgICAgICAgIDw8IC9GMQogICAgICAgICAgICAgICA8PCAvVHlwZSAvRm9udAogICAgICAgICAgICAgICAgICAvU3VidHlwZSAvVHlwZTEKICAgICAgICAgICAgICAgICAgL0Jhc2VGb250IC9UaW1lcy1Sb21hbgogICAgICAgICAgICAgICA+PgogICAgICAgICAgID4+CiAgICAgICA+PgogICAgICAvQ29udGVudHMgNCAwIFIKICA+PgplbmRvYmoKCjQgMCBvYmoKICA8PCAvTGVuZ3RoIDU1ID4+CnN0cmVhbQogIEJUCiAgICAvRjEgMTggVGYKICAgIDAgMCBUZAogICAgKEhlbGxvIFdvcmxkKSBUagogIEVUCmVuZHN0cmVhbQplbmRvYmoKCnhyZWYKMCA1CjAwMDAwMDAwMDAgNjU1MzUgZiAKMDAwMDAwMDAxOCAwMDAwMCBuIAowMDAwMDAwMDc3IDAwMDAwIG4gCjAwMDAwMDAxNzggMDAwMDAgbiAKMDAwMDAwMDQ1NyAwMDAwMCBuIAp0cmFpbGVyCiAgPDwgIC9Sb290IDEgMCBSCiAgICAgIC9TaXplIDUKICA+PgpzdGFydHhyZWYKNTY1CiUlRU9GCg=="
	pswd := "secure123" // FIXME: read this value from env

	// create an instance of pdfcpu handler
	handler := pdfutils.NewPDFCPUHandler()

	pdfProtector := pdfutils.NewPDFProtector(handler, false)

	pdfFileBytes, err := helpers.DecodeBase64ToBytes(helloWorldPdfBase64)
	if err != nil {
		log.Fatal("Invalid pdf file")
	}

	// 1. secure pdf
	if pdfBytes, err := pdfProtector.ProtectPDFWithPassword(pdfFileBytes, pswd, pswd); err == nil {
		log.Println("PDF protected with given password. Find the file at : ", securedPdfPath)
		// write the locked pdf bytes to a file
		helpers.WriteBytesToFile(pdfBytes, securedPdfPath)
	} else {
		log.Println(err)
	}

	// 2. unlock pdf
	securedPdfBytes, err := os.ReadFile(securedPdfPath)
	if err != nil {
		log.Fatal("File not found")
	}

	if unlockedPdfBytes, err := pdfProtector.UnlockPDF(securedPdfBytes, pswd); err == nil {
		log.Println("PDF unlocked successfully. Find the file at path : ", unlockedPdfPath)
		helpers.WriteBytesToFile(unlockedPdfBytes, unlockedPdfPath)
	} else {
		log.Fatal(err)
	}
}
