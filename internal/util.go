package internal

func mapFilter[V any, R any](a []V, f func(V) (R, bool)) []R {
	r := make([]R, 0, len(a))
	for _, v := range a {
		tmp, ok := f(v)
		if ok {
			r = append(r, tmp)
		}
	}
	return r
}
