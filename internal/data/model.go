package data

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	UpdatedAt int64
	DeletedAt *int64 `gorm:"index"`
}