package main

import (
	"github.com/guelfey/go.dbus"
	"log"
)

// DBUS_NAME = "org.freedesktop.Avahi"
// DBUS_INTERFACE_SERVER = DBUS_NAME + ".Server"
// DBUS_PATH_SERVER = "/"
// DBUS_INTERFACE_ENTRY_GROUP = DBUS_NAME + ".EntryGroup"
// DBUS_INTERFACE_DOMAIN_BROWSER = DBUS_NAME + ".DomainBrowser"
// DBUS_INTERFACE_SERVICE_TYPE_BROWSER = DBUS_NAME + ".ServiceTypeBrowser"
// DBUS_INTERFACE_SERVICE_BROWSER = DBUS_NAME + ".ServiceBrowser"
// DBUS_INTERFACE_ADDRESS_RESOLVER = DBUS_NAME + ".AddressResolver"
// DBUS_INTERFACE_HOST_NAME_RESOLVER = DBUS_NAME + ".HostNameResolver"
// DBUS_INTERFACE_SERVICE_RESOLVER = DBUS_NAME + ".ServiceResolver"
// DBUS_INTERFACE_RECORD_BROWSER = DBUS_NAME + ".RecordBrowser"

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

	obj = dconn.Object("org.freedesktop.Avahi", "/")
	obj.Call("org.freedesktop.Avahi.Server.EntryGroupNew", 0).Store(&path)
	obj = dconn.Object("org.freedesktop.Avahi", path)

	var AAY [][]byte
	for _, s := range []string{"email=lemenkov@gmail.com", "jid=lemenkov@gmail.com", "status=avail"} {
		AAY = append(AAY, []byte(s))
	}

	// http://www.dns-sd.org/ServiceTypes.html
	obj.Call("org.freedesktop.Avahi.EntryGroup.AddService", 0,
		int32(-1), // avahi.IF_UNSPEC
		int32(-1), // PROTO_UNSPEC, PROTO_INET, PROTO_INET6  = -1, 0, 1
		uint32(0), // flags
		name, // sname
		"_epmd._tcp", // stype
		"local", // sdomain
		host, // shost
		port, // port
		AAY) // text record

	obj.Call("org.freedesktop.Avahi.EntryGroup.Commit", 0)

	return obj
}

func avahiUnregister(obj *dbus.Object) {
	obj.Call("org.freedesktop.Avahi.EntryGroup.Free", 0)
}
