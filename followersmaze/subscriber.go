// followersmaze package implement subscriber.Subscriber interface
// using User data type.

// Supported events:
//      'F' Follow: Add event source to event recipient follower list,
//                  notify recipient
//      'U' Unfollow: Remove event source from event recipient follower list
//      'P' Private Message: Notify event recipient of a new private message
//      'S' Status Update: Notify all followers of event source
package followersmaze

import (
	"fmt"
	"io"
	"log"

	"github.com/andreadipersio/efr/event"
)

const (
	FOLLOW_ETYPE          = "F"
	UNFOLLOW_ETYPE        = "U"
	PRIVATE_MESSAGE_ETYPE = "P"
	STATUS_UPDATE_ETYPE   = "S"
)

type User struct {
	id        string
	followers map[string]event.Subscriber
	conn      io.WriteCloser
}

func (u *User) Connect(c io.WriteCloser) {
	u.conn = c
}

func (u *User) Disconnect() {
	if u.conn != nil {
		u.conn.Close()
	}
}

func (u *User) IsConnected() bool {
	return u.conn != nil
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) SetID(ID string) {
	u.id = ID
}

func (sender *User) HandleEvent(e event.Event, recipient event.Subscriber) error {
	switch e.EventType() {
	// Unfollow
	case UNFOLLOW_ETYPE:
		recipient.RemoveFollower(e.SenderID())

	// Follow
	case FOLLOW_ETYPE:
		recipient.NewFollower(sender)
		recipient.SendEvent(e)

	// Private message
	case PRIVATE_MESSAGE_ETYPE:
		recipient.SendEvent(e)

	// Status Update
	case STATUS_UPDATE_ETYPE:
		sender.followersBroadcast(e)
	default:
		return fmt.Errorf("Unsupported event %v", e)
	}

	return nil
}

func (u *User) SendEvent(e event.Event) {
	// user is not connected, ignore silently
	if u.conn == nil {
		return
	}

	_, err := fmt.Fprintf(u.conn, "%v\n", e)

	if err != nil {
		log.Printf("*** Cannot send notification %v: %v", e, err)
	}
}

func (u *User) String() string {
	return u.id
}

func (u *User) GetFollowers() []event.Subscriber {
	followers := []event.Subscriber{}

	for _, f := range u.followers {
		followers = append(followers, f)
	}

	return followers
}

func (u *User) NewFollower(follower event.Subscriber) {
	u.followers[follower.GetID()] = follower
}

func (u *User) RemoveFollower(followerID string) {
	delete(u.followers, followerID)
}

func (u *User) followersBroadcast(e event.Event) {
	for _, follower := range u.GetFollowers() {
		follower.SendEvent(e)
	}
}

func (u *User) Init() {
	u.followers = map[string]event.Subscriber{}
}

func NewUser(ID string) event.Subscriber {
	u := &User{}

	u.id = ID
	u.Init()

	return u
}
