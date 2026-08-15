package dice_animation

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/pomquest/internal/adapters/icon"
)

func TestNewAnimatorDefaults(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	anim := NewAnimator(6, rng)

	if anim.FPS != 60 {
		t.Errorf("expected DefaultFPS 60, got %d", anim.FPS)
	}
	if anim.DecelFactor != 3.0 {
		t.Errorf("expected DefaultDecelFactor 3.0, got %v", anim.DecelFactor)
	}
	if anim.BaseDuration != 1500*time.Millisecond {
		t.Errorf("expected BaseDuration 1.5s, got %v", anim.BaseDuration)
	}
	if anim.VarianceDuration != 1000*time.Millisecond {
		t.Errorf("expected VarianceDuration 1.0s, got %v", anim.VarianceDuration)
	}
	if anim.MaxSide != 6 {
		t.Errorf("expected MaxSide 6, got %d", anim.MaxSide)
	}
	if anim.Target != 6 {
		t.Errorf("expected Target 6, got %d", anim.Target)
	}
	if anim.TotalDuration < 500*time.Millisecond || anim.TotalDuration > 2500*time.Millisecond {
		t.Errorf("expected TotalDuration between 0.5s and 2.5s, got %v", anim.TotalDuration)
	}
}

func TestSetFPS(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	anim := NewAnimator(6, rng)

	anim.SetFPS(20)
	if anim.FrameInterval() != 50*time.Millisecond {
		t.Errorf("expected FrameInterval 50ms for 20 FPS, got %v", anim.FrameInterval())
	}

	la := NewLineAnimator("Arsel rolls 󰇏", []int{6}, 0, rng)
	la.SetFPS(10)
	if la.FrameInterval() != 100*time.Millisecond {
		t.Errorf("expected LineAnimator FrameInterval 100ms for 10 FPS, got %v", la.FrameInterval())
	}
}

func TestAnimatorTickAndCompletion(t *testing.T) {
	rng := rand.New(rand.NewSource(100))
	anim := NewAnimator(5, rng)
	anim.TotalDuration = 1000 * time.Millisecond

	if anim.IsFinished() {
		t.Errorf("expected newly created animator not to be finished")
	}

	// Advance time partially
	anim.Tick(500 * time.Millisecond)
	if anim.IsFinished() {
		t.Errorf("expected animator not to be finished at 500ms")
	}
	if anim.CurrentValue < 1 || anim.CurrentValue > anim.MaxSide {
		t.Errorf("expected CurrentValue between 1 and %d during shuffle, got %d", anim.MaxSide, anim.CurrentValue)
	}

	// Advance time past total duration
	anim.Tick(600 * time.Millisecond)
	if !anim.IsFinished() {
		t.Errorf("expected animator to be finished after total duration")
	}
	if anim.CurrentValue != 5 {
		t.Errorf("expected final CurrentValue to be target 5, got %d", anim.CurrentValue)
	}
	expectedIcon := icon.GetD6Icon(5, icon.ModeNerdFont)
	if anim.CurrentIcon() != expectedIcon {
		t.Errorf("expected CurrentIcon to be %s, got %s", expectedIcon, anim.CurrentIcon())
	}
}

func TestLineAnimator(t *testing.T) {
	rng := rand.New(rand.NewSource(777))
	logLine := " Arsel rolls 󰇏 + 󰇎 + 3 = 14 󰓥"
	targets := []int{6, 5}

	la := NewLineAnimator(logLine, targets, 2, rng)
	if la.TargetLineIdx != 2 {
		t.Errorf("expected TargetLineIdx 2, got %d", la.TargetLineIdx)
	}

	for _, anim := range la.Animators {
		anim.TotalDuration = 500 * time.Millisecond
	}

	if la.Finished {
		t.Errorf("expected LineAnimator not to be finished initially")
	}

	// Verify that before completion, nothing after the '=' sign is printed
	initialRender := la.Render()
	if !strings.HasSuffix(initialRender, "=") {
		t.Errorf("expected initial render to end at '=', got: %q", initialRender)
	}
	if strings.Contains(initialRender, "14") {
		t.Errorf("expected nothing after '=' to be printed during animation, got: %q", initialRender)
	}

	// Tick until completed
	la.Tick(600 * time.Millisecond)
	if !la.Finished {
		t.Errorf("expected LineAnimator to be finished")
	}

	rendered := la.Render()
	expectedTargetIcons := icon.GetD6Icon(6, icon.ModeNerdFont) + " + " + icon.GetD6Icon(5, icon.ModeNerdFont)
	if !strings.Contains(rendered, expectedTargetIcons) {
		t.Errorf("expected rendered line to contain target icons %q, got %q", expectedTargetIcons, rendered)
	}
	if !strings.Contains(rendered, "= 14") {
		t.Errorf("expected completed line to contain text after '=', got %q", rendered)
	}

}

