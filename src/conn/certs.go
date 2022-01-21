package conn

import (
	"os"
	"path"
)

// Certificates in this file are created for a host named `localhost`. They are aimed only for
// local development.
//
// To regenerate them, you can run the following command from your $GOROOT/src/crypto/tls folder:
//
// go run generate_cert.go  --rsa-bits 1024 --host localhost --ca --start-date "Jan 1 00:00:00 2022" --duration=1000000h
//
// and replace below the contents of the cert.pem and key.pem files

const localCertificate = `-----BEGIN CERTIFICATE-----
MIICGDCCAYGgAwIBAgIQVibt2IqkUR7Coq10+9Y+azANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTIyMDEwMTAwMDAwMFoYDzIxMzYwMTMwMTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQC8KvxHubFrqMPo9/AhPv15KSr73vGCeee0GhWbb3sRg8tlYHvogI8m4gQq
U5DpoX5k+S5irFCFFcJh6ZWCUujQdRGN7WTFxV8I1IKH/5mhv4/uizYk0+b8Ouqc
oDOJ79nfCGU7ps8Hkx4DXiPkIh7bsqme3Oev89zYTDTkIxt/hQIDAQABo20wazAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQU2M3UwO37ic/1HOlPclLI3zbfxwYwFAYDVR0RBA0wC4IJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4GBAIbZwG5XEaZx7eZn8y0Nc/vFwcd+
NNDM/OhXjGK7l+fAO+CzDX+t7AqOz+Bo1cL3dADsDuqMc5ZzWEClbtVoFlHoGBzX
xZJ2TulCOX5mvHGzv/gTnQuYQBpx/T7KThuUia11Bk+EzD9pq5lHIo1m0rwxWRy0
WVvQJZBOS1D3Skqy
-----END CERTIFICATE-----
`

const localKey = `-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBALwq/Ee5sWuow+j3
8CE+/XkpKvve8YJ557QaFZtvexGDy2Vge+iAjybiBCpTkOmhfmT5LmKsUIUVwmHp
lYJS6NB1EY3tZMXFXwjUgof/maG/j+6LNiTT5vw66pygM4nv2d8IZTumzweTHgNe
I+QiHtuyqZ7c56/z3NhMNOQjG3+FAgMBAAECgYAR+22mkRlid3tZbTBWjQV+KbAA
5/pehLXe4UtFUm8JanXql0DgJEEJ7zmErf3ARf2lOqbzKRJ81WqBHuh5zuCOuTFA
NzygRNiaWVwzLU3u+VSj73Sq7SEARIIOakaog8OZS1ryqpBlBWN9zwUhJEpf9Dk1
ZskhfOVhGHl9RKpGwQJBAMhGEbITlCyVitDdJE2v2YCw8RLmtw1kJSO01DeZxT9d
/jVRCGjZA0g69KcwDDrboXPOX/2vJD08cL2yeZYjgTUCQQDwhpghH22fdHbc1VBX
Ko7DuwJe++fWINi+6CoWEjHO7l3c63zguQHjx3/q3cOnML32DStVfcbbvxoaxLnX
gJ8RAkEAliOJYpGw9JeLQLd4XtEk4ohDwiK6OkzIVvNaYPBjYfTp/ThpcIi4IC8q
eCfaE0nRyMp/ReRF665i6qNg6UBmvQJBAIL++unHQSAAASCCuO/QSNLHDiKHFZvk
ZceLkChXHnNyFQLV6jxF5oaUx9E1mHJ9NGhGgdxc1SonKWN80y5QadECQQDGj3j1
/ezxgTXFR4OCdPjiCI7xAYzhDp5MnwpfdDToCZzp10GcSDgJ86/enKFWXrjLtXVz
fTbjelQnStOOZKsk
-----END PRIVATE KEY-----
`

func createLocalCerts() (tmpLocalCert, tmpLocalKey string) {
	var panicOn = func(err error) {
		if err != nil {
			panic("creating local certificates for development: " + err.Error())
		}
	}
	fldr, err := os.MkdirTemp("", "local_dev_certs")
	panicOn(err)
	tmpLocalCert = path.Join(fldr, "cert.pem")
	panicOn(os.WriteFile(tmpLocalCert, []byte(localCertificate), 0600))
	tmpLocalKey = path.Join(fldr, "key.pem")
	panicOn(os.WriteFile(tmpLocalKey, []byte(localKey), 0600))
	return tmpLocalCert, tmpLocalKey
}
