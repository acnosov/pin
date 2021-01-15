create or alter proc dbo.uspSaveSpecialLines @TVP dbo.SpecialLineType READONLY as
begin
    set nocount on

    MERGE dbo.SpecialLine AS t
    USING @TVP s
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET t.SpecialId = s.SpecialId,
            t.LineId    = s.LineId,
            t.MaxBet    = s.MaxBet,
            t.Price     = s.Price,
            t.Handicap  = s.Handicap,
            t.UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id,
                SpecialId,
                LineId,
                MaxBet,
                Price,
                Handicap)
        VALUES (s.Id,
                s.SpecialId,
                s.LineId,
                s.MaxBet,
                s.Price,
                s.Handicap);
end
