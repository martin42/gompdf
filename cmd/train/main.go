package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/mazzegi/gompdf"
)

func main() {
	train := BuildTrain()

	src, err := ioutil.ReadFile("train.xml")
	if err != nil {
		panic(err)
	}

	t, err := template.New("train").Funcs(funcs()).Parse(string(src))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "train", train)
	if err != nil {
		panic(err)
	}

	//no we've got the xml inside the buffer
	doc, err := gompdf.Load(&buf)
	if err != nil {
		panic(err)
	}
	p, err := gompdf.NewProcessor(doc)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("train.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = p.Process(f)
	if err != nil {
		panic(err)
	}
}

func funcs() template.FuncMap {
	fm := template.FuncMap{
		"number": func(i int) int {
			return i + 1
		},
		"scaleTemp": func(tmp Temp, maxwidth int) int {
			d := float64(tmp.Value) / float64(tmp.Max-tmp.Min)
			return int(d * float64(maxwidth))
		},
		"tempAlarmColor": func(tmp Temp) string {
			switch tmp.Alarm {
			case 1:
				return "#eeee00"
			case 2:
				return "#dd0000"
			default:
				return "#dddddd"
			}
		},
	}
	return fm
}

//
type Temp struct {
	Value int
	Alarm int
	Min   int
	Max   int
}

type Axle struct {
	Number       int
	BearingLeft  Temp
	BearingRight Temp
	Disk         Temp
	Brake        Temp
}

type Vehicle struct {
	Number   int
	ID       string
	Type     string
	Operator string
	Axles    []Axle
}

type Train struct {
	System    string
	Started   string
	Finished  string
	Number    string
	Operator  string
	Track     string
	Direction string
	Ambient   int
	Speed     int
	Length    int
	Axles     int
	Vehicles  []Vehicle
}

func BuildTrain() *Train {
	t := Train{
		System:    "Wuhan<",
		Started:   "2020-03-22 14:11:12",
		Finished:  "2020-03-22 14:17:43",
		Number:    ">FXE123",
		Operator:  "FX",
		Track:     "A",
		Direction: "/West",
		Ambient:   27,
		Speed:     83,
		Length:    42,
		Axles:     16,
	}
	for i := 0; i < 4; i++ {
		v := Vehicle{
			Number:   i + 1,
			ID:       fmt.Sprintf("XCB-22-%d", i+42),
			Type:     "Locomotive",
			Operator: "FX",
		}
		for a := 0; a < 4; a++ {
			v.Axles = append(v.Axles, Axle{
				Number:       i*4 + a + 1,
				BearingLeft:  bearingTemp(0),
				BearingRight: bearingTemp(0),
				Disk:         wheelTemp(0),
				Brake:        wheelTemp(0),
			})
		}
		t.Vehicles = append(t.Vehicles, v)
	}
	t.Vehicles[0].Axles[1].BearingLeft = bearingTemp(1)
	t.Vehicles[0].Axles[3].BearingRight = bearingTemp(2)
	t.Vehicles[1].Axles[0].Disk = wheelTemp(2)
	t.Vehicles[1].Axles[2].Brake = wheelTemp(1)
	return &t
}

func bearingTemp(alarm int) Temp {
	switch alarm {
	case 1:
		return Temp{
			Value: 60 + rand.Intn(30),
			Alarm: alarm,
			Min:   0,
			Max:   150,
		}
	case 2:
		return Temp{
			Value: 90 + rand.Intn(40),
			Alarm: alarm,
			Min:   0,
			Max:   150,
		}
	default:
		return Temp{
			Value: 20 + rand.Intn(40),
			Alarm: alarm,
			Min:   0,
			Max:   150,
		}
	}
}

func wheelTemp(alarm int) Temp {
	switch alarm {
	case 1:
		return Temp{
			Value: 200 + rand.Intn(200),
			Alarm: alarm,
			Min:   80,
			Max:   650,
		}
	case 2:
		return Temp{
			Value: 400 + rand.Intn(150),
			Alarm: alarm,
			Min:   80,
			Max:   650,
		}
	default:
		return Temp{
			Value: 70 + rand.Intn(130),
			Alarm: alarm,
			Min:   80,
			Max:   650,
		}
	}
}
