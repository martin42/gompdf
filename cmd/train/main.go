package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/mazzegi/gompdf"
)

func main() {
	train := BuildTrain()

	start := time.Now()
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
		//fmt.Printf("\n%s\n", buf.String())
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
	fmt.Printf("build pdf in (%s)\n", time.Since(start))
}

func funcs() template.FuncMap {
	fm := template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"log": func(s string) error {
			fmt.Printf("template: %s\n", s)
			return nil
		},
		"number": func(i int) int {
			return i + 1
		},
		"scaleTemp": func(tmp Temp, maxwidth int) int {
			d := float64(tmp.Value) / float64(tmp.Max-tmp.Min)
			return int(d * float64(maxwidth))
		},
		"tempAlarmColor": func(alm int) string {
			switch alm {
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

type Stamp struct {
	Message string
	Date    string
	By      string
}

type Alarm struct {
	Severity       int
	Entity         string
	VehicleID      string
	Description    string
	Value          string
	CaseState      Stamp
	CaseAssessment Stamp
	CaseTreatment  Stamp
	CaseRemark     Stamp
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
	Alarms    []Alarm
}

func BuildTrain() *Train {
	t := Train{
		System:    "Cheddar?Creek<",
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
			Operator: "ÜÄÖ",
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

	t.Alarms = append(t.Alarms, Alarm{
		Severity:    2,
		Entity:      "Wheelset 4",
		VehicleID:   "XCB-22-42",
		Description: "Absolute Hot - Right Bearing",
		Value:       "118 °C",
		CaseState: Stamp{
			Message: "Processed",
			Date:    "2020-03-22 15:10:12",
			By:      "Senior Ramos",
		},
		CaseAssessment: Stamp{
			Message: "Defect",
			Date:    "2020-03-22 15:08:54",
			By:      "Senior Ramos",
		},
		CaseTreatment: Stamp{
			Message: "Vehicle Removed",
			Date:    "2020-03-22 15:09:17",
			By:      "Senior Ramos",
		},
		CaseRemark: Stamp{
			Message: "Knocking off now",
			Date:    "2020-03-22 16:00:00",
			By:      "Senior Ramos",
		},
	})

	t.Alarms = append(t.Alarms, Alarm{
		Severity:    1,
		Entity:      "Wheelset 5",
		VehicleID:   "XCB-22-43",
		Description: "Absolute Hot - Disk Bearing",
		Value:       "420°C °C",
		CaseState: Stamp{
			Message: "Acknowledged",
			Date:    "2020-03-22 15:33:02",
			By:      "Senior Ramos",
		},
		CaseAssessment: Stamp{
			Message: "Too hot",
			Date:    "2020-03-22 15:33:20",
			By:      "Senior Ramos",
		},
		CaseRemark: Stamp{
			Message: "Still investigating",
			Date:    "2020-03-22 15:45:35",
			By:      "Senior Ramos",
		},
	})

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
