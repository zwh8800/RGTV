package video_box

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	evbus "github.com/asaskevich/EventBus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/model"
	"github.com/zwh8800/RGTV/util"
)

const (
	videoBufSize = 640 * 480 * 3
	audioBufSize = 4096 * 2 * 4 // 4倍sample buffer

	eventLag = "VideoBox:Lag"
)

type VideoBox struct {
	url string

	cmd *exec.Cmd

	rawVideoStream *os.File
	rawAudioStream *os.File

	once    sync.Once
	texture *sdl.Texture

	videoBuf    *[2][videoBufSize]byte
	videoBufIdx atomic.Int32

	audioBuf *util.RingBuffer

	audioVolume int

	lastFrame time.Time
	stopChan  chan struct{}

	eventBus evbus.Bus

	pinner runtime.Pinner
}

func New(url string) (*VideoBox, error) {
	initAudio()

	v := &VideoBox{
		url:         url,
		audioVolume: 5,
		eventBus:    evbus.New(),
		stopChan:    make(chan struct{}),
	}

	v.pinner.Pin(v)
	v.videoBuf = new([2][videoBufSize]byte)
	v.pinner.Pin(v.videoBuf)
	v.audioBuf = util.NewLimitedBuffer(2 * audioBufSize)

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

	go v.runFFMPEG(pw1, pw2)
	go v.asyncReadVideo()
	go v.asyncReadAudio()
	go v.monitorLag()

	vf.Set(v)
	return v, nil
}

func (v *VideoBox) runFFMPEG(pw1, pw2 *os.File) {
	defer pw1.Close()
	defer pw2.Close()

	i := ffmpeg.Input(v.url, ffmpeg.KwArgs{
		"re":      "",
		"headers": "User-Agent: Lavf",
	})

	out1 := i.Get("v").
		Filter("scale", ffmpeg.Args{"640:480:force_original_aspect_ratio=decrease"}).
		Filter("pad", ffmpeg.Args{"640:480:(ow-iw)/2:(oh-ih)/2"}).
		Filter("fps", ffmpeg.Args{"30"}).
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

	v.cmd = ffmpeg.MergeOutputs(out1, out2).
		ErrorToStdOut().
		SetFfmpegPath(conf.GetConfig().FFMPEGPath).
		Compile()

	v.cmd.ExtraFiles = []*os.File{
		pw1, pw2,
	}

	err := v.cmd.Run()
	if err != nil && !errors.As(err, new(*exec.ExitError)) {
		log.Printf("ffmpeg error: %s", err)
	}
}

func (v *VideoBox) asyncReadVideo() {
	for {
		bufSize := videoBufSize
		i := v.videoBufIdx.Load()
		nextIdx := (i + 1) % int32(len(v.videoBuf))

		n, err := io.ReadFull(v.rawVideoStream, v.videoBuf[nextIdx][0:bufSize])
		if n == 0 || err == io.EOF {
			return
		} else if n != bufSize || err != nil {
			fmt.Printf("read rawVideoStream error: %d, %s\n", n, err)
			return
		}
		v.videoBufIdx.Store(nextIdx)
		v.lastFrame = time.Now()
	}
}

func (v *VideoBox) asyncReadAudio() {
	_, err := io.Copy(v.audioBuf, v.rawAudioStream)
	if err != nil {
		fmt.Printf("read rawAudioStream error: %s\n", err)
		return
	}
}

func (v *VideoBox) monitorLag() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if time.Since(v.lastFrame) >= 100*time.Millisecond {
				v.eventBus.Publish(eventLag, v)
			}
		case <-v.stopChan:
			return
		}
	}
}

func (v *VideoBox) HandleEvent(e sdl.Event) {}

func (v *VideoBox) Draw(renderer *sdl.Renderer) {
	v.once.Do(func() {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, 640, 480)
		if err != nil {
			panic(err)
		}
		v.texture = texture
	})

	v.texture.Update(nil, unsafe.Pointer(&v.videoBuf[v.videoBufIdx.Load()][0]), 640*3)
	renderer.Copy(v.texture, nil, nil)
}

func (v *VideoBox) Dispose() {
	vf.Set(nil)
	v.pinner.Unpin()
	v.rawVideoStream.Close()
	v.rawAudioStream.Close()
	v.texture.Destroy()
	if v.cmd != nil && v.cmd.Process != nil {
		v.cmd.Process.Kill()
	}
	close(v.stopChan)
}

func (v *VideoBox) SetVolume(volume int) {
	v.audioVolume = volume
}

func (v *VideoBox) OnVideoLag(handler model.EventHandler) {
	v.eventBus.Subscribe(eventLag, handler)
}

var _ component.Component = (*VideoBox)(nil)
