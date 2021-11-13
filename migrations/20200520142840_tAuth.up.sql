create table dbo.Auth
(
    Session     varchar(50)                                not null,
    ApiKey      varchar(50)                                not null,
    Device      varchar(50)                                not null,
    TrustCode   varchar(100)                               not null,
    CreatedAt   datetimeoffset default sysdatetimeoffset() not null,
    LastCheckAt datetimeoffset default sysdatetimeoffset() not null,
    constraint PK_Auth primary key (CreatedAt),
)
