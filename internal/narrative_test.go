package narrative

import (
	"testing"
)

func loadMarkers() *CommentMarkers {
	markers, err := ParseMarkerConfig("../testdata/narrative.yaml")
	if err != nil {
		panic(`No marker config`)
	}
	return markers
}

func TestFindMarkerSh(t *testing.T) {
	markers := loadMarkers()
	marker, err := FindMarker(markers, ".sh")
	if err != nil {
		t.Fatalf(`Could not find marker`)
	}
	if marker.Start != "<< comment" {
		t.Fatalf(`FindMarker(".sh") = %v, %v`, marker, err)
	}
}

func TestFindMarkerTf(t *testing.T) {
	markers := loadMarkers()
	marker, err := FindMarker(markers, ".tf")
	if err != nil {
		t.Fatalf(`Could not find marker`)
	}
	if marker.Start != "/*" {
		t.Fatalf(`FindMarker(".tf") = %v, %v`, marker, err)
	}
}

func TestFindMarkerGo(t *testing.T) {
	markers := loadMarkers()
	marker, err := FindMarker(markers, ".go")
	if err != nil {
		t.Fatalf(`Could not find marker`)
	}
	if marker.Start != "/*" {
		t.Fatalf(`FindMarker(".go") = %v, %v`, marker, err)
	}
}

func TestFindBrokenExt(t *testing.T) {
	markers := loadMarkers()
	_, err := FindMarker(markers, "0")
	if err == nil {
		t.Fatalf(`Should throw error on short extension`)
	}
}

func TestFindBrokenExt2(t *testing.T) {
	markers := loadMarkers()
	_, err := FindMarker(markers, "")
	if err == nil {
		t.Fatalf(`Should throw error on blank extension`)
	}
}

func TestFindBrokenExt3(t *testing.T) {
	markers := loadMarkers()
	_, err := FindMarker(markers, ".你好")
	if err == nil {
		t.Fatalf(`Should throw error on not found utf8 extension`)
	}
}
