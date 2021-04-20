package dao

import (
	"log"
	"crypto/tls"
	"net"

	"gopkg.in/mgo.v2"
)

const (
	COLLECTION_REVIEWS	=	"reviews"
	COLLECTION_CLASSES	=	"classes"
	COLLECTION_REPORTS	=	"reported"
	COLLECTION_QUESTION	=	"questions"
	COLLECTION_ANSWERS	=	"answers"
	COLLECTION_RECAPS	=	"recaps"
)

// SessionDAO is struct for allocate info for create connection with mongoDB
type SessionDAO struct {
	Server   string
	Database string
}

var session *mgo.Session

// Connect is Establish a connection to database
func (m *SessionDAO) Connect() {

	tlsConfig		:=	&tls.Config{}

	dialInfo, err	:=	mgo.ParseURL(m.Server)
	if err != nil {
		log.Fatal(err)
	}

	dialInfo.DialServer	=	func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err	:=	tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err	=	mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)

	log.Println("MGO: Mongo has connected, Server get origin session. 🎉")
}