package protos

import (
	"bytes"
	bin "encoding/binary"
	"github.com/abchain/fabric/core/util"
)

//the max size is limited by the msg size in grpc, so no overflow problem
//it was just an estimation of the whole size of package
func (m *GossipMsg) EstimateSize() (total int) {

	if m == nil {
		return
	}

	switch v := m.GetM().(type) {
	case (*GossipMsg_Dig):
		for _, i := range v.Dig.Data {
			total = total + len(i.State) + 8 //the bytes of a num
		}
		total = total + len(v.Dig.Epoch)
	case (*GossipMsg_Ud):
		total = len(v.Ud.Payload)
	default:
		return
	}

	return
}

// generate a consistent byte represent
func (s *PeerTxState) MsgEncode(id string) ([]byte, error) {

	var buf bytes.Buffer
	w := util.NewConWriter(&buf)

	if err := w.Write([]byte(id)).Write(s.Digest).Write(s.Endorsement).Error(); err != nil {
		return nil, err
	}

	if err := bin.Write(&buf, bin.BigEndian, s.Num); err != nil {
		return nil, err
	}

	if err := bin.Write(&buf, bin.BigEndian, s.EndorsementVer); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
