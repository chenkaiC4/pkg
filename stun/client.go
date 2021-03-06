package stun

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

var (
	maxMessageSize = 1280

	// ErrResponseTooBig is returned if more than maxMessageSize bytes are returned in the response
	// see https://tools.ietf.org/html/rfc5389#section-7 for the size limit
	ErrResponseTooBig = errors.New("received too much data")
)

// Client is a STUN client that sents STUN requests and receives STUN responses
type Client struct {
	conn net.Conn
}

// NewClient creates a configured STUN client
func NewClient(protocol, server string, deadline time.Duration) (*Client, error) {
	dialer := &net.Dialer{
		Timeout: deadline,
	}
	conn, err := dialer.Dial(protocol, server)
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(deadline))
	conn.SetWriteDeadline(time.Now().Add(deadline))
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// Close disconnects the client
func (c *Client) Close() error {
	return c.conn.Close()
}

func generateSTUNTransactionID() []byte {
	return GenerateTransactionId()[:TransactionIDSize]
}

// Request executes a STUN request against the clients server
func (c *Client) Request() (*Message, error) {
	host, port, err := netAddrIPPort(c.conn.RemoteAddr())
	if err != nil {
		return nil, err
	}

	req, err := Build(ClassRequest, MethodBinding, GenerateTransactionId(), &XorMappedAddress{
		XorAddress: XorAddress{
			IP:   host,
			Port: port,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write(req.Pack())
	if err != nil {
		return nil, err
	}

	bs := make([]byte, maxMessageSize)
	n, err := c.conn.Read(bs)
	if err != nil {
		return nil, err
	}
	if n > maxMessageSize {
		return nil, ErrResponseTooBig
	}

	return NewMessage(bs[:n])
}
