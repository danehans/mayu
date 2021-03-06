package rpcout

import (
	"errors"
	"fmt"

	"github.com/coreos/ignition/third_party/github.com/sigma/vmw-guestinfo/message"
)

// ErrRpciFormat represents an invalid result format
var ErrRpciFormat = errors.New("invalid format for RPCI command result")

const rpciProtocolNum uint32 = 0x49435052

// SendOne is a command-oriented wrapper for SendOneRaw
func SendOne(format string, a ...interface{}) (reply []byte, ok bool, err error) {
	request := fmt.Sprintf(format, a...)
	return SendOneRaw([]byte(request))
}

// SendOneRaw uses a throw-away RPCOut to send a request
func SendOneRaw(request []byte) (reply []byte, ok bool, err error) {
	out := &RPCOut{}
	if err = out.Start(); err != nil {
		return
	}
	if reply, ok, err = out.Send(request); err != nil {
		return
	}
	if err = out.Stop(); err != nil {
		return
	}
	return
}

// RPCOut is an ougoing connection from the VM to the hypervisor
type RPCOut struct {
	channel *message.Channel
}

// Start opens the connection
func (out *RPCOut) Start() error {
	channel, err := message.NewChannel(rpciProtocolNum)
	if err != nil {
		return err
	}
	out.channel = channel
	return nil
}

// Stop closes the connection
func (out *RPCOut) Stop() error {
	err := out.channel.Close()
	out.channel = nil
	return err
}

// Send emits a request and receives a response
func (out *RPCOut) Send(request []byte) (reply []byte, ok bool, err error) {
	if err = out.channel.Send(request); err != nil {
		return
	}

	var resp []byte
	if resp, err = out.channel.Receive(); err != nil {
		return
	}

	switch string(resp[:2]) {
	case "0 ":
		reply = resp[2:]
	case "1 ":
		reply = resp[2:]
		ok = true
	default:
		err = ErrRpciFormat
	}
	return
}
