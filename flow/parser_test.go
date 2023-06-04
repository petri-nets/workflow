package flow

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/petri-nets/workflow/wfmod"
)

func TestParser_Arrangement(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "eif",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 1",
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 2",
					},
				},
			},
			placeCnt:      5,
			transitionCnt: 4,
		},
		{
			name: "eif-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 1",
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 2",
					},
				},
				Then: &task{
					Name:    "t03",
					Trigger: wfmod.TransitionTriggerAuto,
				},
			},
			placeCnt:      6,
			transitionCnt: 5,
		},
		{
			name: "seq",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 1",
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 2",
					},
				},
			},
			placeCnt:      4,
			transitionCnt: 3,
		},
		{
			name: "seq-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 1",
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       22,
						Condition: "TDJOB_SHARD_NO == 2",
					},
				},
				Then: &task{
					Name:      "t03",
					Trigger:   wfmod.TransitionTriggerAuto,
					Job:       22,
					Condition: "TDJOB_SHARD_NO == 1",
				},
			},
			placeCnt:      5,
			transitionCnt: 4,
		},
		{
			name: "par",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     22,
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     22,
					},
				},
			},
			placeCnt:      6,
			transitionCnt: 4,
		},
		{
			name: "par-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     22,
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     22,
					},
				},
				Then: &task{
					Name:    "t03",
					Trigger: wfmod.TransitionTriggerAuto,
					Job:     22,
				},
			},
			placeCnt:      6,
			transitionCnt: 4,
		},
		{
			name: "vie",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerTime,
						Job:     22,
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerMsg,
						Job:     22,
					},
				},
			},
			placeCnt:      3,
			transitionCnt: 3,
		},
		{
			name: "vie-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerTime,
						Job:     22,
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerMsg,
						Job:     22,
					},
				},
				Then: &task{
					Name:    "t03",
					Trigger: wfmod.TransitionTriggerAuto,
					Job:     22,
				},
			},
			placeCnt:      4,
			transitionCnt: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_EIF_EIF(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "eif-eif",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       10,
						Condition: "TDJOB_SHARD_NO == 1",
						Router:    RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       20,
						Condition: "TDJOB_SHARD_NO == 2",
						Router:    RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      11,
			transitionCnt: 10,
		},
		{
			name: "eif-eif-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       10,
						Condition: "TDJOB_SHARD_NO == 1",
						Router:    RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t03",
							Trigger: wfmod.TransitionTriggerAuto,
							Job:     23,
						},
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       20,
						Condition: "TDJOB_SHARD_NO == 2",
						Router:    RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t04",
							Trigger: wfmod.TransitionTriggerAuto,
							Job:     24,
						},
					},
				},
			},
			placeCnt:      13,
			transitionCnt: 12,
		},
		{
			name: "eif-eif-then-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       10,
						Condition: "TDJOB_SHARD_NO == 1",
						Router:    RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t03",
							Trigger: wfmod.TransitionTriggerAuto,
							Job:     23,
						},
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       20,
						Condition: "TDJOB_SHARD_NO == 2",
						Router:    RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t04",
							Trigger: wfmod.TransitionTriggerAuto,
							Job:     24,
						},
					},
				},
				Then: &task{
					Name:    "t05",
					Trigger: wfmod.TransitionTriggerAuto,
					Job:     25,
				},
			},
			placeCnt:      14,
			transitionCnt: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_EIF_SEQ(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "eif-seq",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       10,
						Condition: "TDJOB_SHARD_NO == 1",
						Router:    RouterSEQ,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerAuto,
								Job:     11,
								Router:  RouterEIF,
								Tasks: []*task{
									{
										Name:      "t0111",
										Trigger:   wfmod.TransitionTriggerAuto,
										Job:       111,
										Condition: "TDJOB_SHARD_NO == 1",
									},
									{
										Name:      "t0112",
										Trigger:   wfmod.TransitionTriggerAuto,
										Job:       112,
										Condition: "TDJOB_SHARD_NO == 1",
									},
								},
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerAuto,
								Job:     12,
							},
						},
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       20,
						Condition: "TDJOB_SHARD_NO == 2",
						Router:    RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerAuto,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
								Router:    RouterEIF,
								Tasks: []*task{
									{
										Name:      "t0221",
										Trigger:   wfmod.TransitionTriggerAuto,
										Job:       111,
										Condition: "TDJOB_SHARD_NO == 111",
									},
									{
										Name:      "t0222",
										Trigger:   wfmod.TransitionTriggerAuto,
										Job:       112,
										Condition: "TDJOB_SHARD_NO == 112",
									},
								},
							},
						},
					},
				},
			},
			placeCnt:      15,
			transitionCnt: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_EIF_VIE(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "eif-vie",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       10,
						Condition: "TDJOB_SHARD_NO == 1",
						Router:    RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       20,
						Condition: "TDJOB_SHARD_NO == 2",
						Router:    RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      7,
			transitionCnt: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_EIF_PAR(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "eif-par",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterEIF,
				Tasks: []*task{
					{
						Name:      "t01",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       10,
						Condition: "TDJOB_SHARD_NO == 1",
						Router:    RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:      "t02",
						Trigger:   wfmod.TransitionTriggerAuto,
						Job:       20,
						Condition: "TDJOB_SHARD_NO == 2",
						Router:    RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      13,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_SEQ_PAR(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "seq-par",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      12,
			transitionCnt: 9,
		},
		{
			name: "seq-par-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
						Then: &task{
							Name: "t031",
							Job:  31,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name: "t032",
							Job:  32,
						},
					},
				},
				Then: &task{
					Name: "t03",
					Job:  3,
				},
			},
			placeCnt:      13,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_SEQ_VIE(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "seq-vie",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      6,
			transitionCnt: 7,
		},
		{
			name: "seq-vie-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      9,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_SEQ_EIF(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "seq-EIF",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 1",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 2",
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      10,
			transitionCnt: 9,
		},
		{
			name: "seq-eif-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      13,
			transitionCnt: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_SEQ_SEQ(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "seq-seq",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 1",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 2",
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      8,
			transitionCnt: 7,
		},
		{
			name: "seq-seq-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterSEQ,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      11,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_PAR_PAR(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "par-par",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      14,
			transitionCnt: 10,
		},
		{
			name: "par-par-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
						Then: &task{
							Name: "t031",
							Job:  3,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name: "t032",
							Job:  3,
						},
					},
				},
				Then: &task{
					Name: "t03",
					Job:  3,
				},
			},
			placeCnt:      14,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_PAR_VIE(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "par-vie",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      8,
			transitionCnt: 8,
		},
		{
			name: "par-vie-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      10,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_PAR_EIF(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "par-eif",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 1",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 2",
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      12,
			transitionCnt: 10,
		},
		{
			name: "par-eif-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      14,
			transitionCnt: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_PAR_SEQ(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "par-seq",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 1",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 2",
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      10,
			transitionCnt: 8,
		},
		{
			name: "par-seq-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterPAR,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      12,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_VIE_PAR(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "vie-par",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      11,
			transitionCnt: 9,
		},
		{
			name: "vie-par-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
						Then: &task{
							Name: "t031",
							Job:  3,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterPAR,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name: "t032",
							Job:  3,
						},
					},
				},
				Then: &task{
					Name: "t03",
					Job:  3,
				},
			},
			placeCnt:      12,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_VIE_VIE(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "vie-vie",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      5,
			transitionCnt: 7,
		},
		{
			name: "vie-vie-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:    "t011",
								Trigger: wfmod.TransitionTriggerTime,
								Job:     11,
							},
							{
								Name:    "t012",
								Trigger: wfmod.TransitionTriggerUser,
								Job:     12,
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterVIE,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      8,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_VIE_EIF(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "vie-eif",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 1",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 2",
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      9,
			transitionCnt: 9,
		},
		{
			name: "vie-eif-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterEIF,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      12,
			transitionCnt: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}

func TestParser_Arrangement_VIE_SEQ(t *testing.T) {
	tests := []struct {
		name          string
		pipeline      Pipeline
		placeCnt      int
		transitionCnt int
	}{
		{
			name: "vie-seq",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 1",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 2",
							},
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
					},
				},
			},
			placeCnt:      7,
			transitionCnt: 7,
		},
		{
			name: "vie-seq-then",
			pipeline: Pipeline{
				Name:       "w01",
				AppID:      123,
				StartJobID: 11,
				Desc:       "desc xxx",
				Router:     RouterVIE,
				Tasks: []*task{
					{
						Name:    "t01",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     10,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t011",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t012",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t032",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     32,
						},
					},
					{
						Name:    "t02",
						Trigger: wfmod.TransitionTriggerAuto,
						Job:     20,
						Router:  RouterSEQ,
						Tasks: []*task{
							{
								Name:      "t021",
								Trigger:   wfmod.TransitionTriggerUser,
								Job:       11,
								Condition: "TDJOB_SHARD_NO == 11",
							},
							{
								Name:      "t022",
								Trigger:   wfmod.TransitionTriggerTime,
								Job:       12,
								Condition: "TDJOB_SHARD_NO == 12",
							},
						},
						Then: &task{
							Name:    "t033",
							Trigger: wfmod.TransitionTriggerUser,
							Job:     33,
						},
					},
				},
				Then: &task{
					Name:    "t031",
					Trigger: wfmod.TransitionTriggerUser,
					Job:     31,
				},
			},
			placeCnt:      10,
			transitionCnt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Pipeline: &tt.pipeline,
			}
			p.Arrangement()
			t.Logf("%#v", p.transitions)
			t.Logf("%#v", p.places)
			assert.Equal(t, len(p.places), tt.placeCnt)
			assert.Equal(t, len(p.transitions), tt.transitionCnt)
		})
	}
}
