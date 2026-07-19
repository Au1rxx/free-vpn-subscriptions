package maintain

type Policy struct {
	RawPayloadDays, ParseErrorDays, AttemptDays, FetchDays, ExportDays int
	PauseColdSources, StoreRawBodies                                   bool
}

func PolicyForUsage(percent float64) Policy {
	policy := Policy{RawPayloadDays: 30, ParseErrorDays: 90, AttemptDays: 180, FetchDays: 365, ExportDays: 365, StoreRawBodies: true}
	switch {
	case percent >= 94:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.FetchDays, policy.ExportDays = 1, 7, 30, 60, 90
		policy.PauseColdSources, policy.StoreRawBodies = true, false
	case percent >= 90:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.FetchDays, policy.ExportDays = 3, 14, 60, 90, 180
		policy.PauseColdSources, policy.StoreRawBodies = true, false
	case percent >= 80:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.FetchDays, policy.ExportDays = 7, 30, 90, 180, 180
	case percent >= 70:
		policy.RawPayloadDays, policy.ParseErrorDays = 14, 30
	}
	return policy
}
