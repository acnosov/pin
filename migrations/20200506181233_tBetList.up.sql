create table dbo.StraightBetList
(
    BetId              bigint                                     not null,
    WagerNumber        int,
    PlacedAt           datetimeoffset,
    BetStatus          varchar(1000),
    BetType            varchar(1000),
    Win                decimal(9, 5),
    Risk               decimal(9, 5),
    WinLoss            decimal(9, 5),
    OddsFormat         varchar(1000),
    CustomerCommission decimal(9, 5),
    CancellationReason varchar(1000),
    UpdateSequence     bigint,

    SportId            int,
    LeagueId           int,
    EventId            bigint,
    Handicap           decimal(9, 5),
    Price              decimal(9, 5),
    TeamName           varchar(1000),
    Side               varchar(1000),
    Pitcher1           varchar(1000),
    Pitcher2           varchar(1000),
    Pitcher1MustStart  bit,
    Pitcher2MustStart  bit,
    Team1              varchar(1000),
    Team2              varchar(1000),
    PeriodNumber       int,
    Team1Score         decimal(9, 5),
    Team2Score         decimal(9, 5),
    FtTeam1Score       decimal(9, 5),
    FtTeam2Score       decimal(9, 5),
    PTeam1Score        decimal(9, 5),
    PTeam2Score        decimal(9, 5),
    EventStartTime     datetimeoffset,
    CreatedAt          datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt          datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_StraightBetList primary key (BetId),
)
create type StraightBetListType as table
(
    BetId              bigint not null,
    WagerNumber        int,
    PlacedAt           datetimeoffset,
    BetStatus          varchar(1000),
    BetType            varchar(1000),
    Win                decimal(9, 5),

    Risk               decimal(9, 5),
    WinLoss            decimal(9, 5),

    OddsFormat         varchar(1000),
    CustomerCommission decimal(9, 5),
--     CancellationReason varchar(1000),
    UpdateSequence     bigint,
    SportId            int,
    LeagueId           int,
    EventId            bigint,
    Handicap           decimal(9, 5),
    Price              decimal(9, 5),
    TeamName           varchar(1000),
    Side               varchar(1000),
    Pitcher1           varchar(1000),
    Pitcher2           varchar(1000),
    Pitcher1MustStart  bit,
    Pitcher2MustStart  bit,
    Team1              varchar(1000),
    Team2              varchar(1000),
    PeriodNumber       int,
    Team1Score         decimal(9, 5),
    Team2Score         decimal(9, 5),
    FtTeam1Score       decimal(9, 5),
    FtTeam2Score       decimal(9, 5),
    PTeam1Score        decimal(9, 5),
    PTeam2Score        decimal(9, 5),
    EventStartTime     datetimeoffset,
    primary key (BetId)
)
