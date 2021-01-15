create or alter proc dbo.uspSaveSpecials @LeagueId int, @TVP dbo.SpecialType READONLY as
begin
    set nocount on

    MERGE dbo.Special AS t
    USING @TVP s
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET BetType    = s.BetType,
            Name       = s.Name,
            Date       = s.Date,
            Cutoff     = s.Cutoff,
            Category   = s.Category,
            Units      = s.Units,
            Status     = s.Status,
            LiveStatus = s.LiveStatus,
            LeagueId   = @LeagueId,
            UpdatedAt  = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id,
                LeagueId,
                BetType,
                Name,
                Date,
                Cutoff,
                Category,
                Units,
                Status,
                LiveStatus)
        VALUES (s.Id,
                @LeagueId,
                s.BetType,
                s.Name,
                s.Date,
                s.Cutoff,
                s.Category,
                s.Units,
                s.Status,
                s.LiveStatus);
end
