package store

type TransferChannelState struct {
	TimeLock	int64
	LifeTime	int64
}

func NewTransferChannelState(timeLock int64, lifeTime int64) *TransferChannelState {
	return &TransferChannelState{
		TimeLock: timeLock,
		LifeTime: lifeTime,
	}
}
