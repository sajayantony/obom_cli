# SBOM to OCI Artiafact

This is a simple tool to convert a SBOM to a OCI Artifact and push to a target registry

# Build 

```bash
go build -ldflags "-s -w" -o obom main.go
```

## Usage

### View the SBOM to validate it. 

```bash
$ ./obom show -f ./examples/SPDXJSONExample-v2.3.spdx.json
Successfully loaded ./examples/SPDXJSONExample-v2.3.spdx.json
================================================================================
Some Attributes of the Document:
Document Name:         SPDX-Tools-v2.0
DataLicense:           CC0-1.0
Document Namespace:    http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301
SPDX Version:          SPDX-2.3
================================================================================
```

### Convert the SBOM to OCI Artifact

```bash
$ ./obom push -f ./examples/SPDXJSONExample-v2.3.spdx.json localhost:5001/spdx/annotations:test
Successfully loaded ./examples/SPDXJSONExample-v2.3.spdx.json
================================================================================
Some Attributes of the Document:
Document Name:         SPDX-Tools-v2.0
DataLicense:           CC0-1.0
Document Namespace:    http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301
SPDX Version:          SPDX-2.3
================================================================================
Adding ./examples/SPDXJSONExample-v2.3.spdx.json: sha256:2de3741a7be1be5f5e54e837524f2ec627fedfb82307dc004ae03b195abc092f
Pushing localhost:5001/spdx/annotations:test sha256:9c517423daba83776785a018190b13689ef93146465a0dd66187e2f62239a805
```

### View the OCI Artifact

```bash
$ oras manifest get localhost:5001/spdx/annotations:test --pretty
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "text/spdx",
    "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
    "size": 2
  },
  "layers": [
    {
      "mediaType": "text/spdx",
      "digest": "sha256:2de3741a7be1be5f5e54e837524f2ec627fedfb82307dc004ae03b195abc092f",
      "size": 21342,
      "annotations": {
        "org.opencontainers.image.title": "./examples/SPDXJSONExample-v2.3.spdx.json"
      }
    }
  ],
  "annotations": {
    "org.opencontainers.image.created": "2023-05-10T19:10:52Z",
    "org.opencontainers.image.spdx.license": "CC0-1.0",
    "org.opencontainers.image.spdx.name": "SPDX-Tools-v2.0",
    "org.opencontainers.image.spdx.namespace": "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301",
    "org.opencontainers.image.spdx.version": "SPDX-2.3"
  }
}
