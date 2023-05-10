package obom

import (
	"fmt"
	"os"
	"strings"

	json "github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

const (
	OCI_ANNOTATION_DOCUMENT_NAME      = "org.opencontainers.image.spdx.name"
	OCI_ANNOTATION_DATA_LICENSE       = "org.opencontainers.image.spdx.license"
	OCI_ANNOTATION_DOCUMENT_NAMESPACE = "org.opencontainers.image.spdx.namespace"
	OCI_ANNOTATION_SPDX_VERSION       = "org.opencontainers.image.spdx.version"
	OCI_ANNOTATION_ANNOTATOR          = "org.opencontainers.image.spdx.annotator"
	OCI_ANNOTATION_ANNOTATION_DATE    = "org.opencontainers.image.spdx.annotation_date"
)

// LoadSBOM loads an SPDX file into memory
func LoadSBOM(filename string) (*v2_3.Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	doc, err := json.Read(file)
	if err != nil {
		fmt.Printf("Error while parsing SPDX file %s: %v\n", filename, err)
		return nil, err
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded %s\n", filename)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Some Attributes of the Document:")
	fmt.Printf("Document Name:         %s\n", doc.DocumentName)
	fmt.Printf("DataLicense:           %s\n", doc.DataLicense)
	fmt.Printf("Document Namespace:    %s\n", doc.DocumentNamespace)
	fmt.Printf("SPDX Version:          %s\n", doc.SPDXVersion)
	fmt.Println(strings.Repeat("=", 80))

	return doc, nil
}

// GetAnnotations returns the annotations from the SBOM
func GetAnnotations(sbom *v2_3.Document) (map[string]string, error) {
	annotations := make(map[string]string)

	annotations[OCI_ANNOTATION_DOCUMENT_NAME] = sbom.DocumentName
	annotations[OCI_ANNOTATION_DATA_LICENSE] = sbom.DataLicense
	annotations[OCI_ANNOTATION_DOCUMENT_NAMESPACE] = sbom.DocumentNamespace
	annotations[OCI_ANNOTATION_SPDX_VERSION] = sbom.SPDXVersion

	return annotations, nil
}
