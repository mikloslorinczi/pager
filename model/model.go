package model

import (
	"time"
)

type Team struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type User struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Teams []Team `json:"teams"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

type TTeam struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type TeamsResponse struct {
	Teams []TTeam `json:"teams"`
}

type EscalationPolicy struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type UUser struct {
	Summary string `json:"summary"`
	ID      string `json:"id"`
}

type Schedule struct {
	Summary string `json:"summary"`
	ID      string `json:"id"`
}

type Oncall struct {
	EscalationPolicy EscalationPolicy `json:"escalation_policy"`
	EscalationLevel  int              `json:"escalation_level"`
	Schedule         Schedule         `json:"schedule"`
	User             UUser            `json:"user"`
	Start            time.Time        `json:"start"`
	End              time.Time        `json:"end"`
}

type OncallsResponse struct {
	Oncalls []Oncall `json:"oncalls"`
}
