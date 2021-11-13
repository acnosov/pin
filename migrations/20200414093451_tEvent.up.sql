create table dbo.Event
(
    Id            int                                        not null,
    LeagueId      int                                        not null,
    ParentId      int,
    Starts        datetimeoffset,

    Home          varchar(300),
    Away          varchar(300),

    RotNum        varchar(300),
    LiveStatus    tinyint,

    HomePitcher   varchar(300),
    AwayPitcher   varchar(300),

    ResultingUnit varchar(300),

    CreatedAt     datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Event_Id primary key (Id),
);
create type dbo.EventType as table
(
    Id                int not null,
    ParentId          int,
    Starts            datetimeoffset,

    Home              varchar(300),
    Away              varchar(300),

    RotNum            varchar(300),
    LiveStatus        tinyint,

    HomePitcher       varchar(300),
    AwayPitcher       varchar(300),
    Status            varchar(300),
    ParlayRestriction int,
    AltTeaser         bit,

    ResultingUnit     varchar(300),
    primary key (Id)
)

--     // Event id.
-- 	Id *int64 `json:"id,omitempty" xml:"id"`
-- 	// If event is linked to another event, parentId will be populated.  Live event would have pre game event as parent id.
--  ParentId *int64 `json:"parentId,omitempty" xml:"parentId"`
-- 	// Start time of the event in UTC.
-- 	Starts *time.Time `json:"starts,omitempty" xml:"starts"`
-- 	// Home team name.
-- 	Home *string `json:"home,omitempty" xml:"home"`
-- 	// Away team name.
-- 	Away *string `json:"away,omitempty" xml:"away"`
-- 	// Team1 rotation number. Please note that in the next version of /fixtures, rotNum property will be decommissioned. ParentId can be used instead to group the related events.
-- 	RotNum *string `json:"rotNum,omitempty" xml:"rotNum"`
-- 	// Indicates live status of the event.   0 = No live betting will be offered on this event,  1 = Live betting event,  2 = Live betting will be offered on this match, but on a different event. Please note that [pre-game and live events are different](https://github.com/pinnacleapi/pinnacleapi-documentation/blob/master/FAQ.md#how-to-find-associated-events) .
-- 	LiveStatus *int32 `json:"liveStatus,omitempty" xml:"liveStatus"`
-- 	// Home team pitcher. Only for Baseball.
-- 	HomePitcher *string `json:"homePitcher,omitempty" xml:"homePitcher"`
-- 	// Away team pitcher. Only for Baseball.
-- 	AwayPitcher *string `json:"awayPitcher,omitempty" xml:"awayPitcher"`
-- 	// This is deprecated parameter, please check period's `status` in the `/odds` endpoint to see if it's open for betting.   O = This is the starting status of a game.    H = This status indicates that the lines are temporarily unavailable for betting,   I = This status indicates that one or more lines have a red circle (lower maximum bet amount).
-- 	Status *string `json:"status,omitempty" xml:"status"`
-- 	//  Parlay status of the event.   0 = Allowed to parlay, without restrictions,  1 = Not allowed to parlay this event,  2 = Allowed to parlay with the restrictions. You cannot have more than one leg from the same event in the parlay. All events with the same rotation number are treated as same event.
-- 	ParlayRestriction *int32 `json:"parlayRestriction,omitempty" xml:"parlayRestriction"`
-- 	// Whether an event is offer with alternative teaser points. Events with alternative teaser points may vary from teaser definition.
-- 	AltTeaser *bool `json:"altTeaser,omitempty" xml:"altTeaser"`
-- 	// Specifies based on what the event will be resulted, e.g. Corners, Bookings
-- 	ResultingUnit *string `json:"resultingUnit,omitempty" xml:"resultingUnit"`