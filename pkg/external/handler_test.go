package external_test

import (
	"context"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/external"
	"github.com/sisukasco/henki/pkg/service"
	dtesting "github.com/sisukasco/henki/pkg/testing"
	"net/url"
	"testing"
)

var (
	svc *service.Service
)

func TestMain(m *testing.M) {
	dtesting.InitService(m, func(s *service.Service) {
		svc = s
	}, func() {})
}

func TestCreatingRedirectURL(t *testing.T) {
	url, err := external.CreateRedirectURL(svc.Konf, "google")
	if err != nil {
		t.Errorf("Error Creating redirect URL %v ", err)
		return
	}
	t.Logf("Google login redirect URL %s", url)
}

func TestParsingStateJwtToken(t *testing.T) {
	strUrl, err := external.CreateRedirectURL(svc.Konf, "google")
	if err != nil {
		t.Errorf("Error Creating redirect URL %v ", err)
		return
	}
	urlObj, err := url.Parse(strUrl)
	if err != nil {
		t.Errorf("Can't parse Redirect URL %v ", err)
		return
	}

	strState := urlObj.Query().Get("state")

	claims, err := external.DecodeJwtClaims(context.Background(), svc.Konf, strState)
	if err != nil {
		t.Errorf("Can't parse Jwt claims %v ", err)
	}

	if claims.Provider != "google" {
		t.Errorf("External JWT claim parsing didn't work")
	}

	t.Logf(" jwt parsed %v", utils.ToJSONString(claims))

}
