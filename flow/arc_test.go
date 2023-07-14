package flow

import (
	"testing"

	"github.com/petri-nets/workflow/wfmod"
)

func TestArc_CheckPass(t *testing.T) {
	type fields struct {
		WfArc wfmod.WfArc
	}
	type args struct {
		wfContext wfmod.WfContextType
	}
	tests := []struct {
		name      string
		wfArc     wfmod.WfArc
		wfContext wfmod.WfContextType
		want      bool
	}{
		{
			name:      "no condition",
			wfArc:     wfmod.WfArc{Type: wfmod.ArcOutExplicitORsplit, Condition: ""},
			wfContext: wfmod.WfContextType{},
			want:      true,
		},
		{
			name:  "cond01",
			wfArc: wfmod.WfArc{Type: wfmod.ArcOutExplicitORsplit, Condition: "name == 'test001' && age == 100 && age > 99 && age < 101 && age <= 100 && !sleep"},
			wfContext: wfmod.WfContextType{
				"name":  "test001",
				"age":   100,
				"sleep": false,
			},
			want: true,
		},
		{
			name: "cal",
			wfArc: wfmod.WfArc{
				Type:      wfmod.ArcOutExplicitORsplit,
				Condition: "get + put == 105 && get * put == 500 && get - put == 95 && get / put == 20 && set == 1003.256 && set > 1003.25 && set < 1003.257",
			},
			wfContext: wfmod.WfContextType{
				"get": 100,
				"put": 5,
				"set": 1003.256,
			},
			want: true,
		},
		{
			name: "cal01",
			wfArc: wfmod.WfArc{
				Type:      wfmod.ArcOutExplicitORsplit,
				Condition: "set * 5 == 5016.28",
			},
			wfContext: wfmod.WfContextType{
				"get": 100,
				"put": 5,
				"set": 1003.256,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Arc{
				WfArc: tt.wfArc,
			}
			if got := a.CheckPass(tt.wfContext); got != tt.want {
				t.Errorf("Arc.CheckPass() = %v, want %v", got, tt.want)
			}
		})
	}
}
