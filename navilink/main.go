package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	gj "github.com/kpawlik/geojson"
)

var (
	devInput string
	info     bool
	output   string
)

func main() {

	flag.StringVar(&devInput, "input", "", "Périphérique d'entrée")
	flag.StringVar(&output, "output", "", "Sortie du fichier")
	flag.BoolVar(&info, "info", false, "Display devices informations")
	flag.Parse()

	if devInput == "" {
		fmt.Printf("Aucun prériphérique d'entrée spécifié :( \n")
		flag.Usage()
		os.Exit(-1)
	}

	if !info && output == "" {
		fmt.Printf("Au moins une action est requise \n")
		flag.Usage()
		os.Exit(-1)
	}

	device, err := Open(devInput)
	if err != nil {
		fmt.Printf("Impossible de se connecter: %s \n", err)
		os.Exit(-1)
	}
	defer device.Close()

	if info {
		printInfo(device)
	}

	err = toGeoJSON(device, output)
	if err != nil {
		panic(err)
	}
}

func printInfo(device *Device) {
	fmt.Printf("Informations sur le périphérique: \n")
	fmt.Printf("\tDevice name: %s \n", device.Infos.Name)
	fmt.Printf("\tDevice serial: %d \n", device.Infos.SerialNumber)
	fmt.Printf("\tDevice protocol version: %d \n", device.Infos.ProtocolVersion)
	fmt.Printf("\tTotal waypoints: %d \n", device.Infos.TotalWaypoint)
	fmt.Printf("\tTotal routes: %d \n", device.Infos.TotalRoute)
	fmt.Printf("\tTotal track: %d \n\n", device.Infos.TotalTrack)
}

func toGeoJSON(d *Device, file string) error {
	wpts, err := d.GetAllWaypoints()
	if err != nil {
		return err
	}

	fc := gj.NewFeatureCollection([]*gj.Feature{})

	for _, pt := range wpts {
		coord := gj.Coordinate{gj.CoordType(pt.Position.Lng), gj.CoordType(pt.Position.Lat), gj.CoordType(pt.Position.Altitude)}
		p := gj.NewPoint(coord)
		props := map[string]interface{}{"Name": pt.Name, "ID": pt.ID, "Date": pt.Date}
		f1 := gj.NewFeature(p, props, nil)
		fc.AddFeatures(f1)
	}
	out, err := gj.Marshal(fc)
	if err != nil {
		return err
	}

	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.WriteString(fd, out)
	return err
}
