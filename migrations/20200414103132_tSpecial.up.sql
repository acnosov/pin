create table dbo.Special
(
    Id         int                                        not null,
    LeagueId   int                                        not null,
    BetType    varchar(300),
    Name       varchar(300),
    Date       datetimeoffset,
    Cutoff     datetimeoffset,
    Category   varchar(300),
    Units      varchar(300),
    Status     varchar(300),
    LiveStatus tinyint,

    CreatedAt  datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt  datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Special_Id primary key (Id),
);
create table dbo.SpecialEvent
(
    Id           bigint                                        not null,
    PeriodNumber int,
    Home         varchar(600),
    Away         varchar(600),
    SpecialId    bigint                                        not null,
    CreatedAt    datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt    datetimeoffset default sysdatetimeoffset() not null,
    constraint PK_SpecialEvent_Id primary key (Id, SpecialId)
)
create type dbo.SpecialEventType as table
(
    Id           bigint not null,
    PeriodNumber int,
    Home         varchar(600),
    Away         varchar(600),
    SpecialId    bigint not null,
    primary key (Id, SpecialId)
)
create table dbo.SpecialContestant
(
    Id        int                                        not null,
    Name      varchar(300),
    RotNum    int,
    SpecialId int                                        not null,

    CreatedAt datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset default sysdatetimeoffset() not null,
    constraint PK_SpecialContestant_Id primary key (Id)
)
create type dbo.SpecialContestantType as table
(
    Id        int not null,
    Name      varchar(300),
    RotNum    int,
    SpecialId int not null,
    primary key (Id)
)

create type dbo.SpecialType as table
(
    Id         int not null,
    BetType    varchar(300),
    Name       varchar(300),
    Date       datetimeoffset,
    Cutoff     datetimeoffset,
    Category   varchar(300),
    Units      varchar(300),
    Status     varchar(300),
    LiveStatus tinyint,
    primary key (Id)
)
