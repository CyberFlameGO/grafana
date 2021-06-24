package notifiers

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana/pkg/services/validations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

func TestReplaceIllegalCharswithUnderscore(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "foobar",
			expected: "foobar",
		},
		{
			input:    `foo.,\][!?#="~*^&+|<>\'bar09_09`,
			expected: "foo____________________bar09_09",
		},
	}

	for _, c := range cases {
		assert.Equal(t, replaceIllegalCharsInLabelname(c.input), c.expected)
	}
}

func TestWhenAlertManagerShouldNotify(t *testing.T) {
	tcs := []struct {
		prevState models.AlertStateType
		newState  models.AlertStateType

		expect bool
	}{
		{
			prevState: models.AlertStatePending,
			newState:  models.AlertStateOK,
			expect:    false,
		},
		{
			prevState: models.AlertStateAlerting,
			newState:  models.AlertStateOK,
			expect:    true,
		},
		{
			prevState: models.AlertStateOK,
			newState:  models.AlertStatePending,
			expect:    false,
		},
		{
			prevState: models.AlertStateUnknown,
			newState:  models.AlertStatePending,
			expect:    false,
		},
	}

	for _, tc := range tcs {
		am := &AlertmanagerNotifier{log: log.New("test.logger")}
		evalContext := alerting.NewEvalContext(context.Background(), &alerting.Rule{
			State: tc.prevState,
		}, &validations.OSSPluginRequestValidator{})

		evalContext.Rule.State = tc.newState

		res := am.ShouldNotify(context.TODO(), evalContext, &models.AlertNotificationState{})
		if res != tc.expect {
			t.Errorf("got %v expected %v", res, tc.expect)
		}
	}
}

//nolint:goconst
func TestAlertmanagerNotifier(t *testing.T) {
	t.Run("Alertmanager notifier tests", func(t *testing.T) {
		t.Run("Parsing alert notification from settings", func(t *testing.T) {
			t.Run("empty settings should return error", func(t *testing.T) {
				json := `{ }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &models.AlertNotification{
					Name:     "alertmanager",
					Type:     "alertmanager",
					Settings: settingsJSON,
				}

				_, err := NewAlertmanagerNotifier(model)
				require.Error(t, err)
			})

			t.Run("from settings", func(t *testing.T) {
				json := `{ "url": "http://127.0.0.1:9093/", "basicAuthUser": "user", "basicAuthPassword": "password" }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &models.AlertNotification{
					Name:     "alertmanager",
					Type:     "alertmanager",
					Settings: settingsJSON,
				}

				not, err := NewAlertmanagerNotifier(model)
				alertmanagerNotifier := not.(*AlertmanagerNotifier)

				require.NoError(t, err)
				require.Equal(t, "user", alertmanagerNotifier.BasicAuthUser)
				require.Equal(t, "password", alertmanagerNotifier.BasicAuthPassword)

				require.True(t, cmp.Equal(alertmanagerNotifier.URL, []string{"http://127.0.0.1:9093/"}))
			})

			t.Run("from settings with multiple alertmanager", func(t *testing.T) {
				json := `{ "url": "http://alertmanager1:9093,http://alertmanager2:9093" }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &models.AlertNotification{
					Name:     "alertmanager",
					Type:     "alertmanager",
					Settings: settingsJSON,
				}

				not, err := NewAlertmanagerNotifier(model)
				alertmanagerNotifier := not.(*AlertmanagerNotifier)

				require.NoError(t, err)
				require.True(t, cmp.Equal(alertmanagerNotifier.URL, []string{"http://alertmanager1:9093", "http://alertmanager2:9093"}))
			})
		})
	})
}
