package eval

import (
	"io/ioutil"
	"net"
	"time"
)

const (
	TCPCONN_OBJ     = "TCP_CONN_OBJ"
	TCPLISTENER_OBJ = "TCP_LISTENER_OBJ"

	UDPCONN_OBJ = "UDP_CONN_OBJ"

	UNIXCONN_OBJ     = "UNIX_CONN_OBJ"
	UNIXLISTENER_OBJ = "UNIX_LISTENER_OBJ"
)

//Net Object
const (
	NET_OBJ = "NET_OBJ"
	net_name = "net"
)
type NetObj struct{}

func NewNetObj() Object {
	ret := &NetObj{}
	SetGlobalObj(net_name, ret)

	return ret
}

func (n *NetObj) Inspect() string  { return "<" + net_name + ">" }
func (n *NetObj) Type() ObjectType { return NET_OBJ }

func (n *NetObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "joinHostPort":
		return n.JoinHostPort(line, args...)
	case "lookupAddr":
		return n.LookupAddr(line, args...)
	case "lookupHost":
		return n.LookupHost(line, args...)
	case "lookupIP":
		return n.LookupIP(line, args...)
	case "lookupPort":
		return n.LookupPort(line, args...)
	case "splitHostPort":
		return n.SplitHostPort(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, n.Type()))
}

func (n *NetObj) JoinHostPort(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	host, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "joinHostPort", "*String", args[0].Type()))
	}

	port, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "joinHostPort", "*String", args[1].Type()))
	}

	str := net.JoinHostPort(host.String, port.String)
	return NewString(str)
}

func (n *NetObj) LookupAddr(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	addr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lookupAddr", "*String", args[0].Type()))
	}

	names, err := net.LookupAddr(addr.String)
	if err != nil {
		return NewNil(err.Error())
	}

	arr := &Array{}
	for _, v := range names {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (n *NetObj) LookupHost(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	host, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lookupHost", "*String", args[0].Type()))
	}

	addrs, err := net.LookupHost(host.String)
	if err != nil {
		return NewNil(err.Error())
	}

	arr := &Array{}
	for _, v := range addrs {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (n *NetObj) LookupIP(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	host, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lookupIP", "*String", args[0].Type()))
	}

	ips, err := net.LookupIP(host.String)
	if err != nil {
		return NewNil(err.Error())
	}

	arr := &Array{}
	for _, v := range ips {
		ip := string(v)
		arr.Members = append(arr.Members, NewString(ip))
	}

	return arr
}

func (n *NetObj) LookupPort(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	network, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lookupPort", "*String", args[0].Type()))
	}

	service, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "lookupPort", "*String", args[1].Type()))
	}

	port, err := net.LookupPort(network.String, service.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(port))
}

func (n *NetObj) SplitHostPort(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	hostport, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "splitHostPort", "*String", args[0].Type()))
	}

	host, port, err := net.SplitHostPort(hostport.String)
	if err != nil {
		return NewNil(err.Error())
	}

	arr := &Array{}
	arr.Members = append(arr.Members, NewString(host))
	arr.Members = append(arr.Members, NewString(port))

	return arr
}

//TCP connection object
type TcpConnObject struct {
	Conn    *net.TCPConn
	Address string
}

// Implement the 'Closeable' interface
func (t *TcpConnObject) close(line string, args ...Object) Object {
	return t.Close(line, args...)
}

func (t *TcpConnObject) Inspect() string  { return t.Address }
func (t *TcpConnObject) Type() ObjectType { return TCPCONN_OBJ }
func (t *TcpConnObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "addr":
		return t.Addr(line, args...)
	case "close":
		return t.Close(line, args...)
	case "closeRead":
		return t.CloseRead(line, args...)
	case "closeWrite":
		return t.CloseWrite(line, args...)
	case "read":
		return t.Read(line, args...)
	case "write":
		return t.Write(line, args...)
	case "setDeadline":
		return t.SetDeadline(line, args...)
	case "setReadDeadline":
		return t.SetReadDeadline(line, args...)
	case "setWriteDeadline":
		return t.SetWriteDeadline(line, args...)
	case "setLinger":
		return t.SetLinger(line, args...)
	case "setNoDelay":
		return t.SetNoDelay(line, args...)
	case "setReadBuffer":
		return t.SetReadBuffer(line, args...)
	case "setWriteBuffer":
		return t.SetWriteBuffer(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, t.Type()))
	}
}

