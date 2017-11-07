/*
tcpserver.go
by umiringo42

此文件定义了服务器模块

*/
package nfnet

import (
	"fmt"
	"net"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfconst"
	"sync"
)

// check err whether is a net error TODO:
func IsNetError(_err interface{}) bool {
	err, ok := _err.(error)
	if !ok {
		return false
	}

	/*
		err = errors.Cause(err)
		if err == io.EOF {
			return true
		}
	*/
	neterr, ok := err.(net.Error)
	if !ok {
		return false
	}
	if neterr.Timeout() {
		return false
	}

	return true
}

type ServerDelegate interface {
	OnAddSession(nfcommon.SessionId)
	OnDelSession(nfcommon.SessionId)
}

type Server interface {
	handleConn(net.Conn)
	addSession(*NetSession)
	delSession(*NetSession)
}

// Server Interface, all server will implement this interface
type ServerOption struct {
	Ip   string
	Port int
}

type TcpServer struct {
	opts       ServerOption                       // server options
	sessionMap map[nfcommon.SessionId]*NetSession // session map
	mutex      sync.RWMutex                       // mutex for map
	delegate   ServerDelegate
}

func NewTcpServer(opt ServerOption, d ServerDelegate) *TcpServer {
	// 组装数据
	return &TcpServer{
		opts:       opt,
		sessionMap: make(map[nfcommon.SessionId]*NetSession),
		delegate:   d,
	}
}

func (ts *TcpServer) handleConn(conn net.Conn) {
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetWriteBuffer(nfconst.STREAM_WRITE_BUFFER_SIZE)
	tcpConn.SetReadBuffer(nfconst.STREAM_READ_BUFFER_SIZE)

	ts.registerConnection(tcpConn)
}

func (ts *TcpServer) registerConnection(conn net.Conn) {
	session := NewNetSession(conn)
	ts.addSession(session)
	defer ts.delSession(session)
	// session begin to work
	session.Run()
}

func (ts *TcpServer) addSession(session *NetSession) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.sessionMap[session.Id] = session
	ts.delegate.OnAddSession(session.Id)
}

func (ts *TcpServer) delSession(session *NetSession) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	if ts.sessionMap[session.Id] != nil {
		delete(ts.sessionMap, session.Id)
	}
	ts.delegate.OnDelSession(session.Id)
}

// open a port wait for connecting
func (ts *TcpServer) Run() error {
	listenAddr := fmt.Sprintf("%s:%d", ts.opts.Ip, ts.opts.Port)
	listener, err := net.Listen("tcp", listenAddr)
	fmt.Printf("Listening on TCP: %s ...\n", listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue
			} else {
				return err
			}
		}
		fmt.Printf("Connection from: %s\n", conn.RemoteAddr())
		// trigger a delegate to add new conn
		go ts.handleConn(conn)
	}
}

func (ts *TcpServer) BroadCast(id uint32, fromType uint32, fromId uint32, toType uint32, toId uint32, data []byte) {
	proto := nfcommon.NewProto()
	proto.Id = id
	proto.FromType = fromType
	proto.FromId = fromId
	proto.ToType = toType
	proto.ToId = toId
	proto.Data = data
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	for _, v := range ts.sessionMap {
		if v != nil {
			v.Send(proto)
		}
	}
}

func (ts *TcpServer) SendProto(sid nfcommon.SessionId, id uint32, fromType uint32, fromId uint32, toType uint32, toId uint32, data []byte) {
	proto := nfcommon.NewProto()
	proto.Id = id
	proto.FromType = fromType
	proto.FromId = fromId
	proto.ToType = toType
	proto.ToId = toId
	proto.Data = data
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	session := ts.sessionMap[sid]
	if session != nil {
		session.Send(proto)
	}
}
