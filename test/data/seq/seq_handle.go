package seqtester

import (
	"fmt"

	"github.com/petri-nets/workflow/flow"
)

type SeqRouter struct {
	preJobID int
}

func (s *SeqRouter) Handle(cs *flow.Case, w *flow.Workitem) {
	if w.WfWorkitem.JobID < s.preJobID {
		panic(fmt.Sprintf("pre job id: %d, current job id: %d",
			s.preJobID, w.WfWorkitem.JobID))
	}
	s.preJobID = w.WfWorkitem.JobID
}
