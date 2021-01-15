create table dbo.Period
(
    Number                     int                                        not null,
    SportId                    int                                        not null,
    Description                varchar(800),
    ShortDescription           varchar(800),
    SpreadDescription          varchar(800),
    MoneylineDescription       varchar(800),
    TotalDescription           varchar(800),
    Team1TotalDescription      varchar(800),
    Team2TotalDescription      varchar(800),
    SpreadShortDescription     varchar(800),
    MoneylineShortDescription  varchar(800),
    TotalShortDescription      varchar(800),
    Team1TotalShortDescription varchar(800),
    Team2TotalShortDescription varchar(800),

    CreatedAt                  datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt                  datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_Period primary key (number, sportId),
)
create type dbo.PeriodType as table
(
    Number                     int not null,
    Description                varchar(800),
    ShortDescription           varchar(800),
    SpreadDescription          varchar(800),
    MoneylineDescription       varchar(800),
    TotalDescription           varchar(800),
    Team1TotalDescription      varchar(800),
    Team2TotalDescription      varchar(800),
    SpreadShortDescription     varchar(800),
    MoneylineShortDescription  varchar(800),
    TotalShortDescription      varchar(800),
    Team1TotalShortDescription varchar(800),
    Team2TotalShortDescription varchar(800)
--     primary key (Id)
)
