package model

type dgType struct {
	TypeId   uint   `gorm:"column:typeid;primarykey;AUTO_INCREMENT" json:"typeid"`
	TypeName string `gorm:"column:typename;type:varchar(255);" json:"typename"`
}
