// +build zmq_3_x

package gozmq

/*
#cgo pkg-config: libzmq
#include <zmq.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"unsafe"
)


const (

	APUB            = SocketType(C.ZMQ_APUB)

    BLOCK_ADDR      = StringSocketOption(C.ZMQ_BLOCK_ADDR)
    UNBLOCK_ADDR    = StringSocketOption(C.ZMQ_UNBLOCK_ADDR)
    APUB_APPROVE    = StringSocketOption(C.ZMQ_APUB_APPROVE)
    APUB_REQ        = IntSocketOption(C.ZMQ_APUB_REQ)

)

func (s *Socket) APubReq() (bool, error) {
	value, err := s.GetSockOptInt(APUB_REQ)
	return value != 0, err
}

func (s *Socket) APubApprove(value string) error {
	return s.SetSockOptString(APUB_APPROVE, value)
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
    var addr **C.char

    if rc, err := C.zmq_last_recv_peer_addr(s.s, addr); rc != 0 {
        return "", casterr(err)
    }
    defer C.free(unsafe.Pointer(*addr))
    return C.GoString(*addr), nil
}
