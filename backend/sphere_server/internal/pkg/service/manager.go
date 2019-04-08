package service

type Manager struct {
	cacheService CacheService
	services  map[string]CoinService
}

var (
	CoinTypeBTC string = "BTC"
)

func NewManager() *Manager {
	services:= make(map[string]CoinService)
	return &Manager{services: services}
}

func (m *Manager) SetCacheService(service CacheService) *Manager {
	m.cacheService = service
	return m
}

func (m *Manager) GetCacheService() CacheService {
	return m.cacheService
}


func (m *Manager) SetCoinService(coinType string, service CoinService) *Manager {
	m.services[coinType] = service
	return m
}

func (m *Manager) GetCoinService(coinType string) CoinService {
	return m.services[coinType]
}