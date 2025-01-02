package types

type TalkSessionCreateRequest struct {
	// TalkMode indicates the chat mode: 1 for private chat, 2 for group chat
	TalkMode int32 `json:"talk_mode" binding:"required,oneof=1 2"`

	// ToFromID represents the ID of the other participant
	ToFromID int32 `json:"to_from_id" binding:"required"`
}
