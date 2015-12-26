# xlTransport_go

Transport library for
[xLattice_go.](https://jddixon.github.io/xlattice_go)

## Basic Abstractions

XLattice
[nodes](https://jddixon.github.io/xlNode_go)
communicate by sending messages over connections.  The primitive
notions common to all xLattice connections are

* **Acceptor **, an EndPoint on which a Node listens
* **Address**, a sequence of bytes distinguishing EndPoints
* **Connection**, a pair of EndPoints enabling communications between
  the Nodes at either end
* **Connector**, an address and other information necessary to establish
  a connection to another Node
* **EndPoint**, an address and a Transport
* **Transport**, a name identifying a communications protocol

### Acceptor

An Acceptor is an EndPoint on which a Node listens.  On receiving a
connection from a far Node, the Acceptor will either reject it or begin
the procedure necessary to establish a connection.  Usually this
involves

* creating a data structure including at least the details of the
  initiating Node
* starting a separate thread (gorouotine) which will be responsible
  for handling the commuication, in particular the first message
  which the far Node sent to begin the session

### Address

An Address is a sequence of bytes sufficient to distinguish a Node within
a set of Addresses.  The Address may include a local part.  For example

* in **IPv4** the Address must include a 32-bit value identifying the host
  and may include a 16-bit port number associated with a thread or process
  on that host.  The 32-bit limit restricts the number of hosts to roughly
  4 billion
* **IPv6** the main part of the Address is a 64-bit value instead, restricting
  the number of hosts to about 16 quintillion (`10^18`)
* a simple **email** address consists of a userID and a fully qualified domain
  name (so in `john.doe@example.com` the local part is `john.doe` and the
  domain name is `example.com`.

In ipv4 and ipv6 it is necessary to distinguish the **serialization** of the
Address from its binary form.  The main part of an ipv4 address, for example,
is a 32-bit binary value in its wire form, but normally written as a
**dotted quad**, `A.B.C.D`, in which each byte is represented by its decimal
value, a number between `0` and `255` inclusive.

### Connection

A Connection is a established relationship between two EndPoints
supporting the passage of messages between the two.  From the point of
view of a participating host, the Connection will have a near EndPoint
and a far EndPoint.  Both EndPoints must use the same protocol.
Usually the Connection will have an associated State, and the Connection
will proceed through an ordered sequence of States when it is being
established, or the Connection will be closed.

### Connector

A Node will normally acquire a number of **Peers** as it is being set
up and with each a **Connector**.  Essentially this is a recipe for
establishing a Connection to the Peer; this recipe will normally
include at least two of the Peer's RSA public keys, one of which (**ck**)
is used for setting up an encrypted AES session with the Peer.

### EndPoint

An EndPoint is identified by the associated Transport and an Address.
The Address may include an explicit local part or there may be a
default local part included with the Transport.  Thus for example
the default local part associated with `http` is `80` and the default
associated with `https` is `443`.

As another example, for email the transport is `smtp` and there is
no default Address or local part.

### Transport

In the current implementation of xLattice, the only information
associated with a Transport is its name, a Unicode string.

## Project Status

A minimal implementation with only TCP/IP and UDP transport supported.

## On-line Documentation

More information on the **xlTransport_go** project can be found
[here](https://jddixon.github.io/xlTransport_go)
