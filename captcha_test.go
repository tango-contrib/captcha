package captcha

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/macaron-contrib/cache"
	"github.com/tango-contrib/renders"

	"github.com/lunny/tango"
)

type CaptchaAction struct {
	Captcha
	renders.Renderer
}

func (c *CaptchaAction) Get() {
	c.Render("captcha.html", renders.T{
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
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	c, _ := cache.NewCacher("memory", cache.Options{
		Interval: 120,
	})
	tg.Use(New(Options{}, c), renders.New())
	tg.Any("/", new(CaptchaAction))

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
