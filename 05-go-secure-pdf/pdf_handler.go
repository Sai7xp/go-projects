package main

type PDFProtectionHandler interface {
	Encrypt(pdfBytes []byte, userPswd, OwnerPswd string) (encryptedPdfBytes []byte, err error)
	Decrypt(encryptedPdfBytes []byte) (decryptedPdfBytes []byte, err error)
}

type PDFProtector struct {
	handler            PDFProtectionHandler
	AllowModifications bool
}

func NewPDFProtector(handler PDFProtectionHandler) *PDFProtector {
	return &PDFProtector{
		handler: handler,
	}
}

func (p *PDFProtector) ProtectPDFWithPassword(pdfBytes []byte, userPswd, ownerPswd string) (lockedPdf []byte, err error) {
	return nil, nil
}

func (p *PDFProtector) UnlockPDF(lockedPdfBytes []byte, password string) (unlockedPdf []byte, err error) {
	return nil, nil
}
