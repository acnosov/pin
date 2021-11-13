create table dbo.Sport
(
    Id                  int                                        not null,
    Name                varchar(180)                               not null,
    HasOfferings        bit,
    LeagueSpecialsCount int                                        not null,
    EventSpecialsCount  int                                        not null,
    EventCount          int                                        not null,
    CreatedAt           datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt           datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_SportId primary key (Id),
)

create type dbo.SportType as table
(
    Id                  int          not null,
    Name                varchar(180) not null,
    HasOfferings        bit,
    LeagueSpecialsCount int,
    EventSpecialsCount  int,
    EventCount          int,
    primary key (Id)
)