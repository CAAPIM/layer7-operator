package graphman

import (
	"crypto/sha1"
	"encoding/base64"
	"testing"
)

func certBundle(thumb string) Bundle {
	der := []byte{1, 2, 3, 4, 5}
	return Bundle{TrustedCerts: []*TrustedCertInput{{
		Name:           "apim-ssg.example",
		CertBase64:     base64.StdEncoding.EncodeToString(der),
		ThumbprintSha1: thumb,
	}}}
}

func wantThumb() string {
	sum := sha1.Sum([]byte{1, 2, 3, 4, 5})
	return base64.StdEncoding.EncodeToString(sum[:])
}

func TestTrustedCertWithoutThumbprintSurvivesMerge(t *testing.T) {
	merges := map[string]func(Bundle, Bundle) (Bundle, error){
		"CombineWithOverwrite":                     CombineWithOverwrite,
		"CombineWithOverwritePreservingDeleteMaps": CombineWithOverwritePreservingDeleteMappings,
	}
	for name, merge := range merges {
		out, err := merge(certBundle(""), Bundle{})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if len(out.TrustedCerts) != 1 {
			t.Fatalf("%s: trusted cert was dropped (got %d)", name, len(out.TrustedCerts))
		}
		if got := out.TrustedCerts[0].ThumbprintSha1; got != wantThumb() {
			t.Fatalf("%s: thumbprint = %q, want %q", name, got, wantThumb())
		}
	}
}

func TestTrustedCertExistingThumbprintUnchanged(t *testing.T) {
	out, err := CombineWithOverwrite(certBundle("keep-me"), Bundle{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.TrustedCerts) != 1 || out.TrustedCerts[0].ThumbprintSha1 != "keep-me" {
		t.Fatalf("existing thumbprint was changed: %+v", out.TrustedCerts)
	}
}

func TestTrustedCertBadBase64LeftAlone(t *testing.T) {
	b := Bundle{TrustedCerts: []*TrustedCertInput{{Name: "x", CertBase64: "%%%not-base64"}}}
	normalizeTrustedCertThumbprints(&b)
	if b.TrustedCerts[0].ThumbprintSha1 != "" {
		t.Fatalf("thumbprint set for undecodable cert: %q", b.TrustedCerts[0].ThumbprintSha1)
	}
}
