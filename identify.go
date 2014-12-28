package main

import (
	"github.com/keybase/go-libkb"
	"github.com/keybase/protocol/go"
	fmprpc "github.com/maxtaco/go-framed-msgpack-rpc"
	"github.com/ugorji/go/codec"
	"net"
	"net/rpc"
)

type RemoteIdentifyUI struct {
	sessionId int
	rpccli    *rpc.Client
	uicli     keybase_1.IdentifyUiClient
}

type IdentifyHandler struct {
	conn net.Conn
	cli  *rpc.Client
}

func (h *IdentifyHandler) GetRpcClient() (cli *rpc.Client) {
	if cli = h.cli; cli == nil {
		var mh codec.MsgpackHandle
		cdc := fmprpc.MsgpackSpecRpc.ClientCodec(h.conn, &mh, true)
		cli = rpc.NewClientWithCodec(cdc)
		h.cli = cli
	}
	return
}

func NewRemoteIdentifyUI(sessionId int, c *rpc.Client) *RemoteIdentifyUI {
	return &RemoteIdentifyUI{
		sessionId: sessionId,
		rpccli:    c,
		uicli:     keybase_1.IdentifyUiClient{c},
	}
}

func (h *IdentifyHandler) NewUi(sessionId int) libkb.IdentifyUI {
	return NewRemoteIdentifyUI(sessionId, h.GetRpcClient())
}

func (u *RemoteIdentifyUI) FinishWebProofCheck(p keybase_1.RemoteProof, lcr keybase_1.LinkCheckResult) {
	var status keybase_1.Status
	u.uicli.FinishWebProofCheck(keybase_1.FinishWebProofCheckArg{
		SessionId: u.sessionId,
		Rp:        p,
		Lcr:       lcr,
	}, &status)
	return
}

func (u *RemoteIdentifyUI) FinishSocialProofCheck(p keybase_1.RemoteProof, lcr keybase_1.LinkCheckResult) {
	var status keybase_1.Status
	u.uicli.FinishSocialProofCheck(keybase_1.FinishSocialProofCheckArg{
		SessionId: u.sessionId,
		Rp:        p,
		Lcr:     lcr,
	}, &status)
	return
}

func (u *RemoteIdentifyUI) FinishAndPrompt(io *keybase_1.IdentifyOutcome) (ret keybase_1.FinishAndPromptRes) {
	err := u.uicli.FinishAndPrompt(keybase_1.FinishAndPromptArg{SessionId: u.sessionId, Outcome: *io, }, &ret)
	if err != nil {
		ret.Status = libkb.ExportErrorAsStatus(err)
	}
	return
}

func (u *RemoteIdentifyUI) DisplayCryptocurrency(c keybase_1.Cryptocurrency) {
	var status keybase_1.Status
	u.uicli.DisplayCryptocurrency(keybase_1.DisplayCryptocurrencyArg{SessionId: u.sessionId, C: c }, &status)
	return
}

func (u *RemoteIdentifyUI) DisplayKey(k keybase_1.FOKID, d *keybase_1.TrackDiff) {
	var status keybase_1.Status
	u.uicli.DisplayKey(keybase_1.DisplayKeyArg{SessionId: u.sessionId, Fokid: k, Diff: d}, &status)
	return
}

func (u *RemoteIdentifyUI) ReportLastTrack(t *keybase_1.TrackSummary) {
	var status keybase_1.Status
	u.uicli.ReportLastTrack(keybase_1.ReportLastTrackArg{SessionId: u.sessionId, Track: t}, &status)
	return
}

func (u *RemoteIdentifyUI) Start() {}

func (u *RemoteIdentifyUI) LaunchNetworkChecks(id *keybase_1.Identity) {
	var status keybase_1.Status
	u.uicli.LaunchNetworkChecks(keybase_1.LaunchNetworkChecksArg{
		SessionId: u.sessionId,
		Id:        *id,
	}, &status)
	return
}

func (h *IdentifyHandler) IdentifySelf(sessionId *int, res *keybase_1.Status) error {
	luarg := libkb.LoadUserArg{}
	u, err := libkb.LoadMe(luarg)
	if _, not_found := err.(libkb.NoKeyError); not_found {
		err = nil
	} else if _, not_selected := err.(libkb.NoSelectedKeyError); not_selected {
		_, err = u.IdentifySelf(h.NewUi(*sessionId))
	}
	status := libkb.ExportErrorAsStatus(err)
	res = &status
	return nil
}
