package peerex

import (
	"fmt"
	"time"

	mspex "hyperledger.abchain.org/client/hyfabric/msp"
	"hyperledger.abchain.org/client/hyfabric/utils"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/msp"
	fcommon "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	protoutils "github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

const (
	defaultTimeout = 30 * time.Second
)

//Verify 检查参数正确性 没有的构建默认值
func (r *RPCBuilder) Verify(invoke bool) error {
	err := r.ChaincodeEnv.verify()
	if err != nil {
		return err
	}
	if !invoke {
		if len(r.Peers) > 1 {
			r.Peers = r.Peers[:1]
			logger.Warning("query 目前只支持单节点 取第一组数据")
		}
	} else {
		err := r.OrderEnv.verify()
		if err != nil {
			return err
		}
	}
	for _, p := range r.Peers {
		err := p.verify()
		if err != nil {
			return err
		}
	}
	logger.Debug("检查参数正确性=======down")
	return nil
}

func (cc *ChaincodeEnv) verify() error {
	if cc.ChannelID == "" {
		return errors.New("channelID 不能为空")
	}
	// we need chaincode name for everything, including deploy
	if cc.ChaincodeName == "" {
		return errors.New("ChaincodeName 不能为空")
	}
	if utils.IsNullOrEmpty(cc.Function) {
		return errors.New("not fund functionName")
	}
	return nil
}

func (node *NodeEnv) verify() error {
	if node == nil {
		return errors.New("节点配置不能为空")
	}
	if utils.IsNullOrEmpty(node.Address) || utils.IsNullOrEmpty(node.HostnameOverride) {
		return errors.New("Address，HostnameOverride  不能为空")
	}
	if node.TLS {
		if utils.IsNullOrEmpty(node.RootCertFile) {
			return errors.New("RootCertFile  不能为空")
		}
	}

	if node.ConnTimeout == time.Duration(0) {
		node.ConnTimeout = defaultTimeout
		logger.Debug("ConnTimeout is 0,use default timeout:", defaultTimeout)
	}

	return nil
}

//InitCrypto 初始化msp 加密信息
func InitCrypto(m *mspex.MspEnv) error {
	_, err := m.InitCrypto()
	if err != nil {
		// Handle errors reading the config file
		logger.Errorf("Cannot run peer because %s", err.Error())
	}
	return err
}

//InitConn 初始化chaincode命令工厂
func (r *RPCBuilder) InitConn(isOrdererRequired bool) error {
	err := r.ChaincodeEnv.verify()
	if err != nil {
		return err
	}
	logger.Debug("========InitConn start:============")
	var peers []*PeerEnv
	if isOrdererRequired {
		peers = r.Peers
		err := r.OrderEnv.ClientConn()
		if err != nil {
			return errors.WithMessage(err, "orderer grpc conn err")
		}
		logger.Debug("----order grpc conn----")
	} else {
		if len(r.Peers) > 1 {
			// r.Peers = r.Peers[:1]
			logger.Warning("query 目前只支持单节点 取第一组数据", r.Peers[0].Address)
		}
		peers = r.Peers[:1]
	}
	for i := 0; i < len(peers); i++ {
		err := r.Peers[i].ClientConn()
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("peer[%d] grpc conn err", i))
		}
		logger.Debugf("----peer[%d] grpc conn----", i)
		fmt.Printf("----peer[%d] grpc conn---- \n", i)
	}
	return nil
}

func (r *RPCBuilder) CloseConn() {
	for i := 0; i < len(r.Peers); i++ {
		r.Peers[i].CloseConn()
		logger.Debugf("----peer[%d] close grpc conn----", i)
	}
	r.OrderEnv.CloseConn()
	logger.Debug("----order close grpc conn----")
}

func (r *RPCBuilder) VerifyConn(isOrdererRequired bool) error {
	logger.Debug("========InitConn start:============")
	var peers []*PeerEnv
	if isOrdererRequired {
		peers = r.Peers
		err := r.OrderEnv.VerifyConn()
		if err != nil {
			return errors.WithMessage(err, "orderer grpc Verifyconn err")
		}
		logger.Debug("----order grpc conn----")
	} else {
		peers = r.Peers[:1]
	}
	for i := 0; i < len(peers); i++ {
		err := r.Peers[i].VerifyConn()
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("peer[%d] grpc Verifyconn err", i))
		}
		logger.Debugf("----peer[%d] grpc conn----", i)
	}
	return nil
}

// getChaincodeSpec get chaincode spec from the  pramameters
func (cc *ChaincodeEnv) getChaincodeSpec(args [][]byte) *pb.ChaincodeSpec {
	spec := &pb.ChaincodeSpec{}
	funcname := cc.Function
	input := &pb.ChaincodeInput{}
	input.Args = append(input.Args, []byte(funcname))

	for _, arg := range args {
		input.Args = append(input.Args, arg)
	}

	logger.Debug("ChaincodeSpec input :", input, " funcname:", funcname)
	fmt.Println("ChaincodeSpec input :", input, " funcname:", funcname)
	var golang = pb.ChaincodeSpec_Type_name[1]
	spec = &pb.ChaincodeSpec{
		Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value[golang]),
		// ChaincodeId: &pb.ChaincodeID{Name: cc.ChaincodeName, Version: cc.ChaincodeVersion},
		ChaincodeId: &pb.ChaincodeID{Name: cc.ChaincodeName},
		Input:       input,
	}
	return spec
}

