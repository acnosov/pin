create or alter proc dbo.uspCreateSport @Id int, @Name varchar(180), @HasOfferings bit, @LeagueSpecialsCount int,
                                        @EventSpecialsCount int, @EventCount int as
begin
    set nocount on

    MERGE dbo.Sport AS target
    USING (SELECT @Id,
                  @Name,
                  @HasOfferings,
                  @LeagueSpecialsCount,
                  @EventSpecialsCount,
                  @EventCount) AS source (Id, Name, HasOfferings, LeagueSpecialsCount, EventSpecialsCount, EventCount)
    ON (target.Id = source.Id)

    WHEN MATCHED THEN
        UPDATE
        SET Name                = source.Name,
            HasOfferings        = source.HasOfferings,
            LeagueSpecialsCount = source.LeagueSpecialsCount,
            EventSpecialsCount  = source.EventSpecialsCount,
            EventCount          = source.EventCount

    WHEN NOT MATCHED THEN
        INSERT (Id, Name, HasOfferings, LeagueSpecialsCount, EventSpecialsCount, EventCount)
        VALUES (@Id, @Name, @HasOfferings, @LeagueSpecialsCount, @EventSpecialsCount, @EventCount);
end
go

create or alter proc dbo.uspCreateSportTVP @TVP dbo.SportType READONLY as
begin
    set nocount on

    MERGE dbo.Sport AS t
    USING @TVP s
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET Name                = s.Name,
            HasOfferings        = s.HasOfferings,
            LeagueSpecialsCount = s.LeagueSpecialsCount,
            EventSpecialsCount  = s.EventSpecialsCount,
            EventCount          = s.EventCount,
            UpdatedAt           =sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id, Name, HasOfferings, LeagueSpecialsCount, EventSpecialsCount, EventCount)
        VALUES (s.Id, s.Name, s.HasOfferings, s.LeagueSpecialsCount, s.EventSpecialsCount, s.EventCount);
end

