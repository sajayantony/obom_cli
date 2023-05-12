package obom

import (
	"fmt"
	"os"
	"strings"

	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/opencontainers/go-digest"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	json "github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

const (
	MEDIATYPE_SPDX                    = "application/spdx+json"
	OCI_ANNOTATION_DOCUMENT_NAME      = "org.spdx.name"
	OCI_ANNOTATION_DATA_LICENSE       = "org.spdx.license"
	OCI_ANNOTATION_DOCUMENT_NAMESPACE = "org.spdx.namespace"
	OCI_ANNOTATION_SPDX_VERSION       = "org.spdx.version"
	OCI_ANNOTATION_CREATION_DATE      = "org.spdx.created"
	OCI_ANNOTATION_ANNOTATOR          = "org.spdx.annotator"
	OCI_ANNOTATION_ANNOTATION_DATE    = "org.spdx.annotation_date"
)

// LoadSBOM loads an SPDX file into memory
func LoadSBOM(filename string) (*v2_3.Document, *oci.Descriptor, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	doc, err := json.Read(file)
	if err != nil {
		fmt.Printf("Error while parsing SPDX file %s: %v\n", filename, err)
		return nil, nil, err
	}

	desc, err := GetFileDescriptor(filename)
	if err != nil {
		return nil, nil, err
	}

	return doc, desc, nil
}

// PrintSBOMSummary returns the SPDX summary from the SBOM
func PrintSBOMSummary(doc *v2_3.Document, desc *oci.Descriptor) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Document Name:         %s\n", doc.DocumentName)
	fmt.Printf("DataLicense:           %s\n", doc.DataLicense)
	fmt.Printf("Document Namespace:    %s\n", doc.DocumentNamespace)
	fmt.Printf("SPDX Version:          %s\n", doc.SPDXVersion)
	fmt.Printf("Creation Date:         %s\n", doc.CreationInfo.Created)
	fmt.Printf("Packages:              %d\n", len(doc.Packages))
	fmt.Printf("Files:                 %d\n", len(doc.Files))
	fmt.Printf("Digest:                %s\n", desc.Digest)
	fmt.Println(strings.Repeat("=", 80))
}

func GetFileDescriptor(filename string) (*oci.Descriptor, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new SHA256 hasher
	hasher := sha256.New()

	// Copy the file's contents into the hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	// Get the resulting hash as a byte slice
	hash := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash)

	d := digest.NewDigestFromHex("sha256", hashString)

	fInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fInfo.Size()
	desc := &oci.Descriptor{
		MediaType: MEDIATYPE_SPDX,
		Digest:    d,
		Size:      fileSize,
	}

	return desc, nil
}

// GetAnnotations returns the annotations from the SBOM
func GetAnnotations(sbom *v2_3.Document) (map[string]string, error) {
	annotations := make(map[string]string)

	annotations[OCI_ANNOTATION_DOCUMENT_NAME] = sbom.DocumentName
	annotations[OCI_ANNOTATION_DATA_LICENSE] = sbom.DataLicense
	annotations[OCI_ANNOTATION_DOCUMENT_NAMESPACE] = sbom.DocumentNamespace
	annotations[OCI_ANNOTATION_SPDX_VERSION] = sbom.SPDXVersion
	annotations[OCI_ANNOTATION_CREATION_DATE] = sbom.CreationInfo.Created

	return annotations, nil
}

// GetPackages returns the packages from the SBOM
func GetPackages(sbom *v2_3.Document) ([]string, error) {
	var packages []string

	for _, pkg := range sbom.Packages {
		if pkg.PackageExternalReferences != nil {
			for _, exRef := range pkg.PackageExternalReferences {
				packages = append(packages, exRef.Locator)
			}
		}
	}

	return packages, nil
}

func GetFiles(sbom *v2_3.Document) ([]string, error) {
	var files []string

	for _, file := range sbom.Files {
		files = append(files, file.FileName)
	}

	return files, nil
}
