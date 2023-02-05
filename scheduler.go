package main

import (
	"log"
	"time"
	// "fmt"
)
type CallbackOnSubject func(*SubjectRunner)


type SubjectStatus int

const (
        Ready SubjectStatus = iota
        Running
        Done
        Skipped
)


type SubjectRunner struct {
	Subject Subject

	// FIXME: paused?
	LastStarted time.Time
 	LastRunTime time.Duration
 	RunTime time.Duration

 	Status SubjectStatus
}

func (sr *SubjectRunner) Start() {
	sr.LastStarted = time.Now()
	sr.Status = Running
}


func (sr *SubjectRunner) UpdateRunTime() {
	if sr.LastStarted.IsZero() {
		return
	}
	newLastRunTime := time.Now().Sub(sr.LastStarted)	
	sr.RunTime = (sr.RunTime - sr.LastRunTime) + newLastRunTime
	sr.LastRunTime = newLastRunTime
}


func (sr *SubjectRunner) Stop() {
	sr.UpdateRunTime()
	sr.LastRunTime = 0
	sr.Status = Done
}


func (sr *SubjectRunner) IsTimeOver() bool {
	return sr.Subject.Duration <= sr.RunTime
}


type ScheduleRunner struct {
	Subjects []SubjectRunner
	CurrentSubject *SubjectRunner
	// CurrentSubjectIdx int

	Running bool   // FIXME: not needed?
	// Timer 
	// OnNextSubject
	// OnSubjectEnd
	CallbackOnTick CallbackOnSubject // FIXME: use channel
	CallbackOnSubjectFinish CallbackOnSubject // FIXME: use channel

	done chan bool
}


func NewScheduleRunner(schedule *Schedule) *ScheduleRunner {
	var subjects []SubjectRunner
	for i := 0; i < len(schedule.Subjects); i++ {
		subjects = append(subjects, SubjectRunner{
			schedule.Subjects[i],
			time.Time{},
			0,
			0,
			Ready,
		})
	}
	return &ScheduleRunner{subjects, nil, false, nil, nil,  make(chan bool)}
}


// FIXME: use channels
func (s *ScheduleRunner) SetCallbackOnTick(cb CallbackOnSubject) {
	s.CallbackOnTick = cb
}


func (s *ScheduleRunner) SetCallbackOnSubjectFinish(cb CallbackOnSubject) {
	s.CallbackOnSubjectFinish = cb
}
// , onSubjectFinish CallbackOnSubject, callbackOnTick CallbackOnSubject


func (s *ScheduleRunner) Start() {
	if s.Running {
		return
	}
	s.Running = true
	go s.runTicker()
}


func (s *ScheduleRunner) Stop() {
	if !s.Running {
		return
	}
	s.done <- true
}

func (s *ScheduleRunner) runTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case t := <-ticker.C:
			s.OnTick(t)
		}
	}
}


func (s *ScheduleRunner) StartSubject(idx int) {
	if (idx < 0 || idx >= len(s.Subjects)) {
		log.Printf("Invalid idx: %d; Total subjects: %d", idx, len(s.Subjects))
		return
	}
	if s.CurrentSubject != nil {
		s.CurrentSubject.Stop()	
	}
	
	
	s.CurrentSubject = &s.Subjects[idx]
	s.CurrentSubject.Start()
}


func (s *ScheduleRunner) OnTick(current time.Time) {
	subj := s.CurrentSubject
	if subj == nil {
		return
	}
	prevStatus := subj.IsTimeOver()
	subj.UpdateRunTime()
	if subj.IsTimeOver() && !prevStatus && s.CallbackOnSubjectFinish != nil {
		s.CallbackOnSubjectFinish(subj)
	}

	if s.CallbackOnTick != nil {
		s.CallbackOnTick(subj)
	}
}


// func (s *ScheduleRunner) NextSubject() {
// 	s.CurrentSubjectIdx += 1
// 	if s.IsCompleted() {
// 		return
// 	}
// 	s.Subjects[s.CurrentSubjectIdx].Start()
// }


// func (s *ScheduleRunner) CurrentSubject() *SubjectRunner {
// 	return &s.Subjects[s.CurrentSubjectIdx]
// }