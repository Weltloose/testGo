package mongodb

type UserInfo struct {
	Username string
	Password string
	Groups   []int
}

type Event struct {
	EventID  int    `json:"eventID"`
	TimeFrom string `json:"timeFrom"`
	TimeTo   string `json:"timeTo"`
	Content  string `json:"content"`
}

type GroupInfo struct {
	GroupID int
	Users   []string
	Events  []Event
}

type EventListOpt struct {
	Data map[int][]Event `json:"data"`
}

type WSopt struct {
	GroupID int     `json:"groupID"`
	Events  []Event `json:"events"`
}
