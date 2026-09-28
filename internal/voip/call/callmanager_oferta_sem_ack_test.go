package call

import (
	"testing"

	"wacalls/internal/voip/core"
)

// AdvBrasil (28/09/2026): a oferta que o WhatsApp não responde (sem `ack`) encerra com `offer_unacked`; a que
// teve `ack` nunca é encerrada pelo prazo.

func TestOfertaSemAckEncerraComOfferUnacked(t *testing.T) {
	m := ringingManager(&fakeSock{}, nil)
	var fim *CallInfo
	m.OnEnded = func(c *CallInfo) { fim = c }

	m.encerrarSeSemAck("CALL1")

	if fim == nil {
		t.Fatalf("esperava OnEnded com offer_unacked, a chamada seguiu viva")
	}
	if got := fim.StateData.EndReason; got != core.EndCallReasonOfferUnacked {
		t.Fatalf("motivo do fim: esperava %q, veio %q", core.EndCallReasonOfferUnacked, got)
	}
}

func TestOfertaComAckNaoEncerraPeloPrazo(t *testing.T) {
	m := ringingManager(&fakeSock{}, nil)
	m.currentCall.RelayData = &core.RelayData{} // o ack chegou
	ended := false
	m.OnEnded = func(*CallInfo) { ended = true }

	m.encerrarSeSemAck("CALL1")

	if ended || m.currentCall.IsEnded() {
		t.Fatalf("a chamada com ack foi encerrada pelo prazo do ack")
	}
}

func TestPrazoDeOutraChamadaNaoEncerraAAtual(t *testing.T) {
	m := ringingManager(&fakeSock{}, nil)
	ended := false
	m.OnEnded = func(*CallInfo) { ended = true }

	m.encerrarSeSemAck("OUTRA")

	if ended {
		t.Fatalf("o prazo de uma chamada antiga encerrou a chamada atual")
	}
}

func TestPeerServerNaoExpoeONumero(t *testing.T) {
	if got := peerServer("5562999994776@s.whatsapp.net"); got != "s.whatsapp.net" {
		t.Fatalf("peerServer: %q", got)
	}
	if got := peerServer("123@lid"); got != "lid" {
		t.Fatalf("peerServer: %q", got)
	}
}
