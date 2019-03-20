package models

import (
	"context"
	"net"
)

type MemberModel struct {
	Hush       int32
	MemberID   int
	ConnHandle net.Conn
	Send       chan *[]byte
	Ctx        context.Context
	Cancel     context.CancelFunc
}
