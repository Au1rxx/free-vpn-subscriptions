package maintain

type Policy struct {
	RawPayloadDays, ParseErrorDays, AttemptDays, BatchDays, FetchDays, ExportDays int
	PauseColdSources                                                              bool
}

func PolicyForUsage(percent float64) Policy {
	policy := Policy{RawPayloadDays: 30, ParseErrorDays: 30, AttemptDays: 14, BatchDays: 14, FetchDays: 90, ExportDays: 30}
	switch {
	case percent >= 94:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.BatchDays, policy.FetchDays, policy.ExportDays = 1, 1, 1, 1, 7, 1
		policy.PauseColdSources = true
	case percent >= 90:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.BatchDays, policy.FetchDays, policy.ExportDays = 3, 3, 2, 2, 14, 3
		policy.PauseColdSources = true
	case percent >= 80:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.BatchDays, policy.FetchDays, policy.ExportDays = 7, 7, 3, 3, 30, 7
	case percent >= 70:
		policy.RawPayloadDays, policy.ParseErrorDays, policy.AttemptDays, policy.BatchDays, policy.FetchDays, policy.ExportDays = 14, 14, 7, 7, 60, 14
	}
	return policy
}
