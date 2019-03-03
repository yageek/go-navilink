package main

// #cgo CPPFLAGS: -I${SRCDIR}/clib/navilink/include -I${SRCDIR}/clib/libserialport/include
// #cgo darwin LDFLAGS: ${SRCDIR}/clib/navilink/lib/libnavilink.a ${SRCDIR}/clib/libserialport/lib/libserialport.a -framework Foundation -framework IOKit
// #cgo LDFLAGS: ${SRCDIR}/clib/navilink/lib/libnavilink.a clib/libserialport/lib/libserialport.a
// #include <navilink/navilink.h>
import "C"
import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

// DeviceInfo représente les informations du périphérique.
type DeviceInfo struct {
	TotalWaypoint    uint
	TotalRoute       uint
	TotalTrack       uint
	ProtocolVersion  uint
	SerialNumber     uint
	TrackpointsCount uint
	Name             string
}

// Position représente la position d'un Waypoint
// dans la projecton WGS84
type Position struct {
	Lat      float64
	Lng      float64
	Altitude float64
}

// Waypoint est un point du chemin parcouru
type Waypoint struct {
	ID       uint
	Name     string
	Date     time.Time
	Position Position
}

func (d *Device) navilinkErrorToGoError() error {
	errorC := C.navilink_get_error_description(&d.navilinkC)
	return errors.New(C.GoString(errorC))
}

// Device représente une interface
// vers le port série connecté au GPS
type Device struct {
	navilinkC C.NavilinkDevice
	Infos     DeviceInfo
}

// Open permettra au développeur d'ouvrir la communication
func Open(path string) (*Device, error) {
	device := &Device{}
	res := C.navilink_open_device_from_name(C.CString(path), &device.navilinkC)
	if res != C.NavilinkOK {
		return nil, device.navilinkErrorToGoError()
	}
	device.Infos = deviceInfosFromCInfos(device.navilinkC.informations)
	return device, nil
}

// MaxWaypointQueryLength est le maximum de points retournés dans une requête
const MaxWaypointQueryLength = 32

// GetAllWaypoints récupère tous les waypoints présents
// dans l'appareil.
func (d *Device) GetAllWaypoints() ([]Waypoint, error) {

	// Tableau C
	tmpPoints := make([]C.NavilinkWaypoint, d.Infos.TotalWaypoint)
	// Slice Go
	points := make([]Waypoint, d.Infos.TotalWaypoint)

	downloadPassCount := int(d.Infos.TotalWaypoint) / MaxWaypointQueryLength
	rest := int(d.Infos.TotalWaypoint) % MaxWaypointQueryLength
	downloadPassCount += rest

	// On télécharge tous les points en C
	for i := 0; i < downloadPassCount; i++ {

		cursorIndex := 0
		res := C.navilink_query_waypoint(&d.navilinkC, C.int(cursorIndex), MaxWaypointQueryLength, &tmpPoints[cursorIndex])
		if res < 0 {
			return []Waypoint{}, fmt.Errorf("Error retrieving %d waypoints from %d. Aborting", MaxWaypointQueryLength, cursorIndex)
		}

		for c := cursorIndex; c < cursorIndex+int(res); c++ {
			points[c].ID = uint(tmpPoints[c].waypointID)
			points[c].Name = C.GoString((*C.char)(unsafe.Pointer(&tmpPoints[c].waypointName[0])))
			dateTime := tmpPoints[c].datetime
			points[c].Date = time.Date(int(dateTime.year)+2000, time.Month(dateTime.month), int(dateTime.day), int(dateTime.hour), int(dateTime.minute), int(dateTime.second), 0, time.UTC)

			lat := float64(tmpPoints[c].position.latitude) / 10000000
			lng := float64(tmpPoints[c].position.longitude) / 10000000
			alt := float64(tmpPoints[c].position.altitude) * 0.3048
			points[c].Position = Position{Lat: lat, Lng: lng, Altitude: alt}
		}
	}

	return points, nil
}

// Close permettra au développeur de fermer la communication
func (d *Device) Close() {
	C.navilink_close_device(&d.navilinkC)
}

func deviceInfosFromCInfos(cInfo C.NavilinkInformation) DeviceInfo {

	return DeviceInfo{
		TotalWaypoint:    uint(cInfo.totalWaypoint),
		TotalRoute:       uint(cInfo.totalRoute),
		TotalTrack:       uint(cInfo.totalTrack),
		ProtocolVersion:  uint(cInfo.protocolVersion),
		SerialNumber:     uint(cInfo.deviceSerialNum),
		TrackpointsCount: uint(cInfo.numOfTrackpoints),
		Name:             C.GoString((*C.char)(unsafe.Pointer(&cInfo.username[0]))),
	}
}
