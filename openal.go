// OpenAL Test

package main

import (
	"fmt"
	"log"
	"math"
	"time"

	al "github.com/timshannon/go-openal/openal"
)

const (
	fps              = 30
	samplesPerSecond = 44100
	maxSampleUpdate  = samplesPerSecond / 15
)

func main() {
	// sets up OpenAL with default options
	device := al.OpenDevice("")
	defer device.CloseDevice()
	context := device.CreateContext()
	defer context.Destroy()
	context.Activate()
	vendor := al.GetVendor()

	// make sure things have gone well
	if err := al.Err(); err != nil {
		fmt.Printf("Failed to setup OpenAL: %v\n", err)
		return
	}
	fmt.Printf("OpenAL vendor: %s\n", vendor)

	source := al.NewSource()
	defer source.Pause()
	source.SetLooping(false)
	source.SetPosition(&al.Vector{0.0, 0.0, -5.0})
	source.SetGain(source.GetMaxGain())

	buffers := al.NewBuffers(4)
	if err := al.Err(); err != nil {
		fmt.Printf("OpenAL buffer creation failed: %v\n", err)
		return
	}
	freqNums := []float64 {
		freqs["e7"],
		freqs["d7"],
		freqs["f#6"],
		freqs["g#6"],
		freqs["c#7"],
	}
	//r8l16o7ed<f+8g+8>c+<bd8e8&e32ba.c+8&c+32l8.ea2&a&a32r8l16o7ed<f+8g+8>c+<bd8e8&e32ba.c+8&c+32l8.ea2&a&a32r8l16o7ed<f+8g+8>c+<bd8e8&e32ba.c+8&c+32l8.ea2&a&a32r8l16o7ed<f+8g+8>c+<bd8e8&e32ba.c+8&c+32l8.ea2&a&a32
	/*freqNums := make([]float64, 0, len(freqs))
	for _, fr := range freqs {
		freqNums = append(freqNums, fr)
	}*/
	//sort.Float64s(freqNums)
	fn := 0
	for cbuf := range buffers {
		buffers[cbuf].SetData(al.FormatMono8, note(defEnv, freqNums[fn]), samplesPerSecond)
		fn++
	}
	source.QueueBuffers(buffers)

	source.Play()

	// load a sound effect
	// https://developer.tizen.org/dev-guide/2.4/org.tizen.tutorials/html/native/multimedia/openal_tutorial_n.htm
	for {
		time.Sleep(10 * time.Millisecond)
		for source.BuffersProcessed() > 0 {
			log.Printf("processed %v\n", source.BuffersProcessed())
			b := source.UnqueueBuffer()
			b.SetData(al.FormatMono8, note(defEnv, freqNums[fn]), samplesPerSecond)
			fn = (fn + 1) % len(freqNums)
			source.QueueBuffer(b)
		}
	}
}

func note(e envelope, freq float64) []byte {
	timeMs := e.AttackTime + e.DecayTime + e.SustainTime + e.ReleaseTime
	wave := squareBytes(freq, timeMs, samplesPerSecond)
	// todo: cache
	env := e.toFloat(samplesPerSecond)
	for i := 0; i < len(env); i++ {
		//wave[i] = byte(float64(wave[i]) * env[i])
	}
	return wave
}

func squareBytes(freq, timeMs float64, samplesPerSec int) []byte {
	s := make([]byte, int(timeMs*float64(samplesPerSec)/1000))
	for i := 0; i < len(s); i++ {
		l := math.Sin(
			float64(i) * 2 * math.Pi / float64(samplesPerSec) * freq)
		s[i] = byte(math.Round(127.5+127.5*l))
		/*if l > 0 {
			s[i] = 255
		} else {
			s[i] = 0
		}*/
	}
	prev := byte(255)
	lastIdx := 0
	measuredFreq := 0
	fmt.Println("freq", freq)
	for i := 0; i < len(s); i++ {
		if s[i] != prev {
			//fmt.Printf("%d (l: %d) -> %d\n", i, i-lastIdx, s[i])
			prev = s[i]
			if prev == 255 && (i-lastIdx) > 1 {
				measuredFreq++
			}
			lastIdx = i
		}
	}
	fmt.Println("***** measured freq",
		float64(measuredFreq)*float64(samplesPerSec)/float64(len(s)))
	return s
}

type envelope struct {
	AttackTime  float64
	AttackLevel float64
	DecayTime   float64
	DecayLevel  float64
	SustainTime float64
	ReleaseTime float64
}

var defEnv = envelope{
	AttackTime:  100,
	AttackLevel: 1,
	DecayTime:   100,
	DecayLevel:  0.6,
	SustainTime: 200,
	ReleaseTime: 100,
}

