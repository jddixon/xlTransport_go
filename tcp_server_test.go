package transport

// xlTransport_go/server_test.go

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
	"time"
)

// Start an acceptor running on a random port.  Create K*N blocks of
// random data.  These will be sent by K clients to the server.  The
// server will reply to each with the SHA1 hash of the block.  Run
// the clients in separate goroutines.  After all clients have sent
// all of their messages, verify that the hashes received back are
// correct.

const (
	// 2017-01-25: MAX_LEN 8192
	//             K=8,N=128 succeeds 0.279s
	//             K=8,N=256 hangs
	//             K=8,N=129 hangs
	//             K=4,N=256 succeeds 0.305s
	//             K=4,N=512 hangs
	//             K=4,N=257 succeeds 0.290s
	//             K=4,N=260 succeeds 0.281s
	//             K=4,N=270 succeeds 0.304s
	//             K=4,N=280 succeeds, repeatable; 0.38, 0.322, 0.306s, etc
	//             K=4,N=282 hangs sometimes; succeeded at 0.327s, 0.340s
	//             K=4,N=285 hangs reliably
	//             K=4,N=290 hangs reliably

	//             K=16,N=64 succeeds
	//			   K=16,N=65 hangs unreliably (say 4/5 times)
	//			   K=16,N=66 hangs unreliably
	//			   K=16,N=70 hangs fairly reliably (9/10 times)

	//             MAX_LEN 2048
	//             K=16,N=64 succeeds
	//             K=32,N=32 succeeds
	//             K=33,N=32 hangs
	//             K=64,N=16 hangs

	//			   MAX_LEN 65536
	//			   K=16,N=64 succeeds reliably (10/10), about 1.100s
	//             K=64,N=16 hangs unreliably (3/10), gets bad hash (7/10)
	//			   K=16,N=64 hangs unreliably (1/10), gets bad hash (9/10)

	// XXX 2013-07-20 test hangs if K=16,N=32 and K increasd to 32 OR
	// N increased from 32 to 64
	// 2016-11-14 test succeeds K=16,N=32; K=32,N=32; K=16,N=64; K=32,N=64
	//     ** K=64,N=128 HANGS ** but K=64,N=64 succeeds
	//     ** K=128,N=64 HANGS **
	K        = 16   // number of clients
	N        = 64   // number of messages for each client
	MIN_LEN  = 1024 // minimum length of message
	MAX_LEN  = 8192 // maximum
	SHA1_LEN = 20
)

var rng = xr.MakeSimpleRNG()

func (s *XLSuite) handleMsg(cnx ConnectionI) error {
	myCnx := cnx.(*TcpConnection)
	defer myCnx.Close()

	buf := make([]byte, MAX_LEN)

	// read the message
	count, err := myCnx.Read(buf)
	buf = buf[:count] // ESSENTIAL
	if err == nil {
		// calculate its hash
		d := sha1.New()
		d.Write(buf)
		digest := d.Sum(nil) // a binary value

		// send the digest as a reply
		count, err = myCnx.Write(digest)

		_ = count // XXX verify length of 20
	}
	// XXX allow the other end to read the reply; it would be
	// better to loop until a 'closed connection' error is returned
	time.Sleep(100 * time.Millisecond)
	return err
}

func (s *XLSuite) TestHashingServer(c *C) {
	SERVER_ADDR := "127.0.0.1:0"

	fmt.Println("TEST_HASHING_SERVER")
	fmt.Printf("    K %d, N %d, MAX_LEN %d\n", K, N, MAX_LEN)

	// -- setup  -----------------------------------------------------
	// fmt.Println("building messages")
	var messages [][][]byte = make([][][]byte, K)
	var hashes [][][]byte = make([][][]byte, K)
	for i := 0; i < K; i++ {
		messages[i] = make([][]byte, N)
		for j := 0; j < N; j++ {
			msgLen := MIN_LEN + rng.Intn(MAX_LEN-MIN_LEN)
			messages[i][j] = make([]byte, msgLen)
			rng.NextBytes(messages[i][j])
		}
		hashes[i] = make([][]byte, N)
		for j := 0; j < N; j++ {
			hashes[i][j] = make([]byte, SHA1_LEN)
		}
	}

	// -- create and start server -----------------------------------
	acc, err := NewTcpAcceptor(SERVER_ADDR)
	c.Assert(err, Equals, nil)
	defer acc.Close()
	accEndPoint := acc.GetEndPoint()
	//fmt.Printf("server_test acceptor listening on %s\n", accEndPoint.String())
	go func() {
		for {
			cnx, err := acc.Accept()
			if err != nil { // ESSENTIAL
				break
			}
			if cnx != nil {
				go func(cnx ConnectionI) {
					_ = s.handleMsg(cnx)
					// c.Assert(err, Equals, nil)
				}(cnx)
			}
		}
	}()

	// -- create K client connectors --------------------------------
	ktors := make([]*TcpConnector, K)
	for i := 0; i < K; i++ {
		ktors[i], err = NewTcpConnector(accEndPoint)
		c.Assert(err, Equals, nil)
	}

	// -- start the clients -----------------------------------------
	var clientDone [K]chan bool
	for i := 0; i < K; i++ {
		clientDone[i] = make(chan bool)
	}
	for i := 0; i < K; i++ {
		go func(i int) {
			for j := 0; j < N; j++ {
				// the client sends N messages, expecting an SHA1 back
				var count int
				cnx, err := ktors[i].Connect(ANY_TCP_END_POINT)
				c.Assert(err, Equals, nil)
				tcpCnx := cnx.(*TcpConnection)
				count, err = tcpCnx.Write(messages[i][j])
				if err != nil {
					fmt.Printf("error writing [%d][%d]: %v\n", i, j, err)
				}
				count, err = tcpCnx.Read(hashes[i][j])
				if err != nil {
					fmt.Printf("error reading [%d][%d]: %v\n", i, j, err)
				}
				cnx.Close()

				_, _ = count, err // XXX
			}
			clientDone[i] <- true
		}(i)
	}
	// -- when all clients have completed, shut down server ---------
	for i := 0; i < K; i++ {
		<-clientDone[i]
	}
	if acc != nil {
		if err = acc.Close(); err != nil {
			fmt.Printf("unexpected error closing acceptor: %v\n", err)
		} else {
			fmt.Printf("acceptor closed successfully\n")
		}
	}
	// -- calculate and verify K*N hashes ---------------------------
	for i := 0; i < K; i++ {
		for j := 0; j < N; j++ {
			d := sha1.New()
			d.Write(messages[i][j])
			digest := d.Sum(nil)                // a binary value
			hashX := hex.EncodeToString(digest) // DEBUG
			hashY := hex.EncodeToString(hashes[i][j])
			c.Assert(hashX, Equals, hashY)
		}
	}
}
