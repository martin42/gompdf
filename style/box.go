package style

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

type Border struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

type Padding struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

type Margin struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

type Box struct {
	Border  Border  `style:"border"`
	Padding Padding `style:"padding"`
	Margin  Margin  `style:"margin"`
}

func (b *Border) UnmarshalStyle(v string) error {
	_, err := fmt.Fscanf(bytes.NewBufferString(v), "%d,%d,%d,%d", &b.Left, &b.Top, &b.Right, &b.Bottom)
	if err != nil {
		return errors.Wrapf(err, "scan border value (%s)", v)
	}
	return nil
}

func (b *Padding) UnmarshalStyle(v string) error {
	_, err := fmt.Fscanf(bytes.NewBufferString(v), "%f,%f,%f,%f", &b.Left, &b.Top, &b.Right, &b.Bottom)
	if err != nil {
		return errors.Wrapf(err, "scan padding value (%s)", v)
	}
	return nil
}

func (b *Margin) UnmarshalStyle(v string) error {
	_, err := fmt.Fscanf(bytes.NewBufferString(v), "%f,%f,%f,%f", &b.Left, &b.Top, &b.Right, &b.Bottom)
	if err != nil {
		return errors.Wrapf(err, "scan margin value (%s)", v)
	}
	return nil
}
