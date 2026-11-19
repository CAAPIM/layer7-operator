package gateway

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func newChecksumSecret(uid string, data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default", UID: types.UID(uid)},
		Data:       data,
	}
}

func TestSecretLabelChecksum(t *testing.T) {
	data := map[string][]byte{"SSG_ADMIN_USERNAME": []byte("admin"), "SSG_ADMIN_PASSWORD": []byte("7layer")}
	changed := map[string][]byte{"SSG_ADMIN_USERNAME": []byte("admin"), "SSG_ADMIN_PASSWORD": []byte("changed")}

	base := SecretLabelChecksum(newChecksumSecret("uid-1", data))

	t.Run("fits in a label value", func(t *testing.T) {
		if base == "" || len(base) > 63 {
			t.Errorf("unexpected checksum length %d", len(base))
		}
	})

	t.Run("is stable for the same secret", func(t *testing.T) {
		if got := SecretLabelChecksum(newChecksumSecret("uid-1", data)); got != base {
			t.Errorf("Expected %s, Actual %s", base, got)
		}
	})

	t.Run("changes when data changes", func(t *testing.T) {
		if SecretLabelChecksum(newChecksumSecret("uid-1", changed)) == base {
			t.Error("checksum did not change when secret data changed")
		}
	})

	t.Run("differs for different secret uids", func(t *testing.T) {
		if SecretLabelChecksum(newChecksumSecret("uid-2", data)) == base {
			t.Error("checksum must be keyed by the secret uid")
		}
	})

	t.Run("is not the unkeyed sha1 annotation value", func(t *testing.T) {
		dataBytes, _ := json.Marshal(data)
		h := sha1.New()
		h.Write(dataBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		if strings.HasPrefix(sha1Sum, base) || strings.HasPrefix(base, sha1Sum) {
			t.Error("checksum must not be derived from secret data alone")
		}
	})

	t.Run("is empty without a uid", func(t *testing.T) {
		if got := SecretLabelChecksum(newChecksumSecret("", data)); got != "" {
			t.Errorf("Expected empty checksum, Actual %s", got)
		}
	})
}
