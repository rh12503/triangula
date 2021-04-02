// +build !sdl2

package draw

import (
	"github.com/gonutz/mixer"
	"github.com/gonutz/mixer/wav"
)

var wavTable = make(map[string]mixer.SoundSource)

func initSound() error {
	return mixer.Init()
}

func closeSound() {
	mixer.Close()
}

func playSoundFile(path string) error {
	if sound, ok := wavTable[path]; ok {
		sound.PlayOnce()
		return nil
	}

	wave, err := wav.LoadFromFile(path)
	if err != nil {
		return err
	}

	sound, err := mixer.NewSoundSource(wave)
	if err != nil {
		return err
	}
	wavTable[path] = sound
	sound.PlayOnce()

	return nil
}
