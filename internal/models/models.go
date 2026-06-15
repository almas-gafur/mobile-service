package models

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMaster Role = "master"
)

type User struct {
	ID        int64
	Username  string
	Password  string
	Role      Role
	CreatedAt time.Time
}

type OrderStatus string

const (
	StatusAccepted   OrderStatus = "accepted"
	StatusDiagnosis  OrderStatus = "diagnosis"
	StatusWaiting    OrderStatus = "waiting_parts"
	StatusInProgress OrderStatus = "in_progress"
	StatusReady      OrderStatus = "ready"
	StatusIssued     OrderStatus = "issued"
)

func (s OrderStatus) Label() string {
	labels := map[OrderStatus]string{
		StatusAccepted:   "Принят",
		StatusDiagnosis:  "В диагностике",
		StatusWaiting:    "Ожидает запчастей",
		StatusInProgress: "В работе",
		StatusReady:      "Готов",
		StatusIssued:     "Выдан",
	}
	if l, ok := labels[s]; ok {
		return l
	}
	return string(s)
}

var OrderStatuses = []OrderStatus{
	StatusAccepted,
	StatusDiagnosis,
	StatusWaiting,
	StatusInProgress,
	StatusReady,
	StatusIssued,
}

type Order struct {
	ID            int64
	ClientName    string
	Phone         string
	Device        string
	Description   string
	EstimatedCost float64
	Status        OrderStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Parts         []OrderPart
}

type Part struct {
	ID            int64
	Name          string
	Quantity      int
	PurchasePrice float64
	CreatedAt     time.Time
}

type OrderPart struct {
	ID        int64
	OrderID   int64
	PartID    int64
	PartName  string
	Quantity  int
	CreatedAt time.Time
}
