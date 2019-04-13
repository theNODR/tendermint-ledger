package store

type BaseChannelState struct {
	Address		string
	Price		Price
	State		*SpecialChannelState
	TimeLock	int64
	LifeTime	int64
}
