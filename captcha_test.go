package captcha

import (
	"testing"

	"github.com/tango-contrib/renders"
	"github.com/macaron-contrib/cache"

	"github.com/lunny/tango"
)

type CaptchaAction struct {
	Captcha
	renders.Renderer
}

func (c *CaptchaAction) Get() {
	c.Renderer.Render("captcha.html", renders.T{
		"captcha": c.CreateHtml(),
	})
}

func (c *CaptchaAction) Post() string {
	if c.Verify() {
		return "true"
	}
	return "false"
}

func TestCaptcha(t *testing.T) {
	tg := tango.Classic()
	c, _ := cache.NewCache("memory", `{"interval":120}`)
	tg.Use(New(Options{}, c), renders.New())
	tg.Any("/", new(CaptchaAction))
	tg.Run()
}
