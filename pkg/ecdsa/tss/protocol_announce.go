package tss

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

const protocolAnnounceTimeout = 120 * time.Second

func AnnounceProtocol(
	parentCtx context.Context,
	publicKey *operator.PublicKey,
	membersCount int,
	broadcastChannel net.BroadcastChannel,
) (
	[]MemberID,
	error,
) {
	logger.Infof("announcing presence")

	ctx, cancel := context.WithTimeout(parentCtx, protocolAnnounceTimeout)
	defer cancel()

	announceInChan := make(chan *AnnounceMessage, membersCount)
	handleAnnounceMessage := func(netMsg net.Message) {
		switch msg := netMsg.Payload().(type) {
		case *AnnounceMessage:
			announceInChan <- msg
		}
	}
	broadcastChannel.Recv(ctx, handleAnnounceMessage)

	receivedMemberIDs := make(map[string]MemberID)

	go func() {

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-announceInChan:
				// Since broadcast channel has an address filter, we can
				// assume each message come from a valid group member.
				receivedMemberIDs[msg.SenderID.String()] = msg.SenderID

				if len(receivedMemberIDs) == membersCount {
					cancel()
				}
			}
		}
	}()

	go func() {
		sendMessage := func() {
			if err := broadcastChannel.Send(ctx,
				&AnnounceMessage{
					SenderID: MemberIDFromPublicKey(publicKey),
				},
			); err != nil {
				logger.Errorf("failed to send announcement: [%v]", err)
			}
		}

		for {
			select {
			case <-ctx.Done():
				// Send the message once again as the member received messages
				// from all peer members but not all peer members could receive
				// the message from the member as some peer member could join
				// the protocol after the member sent the last message.
				sendMessage()
				return
			default:
				sendMessage()
				time.Sleep(1 * time.Second)
			}
		}
	}()

	<-ctx.Done()
	switch ctx.Err() {
	case context.DeadlineExceeded:
		return nil, fmt.Errorf(
			"waiting for announcements timed out after: [%v]",
			protocolAnnounceTimeout,
		)
	case context.Canceled:
		logger.Infof("announce protocol completed successfully")

		memberIDs := make([]MemberID, 0)
		for _, memberID := range receivedMemberIDs {
			memberIDs = append(memberIDs, memberID)
		}

		return memberIDs, nil
	default:
		return nil, fmt.Errorf("unexpected context error: [%v]", ctx.Err())
	}
}