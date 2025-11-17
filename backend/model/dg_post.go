package model

type DgPost struct {
	ID          uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	PublisherId uint   `gorm:"column:publisherid" json:"publisherid"`
	AccepterId  uint   `gorm:"column:accepterid" json:"accepterid"`
	DormId      uint   `gorm:"column:dormid" json:"dormid"`
	Dorm        DgDorm `gorm:"foreignkey:DormId;references:DormId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Title       string `gorm:"column:title;type:varchar(255);" json:"title"`
	Content     string `gorm:"column:content;type:text;" json:"content"`
	TypeId      uint   `gorm:"column:typeid" json:"typeid"`
	Type        DgType `gorm:"foreignkey:TypeId;references:TypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Status      string `gorm:"column:status;type:varchar(50);" json:"status"`
	CreatedAt   int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at" json:"updated_at"`
}
