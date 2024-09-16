package httpserver

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/internal/container"
)

func SetupAsRecorderWithIDField(method, target, id, body, iDField string) (echo.Context, *httptest.ResponseRecorder) {
	mockServer := NewServer(container.Dependencies{})

	if strings.Contains(target, "?") {
		target = strings.Replace(target, "?", fmt.Sprintf("/%s?", id), 1)
	} else {
		target = fmt.Sprintf("%s/%s", target, id)
	}

	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	context := mockServer.NewServerContext(request, recorder)
	context.SetPath(target)

	if id != "" {
		context.SetPath(fmt.Sprintf("%s/%s", strings.TrimSuffix(target, "/"+id), fmt.Sprintf(":%s", iDField)))
		context.SetParamNames(iDField)
		context.SetParamValues(id)
	}

	return context, recorder
}

func SetupAsRecorderWithDynamicQueryParams(method, baseURL, userID string, queryParams map[string]string,
	body string) (echo.Context, *httptest.ResponseRecorder) {
	mockServer := NewServer(container.Dependencies{})

	queryString := "?"
	for key, value := range queryParams {
		if value != "" {
			queryString += fmt.Sprintf("%s=%s&", key, value)
		}
	}
	queryString = strings.TrimSuffix(queryString, "&")

	target := fmt.Sprintf("%s/%s", baseURL, userID)
	if len(queryParams) > 0 {
		target += queryString
	}

	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	context := mockServer.NewServerContext(request, recorder)
	context.SetPath(fmt.Sprintf("%s/:user_id", baseURL))
	context.SetParamNames("user_id")
	context.SetParamValues(userID)

	return context, recorder
}

func SetupAsRecorder(method, target, id, body string) (echo.Context, *httptest.ResponseRecorder) {
	return SetupAsRecorderWithIDField(method, target, id, body, "id")
}
