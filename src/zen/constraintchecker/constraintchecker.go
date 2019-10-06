package constraintchecker

import "zen/boundschecking"

type ConstraintCheckerState struct {
	knownConstraints *boundschecking.KnownConstraints
}
