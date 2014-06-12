// +build zmq_3_x

package gozmq

/*
#cgo pkg-config: libzmq
#include <zmq.h>
#include <stdlib.h>
#include <string.h>


char *ev_get_addr(zmq_event_t *ev) {
  switch (ev->event) {
  case ZMQ_EVENT_CONNECTED:
    return ev->data.connected.addr;
  case ZMQ_EVENT_CONNECT_DELAYED:
    return ev->data.connect_delayed.addr;
  case ZMQ_EVENT_CONNECT_RETRIED:
    return ev->data.connect_retried.addr;
  case ZMQ_EVENT_LISTENING:
    return ev->data.listening.addr;
  case ZMQ_EVENT_BIND_FAILED:
    return ev->data.bind_failed.addr;
  case ZMQ_EVENT_ACCEPTED:
    return ev->data.accepted.addr;
  case ZMQ_EVENT_ACCEPT_FAILED:
    return ev->data.accept_failed.addr;
  case ZMQ_EVENT_CLOSED:
    return ev->data.closed.addr;
  case ZMQ_EVENT_CLOSE_FAILED:
    return ev->data.close_failed.addr;
  case ZMQ_EVENT_DISCONNECTED:
    return ev->data.disconnected.addr;
  case ZMQ_EVENT_PEER_ATTACHED:
    return ev->data.peer_attached.addr;
  case ZMQ_EVENT_PEER_DETACHED:
    return ev->data.peer_detached.addr;
  }
}

size_t ev_get_extra(zmq_event_t *ev) {
  switch (ev->event) {
  case ZMQ_EVENT_CONNECTED:
    return ev->data.connected.fd;
  case ZMQ_EVENT_CONNECT_DELAYED:
    return ev->data.connect_delayed.err;
  case ZMQ_EVENT_CONNECT_RETRIED:
    return ev->data.connect_retried.interval;
  case ZMQ_EVENT_LISTENING:
    return ev->data.listening.fd;
  case ZMQ_EVENT_BIND_FAILED:
    return ev->data.bind_failed.err;
  case ZMQ_EVENT_ACCEPTED:
    return ev->data.accepted.fd;
  case ZMQ_EVENT_ACCEPT_FAILED:
    return ev->data.accept_failed.err;
  case ZMQ_EVENT_CLOSED:
    return ev->data.closed.fd;
  case ZMQ_EVENT_CLOSE_FAILED:
    return ev->data.close_failed.err;
  case ZMQ_EVENT_DISCONNECTED:
    return ev->data.disconnected.fd;
  case ZMQ_EVENT_PEER_ATTACHED:
    return ev->data.peer_attached.id;
  case ZMQ_EVENT_PEER_DETACHED:
    return ev->data.peer_detached.id;
  }
}


*/
import "C"
import (
	"unsafe"
	"errors"
)


const (

  APUB                  = SocketType(C.ZMQ_APUB)

  BLOCK_ADDR            = StringSocketOption(C.ZMQ_BLOCK_ADDR)
  UNBLOCK_ADDR          = StringSocketOption(C.ZMQ_UNBLOCK_ADDR)
  APUB_APPROVE          = StringSocketOption(C.ZMQ_APUB_APPROVE)
  LAST_PEER_ADDR        = StringSocketOption(C.ZMQ_LAST_PEER_ADDR)
  APUB_REQ              = IntSocketOption(C.ZMQ_APUB_REQ)
  LAST_PEER_UNIQ_ID     = IntSocketOption(C.ZMQ_LAST_PEER_UNIQ_ID)
  DISCONNECT_PEER_BY_ID = IntSocketOption(C.ZMQ_DISCONNECT_PEER_BY_ID)

  EVENT_PEER_ATTACHED = Event(C.ZMQ_EVENT_PEER_ATTACHED)
  EVENT_PEER_DETACHED = Event(C.ZMQ_EVENT_PEER_DETACHED)
)

func (s *Socket) GetSockOptUInt(option IntSocketOption) (value uint32, err error) {
	size := C.size_t(unsafe.Sizeof(value))
	var rc C.int
	if rc, err = C.zmq_getsockopt(s.s, C.int(option), unsafe.Pointer(&value), &size); rc != 0 {
		err = casterr(err)
		return
	}
	return
}

func (s *Socket) APubReq() (bool, error) {
	value, err := s.GetSockOptInt(APUB_REQ)
	return value != 0, err
}

func (s *Socket) DisconnectPeerById(id uint32) error {
	if rc, err := C.zmq_setsockopt(s.s, C.int(DISCONNECT_PEER_BY_ID), unsafe.Pointer(nil), C.size_t(id)); rc != 0 {
		return casterr(err)
	}
	return nil
}

func (s *Socket) GetLastPeerUniqueID() (uint32, error) {
	return s.GetSockOptUInt(LAST_PEER_UNIQ_ID)
}

func (s *Socket) APubApprove(value string, peerid uint32) error {
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	if rc, err := C.zmq_setsockopt(s.s, C.int(APUB_APPROVE), unsafe.Pointer(v), C.size_t(peerid)); rc != 0 {
		return casterr(err)
	}
	return nil
}

func (s *Socket) BlockAddr(value string) error {
	return s.SetSockOptString(BLOCK_ADDR, value)
}

func (s *Socket) UnblockAddr(value string) error {
	return s.SetSockOptString(UNBLOCK_ADDR, value)
}

func (s *Socket) DisconnectLastRecvPeer() error {
	if rc, err := C.zmq_disconnect_last_recv_peer(s.s); rc != 0 {
		return casterr(err)
	}
	return nil
}

func (s *Socket) BlockLastRecvPeer() error {
	if rc, err := C.zmq_block_last_recv_peer(s.s); rc != 0 {
		return casterr(err)
	}
	return nil
}

func (s *Socket) GetLastRecvPeerAddr() (string, error) {
    var addr *C.char

    if rc, err := C.zmq_last_recv_peer_addr(s.s, &addr); rc != 0 {
        return "", casterr(err)
    }
    defer C.free(unsafe.Pointer(addr))
    return C.GoString(addr), nil
}

func (s *Socket) RecvEvent(flags SendRecvOption) (ty Event, addr string, ex int64, err error) {
	// Allocate and initialise a new zmq_msg_t
	var m C.zmq_msg_t
	var rc C.int
	if rc, err = C.zmq_msg_init(&m); rc != 0 {
		err = casterr(err)
		return
	}
	defer C.zmq_msg_close(&m)
	// Receive into message
	if rc, err = C.zmq_recvmsg(s.s, &m, C.int(flags)); rc == -1 {
		err = casterr(err)
		return
	}
  var ev *C.zmq_event_t
	size := C.zmq_msg_size(&m)
  if uintptr(size) != unsafe.Sizeof(*ev) {
    err = errors.New("Invalid event message received")
    return
  }
	err = nil

  ev = (*C.zmq_event_t)(C.zmq_msg_data(&m))
  ty = Event(ev.event)
  addr = C.GoString(C.ev_get_addr(ev))
  ex = int64(C.ev_get_extra(ev))
  return
}
