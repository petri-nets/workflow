package datatester

import "github.com/petri-nets/workflow/flow"

type DTT interface {
	Handle(cs *flow.Case, w *flow.Workitem)
}
