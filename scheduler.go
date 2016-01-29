package main

import (
	"errors"
	"github.com/amundi/escheck.v2/config"
	"time"
)

type scheduler struct {
	isAlertOnlyOnce bool
	isAlertEndMsg   bool
	alertState      bool
	alertSchedule   time.Duration
	waitSchedule    time.Duration
}

func (s *scheduler) initScheduler(info *config.Query) error {
	var err error

	if info == nil {
		s.initSchedulerDefault()
		return errors.New("Error while parsing scheduler, request will have default values")
	}
	s.waitSchedule, err = time.ParseDuration(info.Schedule)
	if err != nil {
		s.initSchedulerDefault()
		return err
	}
	s.isAlertOnlyOnce = info.Alert_onlyonce
	s.isAlertEndMsg = info.Alert_endmsg
	s.alertState = false
	return nil
}

func (s *scheduler) initSchedulerDefault() {
	s.isAlertOnlyOnce = true
	s.alertState = false
	s.alertSchedule = 10 * time.Minute
	s.waitSchedule = 10 * time.Minute
	s.isAlertEndMsg = false
}

func (s *scheduler) wait() {
	time.Sleep(s.waitSchedule)
}
