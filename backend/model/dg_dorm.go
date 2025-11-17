package model

type DgDorm struct {
	DormId   uint   `gorm:"column:dormid;primarykey;AUTO_INCREMENT" json:"dormid"`
	DormName string `gorm:"column:dormname;type:varchar(50);" json:"dormname"`
	//constraint:OnUpdate:CASCADE,OnDelete:SET NULL 表示更新或删除宿舍时，对应用户的 DormId 自动更新或设为 NULL
	Users []DgUser `gorm:"foreignkey:DormId;references:DormId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
}
