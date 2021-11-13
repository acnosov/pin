create or alter proc dbo.uspSaveSpreads @TVP dbo.SpreadType READONLY as
begin
    set nocount on

    MERGE dbo.Spread AS t
    USING @TVP s
    ON (t.LineId = s.LineId and t.Hdp = s.Hdp)

    WHEN MATCHED THEN
        UPDATE
        SET t.AltLineId = s.AltLineId,
            t.Home      = s.Home,
            t.Away      = s.Away,
            t.UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (LineId,
                AltLineId,
                Hdp,
                Home,
                Away)
        VALUES (s.LineId,
                s.AltLineId,
                s.Hdp,
                s.Home,
                s.Away);
end