func (e *envelope) toFloat(sampleRate int) []float64 {
	yVal := func(x, x1, y1, x2, y2 float64) float64 {
		return (x-x1)/(x2-x1)*(y2-y1) + y1
	}
	timeMs := e.AttackTime + e.DecayTime + e.SustainTime + e.ReleaseTime
	samples := make([]float64, int(math.Round(float64(sampleRate)*timeMs/1000)))
	// attack slope
	attackSamples := int(math.Round(e.AttackTime * float64(len(samples)) / timeMs))
	sn := 0
	for i := 0; i < attackSamples; i++ {
		x := float64(i) / float64(attackSamples) * e.AttackTime
		samples[sn] = yVal(x, 0, 0, e.AttackTime, e.AttackLevel)
		sn++
	}
	// decay slope
	decaySamples := int(math.Round(e.DecayTime) * float64(len(samples)) / timeMs)
	for i := 0; i < decaySamples; i++ {
		x := e.AttackTime + float64(i)/float64(decaySamples)*e.DecayTime
		samples[sn] = yVal(x, e.AttackTime, e.AttackLevel, e.AttackTime+e.DecayTime, e.DecayLevel)
		sn++
	}
	// sustain
	sustainSamples := int(math.Round(e.SustainTime * float64(len(samples)) / timeMs))
	for i := 0; i < sustainSamples; i++ {
		samples[sn] = e.DecayLevel
		sn++
	}
	// release
	releaseSamples := int(math.Round(e.ReleaseTime * float64(len(samples)) / timeMs))
	for i := 1; i < releaseSamples; i++ {
		x := e.AttackTime + e.DecayTime + e.SustainTime +
			float64(i)/float64(releaseSamples)*e.ReleaseTime
		// subtract 1 to make sure the very last is 0
		samples[sn-1] = yVal(x,
			e.SustainTime+e.DecayTime+e.AttackTime, e.DecayLevel,
			e.ReleaseTime+e.SustainTime+e.DecayTime+e.AttackTime, 0)
		sn++
	}
	return samples
}

var freqs = map[string]float64{
	"c0":  16.35,
	"c#0": 17.32,
	"d0":  18.35,
	"d#0": 19.45,
	"e0":  20.60,
	"f0":  21.83,
	"f#0": 23.12,
	"g0":  24.50,
	"g#0": 25.96,
	"a0":  27.50,
	"a#0": 29.14,
	"b0":  30.87,
	"c1":  32.70,
	"c#1": 34.65,
	"d1":  36.71,
	"d#1": 38.89,
	"e1":  41.20,
	"f1":  43.65,
	"f#1": 46.25,
	"g1":  49.00,
	"g#1": 51.91,
	"a1":  55.00,
	"a#1": 58.27,
	"b1":  61.74,
	"c2":  65.41,
	"c#2": 69.30,
	"d2":  73.42,
	"d#2": 77.78,
	"e2":  82.41,
	"f2":  87.31,
	"f#2": 92.50,
	"g2":  98.00,
	"g#2": 103.83,
	"a2":  110.00,
	"a#2": 116.54,
	"b2":  123.47,
	"c3":  130.81,
	"c#3": 138.59,
	"d3":  146.83,
	"d#3": 155.56,
	"e3":  164.81,
	"f3":  174.61,
	"f#3": 185.00,
	"g3":  196.00,
	"g#3": 207.65,
	"a3":  220.00,
	"a#3": 233.08,
	"b3":  246.94,
	"c4":  261.63,
	"c#4": 277.18,
	"d4":  293.66,
	"d#4": 311.13,
	"e4":  329.63,
	"f4":  349.23,
	"f#4": 369.99,
	"g4":  392.00,
	"g#4": 415.30,
	"a4":  440.00,
	"a#4": 466.16,
	"b4":  493.88,
	"c5":  523.25,
	"c#5": 554.37,
	"d5":  587.33,
	"d#5": 622.25,
	"e5":  659.25,
	"f5":  698.46,
	"f#5": 739.99,
	"g5":  783.99,
	"g#5": 830.61,
	"a5":  880.00,
	"a#5": 932.33,
	"b5":  987.77,
	"c6":  1046.50,
	"c#6": 1108.73,
	"d6":  1174.66,
	"d#6": 1244.51,
	"e6":  1318.51,
	"f6":  1396.91,
	"f#6": 1479.98,
	"g6":  1567.98,
	"g#6": 1661.22,
	"a6":  1760.00,
	"a#6": 1864.66,
	"b6":  1975.53,
	"c7":  2093.00,
	"c#7": 2217.46,
	"d7":  2349.32,
	"d#7": 2489.02,
	"e7":  2637.02,
	"f7":  2793.83,
	"f#7": 2959.96,
	"g7":  3135.96,
	"g#7": 3322.44,
	"a7":  3520.00,
	"a#7": 3729.31,
	"b7":  3951.07,
	"c8":  4186.01,
	"c#8": 4434.92,
	"d8":  4698.63,
	"d#8": 4978.03,
	"e8":  5274.04,
	"f8":  5587.65,
	"f#8": 5919.91,
	"g8":  6271.93,
	"g#8": 6644.88,
	"a8":  7040.00,
	"a#8": 7458.62,
	"b8":  7902.13,
}
