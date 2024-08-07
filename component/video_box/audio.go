package video_box

/*
#include <string.h>
void audioCallback(void *userdata, char *stream, int len);
*/
import "C"
import (
	"log"
	"math"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	once sync.Once

	volumeMap = map[int]float64{
		1:  -20,
		2:  -14,
		3:  -10,
		4:  -7.5,
		5:  -6,
		6:  -4.5,
		7:  -3,
		8:  -2,
		9:  -1,
		10: 0,
	}
)

type videoBoxRef struct {
	p atomic.Pointer[VideoBox]
}

func (r *videoBoxRef) Get() *VideoBox {
	return r.p.Load()
}

func (r *videoBoxRef) Set(p *VideoBox) {
	r.p.Store(p)
}

var vf = &videoBoxRef{}

//export audioCallback
func audioCallback(userdata unsafe.Pointer, stream *C.char, len C.int) {
	C.memset(unsafe.Pointer(stream), 0, C.ulong(len))

	vf := (*videoBoxRef)(userdata)
	v := vf.Get()
	if v == nil {
		return
	}

	if v.audioVolume == 0 {
		return
	}

	data := (*[1 << 30]byte)(unsafe.Pointer(stream))[:len:len]
	_, err := v.audioBuf.Read(data)
	if err != nil {
		return
	}

	linearVolume := math.Pow(10, volumeMap[v.audioVolume]/10)

	i16data := (*[1 << 30]int16)(unsafe.Pointer(stream))[: len/2 : len/2]
	for i := range i16data {
		i16data[i] = int16(float64(i16data[i]) * linearVolume)
	}
}

func initAudio() {
	once.Do(func() {
		desired := &sdl.AudioSpec{
			Freq:     44100,
			Format:   sdl.AUDIO_S16LSB,
			Channels: 2,
			Samples:  4096,
			Callback: sdl.AudioCallback(C.audioCallback),
			UserData: unsafe.Pointer(vf),
		}
		obtained := &sdl.AudioSpec{}
		err := sdl.OpenAudio(desired, obtained)
		if err != nil {
			panic(err)
		}
		log.Println("obtained audio spec:", obtained)
		sdl.PauseAudio(false)
	})
}
