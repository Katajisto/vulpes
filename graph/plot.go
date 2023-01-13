package graph

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
	"vulpes.ktj.st/models"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

type xySet struct {
	name string
	xy   plotter.XYs
}

type xyPair struct {
	time  time.Time
	value float64
}

func formSensorXySets(points []models.DataPoint) []xySet {
	sensorMap := make(map[string][]xyPair)

	for _, point := range points {
		for _, tempData := range point.TemperatureData {
			sensorMap[tempData.Sensor] = append(sensorMap[tempData.Sensor], xyPair{time: point.Timestamp, value: tempData.Value})
		}
	}

	xySets := make([]xySet, 0)
	for k, v := range sensorMap {
		curXySet := xySet{name: k}
		curXySet.xy = make(plotter.XYs, len(v))
		for i, pair := range v {
			curXySet.xy[i].X = float64(pair.time.Unix())
			curXySet.xy[i].Y = pair.value
		}
		xySets = append(xySets, curXySet)
	}
	return xySets
}

func histPlot(points []xySet) string {
	p := plot.New()
	p.Title.Text = "Lämpötilat"

	c := vgsvg.NewWith(
		vgsvg.UseWH(30*vg.Centimeter, 15*vg.Centimeter),
	)

	for _, v := range points {
		err := plotutil.AddLinePoints(p, v.name, v.xy)
		if err != nil {
			log.Println("Error in adding plot points: ", err)
		}
	}

	var buf bytes.Buffer
	dc := draw.New(c)
	p.Draw(dc)
	c.WriteTo(&buf)
	return toBase64(buf.Bytes())
}

func GetTemperaturePlotImageDataUrl(dataPoints []models.DataPoint, err error) template.URL {
	if err != nil {
		return ""
	}

	b64png := histPlot(formSensorXySets(dataPoints))
	return template.URL(fmt.Sprintf("data:image/svg+xml;base64,%v", b64png))
}