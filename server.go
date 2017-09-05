package foog

import (
	"net"
)

type IServer interface{
	Run(net.Listener, func(IConn))
}