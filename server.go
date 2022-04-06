package govader_backend

//5 minutes
import (
	"github.com/jonreiter/govader"
	echo "github.com/labstack/echo/v4"
)

type RequestType struct {
	Text string `json:"text"`
}

type Handler struct {
	analyzer *govader.SentimentIntensityAnalyzer
}

func (h Handler) HandleGetRequest(c echo.Context) error {

	// var request RequestType
	// c.Bind(request)
	text := c.QueryParam("text")
	if text == "" {
		return c.JSON(400, map[string]string{"error": "?text= parameter is required"})
	}
	score := h.analyzer.PolarityScores(text)
	return c.JSON(200, score)
}

func (h Handler) HandlePostRequest(c echo.Context) error {
	request := new(RequestType)
	c.Bind(request)

	if request.Text == "" {
		return c.JSON(400, map[string]string{"error": "text required"})
	}
	score := h.analyzer.PolarityScores(request.Text)
	return c.JSON(200, score)
}

func Serve() {
	e := echo.New()
	handler := Handler{
		analyzer: govader.NewSentimentIntensityAnalyzer(),
	}
	e.GET("/", handler.HandleGetRequest)
	e.POST("/", handler.HandlePostRequest)
	e.Logger.Fatal(e.Start(":8080"))
}
