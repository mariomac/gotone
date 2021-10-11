package main

/*
const (
	fps = 30
	samplesPerSecond = 22050
	maxSampleUpdate = samplesPerSecond / 15
)

func init() {
	runtime.LockOSThread()
}

func main2() {
	rl.InitWindow(800, 450, "raylib [audio] example - raw audio streaming")

	rl.InitAudioDevice()

	// Init raw audio stream (sample rate: 22050, sample size: 32bit-float, channels: 1-mono)
	stream := rl.InitAudioStream(samplesPerSecond, 32, 1)
	musicBuffer := note(defEnv)

	// NOTE: The generated MAX_SAMPLES do not fit to close a perfect loop
	// for that reason, there is a clip everytime audio stream is looped
	rl.PlayAudioStream(stream)


    block := sync.Mutex{}
	go func() {
		for {
			block.Lock()
			if len(musicBuffer) < maxSampleUpdate {
				musicBuffer = append(musicBuffer, note(defEnv)...)
			}
			block.Unlock()
		}
	}()

	rl.SetTargetFPS(fps)

	for !rl.WindowShouldClose() {
		// Refill audio stream if required
		if rl.IsAudioStreamProcessed(stream) {
			block.Lock()
			bl := int(math.Min(maxSampleUpdate, float64(len(musicBuffer))))
			var send []float32
			send, musicBuffer = musicBuffer[:bl], musicBuffer[bl:]
			block.Unlock()
			rl.UpdateAudioStream(stream, send, int32(bl))
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("SINE WAVE SHOULD BE PLAYING!", 240, 140, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseAudioStream(stream) // Close raw audio stream and delete buffers from RAM

	rl.CloseAudioDevice() // Close audio device (music streaming is automatically stopped)

	rl.CloseWindow()
}




func squareWave(wl int, samples int) []float32 {
	s := make([]float32, samples)
	x := 0
	for {
		for i := 0 ; i < wl / 2 ; i++ {
			if x >= samples {
				return s
			}
			s[x] = 1
			x++
		}
		for i := wl / 2 ; i < wl ; i++ {
			if x >= samples {
				return s
			}
			s[x] = 0
			x++
		}
	}
}
*/
