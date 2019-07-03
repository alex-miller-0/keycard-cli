package main

import (
	keycard "github.com/alex-miller-0/keycard-go"
	"github.com/alex-miller-0/keycard-go/io"
	"github.com/alex-miller-0/keycard-go/types"
	"math/rand"
)

type Certifier struct {
	c types.Channel
}

func NewCertifier(t io.Transmitter) *Certifier {
	return &Certifier{
		c: io.NewNormalChannel(t),
	}
}

// Get the public key corresponding to the card's identity
func (c *Certifier) GetId() ([]byte, error) {
	// Send some random data to `authenticate` and get a signature template back
	cmdSet := keycard.NewCommandSet(c.c)
	
	challenge := make([]byte, 32)
	rand.Read(challenge)
	data, err := cmdSet.GenericCommand(0x00, 0xEE, 0x00, 0x00, challenge)
	if err != nil {
		return nil, err
	}
	sig, err := types.ParseSignature(challenge, data)
	if err != nil {
		return nil, err
	}
	return sig.PubKey(), nil
}

// Put a signature onto the card which contains the certs
func (c *Certifier) PutCert(cert []byte) (error) {
	cmdSet := keycard.NewCommandSet(c.c)
	_, err := cmdSet.GenericCommand(0x80, 0xFA, 0x00, 0x00, cert)
	return err
}

// Get the cert from the card
func (c *Certifier) GetCert() ([]byte, error) {
	cmdSet := keycard.NewCommandSet(c.c)
	data, err := cmdSet.GenericCommand(0x80, 0xFB, 0x00, 0x00, []byte{})
	if err != nil {
		return nil, err
	}
	return data, nil
}

