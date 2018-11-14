package protos

import "fmt"
import "github.com/golang/protobuf/proto"

// Bytes returns this GlobalState as an array of bytes.
func (gs *GlobalState) Bytes() ([]byte, error) {
	data, err := proto.Marshal(gs)
	if err != nil {
		logger.Errorf("Error marshalling GlobalState: %s", err)
		return nil, fmt.Errorf("Could not marshal GlobalState: %s", err)
	}
	return data, nil
}

// Test state is a branch (or a node in the state graphic)
func (gs *GlobalState) Branched() bool {
	return len(gs.NextNodeStateHash) > 1 ||
		len(gs.ParentNodeStateHash) > 1
}

func (gs *GlobalState) NextNode() []byte {

	if gs == nil || len(gs.NextNodeStateHash) != 1 {
		return nil
	}

	return gs.NextNodeStateHash[0]
}

func (gs *GlobalState) ParentNode() []byte {
	if gs == nil || len(gs.ParentNodeStateHash) != 1 {
		return nil
	}

	return gs.ParentNodeStateHash[0]
}

// NewBlock creates a new Block given the input parameters.
func NewGlobalState() *GlobalState {
	gs := new(GlobalState)
	return gs
}

func UnmarshallGS(gsBytes []byte) (*GlobalState, error) {
	gs := &GlobalState{}
	err := proto.Unmarshal(gsBytes, gs)
	if err != nil {
		logger.Errorf("Error unmarshalling GlobalState: %s", err)
		return nil, fmt.Errorf("Could not unmarshal GlobalState: %s", err)
	}
	return gs, nil
}
