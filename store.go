package main

type PotionStore interface {
	GetAllPotions() ([]Potion, error)
}
