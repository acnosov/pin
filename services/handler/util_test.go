package handler

import (
	"context"
	"github.com/aibotsoft/micro/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler_CalcSpread(t *testing.T) {
	got, err := h.CalcHandicapMiddleMargin(context.Background(), "1127970328", 1)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
	}
}

func TestHandler_CalcSpread_System(t *testing.T) {
	var sum, count, max float64
	var min float64 = 100
	// Corners, Regular, Booking
	// 29 soccer, 4 Basketball
	events, err2 := h.store.SelectCurrentEvents(context.Background(), 4, 200, "Regular")
	if assert.NoError(t, err2) {
		for i := range events {
			//t.Log(events[i])
			got, err := h.CalcHandicapMiddleMargin(context.Background(), events[i].Id, 3)
			if assert.NoError(t, err) {
				//assert.NotEmpty(t, got)
				if got > 0 {
					sum += got
					count += 1
					//t.Log(got)
				}
				if got > 0 && got < min {
					min = got
				}
				if got > max {
					max = got
				}
			}
		}
		if count > 0 {
			h.log.Infow("result", "avg:", util.TruncateFloat(sum/count, 2), "count:", count, "min", min, "max", max)
		}
	}
}
func TestHandler_CalcTotalMiddleMargin(t *testing.T) {
	got, err := h.CalcTotalMiddleMargin(context.Background(), "1127970329", 0)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
	}
}

func TestHandler_CalcTotalMiddleMargin_System(t *testing.T) {
	var sum, count float64
	// Corners, Regular, Booking
	events, err2 := h.store.SelectCurrentEvents(context.Background(), 29, 300, "Regular")
	if assert.NoError(t, err2) {
		//t.Log(events)
		for i := range events {
			//t.Log(events[i])
			got, err := h.CalcTotalMiddleMargin(context.Background(), events[i].Id, 0)
			if assert.NoError(t, err) {
				//assert.NotEmpty(t, got)
				if got > 0 {
					sum += got
					count += 1
					//t.Log(got)
				}
			}
		}
		t.Log("avg:", sum/count, "count:", count)
	}
}
