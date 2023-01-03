// SPDX-License-Identifier: Apache-2.0
package connect

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ConnectionArgs struct {
	Key       string
	Cert      string
	Endpoint  string
	CAFile    string
	ThingName string
	Port      int
}

var maxClientIdLen = 8

// getRandomClientId returns randomized ClientId.
func getRandomClientId() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, maxClientIdLen)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return "test-" + string(bytes)
}

func getCertPool(pemPath string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()

	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	return certs, nil
}

func getTLSConfig(args ConnectionArgs) (*tls.Config, error) {
	if args.CAFile == "" {
		return nil, fmt.Errorf("please specify CA file")
	}
	fmt.Printf("CA: %s\n", args.CAFile)
	fmt.Printf("cert: %s\n", args.Cert)
	fmt.Printf("key: %s\n", args.Key)

	// CA
	caPool, err := getCertPool(args.CAFile)
	if err != nil {
		return nil, err
	}

	certPool, err := getCertPool(args.Cert)
	if err != nil {
		return nil, err
	}

	// Cert and Private Key
	cert, err := tls.LoadX509KeyPair(args.Cert, args.Key)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:            caPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	return tlsConfig, nil
}

func getClientOptions(args ConnectionArgs) (*mqtt.ClientOptions, error) {
	opt := mqtt.NewClientOptions()
	opt.SetAutoReconnect(true)
	opt.SetClientID(getRandomClientId())
	opt.AddBroker("ssl://" + args.Endpoint + ":8883")

	tls, err := getTLSConfig(args)
	if err != nil {
		return nil, err
	}
	opt.SetTLSConfig(tls)

	return opt, nil
}

func Connect(args ConnectionArgs) (mqtt.Client, error) {
	// 1. Get connection options
	opts, err := getClientOptions(args)
	if err != nil {
		return nil, err
	}

	// 2. Set Paho MQTT client
	mc := mqtt.NewClient(opts)

	// 4. Connect to MQTT
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return mc, nil
}
