package client

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/mikloslorinczi/pager/logger"
	"github.com/mikloslorinczi/pager/model"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/h2non/gentleman.v2"
)

var (
	defaultClient *gentleman.Client
	once          sync.Once
	log           = logger.NewLogger("HTTP Client")
)

// Setup sets up the defaultClient with the Pager Duty URL and Token provided by Viper
func Setup() {
	once.Do(func() {
		defaultClient = gentleman.New()
		defaultClient.URL(viper.GetString("api_url"))
		defaultClient.SetHeader("Accept", "application/json")
		defaultClient.SetHeader("Authorization", fmt.Sprintf("Token token=%s", viper.GetString("api_token")))
	})
}

// GET creates a new GET Request with the defaultClient and returns it
func GET() *gentleman.Request {
	return defaultClient.Get()
}

// GetUserID returns the user's ID by username or part of username
func GetUserID(query string) string {
	req := GET()
	req.Path("/users")
	res, err := req.Do()
	if err != nil {
		log.WithError(err).Fatal("HTTP Client failed to GET Response from PagerDuty API")
	}
	var usersResp model.UsersResponse
	if err := res.JSON(&usersResp); err != nil {
		log.WithError(err).Fatal("Failed to JSON Parse Response")
	}
	for _, user := range usersResp.Users {
		if query == user.ID {
			return query
		}
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) {
			return user.ID
		}
	}
	return ""
}

// GetAllUsers returns the list of all users
func GetAllUser() (model.UsersResponse, error) {
	var usersResp model.UsersResponse
	req := GET()
	req.Path("/users")
	res, err := req.Do()
	if err != nil {
		return usersResp, errors.Wrap(err, "HTTP Client failed to GET Response from PagerDuty API")
	}
	if !res.Ok {
		return usersResp, errors.Errorf("HTTP Error. Code: %d Body: %s", res.StatusCode, res.String())
	}
	if err := res.JSON(&usersResp); err != nil {
		return usersResp, errors.Wrap(err, "Failed to JSON Parse Response")
	}
	return usersResp, nil
}

// GetAllUserID returns all user IDs belonging to the team (ID or name or part of name)
func GetAllUserID(team string) []string {
	var foundIDs []string
	teamsResp, err := GetTeamsList()
	if err != nil {
		log.WithError(err).Fatal("Failed to get list of teams")
	}
	for _, aTeam := range teamsResp.Teams {
		if strings.Contains(strings.ToLower(aTeam.Name), strings.ToLower(team)) || team == aTeam.ID {
			membersResp, err := GetMembersList(aTeam.ID)
			if err != nil {
				log.WithError(err).Fatalf("Failed to get members of the team %s", aTeam.Name)
			}
			for _, member := range membersResp.Members {
				foundIDs = append(foundIDs, member.User.ID)
			}
			return foundIDs
		}
	}
	return nil
}

// GetTeamsList returns a list of all teams
func GetTeamsList() (model.TeamsResponse, error) {
	var teamsResp model.TeamsResponse
	req := GET()
	req.Path("/teams")
	res, err := req.Do()
	if err != nil {
		return teamsResp, errors.Wrap(err, "HTTP Client failed to GET Response from PagerDuty API")
	}
	if !res.Ok {
		return teamsResp, errors.Wrapf(err, "HTTP Error. Code: %d Body: %s", res.StatusCode, res.String())
	}
	if err := res.JSON(&teamsResp); err != nil {
		return teamsResp, errors.Wrap(err, "Failed to JSON Parse Response")
	}
	return teamsResp, nil
}

// GetMembersList returns all members of a team (by team ID)
func GetMembersList(teamID string) (model.MembersResponse, error) {
	var membersResp model.MembersResponse
	req := GET()
	req.Path(path.Join("/teams", teamID, "members"))
	res, err := req.Do()
	if err != nil {
		return membersResp, errors.Wrap(err, "HTTP Client failed to GET Response from PagerDuty API")
	}
	if !res.Ok {
		return membersResp, errors.Wrapf(err, "HTTP Error. Code: %d Body: %s", res.StatusCode, res.String())
	}
	if err := res.JSON(&membersResp); err != nil {
		return membersResp, errors.Wrap(err, "Failed to JSON Parse Response")
	}
	return membersResp, nil
}
