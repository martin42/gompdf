package gompdf

import "testing"

func TestStyleClasses(t *testing.T) {
	sc := `
	body{
		border: 1,1,1,1;
		margin: 20, 30, 40, 60;		
	}

	launcher{
		padding: 3,4,5,5;
	}
	`
	scs, err := ParseClasses([]byte(sc))
	if err != nil {
		t.Fatalf("parse classes: %v", err)
	}
	t.Logf("scs: %#v", scs)
}