// CreateChaincodeProposalWithTxIDAndTransient creates a proposal from given input
// It returns the proposal and the transaction id associated to the proposal
func CreateChaincodeProposalWithTxIDAndTransient(chainID string, spec *pb.ChaincodeSpec, creator []byte, transientMap map[string][]byte) (*pb.Proposal, string, error) {
	// generate a random nonce
	nonce, err := utils.GetRandomNonce()
	if err != nil {
		return nil, "", err
	}
	txid, err := protoutils.ComputeProposalTxID(nonce, creator)
	if err != nil {
		return nil, "", err
	}
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}
	ccHdrExt := &pb.ChaincodeHeaderExtension{ChaincodeId: spec.ChaincodeId}

	ccHdrExtBytes, err := protoutils.Marshal(ccHdrExt)
	if err != nil {
		return nil, "", err
	}

	cisBytes, err := protoutils.Marshal(invocation)
	if err != nil {
		return nil, "", err
	}

	ccPropPayload := &pb.ChaincodeProposalPayload{Input: cisBytes, TransientMap: transientMap}
	ccPropPayloadBytes, err := protoutils.Marshal(ccPropPayload)
	if err != nil {
		return nil, "", err
	}

	// TODO: epoch is now set to zero. This must be changed once we
	// get a more appropriate mechanism to handle it in.
	var (
		epoch     uint64
		timestamp = util.CreateUtcTimestamp()
		typ       = int32(fcommon.HeaderType_ENDORSER_TRANSACTION)
	)

	channelHeader, err := protoutils.Marshal(&fcommon.ChannelHeader{
		Type:      typ,
		TxId:      txid,
		Timestamp: timestamp,
		ChannelId: chainID,
		Extension: ccHdrExtBytes,
		Epoch:     epoch,
	})
	if err != nil {
		return nil, "", err
	}
	signatureHeader, err := protoutils.Marshal(&fcommon.SignatureHeader{
		Nonce:   nonce,
		Creator: creator,
	})

	if err != nil {
		return nil, "", err
	}

	hdr := &fcommon.Header{
		ChannelHeader:   channelHeader,
		SignatureHeader: signatureHeader,
	}

	hdrBytes, err := protoutils.Marshal(hdr)
	if err != nil {
		return nil, "", err
	}
	return &pb.Proposal{Header: hdrBytes, Payload: ccPropPayloadBytes}, txid, nil
}

// GetSignedProposal returns a signed proposal given a Proposal message and a signing identity
func GetSignedProposal(prop *pb.Proposal, signer msp.SigningIdentity) (*pb.SignedProposal, error) {
	// check for nil argument
	if prop == nil || signer == nil {
		return nil, fmt.Errorf("Nil arguments")
	}

	propBytes, err := protoutils.Marshal(prop)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(propBytes)
	if err != nil {
		return nil, err
	}

	return &pb.SignedProposal{ProposalBytes: propBytes, Signature: signature}, nil
}

// Serialize returns a byte array representation of this identity
// func (id *identity) Serialize() ([]byte, error) {
// 	fmt.Println(`F:\virtualMachineShare\src\github.com\hyperledger\fabric\msp\identities.go Serialize()`, id.id.Mspid)
// 	pb := &pem.Block{Bytes: id.cert.Raw, Type: "CERTIFICATE"}
// 应该是msp/signcerts 的读取
// 	pemBytes := pem.EncodeToMemory(pb)
// 	if pemBytes == nil {
// 		return nil, errors.New("encoding of identity failed")
// 	}

// 	// We serialize identities by prepending the MSPID and appending the ASN.1 DER content of the cert
// 	sId := &msp.SerializedIdentity{Mspid: id.id.Mspid, IdBytes: pemBytes}
// 	idBytes, err := proto.Marshal(sId)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "could not marshal a SerializedIdentity structure for identity %s", id.id)
// 	}

// 	return idBytes, nil
// }

// Sign produces a signature over msg, signed by this instance
// func (id *signingidentity) Sign(msg []byte) ([]byte, error) {
// 	//mspIdentityLogger.Infof("Signing message")
// 	//fmt.Println(`F:\virtualMachineShare\src\github.com\hyperledger\fabric\msp\identities`)
// 	// Compute Hash
// 	hashOpt, err := id.getHashOpt(id.msp.cryptoConfig.SignatureHashFamily)
// 	if err != nil {
// 		return nil, errors.WithMessage(err, "failed getting hash function options")
// 	}

// 	digest, err := id.msp.bccsp.Hash(msg, hashOpt)
// 	if err != nil {
// 		return nil, errors.WithMessage(err, "failed computing digest")
// 	}

// 	if len(msg) < 32 {
// 		mspIdentityLogger.Debugf("Sign: plaintext: %X \n", msg)
// 	} else {
// 		mspIdentityLogger.Debugf("Sign: plaintext: %X...%X \n", msg[0:16], msg[len(msg)-16:])
// 	}
// 	mspIdentityLogger.Debugf("Sign: digest: %X \n", digest)

// 	// Sign
// 	return id.signer.Sign(rand.Reader, digest, nil)
// }
