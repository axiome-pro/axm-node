package types

func (r *EmissionRange) Equal(r2 *EmissionRange) bool {
	return r.Rate.Equal(r2.Rate) && r.Start.Equal(r2.Start)
}
