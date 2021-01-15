create table dbo.Line
(
    LineId       int                                        not null,
    EventId      int                                        not null,
    Number       int,
    Cutoff       datetimeoffset,
    Status       tinyint,

    MaxSpread    real,
    MaxMoneyline real,
    MaxTotal     real,
    MaxTeamTotal real,

    CreatedAt    datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt    datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Line_Id primary key (LineId),
)
create table dbo.Moneyline
(
    LineId    int                                        not null,
    Home      real,
    Away      real,
    Draw      real,
    CreatedAt datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Moneyline_Id primary key (LineId),
)
create table dbo.Total
(
    LineId    int                                        not null,
    AltLineId bigint,
    Points    decimal(9, 6)                              not null,
    [Over]    real,
    Under     real,
    CreatedAt datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Total primary key (LineId, Points),
)
create table dbo.Spread
(
    LineId    int                                        not null,
    AltLineId bigint,
    Hdp       real                                       not null,
    Home      real,
    Away      real,
    CreatedAt datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Spread primary key (LineId, Hdp),
)

create type dbo.LineType as table
(
    LineId       int not null,
    EventId      int not null,
    Number       int,
    Cutoff       datetimeoffset,
    Status       tinyint,

    MaxSpread    real,
    MaxMoneyline real,
    MaxTotal     real,
    MaxTeamTotal real,
    primary key (LineId)
)
create type dbo.MoneylineType as table
(
    LineId int not null,
    Home   real,
    Away   real,
    Draw   real,
    primary key (LineId)
)
create type dbo.TotalType as table
(
    LineId    int           not null,
    AltLineId bigint,
    Points    decimal(9, 6) not null,
    [Over]    real,
    Under     real,
    primary key (LineId, Points)
)
create type dbo.SpreadType as table
(
    LineId    int  not null,
    AltLineId bigint,
    Hdp       real not null,
    Home      real,
    Away      real,
    primary key (LineId, Hdp)
)
