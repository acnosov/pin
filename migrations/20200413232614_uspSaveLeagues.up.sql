create or alter proc dbo.uspSaveLeagues @SportId int, @TVP dbo.LeagueType READONLY as
begin
    set nocount on

    MERGE dbo.League AS t
    USING @TVP s
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET Name                = s.Name,
            HasOfferings        = s.HasOfferings,
            HomeTeamType        = s.HomeTeamType,
            Container           = s.Container,
            LeagueSpecialsCount = s.LeagueSpecialsCount,
            EventSpecialsCount  = s.EventSpecialsCount,
            EventCount          = s.EventCount,
            SportId             = @SportId,
            UpdatedAt           = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id, SportId, Name, HasOfferings, HomeTeamType, Container, LeagueSpecialsCount,
                EventSpecialsCount, EventCount)
        VALUES (s.Id, @SportId, s.Name, s.HasOfferings, s.HomeTeamType, s.Container, s.LeagueSpecialsCount,
                s.EventSpecialsCount, s.EventCount);
end