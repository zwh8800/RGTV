package component

// void audioCallback(void *userdata, char *stream, int len);
import "C"
import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"unsafe"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/veandco/go-sdl2/sdl"
)

//export audioCallback
func audioCallback(userdata unsafe.Pointer, stream *C.char, len C.int) {
	v := (*VideoBox)(userdata)
	if v.rawAudioStream == nil {
		return
	}
	buf := make([]byte, len)
	n, err := io.ReadFull(v.rawAudioStream, buf)
	if n == 0 || err == io.EOF {
		return
	} else if n != int(len) || err != nil {
		panic(fmt.Sprintf("read error: %d, %s", n, err))
	}

	data := (*[1 << 30]byte)(unsafe.Pointer(stream))[:len:len]
	copy(data, buf)
}

type VideoBox struct {
	rawVideoStream *os.File
	rawAudioStream *os.File

	once    sync.Once
	texture *sdl.Texture

	pinner runtime.Pinner
}

func NewVideoBox(url string) (*VideoBox, error) {
	v := &VideoBox{}
	v.pinner.Pin(v)

	desired := &sdl.AudioSpec{
		Freq:     44100,
		Format:   sdl.AUDIO_S16LSB,
		Channels: 2,
		Samples:  4096,
		Callback: sdl.AudioCallback(C.audioCallback),
		UserData: unsafe.Pointer(v),
	}
	obtained := &sdl.AudioSpec{}
	err := sdl.OpenAudio(desired, obtained)
	if err != nil {
		return nil, err
	}
	log.Println("obtained audio spec:", obtained)
	sdl.PauseAudio(false)

	pr1, pw1, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	v.rawVideoStream = pr1
	pr2, pw2, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	v.rawAudioStream = pr2

	go func() {
		defer pw1.Close()
		defer pw2.Close()

		i := ffmpeg.Input(url, ffmpeg.KwArgs{
			"re":      "",
			"headers": "User-Agent: Lavf",
		})

		out1 := i.Get("v").
			Filter("scale", ffmpeg.Args{"640:480:force_original_aspect_ratio=decrease"}).
			Filter("pad", ffmpeg.Args{"640:480:(ow-iw)/2:(oh-ih)/2"}).
			Output("pipe:3",
				ffmpeg.KwArgs{
					"format": "rawvideo", "pix_fmt": "rgb24",
				})
		out2 := i.Get("a").
			Output("pipe:4",
				ffmpeg.KwArgs{
					"format": "s16le",
					"ar":     44100,
					"ac":     2,
				})

		cmd := ffmpeg.MergeOutputs(out1, out2).
			WithOutput(pw1, pw2).
			ErrorToStdOut().
			//SetFfmpegPath("/root/code/go/rgbili/ffmpeg").
			Compile()

		cmd.ExtraFiles = []*os.File{
			pw1, pw2,
		}

		err = cmd.Run()
		if err != nil && !errors.As(err, new(*exec.ExitError)) {
			log.Printf("ffmpeg error: %s", err)
		}
	}()

	return v, nil
}

func (v *VideoBox) HandleEvent(e sdl.Event) {

}

func (v *VideoBox) Draw(renderer *sdl.Renderer) {
	v.once.Do(func() {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, 640, 480)
		if err != nil {
			panic(err)
		}
		v.texture = texture
	})

	frameSize := 640 * 480 * 3
	buf := make([]byte, frameSize)
	n, err := io.ReadFull(v.rawVideoStream, buf)
	if n == 0 || err == io.EOF {
		return
	} else if n != frameSize || err != nil {
		panic(fmt.Sprintf("read error: %d, %s", n, err))
	}

	err = v.texture.Update(nil, unsafe.Pointer(&buf[0]), 640*3)
	if err != nil {
		panic(err)
	}
	err = renderer.Copy(v.texture, nil, nil)
	if err != nil {
		panic(err)
	}
}

func (v *VideoBox) Dispose() {
	v.pinner.Unpin()
	sdl.CloseAudio()
	v.rawVideoStream.Close()
	v.rawAudioStream.Close()
	v.texture.Destroy()
}
