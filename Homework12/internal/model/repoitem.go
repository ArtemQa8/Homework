package model

// StorageObject – интерфейс для всех сущностей, которые можно
// сохранять в репозитории (Game, Player, Move).
type StorageObject interface {
	ObjectType() string
}
