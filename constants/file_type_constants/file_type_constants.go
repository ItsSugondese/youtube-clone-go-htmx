package file_type_constants

type FileType string

const (
	IMAGE FileType = "IMAGE"
	DOC   FileType = "DOC"
	EXCEL FileType = "EXCEL"
	PDF   FileType = "PDF"
	TXT   FileType = "TXT"
)

// Maps to store file types
var (
	ImageType = map[string]FileType{
		"JPEG": IMAGE,
		"JPG":  IMAGE,
		"PNG":  IMAGE,
		"SVG":  IMAGE,
	}

	DocumentType = map[string]FileType{
		"DOC":  DOC,
		"DOCX": DOC,
	}

	PdfType = map[string]FileType{
		"PDF": PDF,
	}

	TxtType = map[string]FileType{
		"TXT": TXT,
	}

	ExcelType = map[string]FileType{
		"XLS":  EXCEL,
		"XLSX": EXCEL,
		"CSV":  EXCEL,
		"ODS":  EXCEL,
	}
)
