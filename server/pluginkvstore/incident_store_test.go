package pluginkvstore

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-plugin-incident-response/server/incident"
	mock_pluginkvstore "github.com/mattermost/mattermost-plugin-incident-response/server/pluginkvstore/mocks"
	mock_plugin "github.com/mattermost/mattermost-plugin-incident-response/server/pluginkvstore/mocks/serverpluginapi"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/require"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
)

var id1 = incident.Incident{
	Header: incident.Header{
		ID:              "id1",
		Name:            "incident 1 - wheel cat aliens wheelbarrow",
		IsActive:        true,
		CommanderUserID: "commander1",
		TeamID:          "team1",
		CreatedAt:       123,
		EndedAt:         440,
	},
}

var id2 = incident.Incident{
	Header: incident.Header{
		ID:              "id2",
		Name:            "incdnt 2 - horse staple battery shotgun mouse shotputmouse",
		IsActive:        true,
		CommanderUserID: "commander2",
		TeamID:          "team1",
		CreatedAt:       145,
		EndedAt:         555,
	},
}

var id3 = incident.Incident{
	Header: incident.Header{
		ID:              "id3",
		Name:            "incident 3 - Horse stapler battery shotgun mouse shotputmouse",
		IsActive:        false,
		CommanderUserID: "commander1",
		TeamID:          "team1",
		CreatedAt:       222,
		EndedAt:         666,
	},
}

var id4 = incident.Incident{
	Header: incident.Header{
		ID:              "id4",
		Name:            "incident 4 - titanic terminator aliens",
		IsActive:        false,
		CommanderUserID: "commander3",
		TeamID:          "team2",
		CreatedAt:       333,
		EndedAt:         444,
	},
}

var id5 = incident.Incident{
	Header: incident.Header{
		ID:              "id5",
		Name:            "incident 5 - ubik high castle electric sheep",
		IsActive:        true,
		CommanderUserID: "commander3",
		TeamID:          "team2",
		CreatedAt:       223,
		EndedAt:         550,
	},
}

var id6 = incident.Incident{
	Header: incident.Header{
		ID:              "id6",
		Name:            "incident 6 - ziggurat!",
		IsActive:        true,
		CommanderUserID: "commander5",
		TeamID:          "team3",
		CreatedAt:       555,
		EndedAt:         777,
	},
}

var id7 = incident.Incident{
	Header: incident.Header{
		ID:              "id7",
		Name:            "incident 7 - Ziggürat!",
		IsActive:        true,
		CommanderUserID: "commander5",
		TeamID:          "team3",
		CreatedAt:       556,
		EndedAt:         778,
	},
}

var dbHeaderMap = idHeaderMap{
	"id1": id1.Header,
	"id2": id2.Header,
	"id3": id3.Header,
	"id4": id4.Header,
	"id5": id5.Header,
	"id6": id6.Header,
	"id7": id7.Header,
}

