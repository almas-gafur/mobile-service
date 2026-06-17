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
	StatusNew        OrderStatus = "new"
	StatusAccepted   OrderStatus = "accepted"
	StatusDiagnosis  OrderStatus = "diagnosis"
	StatusWaiting    OrderStatus = "waiting_parts"
	StatusInProgress OrderStatus = "in_progress"
	StatusReady      OrderStatus = "ready"
	StatusIssued     OrderStatus = "issued"
	StatusRejected   OrderStatus = "rejected"
)

func (s OrderStatus) Label() string {
	labels := map[OrderStatus]string{
		StatusNew:        "Новая заявка",
		StatusAccepted:   "Принят",
		StatusDiagnosis:  "В диагностике",
		StatusWaiting:    "Ожидает запчастей",
		StatusInProgress: "В работе",
		StatusReady:      "Готов",
		StatusIssued:     "Выдан",
	StatusRejected:   "Отклонён",
	}
	if l, ok := labels[s]; ok {
		return l
	}
	return string(s)
}

var OrderStatuses = []OrderStatus{
	StatusNew,
	StatusAccepted,
	StatusDiagnosis,
	StatusWaiting,
	StatusInProgress,
	StatusReady,
	StatusIssued,
	StatusRejected,
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
