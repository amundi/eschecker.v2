package esslack

import (
	"github.com/amundi/escheck.v2/eslog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewSlackMsg(t *testing.T) {
	test := NewSlackMsg("salut", "Bobot", "#general")
	assert.NotNil(t, test)
	assert.Equal(t, "salut", test.text)
	assert.Equal(t, "Bobot", test.user)
	assert.Equal(t, "#general", test.channel)
}

func Test_getRequest(t *testing.T) {
	eslog.Init()
	token := "hellongig564"
	g_slack.token = token
	expectedRequest := "https://slack.com/api/chat.postMessage?username=test&token=" +
		token + "&channel=%23testchannel&pretty=1&text=Salut+les+copains+c%27est+moi"
	p := NewSlackMsg("Salut les copains c'est moi", "test", "#testchannel")
	assert.Equal(t, expectedRequest, p.GetSlackRequest())

	expectedRequest = "https://slack.com/api/chat.postMessage?username=Roberto&token=" +
		token + "&channel=%40lolo&pretty=1&text=Les+sanglots+longs+des+violons"
	p = NewSlackMsg("Les sanglots longs des violons", "Roberto", "@lolo")
	assert.Equal(t, expectedRequest, p.GetSlackRequest())
}