func Test_incidentStore_GetIncidents(t *testing.T) {
	tests := []struct {
		name    string
		options incident.HeaderFilterOptions
		want    incident.GetIncidentsResults
		wantErr bool
	}{
		{
			name:    "simple getHeaders, no options",
			options: incident.HeaderFilterOptions{},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id7, id6, id4, id5, id3, id2, id1},
			},
		},
		{
			name: "team1 only, ascending",
			options: incident.HeaderFilterOptions{
				TeamID: "team1",
				Order:  incident.Asc,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 3,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id1, id2, id3},
			},
		},
		{
			name: "sort by ended_at",
			options: incident.HeaderFilterOptions{
				Sort: incident.EndedAt,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id7, id6, id3, id2, id5, id4, id1},
			},
		},
		{
			name: "no options, paged by 1",
			options: incident.HeaderFilterOptions{
				Page:    0,
				PerPage: 1,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  7,
				HasMore:    true,
				Items:      []incident.Incident{id7},
			},
		},
		{
			name: "no options, paged by 3",
			options: incident.HeaderFilterOptions{
				Page:    0,
				PerPage: 3,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  3,
				HasMore:    true,
				Items:      []incident.Incident{id7, id6, id4},
			},
		},
		{
			name: "no options, page 4 by 2",
			options: incident.HeaderFilterOptions{
				Page:    4,
				PerPage: 2,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  4,
				HasMore:    false,
				Items:      nil,
			},
		},
		{
			name: "no options, page 999 by 2",
			options: incident.HeaderFilterOptions{
				Page:    999,
				PerPage: 2,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  4,
				HasMore:    false,
				Items:      nil,
			},
		},
		{
			name: "no options, page 1 by 2",
			options: incident.HeaderFilterOptions{
				Page:    1,
				PerPage: 2,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  4,
				HasMore:    true,
				Items:      []incident.Incident{id4, id5},
			},
		},
		{
			name: "no options, page 1 by 3",
			options: incident.HeaderFilterOptions{
				Page:    1,
				PerPage: 3,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  3,
				HasMore:    true,
				Items:      []incident.Incident{id5, id3, id2},
			},
		},
		{
			name: "no options, page 1 by 5",
			options: incident.HeaderFilterOptions{
				Page:    1,
				PerPage: 5,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  2,
				HasMore:    false,
				Items:      []incident.Incident{id2, id1},
			},
		},
		{
			name: "sorted by ended, ascending, page 1 by 2",
			options: incident.HeaderFilterOptions{
				Sort:    incident.EndedAt,
				Order:   incident.Asc,
				Page:    1,
				PerPage: 2,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 7,
				PageCount:  4,
				HasMore:    true,
				Items:      []incident.Incident{id5, id2},
			},
		},
		{
			name: "only active, page 1 by 2",
			options: incident.HeaderFilterOptions{
				Page:    1,
				PerPage: 2,
				Status:  incident.Ongoing,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 5,
				PageCount:  3,
				HasMore:    true,
				Items:      []incident.Incident{id5, id2},
			},
		},
		{
			name: "active, commander3, asc",
			options: incident.HeaderFilterOptions{
				Status:      incident.Ongoing,
				CommanderID: "commander3",
				Order:       incident.Asc,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 1,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id5},
			},
		},
		{
			name: "commander1, asc, by ended_at",
			options: incident.HeaderFilterOptions{
				CommanderID: "commander1",
				Order:       incident.Asc,
				Sort:        incident.EndedAt,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 2,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id1, id3},
			},
		},
		{
			name: "search for horse",
			options: incident.HeaderFilterOptions{
				SearchTerm: "horse",
			},
			want: incident.GetIncidentsResults{
				TotalCount: 2,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id3, id2},
			},
		},
		{
			name: "search for aliens & commander3",
			options: incident.HeaderFilterOptions{
				CommanderID: "commander3",
				SearchTerm:  "aliens",
			},
			want: incident.GetIncidentsResults{
				TotalCount: 1,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id4},
			},
		},
		{
			name: "fuzzy search using starting characters -- not implemented",
			options: incident.HeaderFilterOptions{
				SearchTerm: "sbsm",
			},
			want: incident.GetIncidentsResults{
				TotalCount: 0,
				PageCount:  0,
				HasMore:    false,
				Items:      nil,
			},
		},
		{
			name: "fuzzy search using starting characters, active -- not implemented",
			options: incident.HeaderFilterOptions{
				SearchTerm: "sbsm",
				Status:     incident.Ongoing,
			},
			want: incident.GetIncidentsResults{
				TotalCount: 0,
				PageCount:  0,
				HasMore:    false,
				Items:      nil,
			},
		},
		{
			name: "case-insensitive and unicode-normalized",
			options: incident.HeaderFilterOptions{
				SearchTerm: "ziggurat",
			},
			want: incident.GetIncidentsResults{
				TotalCount: 2,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id7, id6},
			},
		},
		{
			name: "case-insensitive and unicode-normalized with unicode search term",
			options: incident.HeaderFilterOptions{
				SearchTerm: "ziggūràt",
			},
			want: incident.GetIncidentsResults{
				TotalCount: 2,
				PageCount:  1,
				HasMore:    false,
				Items:      []incident.Incident{id7, id6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			kvAPI := mock_pluginkvstore.NewMockKVAPI(mockCtrl)
			kvAPI.EXPECT().
				Get(IncidentHeadersKey, gomock.Any()).
				SetArg(1, dbHeaderMap).
				Times(1)

			for _, i := range []incident.Incident{id1, id2, id3, id4, id5, id6, id7} {
				kvAPI.EXPECT().
					Get(fmt.Sprintf(IncidentKey+"%s", i.ID), gomock.Any()).
					SetArg(1, i).
					AnyTimes()
			}

			s := &incidentStore{
				pluginAPI: PluginAPIClient{
					KV: kvAPI,
				},
			}
			got, err := s.GetIncidents(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIncidents() error = %v\nwantErr = %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUpdateHeaders(t *testing.T) {
	t.Run("Update empty headers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		pluginAPI := mock_plugin.NewMockAPI(mockCtrl)

		// Make KVGet return an empty value to simulate that the key is not set yet
		pluginAPI.EXPECT().
			KVGet(IncidentHeadersKey).
			Return([]byte{}, nil).
			Times(1)

		// Verify that KVSet is called to set the first value, proving that
		// SetAtomicWithRetries was called inside updateHeader
		value, err := json.Marshal(idHeaderMap{id1.ID: id1.Header})
		require.NoError(t, err)
		kvSetOptions := model.PluginKVSetOptions{
			Atomic:          true,
			OldValue:        nil,
			ExpireInSeconds: 0,
		}
		pluginAPI.EXPECT().
			KVSetWithOptions(IncidentHeadersKey, value, kvSetOptions).
			Return(true, nil).
			Times(1)

		// Set the wrapped plugin API client with the mocked underlying plugin API
		// and assign it to the store
		pluginAPIClient := pluginapi.NewClient(pluginAPI)
		s := &incidentStore{
			pluginAPI: PluginAPIClient{
				KV: &pluginAPIClient.KV,
			},
		}

		err = s.updateHeader(&id1)
		require.NoError(t, err)
	})
}
