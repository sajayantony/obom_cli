package obom

import (
	"context"
	"fmt"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

// PushFiles pushes the SPDX SBOM file to the registry
func PushFiles(filename string, reference string, spdx_annotations map[string]string, username string, password string) error {

	// 0. Create a file store
	fs, err := file.New(".")
	if err != nil {
		return err
	}
	defer fs.Close()
	ctx := context.Background()

	// 1. Add files to a file store
	mediaType := "text/spdx"
	fileNames := []string{filename}
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))
	for _, name := range fileNames {
		fileDescriptor, err := fs.Add(ctx, name, mediaType, "")
		if err != nil {
			return err
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
		fmt.Printf("Adding %s: %s\n", name, fileDescriptor.Digest)
	}

	annotations := make(map[string]string)
	for k, v := range spdx_annotations {
		annotations[k] = v
	}

	// 2. Pack the files and tag the packed manifest
	artifactType := "text/spdx"
	manifestDescriptor, err := oras.Pack(ctx, fs, artifactType, fileDescriptors, oras.PackOptions{
		PackImageManifest:   true,
		ManifestAnnotations: annotations,
	})
	if err != nil {
		return err
	}

	// Use the latest tag isf no tag is specified
	tag := "latest"
	ref, err := registry.ParseReference(reference)
	if err != nil {
		return err
	}

	if ref.Reference != "" {
		tag = ref.Reference
	}
	fmt.Printf("Pushing %s/%s:%s %s\n", ref.Registry, ref.Repository, tag, manifestDescriptor.Digest)
	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		return err
	}

	// 3. Connect to a remote repository
	repo, err := remote.NewRepository(reference)
	if err != nil {
		panic(err)
	}

	// Check if registry has is localhost or starts with localhost:
	reg := repo.Reference.Registry
	if reg == "localhost" || reg[:10] == "localhost:" {
		repo.PlainHTTP = true
	}

	// Note: The below code can be omitted if authentication is not required
	// Check if username and passowrd are provided
	if username == "" || password == "" {
		repo.Client = &auth.Client{
			Client: retry.DefaultClient,
			Cache:  auth.DefaultCache,
			Credential: auth.StaticCredential(reg, auth.Credential{
				Username: "username",
				Password: "password",
			}),
		}
	}

	// 3. Copy from the file store to the remote repository

	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	return err
}
