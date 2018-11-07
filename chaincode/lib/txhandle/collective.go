package tx

import (
	"fmt"
)

type CollectiveTxs map[string]*ChaincodeTx

func NewCollectiveTxs() CollectiveTxs { return CollectiveTxs(make(map[string]*ChaincodeTx)) }

func (s CollectiveTxs) mergeone(in CollectiveTxs) error {

	for k, v := range in {
		if _, ok := s[k]; ok {
			return fmt.Errorf("Method [%s] is collision", k)
		}

		s[k] = v
	}

	return nil

}

func (s CollectiveTxs) Merge(ins ...CollectiveTxs) (CollectiveTxs, error) {

	for _, in := range ins {
		if err := s.mergeone(in); err != nil {
			return s, err
		}
	}

	return s, nil

}
