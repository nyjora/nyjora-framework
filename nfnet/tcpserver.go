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
	"nyjora-framework/nflog"
	"sync"
	"time"
)

const (
	RESTART_TCP_SERVER_INTERVAL = 5 * time.Second
	MAX_RETRY_TIME_INTERVAL     = 5 * time.Second
)

type ServerDelegate interface {
	OnAddSession(*NetSession)
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
	opts       ServerOption // server options
	sessionMap *sync.Map    // session map
	mutex      sync.RWMutex // mutex for map
	delegate   ServerDelegate
	wg         *sync.WaitGroup
	listener   net.Listener
}

func NewTcpServer(opt ServerOption, d ServerDelegate) *TcpServer {
	// 组装数据
	s := &TcpServer{
		opts:       opt,
		sessionMap: &sync.Map{},
		delegate:   d,
		wg:         &sync.WaitGroup{},
	}
	return s
}

func (ts *TcpServer) handleConn(conn net.Conn) {
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetWriteBuffer(STREAM_WRITE_BUFFER_SIZE)
	tcpConn.SetReadBuffer(STREAM_READ_BUFFER_SIZE)

	ts.registerConnection(tcpConn)
}

func (ts *TcpServer) registerConnection(conn net.Conn) {
	session := NewNetSession(conn)
	ts.addSession(session)
	defer ts.delSession(session)
	// session begin to work
	session.Run(ts.wg)
}

func (ts *TcpServer) addSession(session *NetSession) {
	ts.sessionMap.Store(session.Id, session)
	ts.delegate.OnAddSession(session)
}

func (ts *TcpServer) delSession(session *NetSession) {
	_, ok := ts.sessionMap.Load(session.Id)
	if ok {
		ts.sessionMap.Delete(session.Id)
		ts.delegate.OnDelSession(session.Id)
	}

}

func (ts *TcpServer) Run(wg *sync.WaitGroup) {
	listenAddr := fmt.Sprintf("%s:%d", ts.opts.Ip, ts.opts.Port)
	l, err := net.Listen("tcp", listenAddr)
	nflog.Info("[TcpServer] Run: listening on TCP: %s ...\n", listenAddr)
	if err != nil {
		nflog.Err("[TcpServer] Run: listen to %s err = %v\n", listenAddr, err)
		return
	}
	ts.listener = l
	wg.Add(1)
	defer func() {
		ts.Close(wg)
	}()
	var delay time.Duration
	for {
		conn, err := ts.listener.Accept()
		// err occurs
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				// retry
				if delay == 0 {
					delay = 5 * time.Millisecond
				} else {
					delay *= 2
				}
				if delay > MAX_RETRY_TIME_INTERVAL {
					delay = MAX_RETRY_TIME_INTERVAL
				}
				nflog.Err("[TcpServer] Server: accept err %v, retry in %d\n", err, delay)
				select {
				case <-time.After(delay):
					nflog.Debug("[TcpServer] Server: time.After(%d)\n", delay)
				}
				continue
			}
			nflog.Info("[TcpServer] Server: accept closed %v\n", err)
			return
		}
		delay = 0
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
	ts.sessionMap.Range(func(k, v interface{}) bool {
		v.(*NetSession).Send(proto)
		return true
	})
}

func (ts *TcpServer) SendProto(sid nfcommon.SessionId, id uint32, fromType uint32, fromId uint32, toType uint32, toId uint32, data []byte) {
	proto := nfcommon.NewProto()
	proto.Id = id
	proto.FromType = fromType
	proto.FromId = fromId
	proto.ToType = toType
	proto.ToId = toId
	proto.Data = data
	session, ok := ts.sessionMap.Load(sid)
	if ok {
		session.(*NetSession).Send(proto)
	}
}

func (ts *TcpServer) Stop(wg *sync.WaitGroup) {
	// break accept loop
	nflog.Debug("[TcpServer] Stop")
	ts.listener.Close()
}

func (ts *TcpServer) Close(wg *sync.WaitGroup) {
	nflog.Debug("[TcpServer] Close.")
	// close all session
	tsm := map[nfcommon.SessionId]*NetSession{}
	ts.sessionMap.Range(func(k, v interface{}) bool {
		tsm[k.(nfcommon.SessionId)] = v.(*NetSession)
		return true
	})
	ts.sessionMap = &sync.Map{}
	for _, c := range tsm {
		c.Close()
	}
	ts.wg.Wait()
	wg.Done()
}

func (ts *TcpServer) StopSession(sid nfcommon.SessionId) {
	nflog.Debug("[TcpServer] StopSession sid = %d\n", sid)
	s, ok := ts.sessionMap.Load(sid)
	if ok {
		s.(*NetSession).Close()
		ts.delSession(s.(*NetSession))
	}
}
