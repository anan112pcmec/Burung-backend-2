package notify_payload

type TableAndAction struct {
	Table  string `json:"table"`
	Action string `json:"action"`
}

type ColumnChangeInfo struct {
	ChangedColumns     map[string]interface{} `json:"changed_columns"`
	ColumnChangedNames []string               `json:"column_change_name"`
}

// ////////////////////////////////////////////////////////////////////////////
// BARANG PAYLOAD
// ////////////////////////////////////////////////////////////////////////////
