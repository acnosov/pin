create table dbo.SpecialLine
(
    Id        int                                        not null,
    SpecialId int                                        not null,
    LineId    int,
    MaxBet    real,
    Price     real,
    Handicap  real,

    CreatedAt datetimeoffset default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset default sysdatetimeoffset() not null,

    constraint PK_SpecialLine_Id primary key (Id),
)

create type dbo.SpecialLineType as table
(
    Id        int not null,
    SpecialId int not null,
    LineId    int,
    MaxBet    real,
    Price     real,
    Handicap  real,
    primary key (Id)
)
