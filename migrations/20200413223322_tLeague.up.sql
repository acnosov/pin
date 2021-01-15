create table dbo.League
(
    Id                  int                                        not null,
    SportId             int                                        not null,
    Name                varchar(300)                               not null,
    HasOfferings        bit,
    HomeTeamType        varchar(300),
    Container           varchar(300),
    LeagueSpecialsCount int,
    EventSpecialsCount  int,
    EventCount          int,
    CreatedAt           datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt           datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_LeagueId primary key (Id),
)
create type dbo.LeagueType as table
(
    Id                  int          not null,
    Name                varchar(300) not null,
    HomeTeamType        varchar(300),
    HasOfferings        bit,
    Container           varchar(300),
    AllowRoundRobins    bit,
    LeagueSpecialsCount int,
    EventSpecialsCount  int,
    EventCount          int,
    primary key (Id)
)
