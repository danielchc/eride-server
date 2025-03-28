package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

type TLSConfigData struct {
	CA         string
	ServerCert string
	ServerKey  string
}

func LoadTLSCredentials(tlsConfigData *TLSConfigData) (credentials.TransportCredentials, error) {

	caCert, err := os.ReadFile(tlsConfigData.CA)

	if err != nil {
		return nil, fmt.Errorf("error reading CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	serverCert, err := tls.LoadX509KeyPair(tlsConfigData.ServerCert, tlsConfigData.ServerKey)
	if err != nil {
		return nil, fmt.Errorf("error loading server key pair: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}
