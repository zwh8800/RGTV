package channel_source

import (
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/model"
	"github.com/zwh8800/RGTV/text"
)

const (
	closeTimeout = 3 * time.Second
)

const (
	eventSourceChange = "ChannelSource:SourceChange"
)

const (
	maxTextWidth    = 6
	maxTextWidthCap = 8
	shadowWidth     = 3
)

var (
	shadowAlpha = [shadowWidth]int{50, 128, 180}
)

type ChannelSource struct {
	channel   *model.Channel
	sourceIdx int

	shown      bool
	closeTimer *time.Timer
	eventBus   evbus.Bus
}

func New(channel *model.Channel) *ChannelSource {
	return &ChannelSource{
		channel:  channel,
		eventBus: evbus.New(),
	}
}

func (c *ChannelSource) HandleEvent(e sdl.Event) {}
func (c *ChannelSource) Dispose()                {}

func (c *ChannelSource) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}

	var drawTexts [3]string
	for i := 0; i < 3; i++ {
		idx := c.sourceIdx + i - 1
		if idx < 0 && idx >= len(c.channel.Sources) {
			continue
		}
		source := []rune(c.channel.Sources[idx].Name)
		if len(source) > maxTextWidth {
			source = append(source[:maxTextWidth], '…', '…')
		}
		drawTexts[i] = string(source)

		// TODO: draw text
	}

}

func (c *ChannelSource) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(closeTimeout, func() {
		c.shown = false
	})
}

func (c *ChannelSource) SetChannel(channel *model.Channel) {
	c.channel = channel
	if c.sourceIdx >= len(c.channel.Sources) {
		c.sourceIdx = len(c.channel.Sources)
	} else if c.sourceIdx < 0 {
		c.sourceIdx = 0
	}
}

func (c *ChannelSource) NextSource() {
	c.sourceIdx++
	if c.sourceIdx >= len(c.channel.Sources) {
		c.sourceIdx = len(c.channel.Sources)
	}
	c.Show()
	c.eventBus.Publish(eventSourceChange, c)
}

func (c *ChannelSource) PrevSource() {
	c.sourceIdx--
	if c.sourceIdx < 0 {
		c.sourceIdx = 0
	}
	c.Show()
	c.eventBus.Publish(eventSourceChange, c)
}

func (c *ChannelSource) GetSource() *model.Source {
	if c.sourceIdx >= len(c.channel.Sources) {
		return nil
	}
	return c.channel.Sources[c.sourceIdx]
}

func (c *ChannelSource) OnSourceChange(f model.EventHandler) {
	c.eventBus.Subscribe(eventSourceChange, f)
}

var _ component.Component = (*ChannelSource)(nil)
