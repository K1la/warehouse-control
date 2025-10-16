package dto

type HistoryResponse struct {
	ItemID   int64  `json:"item_id"`
	Action   string `json:"action"`
	UserID   int64  `json:"user_id"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
	Date     string `json:"date"`
}
