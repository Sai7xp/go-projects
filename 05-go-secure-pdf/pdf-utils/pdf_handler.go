package pdfutils

import (
	"bytes"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFProtectionHandler interface {
	Encrypt(pdfBytes []byte, userPswd, OwnerPswd string) (encryptedPdfBytes []byte, err error)
	Decrypt(encryptedPdfBytes []byte, password string) (decryptedPdfBytes []byte, err error)
}

type PDFProtector struct {
	handler            PDFProtectionHandler
	allowModifications bool
}

func NewPDFProtector(handler PDFProtectionHandler, allowModifications bool) *PDFProtector {
	return &PDFProtector{
		handler:            handler,
		allowModifications: allowModifications,
	}
}

func (p *PDFProtector) ProtectPDFWithPassword(pdfBytes []byte, userPswd, ownerPswd string) (lockedPdf []byte, err error) {
	encryptedPdfBytes, err := p.handler.Encrypt(pdfBytes, userPswd, ownerPswd)
	if err != nil {
		log.Println("Error while adding password to PDF: ", err)
		return nil, err
	}
	return encryptedPdfBytes, nil
}

func (p *PDFProtector) UnlockPDF(lockedPdfBytes []byte, password string) (unlockedPdf []byte, err error) {
	decryptedPDFBytes, err := p.handler.Decrypt(lockedPdfBytes, password)
	if err != nil {
		log.Println("Error while unlocking PDF: ", err)
		return nil, err
	}
	return decryptedPDFBytes, nil
}

// Dependencies Implementation

type PDFCPUHandler struct {
	config *model.Configuration
}

func NewPDFCPUHandler() *PDFCPUHandler {
	return &PDFCPUHandler{
		config: model.NewDefaultConfiguration(),
	}
}

func (c *PDFCPUHandler) Encrypt(pdfBytes []byte, userPswd, OwnerPswd string) (encryptedPdfBytes []byte, err error) {
	c.config.OwnerPW = OwnerPswd
	c.config.UserPW = userPswd
	c.config.Cmd = model.ENCRYPT

	r := bytes.NewReader(pdfBytes)
	w := bytes.NewBuffer(nil)

	if err := api.Encrypt(r, w, c.config); err != nil {
		log.Println("Error while encrypting the pdf using cpupdf: ", err)
		return nil, err
	}

	return w.Bytes(), nil
}

func (c *PDFCPUHandler) Decrypt(encryptedPdfBytes []byte, password string) (decryptedPdfBytes []byte, err error) {
	c.config.UserPW = password
	r := bytes.NewReader(encryptedPdfBytes)
	w := bytes.NewBuffer(nil)

	if err := api.Decrypt(r, w, c.config); err != nil {
		log.Println("Error while decypting the pdf using cpupdf: ", err)
	}
	
	return w.Bytes(), nil
}
