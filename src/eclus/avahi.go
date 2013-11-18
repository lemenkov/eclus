package main

import (
	"github.com/guelfey/go.dbus"
	"log"
)

const DBUS_NAME = "org.freedesktop.Avahi"
const DBUS_INTERFACE_SERVER = DBUS_NAME + ".Server"
const DBUS_PATH_SERVER = "/"
const DBUS_INTERFACE_ENTRY_GROUP = DBUS_NAME + ".EntryGroup"
const DBUS_INTERFACE_DOMAIN_BROWSER = DBUS_NAME + ".DomainBrowser"
const DBUS_INTERFACE_SERVICE_TYPE_BROWSER = DBUS_NAME + ".ServiceTypeBrowser"
const DBUS_INTERFACE_SERVICE_BROWSER = DBUS_NAME + ".ServiceBrowser"
const DBUS_INTERFACE_ADDRESS_RESOLVER = DBUS_NAME + ".AddressResolver"
const DBUS_INTERFACE_HOST_NAME_RESOLVER = DBUS_NAME + ".HostNameResolver"
const DBUS_INTERFACE_SERVICE_RESOLVER = DBUS_NAME + ".ServiceResolver"
const DBUS_INTERFACE_RECORD_BROWSER = DBUS_NAME + ".RecordBrowser"

const PROTO_UNSPEC = int32(-1)
const PROTO_INET = int32(0)
const PROTO_INET6 = int32(1)

const IF_UNSPEC = int32(-1)

func dbusConnect (isDbusSystemWide bool) *dbus.Conn {
	var err error
	var dconn *dbus.Conn
	var visibility string
	if isDbusSystemWide {
		dconn, err = dbus.SystemBus()
		visibility = "system-wide"
	} else {
		dconn, err = dbus.SessionBus()
		visibility = "local/session"
	}
	if err == nil {
		for _, name := range []string{"org.goerlang.Eclus", "org.erlang.Epmd"} {
			log.Printf("Registering at D-Bus: %s (%s)", name, visibility)
			reply, err := dconn.RequestName(name, dbus.NameFlagDoNotQueue)
			if reply != dbus.RequestNameReplyPrimaryOwner {
				// Shall we continue as is instead?
				log.Fatal("Registering at D-Bus: %s name failed. Aready taken,  %+v,  %+v", name, reply, err)
			}
			if err != nil {
				log.Fatal("Registering at D-Bus: %s name failed. Other error, %+v,  %+v", name, reply, err)
			}
		}
		return dconn
	}
	return nil
}

func avahiRegister(dconn *dbus.Conn, name string, host string, port uint16) *dbus.Object {
	var obj *dbus.Object
	var path dbus.ObjectPath

	obj = dconn.Object(DBUS_NAME, DBUS_PATH_SERVER)
	obj.Call(DBUS_INTERFACE_SERVER + ".EntryGroupNew", 0).Store(&path)
	obj = dconn.Object(DBUS_NAME, path)

	var AAY [][]byte
	for _, s := range []string{"email=lemenkov@gmail.com", "jid=lemenkov@gmail.com", "status=avail"} {
		AAY = append(AAY, []byte(s))
	}

	// http://www.dns-sd.org/ServiceTypes.html
	obj.Call(DBUS_INTERFACE_ENTRY_GROUP + ".AddService", 0,
		IF_UNSPEC,
		PROTO_UNSPEC,
		uint32(0), // flags
		name, // sname
		"_epmd._tcp", // stype
		"local", // sdomain
		host, // shost
		port, // port
		AAY) // text record

	obj.Call(DBUS_INTERFACE_ENTRY_GROUP + ".Commit", 0)

	return obj
}

func avahiUnregister(obj *dbus.Object) {
	obj.Call(DBUS_INTERFACE_ENTRY_GROUP + ".Free", 0)
}
