package dice_animation

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/pomquest/internal/adapters/icon"
)

// Default configuration variables as requested in task 8.3.1
var (
	DefaultFPS              = 60
	DefaultDecelFactor      = 3.0
	DefaultBaseDuration     = 1500 * time.Millisecond
	DefaultVarianceDuration = 1000 * time.Millisecond
	DefaultMaxSide          = 6
)


// Animator manages the frame animation for a single die with physical deceleration.
type Animator struct {
	FPS              int
	DecelFactor      float64
	BaseDuration     time.Duration
	VarianceDuration time.Duration
	MaxSide          int
	Target           int
	CurrentValue     int
	TotalDuration    time.Duration
	ElapsedTime      time.Duration
	TimeSinceChange  time.Duration
	Finished         bool
	Mode             icon.Mode
	rng              *rand.Rand
}

// NewAnimator constructs a die animator targeting a specific final face value.
func NewAnimator(target int, rng *rand.Rand) *Animator {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	maxSide := DefaultMaxSide
	if target > maxSide {
		maxSide = target
	}
	if maxSide < 1 {
		maxSide = 1
	}

	// Variance ranges from -VarianceDuration to +VarianceDuration (-1.0s to +1.0s)
	// Output duration: between 0.5s and 2.5s (BaseDuration = 1.5s)
	varFactor := (rng.Float64() * 2.0) - 1.0
	varianceOffset := time.Duration(varFactor * float64(DefaultVarianceDuration))
	totalDur := DefaultBaseDuration + varianceOffset
	if totalDur < 500*time.Millisecond {
		totalDur = 500 * time.Millisecond
	}
	if totalDur > 2500*time.Millisecond {
		totalDur = 2500 * time.Millisecond
	}

	initialVal := rng.Intn(maxSide) + 1

	return &Animator{
		FPS:              DefaultFPS,
		DecelFactor:      DefaultDecelFactor,
		BaseDuration:     DefaultBaseDuration,
		VarianceDuration: DefaultVarianceDuration,
		MaxSide:          maxSide,
		Target:           target,
		CurrentValue:     initialVal,
		TotalDuration:    totalDur,
		ElapsedTime:      0,
		TimeSinceChange:  0,
		Finished:         false,
		Mode:             icon.ModeNerdFont,
		rng:              rng,
	}
}

// SetFPS sets the animation rate by target FPS.
func (a *Animator) SetFPS(fps int) {
	if fps <= 0 {
		fps = 60
	}
	a.FPS = fps
}

// FrameInterval returns the current duration per frame interval based on FPS.
func (a *Animator) FrameInterval() time.Duration {
	if a.FPS <= 0 {
		return time.Second / 60
	}
	return time.Second / time.Duration(a.FPS)
}

// Step advances the animation state by one frame interval.
func (a *Animator) Step() {
	a.Tick(a.FrameInterval())
}

// Tick advances the animation state by the specified delta time.
func (a *Animator) Tick(delta time.Duration) {
	if a.Finished {
		return
	}

	a.ElapsedTime += delta
	a.TimeSinceChange += delta

	if a.ElapsedTime >= a.TotalDuration {
		a.Finished = true
		a.CurrentValue = a.Target
		return
	}

	// Physics-based deceleration easing curve for weighty die roll feel:
	// As progress goes from 0 to 1, friction slows down face changes based on DecelFactor.
	progress := float64(a.ElapsedTime) / float64(a.TotalDuration)
	baseInterval := a.FrameInterval()
	decelFactor := 1.0 + a.DecelFactor*(progress*progress)
	currentInterval := time.Duration(float64(baseInterval) * decelFactor)

	if a.TimeSinceChange >= currentInterval {
		a.TimeSinceChange = 0
		nextVal := a.rng.Intn(a.MaxSide) + 1
		if a.MaxSide > 1 && nextVal == a.CurrentValue {
			nextVal = (nextVal % a.MaxSide) + 1
		}
		a.CurrentValue = nextVal
	}
}

// IsFinished returns whether the die animation has landed on the target value.
func (a *Animator) IsFinished() bool {
	return a.Finished
}

// CurrentIcon returns the icon string corresponding to the current face value.
func (a *Animator) CurrentIcon() string {
	return icon.GetD6Icon(a.CurrentValue, a.Mode)
}

// TargetIcon returns the icon string corresponding to the target face value.
func (a *Animator) TargetIcon() string {
	return icon.GetD6Icon(a.Target, a.Mode)
}

// LineAnimator animates multiple die icons embedded within a specific quest log line.
type LineAnimator struct {
	OriginalLine  string
	Animators     []*Animator
	Finished      bool
	TargetLineIdx int
}

var d6IconRegex = regexp.MustCompile(`[󰇋󰇌󰇍󰇎󰇏⚀⚁⚂⚃⚄⚅]`)

// NewLineAnimator parses a log line and creates animators for each target die value.
func NewLineAnimator(line string, diceTargets []int, targetLineIdx int, rng *rand.Rand) *LineAnimator {
	var animators []*Animator
	for _, target := range diceTargets {
		animators = append(animators, NewAnimator(target, rng))
	}

	return &LineAnimator{
		OriginalLine:  line,
		Animators:     animators,
		Finished:      len(animators) == 0,
		TargetLineIdx: targetLineIdx,
	}
}

// SetFPS sets the animation rate by target FPS for all die animators in the line.
func (la *LineAnimator) SetFPS(fps int) {
	for _, anim := range la.Animators {
		anim.SetFPS(fps)
	}
}

// FrameInterval returns the frame interval of the first animator, or default interval.
func (la *LineAnimator) FrameInterval() time.Duration {
	if len(la.Animators) > 0 {
		return la.Animators[0].FrameInterval()
	}
	return time.Second / time.Duration(DefaultFPS)
}

// IsFinished returns whether all die animators in the line have finished.
func (la *LineAnimator) IsFinished() bool {
	return la.Finished
}

// Finish immediately completes all die animators in the line.
func (la *LineAnimator) Finish() {
	for _, anim := range la.Animators {
		anim.Tick(anim.TotalDuration)
	}
	la.Finished = true
}

// Tick advances all die animators for this line by delta time.
func (la *LineAnimator) Tick(delta time.Duration) {
	if la.Finished {
		return
	}

	allFinished := true
	for _, anim := range la.Animators {
		anim.Tick(delta)
		if !anim.IsFinished() {
			allFinished = false
		}
	}
	la.Finished = allFinished
}

// Render returns the log line with die icons replaced by their current animated frame icons.
// While animation is in progress (!la.Finished), nothing after the '=' sign is printed until completion.
func (la *LineAnimator) Render() string {
	if len(la.Animators) == 0 {
		return la.OriginalLine
	}

	targetLine := la.OriginalLine
	if !la.Finished {
		if eqIdx := strings.Index(la.OriginalLine, "="); eqIdx != -1 {
			targetLine = la.OriginalLine[:eqIdx+1]
		} else {
			matches := d6IconRegex.FindAllStringIndex(la.OriginalLine, -1)
			numAnim := len(la.Animators)
			if numAnim > 0 && len(matches) >= numAnim {
				lastEnd := matches[numAnim-1][1]
				targetLine = la.OriginalLine[:lastEnd]
			}
		}
	}

	animIdx := 0
	return d6IconRegex.ReplaceAllStringFunc(targetLine, func(match string) string {
		if animIdx < len(la.Animators) {
			iconStr := la.Animators[animIdx].CurrentIcon()
			animIdx++
			return iconStr
		}
		return match
	})
}


