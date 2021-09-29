package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/wkt"
	"io"
	"log"
	"os"
)

func convert(r *csv.Reader, w *bufio.Writer) error {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		g, err := ewkbhex.Decode(record[0])
		if err != nil {
			return err
		}
		switch t := g.(type) {
		case *geom.LineString:
			g2 := geom.LineString(*t)
			coords := [][][]geom.Coord{[][]geom.Coord{g2.Coords()}}
			g = geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)
		case *geom.Polygon:
			g2 := geom.Polygon(*t)
			coords := [][][]geom.Coord{g2.Coords()}
			g = geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)
		case *geom.MultiPolygon: // ignore
		default:
			return fmt.Errorf("cannot handle geometry %v", t)
		}
		h, err := wkt.Marshal(g)
		if err != nil {
			return err
		}
		_, err = w.WriteString(h + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func main() {
	r := csv.NewReader(os.Stdin)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.ReuseRecord = true

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	checkErr(convert(r, w))
}
