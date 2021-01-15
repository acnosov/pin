create or alter proc dbo.uspSaveEvents @LeagueId int, @TVP dbo.EventType READONLY as
begin
    set nocount on

    MERGE dbo.Event AS t
    USING @TVP s
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET ParentId      = s.ParentId,
            Starts        = s.Starts,
            Home          = s.Home,
            Away          = s.Away,
            RotNum        = s.RotNum,
            LiveStatus    = s.LiveStatus,
            HomePitcher   = s.HomePitcher,
            AwayPitcher   = s.AwayPitcher,
            ResultingUnit = s.ResultingUnit,
            LeagueId      = @LeagueId,
            UpdatedAt     = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id,
                LeagueId,
                ParentId,
                Starts,
                Home,
                Away,
                RotNum,
                LiveStatus,
                HomePitcher,
                AwayPitcher,
                ResultingUnit)
        VALUES (s.Id,
                @LeagueId,
                s.ParentId,
                s.Starts,
                s.Home,
                s.Away,
                s.RotNum,
                s.LiveStatus,
                s.HomePitcher,
                s.AwayPitcher,
                s.ResultingUnit);
end
