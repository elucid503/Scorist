package models

import (
	"math/rand"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type ChannelType int

const (

	ChannelTypeUpdates ChannelType = iota

)

type UpdateMode string

const (

	UpdateModeFinalScores UpdateMode = "final_scores"
	UpdateModeScoreChanges UpdateMode = "score_changes"

)

type DelayMode string

const (

	Delay1Min DelayMode = "1_min"
	Delay5Min DelayMode = "5_min"
	Delay30Min DelayMode = "30_min"
	DelayRandom DelayMode = "random"

)

func (d DelayMode) Duration() time.Duration {

	switch d {

		case Delay1Min:

			return time.Minute

		case Delay5Min:

			return 5 * time.Minute

		case Delay30Min:

			return 30 * time.Minute

		case DelayRandom:

			return time.Minute + time.Duration(rand.Int63n(int64(29*time.Minute)))

		default:

			return time.Minute

	}

}

type Guild struct {

	ID snowflake.ID `bson:"_id"`

	Channels Channels `bson:"channels"`
	Preferences Preferences `bson:"preferences"`

}

type Channels struct {

	Updates snowflake.ID `bson:"updates"`

}

type Preferences struct {

	Updates UpdateMode `bson:"updates"`
	Delay DelayMode `bson:"delay"`

}

func NewGuild(ID snowflake.ID) *Guild {

	return &Guild{

		ID: ID,

		Channels: Channels{

			Updates: 0,

		},

		Preferences: Preferences{

			Updates: UpdateModeFinalScores,
			Delay: Delay1Min,

		},

	}

}

func (g *Guild) SetChannel(channelType ChannelType, channelID snowflake.ID) {

	switch channelType {

		case ChannelTypeUpdates:

			g.Channels.Updates = channelID

	}

}

func (g *Guild) SetUpdates(mode UpdateMode) {

	g.Preferences.Updates = mode

}

func (g *Guild) SetDelay(delay DelayMode) {

	g.Preferences.Delay = delay

}

func (g *Guild) Configured() bool {

	return g.Channels.Updates != 0

}