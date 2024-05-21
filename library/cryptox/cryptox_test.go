package cryptox

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMd5(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"non-empty string", _cipherkey, "996ce17f6abc9fe126b57aa5f1d8c92c"},
		{"empty string", []byte(""), "d41d8cd98f00b204e9800998ecf8427e"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Md5(tc.input)
			if result != tc.expected {
				t.Errorf("Md5(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSha1(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"non-empty string", _cipherkey, "ff998abc1ce6d8f01a675fa197368e44c8916e9c"},
		{"empty string", []byte(""), "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Sha1(tc.input)
			if result != tc.expected {
				t.Errorf("Sha1(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"non-empty string", _cipherkey, "8e9916c5340c43fa003fe2dd54cd4e3027affbfc0d631e4cd858f64ec09fa9ed"},
		{"empty string", []byte(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Sha256(tc.input)
			if result != tc.expected {
				t.Errorf("Sha256(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSha512(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"non-empty string", _cipherkey, "6df7de339b39ae1125f181c9229cdb904fe31eac219aa2335059240101939083495221f7fbe8c32d73f8ee81dc68539c98c143f15700d944c8c0eafb27567a7a"},
		{"empty string", []byte(""), "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Sha512(tc.input)
			if result != tc.expected {
				t.Errorf("Sha512(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestFnv1aToUint64(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint64
	}{
		{"non-empty string", _cipherkey, 0x30f1a9bc9e896233},
		{"empty string", []byte(""), 0xcbf29ce484222325},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Fnv1aToUint64(tc.input)
			if result != tc.expected {
				t.Errorf("Fnv1aToUint64(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestFnv1aToUint32(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint32
	}{
		{"non-empty string", _cipherkey, 0x7f19f353},
		{"empty string", []byte(""), 0x811c9dc5},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Fnv1aToUint32(tc.input)
			if result != tc.expected {
				t.Errorf("Fnv1aToUint32(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
	}{
		{"non-empty string", []byte("thisis32bitlongpassphraseimusing"), _cipherkey},
		{"empty string", []byte("thisis32bitlongpassphraseimusing"), []byte("")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext := AESEncrypt(tc.key, tc.plaintext)
			decrypted, err := AESDecrypt(tc.key, ciphertext)
			if err != nil {
				t.Fatalf("AESDecrypt error: %v", err)
			}
			if !bytes.Equal(tc.plaintext, decrypted) {
				t.Errorf("AESDecrypt = %v; expected %v", decrypted, tc.plaintext)
			}
		})
	}
}

func TestAESCBCEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
	}{
		{"non-empty string", []byte("thisis32bitlongpassphraseimusing"), _cipherkey},
		{"empty string", []byte("thisis32bitlongpassphraseimusing"), []byte("")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext := AESCBCEncrypt(tc.key, tc.plaintext)
			decrypted, err := AESCBCDecrypt(tc.key, ciphertext)
			if err != nil {
				t.Fatalf("AESCBCDecrypt error: %v", err)
			}
			if !bytes.Equal(tc.plaintext, decrypted) {
				t.Errorf("AESCBCDecrypt = %v; expected %v", decrypted, tc.plaintext)
			}
		})
	}
}

func TestAESCTREncryptDecrypt(t *testing.T) {
	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
	}{
		{"non-empty string", []byte("thisis32bitlongpassphraseimusing"), _cipherkey},
		{"empty string", []byte("thisis32bitlongpassphraseimusing"), []byte("")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext := AESCTREncrypt(tc.key, tc.plaintext)
			decrypted, err := AESCTRDecrypt(tc.key, ciphertext)
			if err != nil {
				t.Fatalf("AESCTRDecrypt error: %v", err)
			}
			if !bytes.Equal(tc.plaintext, decrypted) {
				t.Errorf("AESCTRDecrypt = %v; expected %v", decrypted, tc.plaintext)
			}
		})
	}
}

var (
	_cipherkey = []byte("1234567890abcdef")
	_plaintext = []byte("text1234")
)

func TestPading(t *testing.T) {
	blockSize := 16
	padded := hexEncode(pkcs5Padding(_plaintext, blockSize))
	t.Log(string(padded))
	r, err := hexDecode(padded)
	assert.NoError(t, err)
	unpaded, err := pkcs5Unpadding(r)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, unpaded)
}