//Return the remote address
func (t *TcpConnObject) Addr(line string, args ...Object) Object {
	return NewString(t.Address)
}

func (t *TcpConnObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := t.Conn.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) CloseRead(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := t.Conn.CloseRead()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) CloseWrite(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := t.Conn.CloseWrite()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) Read(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	bytes, err := ioutil.ReadAll(t.Conn)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(string(bytes))
}

func (t *TcpConnObject) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	str, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	n, err := t.Conn.Write([]byte(str.String))
	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}

func (t *TcpConnObject) SetDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := t.Conn.SetDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) SetReadDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := t.Conn.SetReadDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) SetWriteDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := t.Conn.SetWriteDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) SetLinger(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setLinger", "*Integer", args[0].Type()))
	}

	err := t.Conn.SetLinger(int(sec.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) SetNoDelay(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	noDelay, ok := args[0].(*Boolean)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setNoDelay", "*Boolean", args[0].Type()))
	}

	err := t.Conn.SetNoDelay(noDelay.Bool)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) SetReadBuffer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	bytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadBuffer", "*Integer", args[0].Type()))
	}

	err := t.Conn.SetReadBuffer(int(bytes.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *TcpConnObject) SetWriteBuffer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	bytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteBuffer", "*Integer", args[0].Type()))
	}

	err := t.Conn.SetWriteBuffer(int(bytes.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//TCP Listener object
type TCPListenerObject struct {
	Listener *net.TCPListener
	Address  string
}

// Implement the 'Closeable' interface
func (l *TCPListenerObject) close(line string, args ...Object) Object {
	return l.Close(line, args...)
}

func (l *TCPListenerObject) Inspect() string  { return l.Address }
func (l *TCPListenerObject) Type() ObjectType { return TCPLISTENER_OBJ }
func (l *TCPListenerObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "addr":
		return l.Addr(line, args...)
	case "close":
		return l.Close(line, args...)
	case "acceptTCP":
		return l.AcceptTCP(line, args...)
	case "setDeadline":
		return l.SetDeadline(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, l.Type()))
	}
}

func (l *TCPListenerObject) Addr(line string, args ...Object) Object {
	return NewString(l.Address)
}

func (l *TCPListenerObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := l.Listener.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (l *TCPListenerObject) AcceptTCP(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	tcpConn, err := l.Listener.AcceptTCP()
	if err != nil {
		return NewNil(err.Error())
	}

	return &TcpConnObject{Conn: tcpConn, Address: l.Address}
}

func (l *TCPListenerObject) SetDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := l.Listener.SetDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//UDP connection object
type UdpConnObject struct {
	Conn    *net.UDPConn
	Address string
}

// Implement the 'Closeable' interface
func (u *UdpConnObject) close(line string, args ...Object) Object {
	return u.Close(line, args...)
}

func (u *UdpConnObject) Inspect() string  { return u.Address }
func (u *UdpConnObject) Type() ObjectType { return UDPCONN_OBJ }
func (u *UdpConnObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "addr":
		return u.Addr(line, args...)
	case "close":
		return u.Close(line, args...)
	case "read":
		return u.Read(line, args...)
	case "write":
		return u.Write(line, args...)
	case "setDeadline":
		return u.SetDeadline(line, args...)
	case "setReadDeadline":
		return u.SetReadDeadline(line, args...)
	case "setWriteDeadline":
		return u.SetWriteDeadline(line, args...)
	case "setReadBuffer":
		return u.SetReadBuffer(line, args...)
	case "setWriteBuffer":
		return u.SetWriteBuffer(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, u.Type()))
	}
}

//Return the remote address
func (u *UdpConnObject) Addr(line string, args ...Object) Object {
	return NewString(u.Address)
}

func (u *UdpConnObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := u.Conn.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UdpConnObject) Read(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	bytes, err := ioutil.ReadAll(u.Conn)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(string(bytes))
}

func (u *UdpConnObject) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	str, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	n, err := u.Conn.Write([]byte(str.String))
	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}

