module github.com/pions/webrtc

replace github.com/pions/webrtc/pkg/quic => ./pkg/quic

require (
	github.com/pions/datachannel v1.1.0
	github.com/pions/dtls v1.0.2
	github.com/pions/sctp v1.1.0
	github.com/pions/sdp v1.1.0
	github.com/pions/stun v0.1.0
	github.com/pions/transport v0.0.0-20190123145644-fbbbdd95131a
	github.com/pions/webrtc/pkg/quic v0.0.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
)
