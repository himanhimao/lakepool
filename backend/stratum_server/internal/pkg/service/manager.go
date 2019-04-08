package service

import "errors"

type Manager struct {
	userService   UserService
	sphereService StratumService
	statsService  StatsService
	logService    LogService
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) GetUserService() UserService {
	return m.userService
}

func (m *Manager) GetSphereService() StratumService {
	return m.sphereService
}

func (m *Manager) GetLogService() LogService {
	return m.logService
}

func (m *Manager) GetStatsService() StatsService {
	return m.statsService
}

func (m *Manager) SetUserService(service UserService) *Manager {
	m.userService = service
	return m
}

func (m *Manager) SetSphereService(service StratumService) *Manager {
	m.sphereService = service
	return m
}

func (m *Manager) SetLogService(service LogService) *Manager {
	m.logService = service
	return m
}

func (m *Manager) SetStatsService(service StatsService) *Manager {
	m.statsService = service
	return m
}

func (m *Manager) Init() error {
	var err error
	if m.statsService == nil {
		return errors.New("service manager not found stats service")
	}

	if m.sphereService == nil {
		return errors.New("service manager not found sphere service")
	}

	if m.userService == nil {
		return errors.New("service manager not found user service")
	}

	if err = m.statsService.Init(); err != nil {
		return err
	}

	if err = m.sphereService.Init(); err != nil {
		return err
	}

	if err = m.userService.Init(); err != nil {
		return err
	}
	return nil
}
