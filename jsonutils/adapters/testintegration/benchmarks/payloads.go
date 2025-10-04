// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import "strconv"

// The definition of these payloads is borrowed from github.com/goccy/go-json.

type SmallPayload struct {
	St   int    `json:"st"`
	Sid  int    `json:"sid"`
	Tt   string `json:"tt"`
	Gr   int    `json:"gr"`
	UUID string `json:"uuid"`
	IP   string `json:"ip"`
	Ua   string `json:"ua"`
	Tz   int    `json:"tz"`
	V    int    `json:"v"`
}

//nolint:mnd
func NewSmallPayload() *SmallPayload {
	return &SmallPayload{
		St:   1,
		Sid:  2,
		Tt:   "TestString",
		Gr:   4,
		UUID: "8f9a65eb-4807-4d57-b6e0-bda5d62f1429",
		IP:   "127.0.0.1",
		Ua:   "Mozilla",
		Tz:   8,
		V:    6,
	}
}

type MediumPayload struct {
	Person  *CBPerson `json:"person"`
	Company string    `json:"company"`
}

const mediumNumberOfAvatars = 8

//nolint:mnd
func NewMediumPayload() *MediumPayload {
	p := &MediumPayload{
		Company: "test",
		Person: &CBPerson{
			Name: &CBName{
				FullName: "test",
			},
			Github: &CBGithub{
				Followers: 100,
			},
			Gravatar: &CBGravatar{},
		},
	}

	avatars := make(Avatars, mediumNumberOfAvatars)
	for i := range mediumNumberOfAvatars {
		avatars[i] = &CBAvatar{
			URL: "http://test.com",
		}
	}
	p.Person.Gravatar.Avatars = avatars

	return p
}

type CBPerson struct {
	Name     *CBName     `json:"name"`
	Github   *CBGithub   `json:"github"`
	Gravatar *CBGravatar `json:"gravatar"`
}

type CBName struct {
	FullName string `json:"fullName"`
}

type CBGithub struct {
	Followers int `json:"followers"`
}

type CBGravatar struct {
	Avatars Avatars `json:"avatars"`
}

type Avatars []*CBAvatar

type CBAvatar struct {
	URL string `json:"url"`
}

type LargePayload struct {
	Users  DSUsers       `json:"users"`
	Topics *DSTopicsList `json:"topics"`
}

const largeNumberOfUsers = 100

func NewLargePayload() *LargePayload {
	dsUsers := make(DSUsers, largeNumberOfUsers)
	dsTopics := make(DSTopics, largeNumberOfUsers)
	for i := range largeNumberOfUsers {
		str := "test" + strconv.Itoa(i)
		dsUsers[i] = &DSUser{
			Username: str,
		}
		dsTopics[i] = &DSTopic{
			ID:   i,
			Slug: str,
		}
	}

	return &LargePayload{
		Users: dsUsers,
		Topics: &DSTopicsList{
			Topics:        dsTopics,
			MoreTopicsURL: "http://test.com",
		},
	}
}

type DSUser struct {
	Username string `json:"username"`
}

type DSUsers []*DSUser

type DSTopicsList struct {
	Topics        DSTopics `json:"topics"`
	MoreTopicsURL string   `json:"moreTopicsURL"`
}

type DSTopic struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type DSTopics []*DSTopic
