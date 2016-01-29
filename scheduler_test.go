package main

import (
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Scheduler(t *testing.T) {
	sched := new(scheduler)
	eslog.InitSilent()
	info := &config.Query{
		Schedule:       "500s",
		Alert_onlyonce: true,
		Alert_endmsg:   true,
	}

	sched.initScheduler(info)
	assert.Equal(t, true, sched.isAlertOnlyOnce)
	assert.Equal(t, true, sched.isAlertEndMsg)
	assert.Equal(t, "8m20s", sched.waitSchedule.String())

	info = &config.Query{
		Schedule:       "30m",
		Alert_onlyonce: false,
		Alert_endmsg:   false,
	}
	sched.initScheduler(info)
	assert.Equal(t, false, sched.isAlertOnlyOnce)
	assert.Equal(t, false, sched.isAlertEndMsg)
	assert.Equal(t, "30m0s", sched.waitSchedule.String())

	info = &config.Query{
		Alert_onlyonce: false,
	}
	sched.initScheduler(info)
	assert.Equal(t, true, sched.isAlertOnlyOnce)
	assert.Equal(t, false, sched.isAlertEndMsg)
	assert.Equal(t, "10m0s", sched.waitSchedule.String())
	assert.Equal(t, "10m0s", sched.alertSchedule.String())

	info = &config.Query{}
	sched.initScheduler(info)
	assert.Equal(t, true, sched.isAlertOnlyOnce)
	assert.Equal(t, false, sched.isAlertEndMsg)
	assert.Equal(t, "10m0s", sched.waitSchedule.String())
	assert.Equal(t, "10m0s", sched.alertSchedule.String())

	info = &config.Query{
		Schedule:       "pouet",
		Alert_onlyonce: false,
	}
	sched.initScheduler(info)
	assert.Equal(t, true, sched.isAlertOnlyOnce)
	assert.Equal(t, false, sched.isAlertEndMsg)
	assert.Equal(t, "10m0s", sched.waitSchedule.String())
	assert.Equal(t, "10m0s", sched.alertSchedule.String())

	info = &config.Query{
		Schedule:       "40z",
		Alert_onlyonce: false,
	}
	assert.Equal(t, true, sched.isAlertOnlyOnce)
	assert.Equal(t, false, sched.isAlertEndMsg)
	assert.Equal(t, "10m0s", sched.waitSchedule.String())
	assert.Equal(t, "10m0s", sched.alertSchedule.String())
	sched.initScheduler(info)
}
