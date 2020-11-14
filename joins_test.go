package turf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomchavakis/turf-go/geojson"
	"github.com/tomchavakis/turf-go/geojson/feature"
	"github.com/tomchavakis/turf-go/geojson/geometry"
	"github.com/tomchavakis/turf-go/utils"
)

const PolyWithHoleFixture = "test-data/poly-with-hole.json"
const MultiPolyWithHoleFixture = "test-data/multipoly-with-hole.json"

func TestPointInPolygon(t *testing.T) {
	type args struct {
		point   geometry.Point
		polygon geometry.Polygon
	}

	poly := geometry.Polygon{
		Coordinates: []geometry.LineString{
			{
				Coordinates: []geometry.Point{
					{
						Lat: 36.171278341935434,
						Lng: -86.76624298095703,
					},
					{
						Lat: 36.170862616662134,
						Lng: -86.74238204956055,
					},
					{
						Lat: 36.19607929145354,
						Lng: -86.74100875854492,
					},
					{
						Lat: 36.2014818084173,
						Lng: -86.77362442016602,
					},
					{
						Lat: 36.171278341935434,
						Lng: -86.76624298095703,
					},
				},
			},
		},
	}

	tests := map[string]struct {
		args    args
		want    bool
		wantErr bool
	}{
		"point in Polygon": {
			args: args{
				point: geometry.Point{
					Lat: 36.185411688981105,
					Lng: -86.76074981689453,
				},
				polygon: poly,
			},
			want:    true,
			wantErr: false,
		},
		"point in Polygon 2": {
			args: args{
				point: geometry.Point{
					Lat: 36.19393203374786,
					Lng: -86.75946235656737,
				},
				polygon: poly,
			},
			want:    true,
			wantErr: false,
		},
		"point out of Polygon": {
			args: args{
				point: geometry.Point{
					Lat: 36.18416473150645,
					Lng: -86.73036575317383,
				},
				polygon: poly,
			},
			want:    false,
			wantErr: false,
		},
		"point out of Polygon - really close to polygon": {
			args: args{
				point: geometry.Point{
					Lat: 36.18200632243299,
					Lng: -86.74175441265106,
				},
				polygon: poly,
			},
			want:    false,
			wantErr: false,
		},
		"point in Polygon - on boundary": {
			args: args{
				point: geometry.Point{
					Lat: 36.171278341935434,
					Lng: -86.76624298095703,
				},
				polygon: poly,
			},
			want:    true,
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := PointInPolygon(tt.args.point, tt.args.polygon)
			if (err != nil) != tt.wantErr {
				t.Errorf("PointInPolygon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PointInPolygon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeatureCollection(t *testing.T) {
	// test for a simple Polygon
	coords := []geometry.Point{
		{
			Lat: 0,
			Lng: 0,
		},
		{
			Lat: 0,
			Lng: 100,
		},
		{
			Lat: 100,
			Lng: 100,
		},
		{
			Lat: 100,
			Lng: 0,
		},
		{
			Lat: 0,
			Lng: 0,
		},
	}

	ml := []geometry.LineString{}
	ln, err := geometry.NewLineString(coords)
	assert.Nil(t, err, "error message %s", err)
	ml = append(ml, *ln)

	poly, err := geometry.NewPolygon(ml)
	assert.Nil(t, err, "error message %s", err)

	ptIn := geometry.Point{
		Lat: 50,
		Lng: 50,
	}

	ptOut := geometry.Point{
		Lat: 140,
		Lng: 150,
	}

	pip, err := PointInPolygon(ptIn, *poly)
	assert.Nil(t, err, "error message %s", err)
	assert.True(t, pip, "Point in not in Polygon")

	pop, err := PointInPolygon(ptOut, *poly)
	assert.Nil(t, err, "error message %s", err)
	assert.False(t, pop, "Point in not in Polygon")
}

func TestPolyWithHole(t *testing.T) {
	ptInHole := geometry.Point{
		Lat: 36.20373274711739,
		Lng: -86.69208526611328,
	}
	ptInPoly := geometry.Point{
		Lat: 36.20258997094334,
		Lng: -86.72229766845702,
	}
	ptOutsidePoly := geometry.Point{
		Lat: 36.18527313913089,
		Lng: -86.75079345703125,
	}

	fix, err := utils.LoadJSONFixture(PolyWithHoleFixture)
	assert.NoError(t, err, "error loading fixture")

	f, err := feature.FromJSON(fix)
	assert.NoError(t, err, "error decoding json to feature")
	assert.NotNil(t, f, "feature is nil")

	assert.Equal(t, f.Type, geojson.Feature, "invalid base type")
	props := map[string]interface{}{
		"name":     "Poly with Hole",
		"value":    float64(3),
		"filename": "poly-with-hole.json",
	}
	assert.Equal(t, f.Properties, props, "invalid properties")
	assert.Equal(t, f.Bbox, []float64{-86.73980712890625, 36.173495506147, -86.67303085327148, 36.23084281427824}, "invalid properties object")
	assert.Equal(t, f.Geometry.GeoJSONType, geojson.Polygon, "invalid geojson type")

	poly, err := f.ToPolygon()
	assert.NoError(t, err, "error converting feature to polygon")

	pih, err := PointInPolygon(ptInHole, *poly)
	assert.Nil(t, err, "error message %s", err)
	assert.False(t, pih, "Point in hole is not in Polygon")

	pip, err := PointInPolygon(ptInPoly, *poly)
	assert.Nil(t, err, "error message %s", err)
	assert.True(t, pip, "Point in poly is not in Polygon")

	pop, err := PointInPolygon(ptOutsidePoly, *poly)
	assert.Nil(t, err, "error message %s", err)
	assert.False(t, pop, "Point in not in Polygon")
}

func TestMultiPolyWithHole(t *testing.T) {
	ptInHole := geometry.Point{
		Lat: 36.20373274711739,
		Lng: -86.69208526611328,
	}
	ptInPoly := geometry.Point{
		Lat: 36.20258997094334,
		Lng: -86.72229766845702,
	}
	ptInPoly2 := geometry.Point{
		Lat: 36.18527313913089,
		Lng: -86.75079345703125,
	}
	ptOutsidePoly := geometry.Point{
		Lat: 36.23015046460186,
		Lng: -86.75302505493164,
	}

	fixture, err := utils.LoadJSONFixture(MultiPolyWithHoleFixture)
	assert.NoError(t, err, "error loading fixture")

	f, err := feature.FromJSON(fixture)
	assert.NoError(t, err, "error decoding json to feature")
	assert.NotNil(t, f, "feature is nil")

	assert.Equal(t, f.Type, geojson.Feature, "invalid base type")
	props := map[string]interface{}{
		"name":     "Poly with Hole",
		"value":    float64(3),
		"filename": "poly-with-hole.json",
	}
	assert.Equal(t, f.Properties, props, "invalid properties")
	assert.Equal(t, f.Bbox, []float64{-86.77362442016602, 36.170862616662134, -86.67303085327148, 36.23084281427824}, "invalid properties object")
	assert.Equal(t, f.Geometry.GeoJSONType, geojson.MultiPolygon, "invalid geojson type")

	poly, err := f.ToMultiPolygon()
	assert.NoError(t, err, "error converting feature to MultiPolygon")

	pih := PointInMultiPolygon(ptInHole, *poly)
	assert.False(t, pih, "Point in hole is not in MultiPolygon")

	pip := PointInMultiPolygon(ptInPoly, *poly)
	assert.True(t, pip, "Point in poly is not in MultiPolygon")

	pip2 := PointInMultiPolygon(ptInPoly2, *poly)
	assert.True(t, pip2, "Point in poly is not in MultiPolygon")

	pop := PointInMultiPolygon(ptOutsidePoly, *poly)
	assert.False(t, pop, "Point in not in MultiPolygon")
}