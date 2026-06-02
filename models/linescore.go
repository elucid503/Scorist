package models

// URL: https://statsapi.mlb.com/api/v1/game/{gamePk}/linescore

type Linescore struct {

	GamePk int `json:"gamePk"` // Not supplied by API, but useful for internal tracking
	CreatedAt int `json:"appeared"` // Not supplied by API, but useful for internal tracking

	Copyright string `json:"copyright"`

	CurrentInning int `json:"currentInning"`
	CurrentInningOrdinal string `json:"currentInningOrdinal"`
	InningState string `json:"inningState"`
	InningHalf string `json:"inningHalf"`

	IsTopInning bool `json:"isTopInning"`

	ScheduledInnings int `json:"scheduledInnings"`

	Innings []Inning `json:"innings"`

	Teams LinescoreTeams `json:"teams"`

	Defense Player `json:"defense"`
	Offense Player `json:"offense"`

	Balls int `json:"balls"`
	Strikes int `json:"strikes"`
	Outs int `json:"outs"`

}

type Inning struct {

	Num int `json:"num"`
	OrdinalNum string `json:"ordinalNum"`

	Home InningStats `json:"home"`
	Away InningStats `json:"away"`

}

type InningStats struct {

	Runs int `json:"runs"`
	Hits int `json:"hits"`
	Errors int `json:"errors"`

	LeftOnBase int `json:"leftOnBase"`

}

type LinescoreTeams struct {

	Home InningStats `json:"home"`
	Away InningStats `json:"away"`

}

type Player struct {

	Pitcher Person `json:"pitcher"`
	Catcher Person `json:"catcher"`

	First Person `json:"first"`
	Second Person `json:"second"`
	Third Person `json:"third"`
	Shortstop Person `json:"shortstop"`

	Left Person `json:"left"`
	Center Person `json:"center"`
	Right Person `json:"right"`

	Batter Person `json:"batter"`
	OnDeck Person `json:"onDeck"`
	InHole Person `json:"inHole"`

	BattingOrder int `json:"battingOrder"`

	Team Team `json:"team"`

}

type Person struct {

	ID int `json:"id"`
	FullName string `json:"fullName"`
	Link string `json:"link"`

}
