create or alter proc dbo.uspSaveTotals @TVP dbo.TotalType READONLY as
begin
    set nocount on
    MERGE dbo.Total AS t
    USING @TVP s
    ON (t.LineId = s.LineId and t.Points = s.Points)

    WHEN MATCHED THEN
        UPDATE
        SET t.AltLineId = s.AltLineId,
            t.[Over]      = s.[Over],
            t.Under     = s.Under,
            t.UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (LineId,
                AltLineId,
                Points,
                [Over],
                Under)
        VALUES (s.LineId,
                s.AltLineId,
                s.Points,
                s.[Over],
                s.Under);
end
go

MERGE dbo.Total AS t
using (select 2, 2, 3, 4, 5) as source (LineId, AltLineId, Points, [Over], Under)
on t.LineId = source.LineId
WHEN MATCHED THEN update set t.UpdatedAt = sysdatetimeoffset()
WHEN NOT MATCHED THEN
    INSERT (LineId,
            AltLineId,
            Points,
            [Over],
            Under)
    VALUES (1,
            2,
            3,
            4,
            5)

    ;

-- INSERT into dbo.Total (LineId,
--                        AltLineId,
--                        Points,
--                        [Over],
--                        Under)
-- VALUES (1,
--         1,
--         1,
--         1,
--         1)