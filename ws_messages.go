package main

// Type を確認する用
type message struct {
	Type string `json:"type"`
}

type registerMessage struct {
	Type          string       `json:"type"`
	RoomID        string       `json:"roomId"`
	ClientID      string       `json:"clientId"`
	AuthnMetadata *interface{} `json:"authnMetadata"`
	SignalingKey  *string      `json:"signalingKey"`
	// 後方互換性対応
	Key *string `json:"key"`
	// Ayame クライアント情報が詰まっている
	AyameClient *string `json:"ayameClient"`
	Libwebrtc   *string `json:"libwebrtc"`
	Environment *string `json:"environment"`
}

type pingMessage struct {
	Type string `json:"type"`
}

type byeMessage struct {
	Type string `json:"type"`
}

type knockMessage struct {
	Type string `json:"type"`
	KnockID string `json:"knockId"`
}

type replyMessage struct {
	Type string `json:"type"`
	Result int `json:"result"`
}

type acceptMessage struct {
	Type         string `json:"type"`
	ConnectionID string `json:"connectionId"`
	// WaitingOffer  bool         `json:"waitingOffer"`
	AuthzMetadata *interface{} `json:"authzMetadata,omitempty"`
	IceServers    *[]iceServer `json:"iceServers,omitempty"`

	// 後方互換性対応
	IsExistClient bool `json:"isExistClient"`
	// 後方互換性対応
	IsExistUser bool `json:"isExistUser"`
}

type rejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type iceServer struct {
	Urls       []string `json:"urls"`
	UserName   *string  `json:"username,omitempty"`
	Credential *string  `json:"credential,omitempty"`
}