func (u *UdpConnObject) SetDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := u.Conn.SetDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UdpConnObject) SetReadDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := u.Conn.SetReadDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UdpConnObject) SetWriteDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := u.Conn.SetWriteDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UdpConnObject) SetReadBuffer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	bytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadBuffer", "*Integer", args[0].Type()))
	}

	err := u.Conn.SetReadBuffer(int(bytes.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UdpConnObject) SetWriteBuffer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	bytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteBuffer", "*Integer", args[0].Type()))
	}

	err := u.Conn.SetWriteBuffer(int(bytes.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//UNIX connection object
type UnixConnObject struct {
	Conn    *net.UnixConn
	Address string
}

// Implement the 'Closeable' interface
func (u *UnixConnObject) close(line string, args ...Object) Object {
	return u.Close(line, args...)
}

func (u *UnixConnObject) Inspect() string  { return u.Address }
func (u *UnixConnObject) Type() ObjectType { return UNIXCONN_OBJ }
func (u *UnixConnObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "addr":
		return u.Addr(line, args...)
	case "close":
		return u.Close(line, args...)
	case "closeRead":
		return u.CloseRead(line, args...)
	case "closeWrite":
		return u.CloseWrite(line, args...)
	case "read":
		return u.Read(line, args...)
	case "write":
		return u.Write(line, args...)
	case "setDeadline":
		return u.SetDeadline(line, args...)
	case "setReadDeadline":
		return u.SetReadDeadline(line, args...)
	case "setWriteDeadline":
		return u.SetWriteDeadline(line, args...)
	case "setReadBuffer":
		return u.SetReadBuffer(line, args...)
	case "setWriteBuffer":
		return u.SetWriteBuffer(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, u.Type()))
	}
}

//Return the remote address
func (u *UnixConnObject) Addr(line string, args ...Object) Object {
	return NewString(u.Address)
}

func (u *UnixConnObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := u.Conn.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UnixConnObject) CloseRead(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := u.Conn.CloseRead()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (u *UnixConnObject) CloseWrite(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := u.Conn.CloseWrite()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (u *UnixConnObject) Read(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	bytes, err := ioutil.ReadAll(u.Conn)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(string(bytes))
}

func (u *UnixConnObject) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	str, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	n, err := u.Conn.Write([]byte(str.String))
	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}

func (u *UnixConnObject) SetDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := u.Conn.SetDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UnixConnObject) SetReadDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := u.Conn.SetReadDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UnixConnObject) SetWriteDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := u.Conn.SetWriteDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UnixConnObject) SetReadBuffer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	bytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadBuffer", "*Integer", args[0].Type()))
	}

	err := u.Conn.SetReadBuffer(int(bytes.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (u *UnixConnObject) SetWriteBuffer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	bytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteBuffer", "*Integer", args[0].Type()))
	}

	err := u.Conn.SetWriteBuffer(int(bytes.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//UNIX Listener object
type UnixListenerObject struct {
	Listener *net.UnixListener
	Address  string
}

// Implement the 'Closeable' interface
func (l *UnixListenerObject) close(line string, args ...Object) Object {
	return l.Close(line, args...)
}

func (l *UnixListenerObject) Inspect() string  { return l.Address }
func (l *UnixListenerObject) Type() ObjectType { return UNIXLISTENER_OBJ }
func (l *UnixListenerObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "addr":
		return l.Addr(line, args...)
	case "close":
		return l.Close(line, args...)
	case "acceptUnix":
		return l.AcceptUnix(line, args...)
	case "setDeadline":
		return l.SetDeadline(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, l.Type()))
	}
}

func (l *UnixListenerObject) Addr(line string, args ...Object) Object {
	return NewString(l.Address)
}

func (l *UnixListenerObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := l.Listener.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (l *UnixListenerObject) AcceptUnix(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	unixConn, err := l.Listener.AcceptUnix()
	if err != nil {
		return NewNil(err.Error())
	}

	return &UnixConnObject{Conn: unixConn, Address: l.Address}
}

func (l *UnixListenerObject) SetDeadline(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sec, ok := args[0].(*Integer) //seconds
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setDeadline", "*Integer", args[0].Type()))
	}

	d := time.Duration(int64(time.Second) * sec.Int64)
	err := l.Listener.SetDeadline(time.Now().Add(d))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}
